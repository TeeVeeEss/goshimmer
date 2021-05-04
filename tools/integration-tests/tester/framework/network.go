package framework

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/iotaledger/hive.go/crypto/ed25519"
	"github.com/iotaledger/hive.go/identity"
	"github.com/mr-tron/base58"

	walletseed "github.com/iotaledger/goshimmer/client/wallet/packages/seed"
)

// Network represents a complete GoShimmer network within Docker.
// Including an entry node and arbitrary many peers.
type Network struct {
	id   string
	name string

	peers  []*Peer
	tester *DockerContainer

	entryNode         *DockerContainer
	entryNodeIdentity *identity.Identity

	partitions []*Partition

	dockerClient *client.Client
}

// newNetwork returns a Network instance, creates its underlying Docker network and adds the tester container to the network.
func newNetwork(dockerClient *client.Client, name string, tester *DockerContainer) (*Network, error) {
	// create Docker network
	resp, err := dockerClient.NetworkCreate(context.Background(), name, types.NetworkCreate{})
	if err != nil {
		return nil, err
	}

	// the tester container needs to join the Docker network in order to communicate with the peers
	err = tester.ConnectToNetwork(resp.ID)
	if err != nil {
		return nil, err
	}

	return &Network{
		id:           resp.ID,
		name:         name,
		tester:       tester,
		dockerClient: dockerClient,
	}, nil
}

// createEntryNode creates the network's entry node.
func (n *Network) createEntryNode() error {
	// create identity
	publicKey, privateKey, err := ed25519.GenerateKey()
	if err != nil {
		return err
	}

	n.entryNodeIdentity = identity.New(publicKey)
	seed := privateKey.Seed().String()

	// create entry node container
	n.entryNode = NewDockerContainer(n.dockerClient)
	err = n.entryNode.CreateGoShimmerEntryNode(n.namePrefix(containerNameEntryNode), seed)
	if err != nil {
		return err
	}
	err = n.entryNode.ConnectToNetwork(n.id)
	if err != nil {
		return err
	}
	err = n.entryNode.Start()
	if err != nil {
		return err
	}

	return nil
}

// CreatePeer creates a new peer/GoShimmer node in the network and returns it.
// Passing bootstrap true enables the bootstrap plugin on the given peer.
func (n *Network) CreatePeer(c GoShimmerConfig) (*Peer, error) {
	name := n.namePrefix(fmt.Sprintf("%s%d", containerNameReplica, len(n.peers)))
	config := c

	// create identity
	var publicKey ed25519.PublicKey
	var privateKey ed25519.PrivateKey
	var err error
	if config.Seed == "" {
		publicKey, privateKey, err = ed25519.GenerateKey()
		if err != nil {
			return nil, err
		}
		seed := privateKey.Seed().String()
		config.Seed = seed
	} else {
		bytes, encodeErr := base58.Decode(config.Seed)
		if encodeErr != nil {
			return nil, encodeErr
		}
		publicKey = ed25519.PrivateKeyFromSeed(bytes).Public()
	}

	config.Name = name
	config.EntryNodeHost = n.namePrefix(containerNameEntryNode)
	config.EntryNodePublicKey = n.entryNodePublicKey()
	config.DisabledPlugins = func() string {
		if !config.SyncBeaconFollower {
			return disabledPluginsPeer + ",SyncBeaconFollower"
		}
		return disabledPluginsPeer
	}()
	config.SnapshotFilePath = snapshotFilePath
	if config.SyncBeaconFollowNodes == "" {
		config.SyncBeaconFollowNodes = syncBeaconPublicKey
	}
	if config.SyncBeaconBroadcastInterval == 0 {
		config.SyncBeaconBroadcastInterval = 5
	}
	if config.FPCRoundInterval == 0 {
		config.FPCRoundInterval = 5
	}
	if config.FPCTotalRoundsFinalization == 0 {
		config.FPCTotalRoundsFinalization = 10
	}

	// create wallet
	var nodeSeed *walletseed.Seed
	if c.Faucet == true {
		nodeSeed = walletseed.NewSeed(genesisSeed)
	} else {
		nodeSeed = walletseed.NewSeed()
	}

	// create Docker container
	container := NewDockerContainer(n.dockerClient)
	err = container.CreateGoShimmerPeer(config)
	if err != nil {
		return nil, err
	}
	err = container.ConnectToNetwork(n.id)
	if err != nil {
		return nil, err
	}
	err = container.Start()
	if err != nil {
		return nil, err
	}

	peer, err := newPeer(name, identity.New(publicKey), container, nodeSeed, n)
	if err != nil {
		return nil, err
	}
	n.peers = append(n.peers, peer)
	return peer, nil
}

