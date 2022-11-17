package blockgadget

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/iotaledger/hive.go/core/debug"
	"github.com/iotaledger/hive.go/core/generics/event"
	"github.com/iotaledger/hive.go/core/generics/options"
	"github.com/iotaledger/hive.go/core/generics/set"
	"github.com/iotaledger/hive.go/core/types/confirmation"
	"github.com/stretchr/testify/assert"

	"github.com/iotaledger/goshimmer/packages/core/epoch"
	"github.com/iotaledger/goshimmer/packages/protocol/engine/eviction"
	"github.com/iotaledger/goshimmer/packages/protocol/engine/sybilprotection/impl"
	"github.com/iotaledger/goshimmer/packages/protocol/engine/tangle"
	"github.com/iotaledger/goshimmer/packages/protocol/engine/tangle/booker/markers"
	"github.com/iotaledger/goshimmer/packages/protocol/ledger"
	"github.com/iotaledger/goshimmer/packages/protocol/ledger/utxo"
	"github.com/iotaledger/goshimmer/packages/protocol/models"
	"github.com/iotaledger/goshimmer/packages/storage"
)

// region TestFramework //////////////////////////////////////////////////////////////////////////////////////////////////////

type TestFramework struct {
	Gadget *Gadget

	test              *testing.T
	acceptedBlocks    uint32
	confirmedBlocks   uint32
	conflictsAccepted uint32
	conflictsRejected uint32
	reorgCount        uint32

	optsGadgetOptions       []options.Option[Gadget]
	optsLedger              *ledger.Ledger
	optsLedgerOptions       []options.Option[ledger.Ledger]
	optsEvictionState       *eviction.State
	optsTangle              *tangle.Tangle
	optsTangleOptions       []options.Option[tangle.Tangle]
	optsTotalWeightCallback func() int64
	optsActiveNodes         *impl.ActiveValidators

	*TangleTestFramework
}

func NewTestFramework(test *testing.T, opts ...options.Option[TestFramework]) (t *TestFramework) {
	return options.Apply(&TestFramework{
		test: test,
		optsTotalWeightCallback: func() int64 {
			return t.TangleTestFramework.ActiveNodes.Weight()
		},
	}, opts, func(t *TestFramework) {
		if t.Gadget == nil {
			if t.optsTangle == nil {
				storageInstance := storage.New(test.TempDir(), 1)
				test.Cleanup(func() {
					t.optsLedger.Shutdown()
					if err := storageInstance.Shutdown(); err != nil {
						test.Fatal(err)
					}
				})

				if t.optsLedger == nil {
					t.optsLedger = ledger.New(storageInstance, t.optsLedgerOptions...)
				}

				if t.optsEvictionState == nil {
					t.optsEvictionState = eviction.NewState(storageInstance)
				}

				if t.optsActiveNodes == nil {
					t.optsActiveNodes = impl.New(time.Now)
				}

				t.optsTangle = tangle.New(t.optsLedger, t.optsEvictionState, t.optsActiveNodes, func() epoch.Index {
					return 0
				}, func(id markers.SequenceID) markers.Index {
					return 1
				}, t.optsTangleOptions...)
			}

			t.Gadget = New(t.optsTangle, t.optsEvictionState, t.optsTotalWeightCallback, t.optsGadgetOptions...)
		}

		if t.TangleTestFramework == nil {
			t.TangleTestFramework = tangle.NewTestFramework(test, tangle.WithTangle(t.optsTangle))
		}
	}, (*TestFramework).setupEvents)
}

func (t *TestFramework) setupEvents() {
	t.Gadget.Events.BlockAccepted.Hook(event.NewClosure(func(metadata *Block) {
		if debug.GetEnabled() {
			t.test.Logf("ACCEPTED: %s", metadata.ID())
		}

		atomic.AddUint32(&(t.acceptedBlocks), 1)
	}))

	t.Gadget.Events.BlockConfirmed.Hook(event.NewClosure(func(metadata *Block) {
		if debug.GetEnabled() {
			t.test.Logf("CONFIRMED: %s", metadata.ID())
		}

		atomic.AddUint32(&(t.confirmedBlocks), 1)
	}))

	t.Gadget.Events.Reorg.Hook(event.NewClosure(func(conflictID utxo.TransactionID) {
		if debug.GetEnabled() {
			t.test.Logf("REORG NEEDED: %s", conflictID)
		}
		atomic.AddUint32(&(t.reorgCount), 1)
	}))

	t.ConflictDAG().Events.ConflictAccepted.Hook(event.NewClosure(func(conflictID utxo.TransactionID) {
		if debug.GetEnabled() {
			t.test.Logf("CONFLICT ACCEPTED: %s", conflictID)
		}
		atomic.AddUint32(&(t.conflictsAccepted), 1)
	}))

	t.ConflictDAG().Events.ConflictRejected.Hook(event.NewClosure(func(conflictID utxo.TransactionID) {
		if debug.GetEnabled() {
			t.test.Logf("CONFLICT REJECTED: %s", conflictID)
		}

		atomic.AddUint32(&(t.conflictsRejected), 1)
	}))
}

