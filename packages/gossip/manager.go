package gossip

import (
	"context"
	"fmt"
	"net"
	"runtime"
	"sync"

	"github.com/cockroachdb/errors"
	"github.com/iotaledger/hive.go/autopeering/peer"
	"github.com/iotaledger/hive.go/events"
	"github.com/iotaledger/hive.go/identity"
	"github.com/iotaledger/hive.go/logger"
	"github.com/iotaledger/hive.go/workerpool"
	"google.golang.org/protobuf/proto"

	pb "github.com/iotaledger/goshimmer/packages/gossip/proto"
	"github.com/iotaledger/goshimmer/packages/gossip/server"
	"github.com/iotaledger/goshimmer/packages/tangle"
)

const (
	// maxPacketSize defines the maximum packet size allowed for gossip and bufferedconn.
	maxPacketSize = 65 * 1024
)

var (
	messageWorkerCount     = runtime.GOMAXPROCS(0) * 4
	messageWorkerQueueSize = 1000

	messageRequestWorkerCount     = runtime.GOMAXPROCS(0)
	messageRequestWorkerQueueSize = 100
)

// LoadMessageFunc defines a function that returns the message for the given id.
type LoadMessageFunc func(messageId tangle.MessageID) ([]byte, error)

// The Manager handles the connected neighbors.
type Manager struct {
	local           *peer.Local
	loadMessageFunc LoadMessageFunc
	log             *logger.Logger
	events          Events
	neighborsEvents map[NeighborsGroup]NeighborsEvents

	wg sync.WaitGroup

	server      *server.TCP
	serverMutex sync.RWMutex

	neighbors         map[identity.ID]*Neighbor
	neighborsMapMutex sync.RWMutex

	// neighborConnectionMutexes is an ever-growing map where the key is a neighbor ID and the value is a mutex.
	// We create a new record in that map for every new unique neighbor and never delete it.
	// Taking in account that amount of possible unique neighbors is pretty small,
	// the size on that map shouldn't be an issue.
	// We need this map to acquire a lock for the particular neighbor ID only during the connection and dropping,
	// because otherwise long-running, hanging operations like accept will block the whole manager.
	neighborConnectionMutexes sync.Map

	// messageWorkerPool defines a worker pool where all incoming messages are processed.
	messageWorkerPool *workerpool.WorkerPool

	messageRequestWorkerPool *workerpool.WorkerPool
}

// NewManager creates a new Manager.
func NewManager(local *peer.Local, f LoadMessageFunc, log *logger.Logger) *Manager {
	m := &Manager{
		local:           local,
		loadMessageFunc: f,
		log:             log,
		events: Events{
			MessageReceived: events.NewEvent(messageReceived),
		},
		neighborsEvents: map[NeighborsGroup]NeighborsEvents{
			NeighborsGroupAuto:   NewNeighborsEvents(),
			NeighborsGroupManual: NewNeighborsEvents(),
		},
		neighbors: map[identity.ID]*Neighbor{},
		server:    nil,
	}

	m.messageWorkerPool = workerpool.New(func(task workerpool.Task) {
		m.processPacketMessage(task.Param(0).([]byte), task.Param(1).(*Neighbor))

		task.Return(nil)
	}, workerpool.WorkerCount(messageWorkerCount), workerpool.QueueSize(messageWorkerQueueSize))

	m.messageRequestWorkerPool = workerpool.New(func(task workerpool.Task) {
		m.processMessageRequest(task.Param(0).([]byte), task.Param(1).(*Neighbor))

		task.Return(nil)
	}, workerpool.WorkerCount(messageRequestWorkerCount), workerpool.QueueSize(messageRequestWorkerQueueSize))

	return m
}

// Start starts the manager for the given TCP server.
func (m *Manager) Start(srv *server.TCP) {
	m.serverMutex.Lock()
	defer m.serverMutex.Unlock()

	m.server = srv

	m.messageWorkerPool.Start()
	m.messageRequestWorkerPool.Start()
}

// Stop stops the manager and closes all established connections.
func (m *Manager) Stop() {
	m.serverMutex.Lock()
	defer m.serverMutex.Unlock()

	m.server = nil

	m.dropAllNeighbors()

	m.messageWorkerPool.Stop()
	m.messageRequestWorkerPool.Stop()
}

func (m *Manager) dropAllNeighbors() {
	neighborsList := m.AllNeighbors()
	for _, nbr := range neighborsList {
		_ = m.dropNeighbor(nbr.ID(), nil /* group */)
	}
}

// Events returns the events related to the gossip protocol.
func (m *Manager) Events() Events {
	return m.events
}