// CreatePeerWithMana creates a new peers/Goshimmer node in the network and returns it.
// It requests funds from the faucet and pledges mana to itself.
func (n *Network) CreatePeerWithMana(c GoShimmerConfig) (*Peer, error) {
	peer, err := n.CreatePeer(c)
	if err != nil {
		return nil, err
	}
	addr := peer.Seed.Address(uint64(0)).Address()
	ID := base58.Encode(peer.ID().Bytes())
	_, err = peer.SendFaucetRequest(addr.Base58(), ID, ID)
	if err != nil {
		_ = peer.Stop()
		return nil, fmt.Errorf("error sending faucet request... shutting down: %w", err)
	}
	err = n.WaitForMana(peer)
	if err != nil {
		return nil, err
	}
	return peer, nil
}

// Shutdown creates logs and removes network and containers.
// Should always be called when a network is not needed anymore!
func (n *Network) Shutdown() error {
	// stop containers
	err := n.entryNode.Stop()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup

	wg.Add(len(n.peers))
	errs := make([]error, len(n.peers))
	for i := range n.peers {
		go func(index int) {
			defer wg.Done()
			err = n.peers[index].Stop()
			if err != nil {
				errs[index] = err
			}
		}(i)
	}

	wg.Wait()

	for _, e := range errs {
		if e != nil {
			return e
		}
	}
	// delete all partitions
	err = n.DeletePartitions()
	if err != nil {
		return err
	}

	// retrieve logs
	logs, err := n.entryNode.Logs()
	if err != nil {
		return err
	}
	err = createLogFile(n.namePrefix(containerNameEntryNode), logs)
	if err != nil {
		return err
	}
	for _, p := range n.peers {
		logs, err = p.Logs()
		if err != nil {
			return err
		}
		err = createLogFile(p.name, logs)
		if err != nil {
			return err
		}
	}

	// save exit status of containers to check at end of shutdown process
	exitStatus := make(map[string]int, len(n.peers)+1)
	exitStatus[containerNameEntryNode], err = n.entryNode.ExitStatus()
	if err != nil {
		return err
	}
	for _, p := range n.peers {
		exitStatus[p.name], err = p.ExitStatus()
		if err != nil {
			return err
		}
	}

	// remove containers
	err = n.entryNode.Remove()
	if err != nil {
		return err
	}
	for _, p := range n.peers {
		err = p.Remove()
		if err != nil {
			return err
		}
	}

	// disconnect tester from network otherwise the network can't be removed
	err = n.tester.DisconnectFromNetwork(n.id)
	if err != nil {
		return err
	}

	// remove network
	err = n.dockerClient.NetworkRemove(context.Background(), n.id)
	if err != nil {
		return err
	}

	// check exit codes of containers
	for name, status := range exitStatus {
		if status != exitStatusSuccessful {
			return fmt.Errorf("container %s exited with code %d", name, status)
		}
	}

	return nil
}