func (t *TestFramework) AssertBlockAccepted(blocksAccepted uint32) {
	assert.Equal(t.test, blocksAccepted, atomic.LoadUint32(&t.acceptedBlocks), "expected %d blocks to be accepted but got %d", blocksAccepted, atomic.LoadUint32(&t.acceptedBlocks))
}

func (t *TestFramework) AssertBlockConfirmed(blocksConfirmed uint32) {
	assert.Equal(t.test, blocksConfirmed, atomic.LoadUint32(&t.confirmedBlocks), "expected %d blocks to be accepted but got %d", blocksConfirmed, atomic.LoadUint32(&t.confirmedBlocks))
}

func (t *TestFramework) AssertConflictsAccepted(conflictsAccepted uint32) {
	assert.Equal(t.test, conflictsAccepted, atomic.LoadUint32(&t.conflictsAccepted), "expected %d conflicts to be accepted but got %d", conflictsAccepted, atomic.LoadUint32(&t.acceptedBlocks))
}

func (t *TestFramework) AssertConflictsRejected(conflictsRejected uint32) {
	assert.Equal(t.test, conflictsRejected, atomic.LoadUint32(&t.conflictsRejected), "expected %d conflicts to be rejected but got %d", conflictsRejected, atomic.LoadUint32(&t.acceptedBlocks))
}

func (t *TestFramework) AssertReorgs(reorgCount uint32) {
	assert.Equal(t.test, reorgCount, atomic.LoadUint32(&t.reorgCount), "expected %d reorgs but got %d", reorgCount, atomic.LoadUint32(&t.reorgCount))
}

func (t *TestFramework) ValidateAcceptedBlocks(expectedAcceptedBlocks map[string]bool) {
	for blockID, blockExpectedAccepted := range expectedAcceptedBlocks {
		actualBlockAccepted := t.Gadget.IsBlockAccepted(t.Block(blockID).ID())
		assert.Equal(t.test, blockExpectedAccepted, actualBlockAccepted, "Block %s should be accepted=%t but is %t", blockID, blockExpectedAccepted, actualBlockAccepted)
	}
}

func (t *TestFramework) ValidateConfirmedBlocks(expectedConfirmedBlocks map[string]bool) {
	for blockID, blockExpectedConfirmed := range expectedConfirmedBlocks {
		actualBlockConfirmed := t.Gadget.isBlockConfirmed(t.Block(blockID).ID())
		assert.Equal(t.test, blockExpectedConfirmed, actualBlockConfirmed, "Block %s should be confirmed=%t but is %t", blockID, blockExpectedConfirmed, actualBlockConfirmed)
	}
}

func (t *TestFramework) ValidateAcceptedMarker(expectedConflictIDs map[markers.Marker]bool) {
	for marker, markerExpectedAccepted := range expectedConflictIDs {
		actualMarkerAccepted := t.Gadget.IsMarkerAccepted(marker)
		assert.Equal(t.test, markerExpectedAccepted, actualMarkerAccepted, "%s should be accepted=%t but is %t", marker, markerExpectedAccepted, actualMarkerAccepted)
	}
}

func (t *TestFramework) ValidateConflictAcceptance(expectedConflictIDs map[string]confirmation.State) {
	for conflictIDAlias, conflictExpectedState := range expectedConflictIDs {
		actualMarkerAccepted := t.ConflictDAG().ConfirmationState(set.NewAdvancedSet(t.Transaction(conflictIDAlias).ID()))
		assert.Equal(t.test, conflictExpectedState, actualMarkerAccepted, "%s should be accepted=%t but is %t", conflictIDAlias, conflictExpectedState, actualMarkerAccepted)
	}
}