// NeighborsEvents returns the events related to the gossip protocol.
func (m *Manager) NeighborsEvents(group NeighborsGroup) NeighborsEvents {
	return m.neighborsEvents[group]
}

// AddOutbound tries to add a neighbor by connecting to that peer.
func (m *Manager) AddOutbound(ctx context.Context, p *peer.Peer, group NeighborsGroup,
	connectOpts ...server.ConnectPeerOption) error {
	return m.addNeighbor(ctx, p, group, m.server.DialPeer, connectOpts)
}

// AddInbound tries to add a neighbor by accepting an incoming connection from that peer.
func (m *Manager) AddInbound(ctx context.Context, p *peer.Peer, group NeighborsGroup,
	connectOpts ...server.ConnectPeerOption) error {
	return m.addNeighbor(ctx, p, group, m.server.AcceptPeer, connectOpts)
}

// DropNeighbor disconnects the neighbor with the given ID.
func (m *Manager) DropNeighbor(id identity.ID, group NeighborsGroup) error {
	return m.dropNeighbor(id, &group)
}

func (m *Manager) dropNeighbor(id identity.ID, group *NeighborsGroup) error {
	nbrConnectionMutex := m.getNeighborConnectionMutex(id)
	nbrConnectionMutex.Lock()
	defer nbrConnectionMutex.Unlock()
	nbr, err := m.getAndDeleteNeighbor(id, group)
	if err != nil {
		return errors.WithStack(err)
	}

	return nbr.Close()
}

func (m *Manager) getAndDeleteNeighbor(id identity.ID, group *NeighborsGroup) (*Neighbor, error) {
	m.neighborsMapMutex.Lock()
	defer m.neighborsMapMutex.Unlock()
	nbr, ok := m.neighbors[id]
	if !ok || (group != nil && nbr.Group != *group) {
		return nil, ErrUnknownNeighbor
	}
	delete(m.neighbors, id)
	return nbr, nil
}

// RequestMessage requests the message with the given id from the neighbors.
// If no peer is provided, all neighbors are queried.
func (m *Manager) RequestMessage(messageID []byte, to ...identity.ID) {
	msgReq := &pb.MessageRequest{Id: messageID}
	m.send(marshal(msgReq), to...)
}

// SendMessage adds the given message the send queue of the neighbors.
// The actual send then happens asynchronously. If no peer is provided, it is send to all neighbors.
func (m *Manager) SendMessage(msgData []byte, to ...identity.ID) {
	msg := &pb.Message{Data: msgData}
	m.send(marshal(msg), to...)
}

// AllNeighbors returns all the neighbors that are currently connected.
func (m *Manager) AllNeighbors() []*Neighbor {
	m.neighborsMapMutex.RLock()
	defer m.neighborsMapMutex.RUnlock()
	result := make([]*Neighbor, 0, len(m.neighbors))
	for _, n := range m.neighbors {
		result = append(result, n)
	}
	return result
}

func (m *Manager) getNeighbors(ids ...identity.ID) []*Neighbor {
	if len(ids) > 0 {
		return m.getNeighborsByID(ids)
	}
	return m.AllNeighbors()
}

func (m *Manager) getNeighborsByID(ids []identity.ID) []*Neighbor {
	result := make([]*Neighbor, 0, len(ids))
	m.neighborsMapMutex.RLock()
	defer m.neighborsMapMutex.RUnlock()
	for _, id := range ids {
		if n, ok := m.neighbors[id]; ok {
			result = append(result, n)
		}
	}
	return result
}

func (m *Manager) send(b []byte, to ...identity.ID) {
	neighbors := m.getNeighbors(to...)

	for _, nbr := range neighbors {
		if _, err := nbr.Write(b); err != nil {
			m.log.Warnw("send error", "peer-id", nbr.ID(), "err", err)
		}
	}
}