func (n *Network) doManualPeering() error {
	for idx, p := range n.peers {
		allOtherPeers := make([]*Peer, 0, len(n.peers)-1)
		allOtherPeers = append(allOtherPeers, n.peers[:idx]...)
		allOtherPeers = append(allOtherPeers, n.peers[idx+1:]...)
		peersToAdd := ToPeerModels(allOtherPeers)
		if err := p.AddManualPeers(peersToAdd); err != nil {
			return errors.Wrap(err, "failed to add manual peers via API")
		}
	}
	return nil
}

// WaitForAutopeering waits until all peers have reached the minimum amount of neighbors.
// Returns error if this minimum is not reached after peeringMaxTries.
func (n *Network) WaitForAutopeering(minimumNeighbors int) error {
	getNeighborsFn := func(p *Peer) (int, error) {
		resp, err := p.GetAutopeeringNeighbors(false)
		if err != nil {
			return 0, errors.Wrap(err, "client failed to return autopeering neighbors")
		}
		return len(resp.Chosen) + len(resp.Accepted), nil
	}
	err := n.waitForPeering(minimumNeighbors, getNeighborsFn)
	return errors.WithStack(err)
}

// WaitForManualpeering waits until all peers have reached together as neighbors.
func (n *Network) WaitForManualpeering() error {
	getNeighborsFn := func(p *Peer) (int, error) {
		peers, err := p.GetManualConnectedPeers()
		if err != nil {
			return 0, errors.Wrap(err, "client failed to return manually connected peers")
		}
		return len(peers), nil
	}
	err := n.waitForPeering(len(n.peers)-1, getNeighborsFn)
	return errors.WithStack(err)
}

type getNeighborsNumberFunc func(p *Peer) (int, error)

