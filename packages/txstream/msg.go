package txstream

import (
	"fmt"

	"github.com/iotaledger/goshimmer/packages/ledgerstate"
	"github.com/iotaledger/hive.go/marshalutil"
)

// MessageType represents the type of a message in the txstream protocol
type MessageType byte

const (
	// FlagClientToServer is set in a message type if the message is client to server
	FlagClientToServer = byte(0x80)
	// FlagServerToClient is set in a message type if the message is server to client
	FlagServerToClient = byte(0x40)

	msgTypeChunk = MessageType((FlagClientToServer | FlagServerToClient) + iota)

	msgTypePostTransaction = MessageType(FlagClientToServer + iota)
	msgTypeSubscribe
	msgTypeGetConfirmedTransaction
	msgTypeGetTxInclusionState
	msgTypeGetBacklog
	msgTypeSetID

	msgTypeTransaction = MessageType(FlagServerToClient + iota)
	msgTypeTxInclusionState
)

// Message is the common interface of all messages in the txstream protocol
type Message interface {
	Write(w *marshalutil.MarshalUtil)
	Read(r *marshalutil.MarshalUtil) error
	Type() MessageType
}

// ChunkMessageHeaderSize is the amount of bytes added by MsgChunk as overhead to each chunk
const ChunkMessageHeaderSize = 3

// MsgChunk is a special message for big data packets chopped into pieces
type MsgChunk struct {
	Data []byte
}

// region client --> server

// MsgPostTransaction is a request from the client to post a
// transaction in the ledger.
// No reply from server.
type MsgPostTransaction struct {
	Tx *ledgerstate.Transaction
}

// MsgUpdateSubscriptions is a request from the client to subscribe to
// requests/transactions for the given addresses. Server will send
// all transactions containing unspent outputs to the address,
// and then whenever a relevant transaction is confirmed
// in the ledger, i will be sent in real-time.
type MsgUpdateSubscriptions struct {
	Addresses []ledgerstate.Address
}

// MsgGetConfirmedTransaction is a request to get a specific confirmed
// transaction from the ledger. Server replies with MsgTransaction.
type MsgGetConfirmedTransaction struct {
	Address ledgerstate.Address
	TxID    ledgerstate.TransactionID
}

// MsgGetTxInclusionState is a request to get the inclusion state for a transaction.
// Server replies with MsgTxInclusionState.
type MsgGetTxInclusionState struct {
	Address ledgerstate.Address
	TxID    ledgerstate.TransactionID
}

// MsgGetBacklog is a request to get the backlog for the given address. Server replies
// sending one MsgTransaction for each transaction with unspent outputs targeted
// to the address.
type MsgGetBacklog struct {
	Address ledgerstate.Address
}

// MsgSetID is a message from client informing its ID, used mostly for tracing/loging.
type MsgSetID struct {
	ClientID string
}

// endregion

// region server --> client

// MsgTransaction informs the client of a given confirmed transaction in the ledger.
type MsgTransaction struct {
	// Address is the address that requested the transaction
	Address ledgerstate.Address
	// Tx is the transaction being sent
	Tx *ledgerstate.Transaction
}

// MsgTxInclusionState informs the client with the inclusion state of a given
// transaction as a response from the given address.
type MsgTxInclusionState struct {
	Address ledgerstate.Address
	TxID    ledgerstate.TransactionID
	State   ledgerstate.InclusionState
}

// endregion

func EncodeMsg(msg Message) []byte {
	m := marshalutil.New()
	m.WriteByte(byte(msg.Type()))
	msg.Write(m)
	return m.Bytes()
}

func DecodeMsg(data []byte, expectedFlags uint8) (interface{}, error) {
	if len(data) < 1 {
		return nil, fmt.Errorf("wrong message")
	}
	var ret Message

	msgType := MessageType(data[0])
	if uint8(msgType)&expectedFlags == 0 {
		return nil, fmt.Errorf("unexpected message")
	}

	switch msgType {
	case msgTypeChunk:
		ret = &MsgChunk{}

	case msgTypePostTransaction:
		ret = &MsgPostTransaction{}

	case msgTypeSubscribe:
		ret = &MsgUpdateSubscriptions{}

	case msgTypeGetConfirmedTransaction:
		ret = &MsgGetConfirmedTransaction{}

	case msgTypeGetTxInclusionState:
		ret = &MsgGetTxInclusionState{}

	case msgTypeGetBacklog:
		ret = &MsgGetBacklog{}

	case msgTypeSetID:
		ret = &MsgSetID{}

	case msgTypeTransaction:
		ret = &MsgTransaction{}

	case msgTypeTxInclusionState:
		ret = &MsgTxInclusionState{}

	default:
		return nil, fmt.Errorf("unknown message type %d", msgType)
	}
	if err := ret.Read(marshalutil.New(data[1:])); err != nil {
		return nil, err
	}
	return ret, nil
}

func (msg *MsgPostTransaction) Write(w *marshalutil.MarshalUtil) {
	w.Write(msg.Tx)
}

func (msg *MsgPostTransaction) Read(m *marshalutil.MarshalUtil) error {
	var err error
	if msg.Tx, err = ledgerstate.TransactionFromMarshalUtil(m); err != nil {
		return err
	}
	return nil
}