func (m *Manager) addNeighbor(ctx context.Context, p *peer.Peer, group NeighborsGroup,
	connectorFunc func(context.Context, *peer.Peer, ...server.ConnectPeerOption) (net.Conn, error),
	connectOpts []server.ConnectPeerOption,
) error {
	if p.ID() == m.local.ID() {
		return ErrLoopbackNeighbor
	}
	m.serverMutex.RLock()
	defer m.serverMutex.RUnlock()
	if m.server == nil {
		return ErrNotRunning
	}
	nbrConnectionMutex := m.getNeighborConnectionMutex(p.ID())
	nbrConnectionMutex.Lock()
	defer nbrConnectionMutex.Unlock()
	if m.neighborExists(p.ID()) {
		m.neighborsEvents[group].ConnectionFailed.Trigger(p, ErrDuplicateNeighbor)
		return ErrDuplicateNeighbor
	}

	conn, err := connectorFunc(ctx, p, connectOpts...)
	if err != nil {
		m.neighborsEvents[group].ConnectionFailed.Trigger(p, err)
		return err
	}

	// create and add the neighbor
	nbr := NewNeighbor(p, group, conn, m.log)
	nbr.Events.Close.Attach(events.NewClosure(func() {
		// assure that the neighbor is removed and notify
		_ = m.DropNeighbor(p.ID(), group)
		m.neighborsEvents[group].NeighborRemoved.Trigger(nbr)
	}))
	nbr.Events.ReceiveMessage.Attach(events.NewClosure(func(data []byte) {
		dataCopy := make([]byte, len(data))
		copy(dataCopy, data)
		if err := m.handlePacket(dataCopy, nbr); err != nil {
			m.log.Debugw("error handling packet", "err", err)
		}
	}))

	m.setNeighbor(nbr)

	nbr.Listen()
	m.neighborsEvents[group].NeighborAdded.Trigger(nbr)

	return nil
}

func (m *Manager) neighborExists(id identity.ID) bool {
	m.neighborsMapMutex.RLock()
	defer m.neighborsMapMutex.RUnlock()
	_, ok := m.neighbors[id]
	return ok
}

func (m *Manager) setNeighbor(neighbor *Neighbor) {
	m.neighborsMapMutex.Lock()
	defer m.neighborsMapMutex.Unlock()
	m.neighbors[neighbor.ID()] = neighbor
}

func (m *Manager) getNeighborConnectionMutex(id identity.ID) *sync.Mutex {
	newMutex := &sync.Mutex{}
	value, _ := m.neighborConnectionMutexes.LoadOrStore(id, newMutex)
	mutex := value.(*sync.Mutex)
	return mutex
}

func (m *Manager) handlePacket(data []byte, nbr *Neighbor) error {
	// ignore empty packages
	if len(data) == 0 {
		return nil
	}

	switch pb.PacketType(data[0]) {
	case pb.PacketMessage:
		if _, added := m.messageWorkerPool.TrySubmit(data, nbr); !added {
			return fmt.Errorf("messageWorkerPool full: packet message discarded")
		}
	case pb.PacketMessageRequest:
		if _, added := m.messageRequestWorkerPool.TrySubmit(data, nbr); !added {
			return fmt.Errorf("messageRequestWorkerPool full: message request discarded")
		}

	default:
		return ErrInvalidPacket
	}

	return nil
}

func marshal(packet pb.Packet) []byte {
	packetType := packet.Type()
	if packetType > 0xFF {
		panic("invalid packet")
	}

	data, err := proto.Marshal(packet)
	if err != nil {
		panic("invalid packet")
	}
	return append([]byte{byte(packetType)}, data...)
}

// MessageWorkerPoolStatus returns the name and the load of the workerpool.
func (m *Manager) MessageWorkerPoolStatus() (name string, load int) {
	return "messageWorkerPool", m.messageWorkerPool.GetPendingQueueSize()
}

// MessageRequestWorkerPoolStatus returns the name and the load of the workerpool.
func (m *Manager) MessageRequestWorkerPoolStatus() (name string, load int) {
	return "messageRequestWorkerPool", m.messageRequestWorkerPool.GetPendingQueueSize()
}

func (m *Manager) processPacketMessage(data []byte, nbr *Neighbor) {
	packet := new(pb.Message)
	if err := proto.Unmarshal(data[1:], packet); err != nil {
		m.log.Debugw("error processing packet", "err", err)
	}
	m.events.MessageReceived.Trigger(&MessageReceivedEvent{Data: packet.GetData(), Peer: nbr.Peer})
}

func (m *Manager) processMessageRequest(data []byte, nbr *Neighbor) {
	packet := new(pb.MessageRequest)
	if err := proto.Unmarshal(data[1:], packet); err != nil {
		m.log.Debugw("invalid packet", "err", err)
	}

	msgID, _, err := tangle.MessageIDFromBytes(packet.GetId())
	if err != nil {
		m.log.Debugw("invalid message id:", "err", err)
	}

	msgBytes, err := m.loadMessageFunc(msgID)
	if err != nil {
		m.log.Debugw("error loading message", "msg-id", msgID, "err", err)
	}

	// send the loaded message directly to the neighbor
	_, _ = nbr.Write(marshal(&pb.Message{Data: msgBytes}))
}