type TangleTestFramework = tangle.TestFramework

// endregion ///////////////////////////////////////////////////////////////////////////////////////////////////////////

// region Options //////////////////////////////////////////////////////////////////////////////////////////////////////

func WithGadget(gadget *Gadget) options.Option[TestFramework] {
	return func(tf *TestFramework) {
		tf.Gadget = gadget
	}
}

func WithTotalWeightCallback(totalWeightCallback func() int64) options.Option[TestFramework] {
	return func(tf *TestFramework) {
		tf.optsTotalWeightCallback = totalWeightCallback
	}
}

func WithGadgetOptions(opts ...options.Option[Gadget]) options.Option[TestFramework] {
	return func(tf *TestFramework) {
		tf.optsGadgetOptions = opts
	}
}

func WithTangle(tangle *tangle.Tangle) options.Option[TestFramework] {
	return func(tf *TestFramework) {
		tf.optsTangle = tangle
	}
}

func WithTangleOptions(opts ...options.Option[tangle.Tangle]) options.Option[TestFramework] {
	return func(tf *TestFramework) {
		tf.optsTangleOptions = opts
	}
}

func WithTangleTestFramework(testFramework *tangle.TestFramework) options.Option[TestFramework] {
	return func(tf *TestFramework) {
		tf.TangleTestFramework = testFramework
	}
}

func WithLedger(ledger *ledger.Ledger) options.Option[TestFramework] {
	return func(tf *TestFramework) {
		tf.optsLedger = ledger
	}
}

func WithLedgerOptions(opts ...options.Option[ledger.Ledger]) options.Option[TestFramework] {
	return func(tf *TestFramework) {
		tf.optsLedgerOptions = opts
	}
}

func WithEvictionState(evictionState *eviction.State) options.Option[TestFramework] {
	return func(t *TestFramework) {
		t.optsEvictionState = evictionState
	}
}

// endregion ///////////////////////////////////////////////////////////////////////////////////////////////////////////

// region Options //////////////////////////////////////////////////////////////////////////////////////////////////////

// MockAcceptanceGadget mocks ConfirmationOracle marking all blocks as confirmed.
type MockAcceptanceGadget struct {
	BlockAcceptedEvent *event.Linkable[*Block]
	AcceptedBlocks     models.BlockIDs
	AcceptedMarkers    *markers.Markers

	mutex sync.RWMutex
}

func NewMockAcceptanceGadget() *MockAcceptanceGadget {
	return &MockAcceptanceGadget{
		BlockAcceptedEvent: event.NewLinkable[*Block](),
		AcceptedBlocks:     models.NewBlockIDs(),
		AcceptedMarkers:    markers.NewMarkers(),
	}
}

func (m *MockAcceptanceGadget) SetBlocksAccepted(blocks models.BlockIDs) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for block := range blocks {
		m.AcceptedBlocks.Add(block)
	}
}

func (m *MockAcceptanceGadget) SetMarkersAccepted(markers ...markers.Marker) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, marker := range markers {
		m.AcceptedMarkers.Set(marker.SequenceID(), marker.Index())
	}
}

// IsBlockAccepted mocks its interface function returning that all blocks are confirmed.
func (m *MockAcceptanceGadget) IsBlockAccepted(blockID models.BlockID) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.AcceptedBlocks.Contains(blockID)
}

func (m *MockAcceptanceGadget) IsMarkerAccepted(marker markers.Marker) (accepted bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if marker.Index() == 0 {
		return true
	}

	if m.AcceptedMarkers == nil || m.AcceptedMarkers.Size() == 0 {
		return false
	}
	acceptedIndex, exists := m.AcceptedMarkers.Get(marker.SequenceID())
	if !exists {
		return false
	}
	return marker.Index() <= acceptedIndex
}

func (m *MockAcceptanceGadget) FirstUnacceptedIndex(sequenceID markers.SequenceID) (firstUnacceptedIndex markers.Index) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	acceptedIndex, exists := m.AcceptedMarkers.Get(sequenceID)
	if exists {
		return acceptedIndex + 1
	}
	return 1
}

func WithActiveNodes(activeNodes *impl.ActiveValidators) options.Option[TestFramework] {
	return func(t *TestFramework) {
		t.optsActiveNodes = activeNodes
	}
}

// endregion ///////////////////////////////////////////////////////////////////////////////////////////////////////////