func (msg *MsgPostTransaction) Type() MessageType {
	return msgTypePostTransaction
}

func (msg *MsgUpdateSubscriptions) Write(w *marshalutil.MarshalUtil) {
	w.WriteUint16(uint16(len(msg.Addresses)))
	for _, addr := range msg.Addresses {
		w.Write(addr)
	}
}

func (msg *MsgUpdateSubscriptions) Read(m *marshalutil.MarshalUtil) error {
	var err error
	var size uint16
	if size, err = m.ReadUint16(); err != nil {
		return err
	}
	msg.Addresses = make([]ledgerstate.Address, size)
	for i := uint16(0); i < size; i++ {
		if msg.Addresses[i], err = ledgerstate.AddressFromMarshalUtil(m); err != nil {
			return err
		}
	}
	return nil
}

func (msg *MsgUpdateSubscriptions) Type() MessageType {
	return msgTypeSubscribe
}

func (msg *MsgGetConfirmedTransaction) Write(w *marshalutil.MarshalUtil) {
	w.Write(msg.Address)
	w.Write(msg.TxID)
}

func (msg *MsgGetConfirmedTransaction) Read(m *marshalutil.MarshalUtil) error {
	var err error
	if msg.Address, err = ledgerstate.AddressFromMarshalUtil(m); err != nil {
		return err
	}
	msg.TxID, err = ledgerstate.TransactionIDFromMarshalUtil(m)
	return err
}

func (msg *MsgGetConfirmedTransaction) Type() MessageType {
	return msgTypeGetConfirmedTransaction
}

func (msg *MsgGetTxInclusionState) Write(w *marshalutil.MarshalUtil) {
	w.Write(msg.Address)
	w.Write(msg.TxID)
}

func (msg *MsgGetTxInclusionState) Read(m *marshalutil.MarshalUtil) error {
	var err error
	if msg.Address, err = ledgerstate.AddressFromMarshalUtil(m); err != nil {
		return err
	}
	if msg.TxID, err = ledgerstate.TransactionIDFromMarshalUtil(m); err != nil {
		return err
	}
	return nil
}

func (msg *MsgGetTxInclusionState) Type() MessageType {
	return msgTypeGetTxInclusionState
}

func (msg *MsgGetBacklog) Write(w *marshalutil.MarshalUtil) {
	w.Write(msg.Address)
}

func (msg *MsgGetBacklog) Read(m *marshalutil.MarshalUtil) error {
	var err error
	msg.Address, err = ledgerstate.AddressFromMarshalUtil(m)
	return err
}

func (msg *MsgGetBacklog) Type() MessageType {
	return msgTypeGetBacklog
}

func (msg *MsgSetID) Write(w *marshalutil.MarshalUtil) {
	w.WriteUint16(uint16(len(msg.ClientID)))
	w.WriteBytes([]byte(msg.ClientID))
}

func (msg *MsgSetID) Read(m *marshalutil.MarshalUtil) error {
	var err error
	var size uint16
	if size, err = m.ReadUint16(); err != nil {
		return err
	}
	var clientID []byte
	if clientID, err = m.ReadBytes(int(size)); err != nil {
		return err
	}
	msg.ClientID = string(clientID)
	return nil
}

func (msg *MsgSetID) Type() MessageType {
	return msgTypeSetID
}

func (msg *MsgTransaction) Write(w *marshalutil.MarshalUtil) {
	w.Write(msg.Address)
	w.Write(msg.Tx)
}

func (msg *MsgTransaction) Read(m *marshalutil.MarshalUtil) error {
	var err error
	if msg.Address, err = ledgerstate.AddressFromMarshalUtil(m); err != nil {
		return err
	}
	if msg.Tx, err = ledgerstate.TransactionFromMarshalUtil(m); err != nil {
		return err
	}
	return nil
}

func (msg *MsgTransaction) Type() MessageType {
	return msgTypeTransaction
}

func (msg *MsgTxInclusionState) Write(w *marshalutil.MarshalUtil) {
	w.Write(msg.Address)
	w.Write(msg.State)
	w.Write(msg.TxID)
}

func (msg *MsgTxInclusionState) Read(m *marshalutil.MarshalUtil) error {
	var err error
	if msg.Address, err = ledgerstate.AddressFromMarshalUtil(m); err != nil {
		return err
	}
	if msg.State, err = ledgerstate.InclusionStateFromMarshalUtil(m); err != nil {
		return err
	}
	if msg.TxID, err = ledgerstate.TransactionIDFromMarshalUtil(m); err != nil {
		return err
	}
	return nil
}

func (msg *MsgTxInclusionState) Type() MessageType {
	return msgTypeTxInclusionState
}

func (msg *MsgChunk) Write(w *marshalutil.MarshalUtil) {
	w.WriteUint16(uint16(len(msg.Data)))
	w.WriteBytes(msg.Data)
}

func (msg *MsgChunk) Read(m *marshalutil.MarshalUtil) error {
	var err error
	var size uint16
	if size, err = m.ReadUint16(); err != nil {
		return err
	}
	msg.Data, err = m.ReadBytes(int(size))
	return err
}

func (msg *MsgChunk) Type() MessageType {
	return msgTypeChunk
}
