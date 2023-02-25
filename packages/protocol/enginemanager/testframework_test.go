package enginemanager_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/goshimmer/packages/core/database"
	"github.com/iotaledger/goshimmer/packages/core/snapshotcreator"
	"github.com/iotaledger/goshimmer/packages/protocol"
	"github.com/iotaledger/goshimmer/packages/protocol/engine"
	"github.com/iotaledger/goshimmer/packages/protocol/engine/sybilprotection/dpos"
	"github.com/iotaledger/goshimmer/packages/protocol/engine/throughputquota/mana1"
	"github.com/iotaledger/goshimmer/packages/protocol/enginemanager"
	"github.com/iotaledger/goshimmer/packages/protocol/ledger"
	"github.com/iotaledger/goshimmer/packages/protocol/ledger/vm/devnetvm"
	"github.com/iotaledger/goshimmer/packages/storage/utils"
	"github.com/iotaledger/hive.go/core/crypto/ed25519"
	"github.com/iotaledger/hive.go/lo"
	"github.com/iotaledger/hive.go/runtime/options"
	"github.com/iotaledger/hive.go/runtime/workerpool"
)

type EngineManagerTestFramework struct {
	EngineManager *enginemanager.EngineManager

	ActiveEngine *enginemanager.EngineInstance
}

func NewEngineManagerTestFramework(t *testing.T, workers *workerpool.Group, identitiesWeights map[ed25519.PublicKey]uint64) *EngineManagerTestFramework {
	tf := &EngineManagerTestFramework{}

	ledgerVM := new(devnetvm.VM)

	snapshotPath := utils.NewDirectory(t.TempDir()).Path("snapshot.bin")

	snapshotcreator.CreateSnapshot(protocol.DatabaseVersion, snapshotPath, 1, make([]byte, 32), identitiesWeights, lo.Keys(identitiesWeights), ledgerVM)

	tf.EngineManager = enginemanager.New(workers.CreateGroup("EngineManager"),
		t.TempDir(),
		protocol.DatabaseVersion,
		[]options.Option[database.Manager]{},
		[]options.Option[engine.Engine]{
			engine.WithLedgerOptions(ledger.WithVM(ledgerVM)),
		},
		dpos.NewProvider(),
		mana1.NewProvider())

	var err error
	tf.ActiveEngine, err = tf.EngineManager.LoadActiveEngine()
	require.NoError(t, err)

	require.NoError(t, tf.ActiveEngine.InitializeWithSnapshot(snapshotPath))

	return tf
}