func (n *Network) waitForPeering(minimumNeighbors int, getNeighborsFn getNeighborsNumberFunc) error {
	log.Printf("Waiting for peering...\n")
	defer log.Printf("Waiting for peering... done\n")

	if minimumNeighbors == 0 {
		return nil
	}

	for i := peeringMaxTries; i > 0; i-- {

		for _, p := range n.peers {
			if neighborsNumber, err := getNeighborsFn(p); err != nil {
				log.Printf("request error: %v\n", err)
			} else {
				p.SetNeighborsNumber(neighborsNumber)
			}
		}

		// verify neighbor requirement
		min := math.MaxInt64
		total := 0
		for _, p := range n.peers {
			neighbors := p.TotalNeighbors()
			if neighbors < min {
				min = neighbors
			}
			total += neighbors
		}
		if min >= minimumNeighbors {
			log.Printf("Neighbors: min=%d avg=%.2f\n", min, float64(total)/float64(len(n.peers)))
			return nil
		}

		log.Println("Not done yet. Try again in 5 seconds...")
		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("peering not successful")
}

// WaitForMana waits until all peers have access mana.
// Returns error if all peers don't have mana after waitForManaMaxTries
func (n *Network) WaitForMana(optionalPeers ...*Peer) error {
	log.Printf("Waiting for nodes to get mana...\n")
	defer log.Printf("Waiting for nodes to get mana... done\n")

	peers := n.peers
	if len(optionalPeers) > 0 {
		peers = optionalPeers
	}
	m := make(map[*Peer]struct{})
	for _, peer := range peers {
		m[peer] = struct{}{}
	}
	for i := waitForManaMaxTries; i > 0; i-- {
		for peer := range m {
			infoRes, err := peer.Info()
			if err == nil && infoRes.Mana.Access > 0.0 {
				delete(m, peer)
			}
		}
		if len(m) == 0 {
			return nil
		}
		log.Println("Not done yet. Try again in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
	return fmt.Errorf("waiting for mana not successful")
}

// namePrefix returns the suffix prefixed with the name.
func (n *Network) namePrefix(suffix string) string {
	return fmt.Sprintf("%s-%s", n.name, suffix)
}

// entryNodePublicKey returns the entry node's public key encoded as base58
func (n *Network) entryNodePublicKey() string {
	return n.entryNodeIdentity.PublicKey().String()
}

// Peers returns all available peers in the network.
func (n *Network) Peers() []*Peer {
	return n.peers
}

// RandomPeer returns a random peer out of the list of peers.
func (n *Network) RandomPeer() *Peer {
	return n.peers[rand.Intn(len(n.peers))]
}

// createPumba creates and starts a Pumba Docker container.
func (n *Network) createPumba(name string, containerName string, targetIPs []string) (*DockerContainer, error) {
	container := NewDockerContainer(n.dockerClient)
	err := container.CreatePumba(name, containerName, targetIPs)
	if err != nil {
		return nil, err
	}
	err = container.Start()
	if err != nil {
		return nil, err
	}

	return container, nil
}

// createPartition creates a partition with the given peers.
// It starts a Pumba container for every peer that blocks traffic to all other partitions.
func (n *Network) createPartition(peers []*Peer) (*Partition, error) {
	peersMap := make(map[string]*Peer)
	for _, peer := range peers {
		peersMap[peer.ID().String()] = peer
	}

	// block all traffic to all other peers except in the current partition
	var targetIPs []string
	for _, peer := range n.peers {
		if _, ok := peersMap[peer.ID().String()]; ok {
			continue
		}
		targetIPs = append(targetIPs, peer.ip)
	}

	partitionName := n.namePrefix(fmt.Sprintf("partition_%d-", len(n.partitions)))

	// create pumba container for every peer in the partition
	pumbas := make([]*DockerContainer, len(peers))
	for i, p := range peers {
		name := partitionName + p.name + containerNameSuffixPumba
		pumba, err := n.createPumba(name, p.name, targetIPs)
		if err != nil {
			return nil, err
		}
		pumbas[i] = pumba
		time.Sleep(1 * time.Second)
	}

	partition := &Partition{
		name:     partitionName,
		peers:    peers,
		peersMap: peersMap,
		pumbas:   pumbas,
	}
	n.partitions = append(n.partitions, partition)

	return partition, nil
}

// DeletePartitions deletes all partitions of the network.
// All nodes can communicate with the full network again.
func (n *Network) DeletePartitions() error {
	for _, p := range n.partitions {
		err := p.deletePartition()
		if err != nil {
			return err
		}
	}
	n.partitions = nil
	return nil
}

// Partitions returns the network's partitions.
func (n *Network) Partitions() []*Partition {
	return n.partitions
}

// Split splits the existing network in given partitions.
func (n *Network) Split(partitions ...[]*Peer) error {
	for _, peers := range partitions {
		_, err := n.createPartition(peers)
		if err != nil {
			return err
		}
	}
	// wait until pumba containers are started and block traffic between partitions
	time.Sleep(5 * time.Second)

	return nil
}

// Partition represents a network partition.
// It contains its peers and the corresponding Pumba instances that block all traffic to peers in other partitions.
type Partition struct {
	name     string
	peers    []*Peer
	peersMap map[string]*Peer
	pumbas   []*DockerContainer
}

// Peers returns the partition's peers.
func (p *Partition) Peers() []*Peer {
	return p.peers
}

// PeersMap returns the partition's peers map.
func (p *Partition) PeersMap() map[string]*Peer {
	return p.peersMap
}

func (p *Partition) String() string {
	return fmt.Sprintf("Partition{%s, %s}", p.name, p.peers)
}

// deletePartition deletes a partition, all its Pumba containers and creates logs for them.
func (p *Partition) deletePartition() error {
	// stop containers
	for _, pumba := range p.pumbas {
		err := pumba.Stop()
		if err != nil {
			return err
		}
	}

	// retrieve logs
	for i, pumba := range p.pumbas {
		logs, err := pumba.Logs()
		if err != nil {
			return err
		}
		err = createLogFile(fmt.Sprintf("%s%s", p.name, p.peers[i].name), logs)
		if err != nil {
			return err
		}
	}

	for _, pumba := range p.pumbas {
		err := pumba.Remove()
		if err != nil {
			return err
		}
	}

	return nil
}
