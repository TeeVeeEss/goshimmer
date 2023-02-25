package metrics

import (
	"fmt"
	"strconv"

	"github.com/iotaledger/goshimmer/packages/app/collector"
	"github.com/iotaledger/goshimmer/packages/protocol/engine/consensus/blockgadget"
	"github.com/iotaledger/goshimmer/packages/protocol/engine/notarization"
	"github.com/iotaledger/goshimmer/packages/protocol/engine/tangle/blockdag"
	"github.com/iotaledger/goshimmer/packages/protocol/engine/tangle/booker/virtualvoting"
	"github.com/iotaledger/goshimmer/packages/protocol/ledger"
	"github.com/iotaledger/goshimmer/packages/protocol/ledger/conflictdag"
	"github.com/iotaledger/goshimmer/packages/protocol/ledger/utxo"
	"github.com/iotaledger/hive.go/runtime/event"
)

const (
	epochNamespace            = "epochs"
	labelName                 = "epoch"
	metricEvictionOffset      = 6
	totalBlocks               = "total_blocks"
	acceptedBlocksInEpoch     = "accepted_blocks"
	orphanedBlocks            = "orphaned_blocks"
	invalidBlocks             = "invalid_blocks"
	subjectivelyInvalidBlocks = "subjectively_invalid_blocks"
	totalAttachments          = "total_attachments"
	rejectedAttachments       = "rejected_attachments"
	acceptedAttachments       = "accepted_attachments"
	orphanedAttachments       = "orphaned_attachments"
	createdConflicts          = "created_conflicts"
	acceptedConflicts         = "accepted_conflicts"
	rejectedConflicts         = "rejected_conflicts"
	notConflictingConflicts   = "not_conflicting_conflicts"
)

var EpochMetrics = collector.NewCollection(epochNamespace,
	collector.WithMetric(collector.NewMetric(totalBlocks,
		collector.WithType(collector.CounterVec),
		collector.WithLabels(labelName),
		collector.WithHelp("Number of blocks seen by the node in an epoch."),
		collector.WithInitFunc(func() {
			deps.Protocol.Events.Engine.Tangle.BlockDAG.BlockAttached.Hook(func(block *blockdag.Block) {
				eventEpoch := int(block.ID().Index())
				deps.Collector.Increment(epochNamespace, totalBlocks, strconv.Itoa(eventEpoch))

				// need to initialize epoch metrics with 0 to have consistent data for each epoch
				for _, metricName := range []string{acceptedBlocksInEpoch, orphanedBlocks, invalidBlocks, subjectivelyInvalidBlocks, totalAttachments, orphanedAttachments, rejectedAttachments, acceptedAttachments, createdConflicts, acceptedConflicts, rejectedConflicts, notConflictingConflicts} {
					deps.Collector.Update(epochNamespace, metricName, map[string]float64{
						strconv.Itoa(eventEpoch): 0,
					})
				}
			}, event.WithWorkerPool(Plugin.WorkerPool))

			// initialize it once and remove committed epoch from all metrics (as they will not change afterwards)
			// in a single attachment instead of multiple ones
			deps.Protocol.Events.Engine.NotarizationManager.EpochCommitted.Hook(func(details *notarization.EpochCommittedDetails) {
				epochToEvict := int(details.Commitment.Index()) - metricEvictionOffset

				// need to remove metrics for old epochs, otherwise they would be stored in memory and always exposed to Prometheus, forever
				for _, metricName := range []string{totalBlocks, acceptedBlocksInEpoch, orphanedBlocks, invalidBlocks, subjectivelyInvalidBlocks, totalAttachments, orphanedAttachments, rejectedAttachments, acceptedAttachments, createdConflicts, acceptedConflicts, rejectedConflicts, notConflictingConflicts} {
					deps.Collector.ResetMetricLabels(epochNamespace, metricName, map[string]string{
						labelName: strconv.Itoa(epochToEvict),
					})
				}
			}, event.WithWorkerPool(Plugin.WorkerPool))
		}),
	)),

	collector.WithMetric(collector.NewMetric(acceptedBlocksInEpoch,
		collector.WithType(collector.CounterVec),
		collector.WithHelp("Number of accepted blocks in an epoch."),
		collector.WithLabels(labelName),
		collector.WithInitFunc(func() {
			deps.Protocol.Events.Engine.Consensus.BlockGadget.BlockAccepted.Hook(func(block *blockgadget.Block) {
				eventEpoch := int(block.ID().Index())
				deps.Collector.Increment(epochNamespace, acceptedBlocksInEpoch, strconv.Itoa(eventEpoch))
			}, event.WithWorkerPool(Plugin.WorkerPool))
		}),
	)),
	collector.WithMetric(collector.NewMetric(orphanedBlocks,
		collector.WithType(collector.CounterVec),
		collector.WithLabels(labelName),
		collector.WithHelp("Number of orphaned blocks in an epoch."),
		collector.WithInitFunc(func() {
			deps.Protocol.Events.Engine.Tangle.BlockDAG.BlockOrphaned.Hook(func(block *blockdag.Block) {
				eventEpoch := int(block.ID().Index())
				deps.Collector.Increment(epochNamespace, orphanedBlocks, strconv.Itoa(eventEpoch))
			}, event.WithWorkerPool(Plugin.WorkerPool))
		}),
	)),
	collector.WithMetric(collector.NewMetric(invalidBlocks,
		collector.WithType(collector.CounterVec),
		collector.WithLabels(labelName),
		collector.WithHelp("Number of invalid blocks in an epoch epoch."),
		collector.WithInitFunc(func() {
			deps.Protocol.Events.Engine.Tangle.BlockDAG.BlockInvalid.Hook(func(blockInvalidEvent *blockdag.BlockInvalidEvent) {
				fmt.Println("block invalid", blockInvalidEvent.Block.ID(), blockInvalidEvent.Reason)
				eventEpoch := int(blockInvalidEvent.Block.ID().Index())
				deps.Collector.Increment(epochNamespace, invalidBlocks, strconv.Itoa(eventEpoch))
			}, event.WithWorkerPool(Plugin.WorkerPool))
		}),
	)),
	collector.WithMetric(collector.NewMetric(subjectivelyInvalidBlocks,
		collector.WithType(collector.CounterVec),
		collector.WithLabels(labelName),
		collector.WithHelp("Number of invalid blocks in an epoch epoch."),
		collector.WithInitFunc(func() {
			deps.Protocol.Events.Engine.Tangle.Booker.VirtualVoting.BlockTracked.Hook(func(block *virtualvoting.Block) {
				if block.IsSubjectivelyInvalid() {
					eventEpoch := int(block.ID().Index())
					deps.Collector.Increment(epochNamespace, subjectivelyInvalidBlocks, strconv.Itoa(eventEpoch))
				}
			}, event.WithWorkerPool(Plugin.WorkerPool))
		}),
	)),

	collector.WithMetric(collector.NewMetric(totalAttachments,
		collector.WithType(collector.CounterVec),
		collector.WithLabels(labelName),
		collector.WithHelp("Number of transaction attachments by the node per epoch."),
		collector.WithInitFunc(func() {
			deps.Protocol.Events.Engine.Tangle.Booker.AttachmentCreated.Hook(func(block *virtualvoting.Block) {
				eventEpoch := int(block.ID().Index())
				deps.Collector.Increment(epochNamespace, totalAttachments, strconv.Itoa(eventEpoch))
			}, event.WithWorkerPool(Plugin.WorkerPool))
		}),
	)),
	collector.WithMetric(collector.NewMetric(orphanedAttachments,
		collector.WithType(collector.CounterVec),
		collector.WithLabels(labelName),
		collector.WithHelp("Number of orphaned attachments by the node per epoch."),
		collector.WithInitFunc(func() {
			deps.Protocol.Events.Engine.Tangle.Booker.AttachmentOrphaned.Hook(func(block *virtualvoting.Block) {
				eventEpoch := int(block.ID().Index())
				deps.Collector.Increment(epochNamespace, orphanedAttachments, strconv.Itoa(eventEpoch))
			}, event.WithWorkerPool(Plugin.WorkerPool))
		}),
	)),
	collector.WithMetric(collector.NewMetric(rejectedAttachments,
		collector.WithType(collector.CounterVec),
		collector.WithLabels(labelName),
		collector.WithHelp("Number of rejected attachments by the node per epoch."),
		collector.WithInitFunc(func() {
			deps.Protocol.Events.Engine.Ledger.TransactionRejected.Hook(func(transactionMetadata *ledger.TransactionMetadata) {
				for it := deps.Protocol.Engine().Tangle.Booker.GetAllAttachments(transactionMetadata.ID()).Iterator(); it.HasNext(); {
					attachmentBlock := it.Next()
					if !attachmentBlock.IsOrphaned() {
						deps.Collector.Increment(epochNamespace, rejectedAttachments, strconv.Itoa(int(attachmentBlock.ID().Index())))
					}
				}
			}, event.WithWorkerPool(Plugin.WorkerPool))
		}),
	)),
	collector.WithMetric(collector.NewMetric(acceptedAttachments,
		collector.WithType(collector.CounterVec),
		collector.WithLabels(labelName),
		collector.WithHelp("Number of accepted attachments by the node per epoch."),
		collector.WithInitFunc(func() {
			deps.Protocol.Events.Engine.Ledger.TransactionAccepted.Hook(func(transactionEvent *ledger.TransactionEvent) {
				for it := deps.Protocol.Engine().Tangle.Booker.GetAllAttachments(transactionEvent.Metadata.ID()).Iterator(); it.HasNext(); {
					attachmentBlock := it.Next()
					if !attachmentBlock.IsOrphaned() {
						deps.Collector.Increment(epochNamespace, acceptedAttachments, strconv.Itoa(int(attachmentBlock.ID().Index())))
					}
				}
			}, event.WithWorkerPool(Plugin.WorkerPool))
		}),
	)),
	collector.WithMetric(collector.NewMetric(createdConflicts,
		collector.WithType(collector.CounterVec),
		collector.WithLabels(labelName),
		collector.WithHelp("Number of conflicts created per epoch."),
		collector.WithInitFunc(func() {
			deps.Protocol.Events.Engine.Ledger.ConflictDAG.ConflictCreated.Hook(func(conflictCreated *conflictdag.Conflict[utxo.TransactionID, utxo.OutputID]) {
				for it := deps.Protocol.Engine().Tangle.Booker.GetAllAttachments(conflictCreated.ID()).Iterator(); it.HasNext(); {
					deps.Collector.Increment(epochNamespace, createdConflicts, strconv.Itoa(int(it.Next().ID().Index())))
				}
			}, event.WithWorkerPool(Plugin.WorkerPool))
		}),
	)),
	collector.WithMetric(collector.NewMetric(acceptedConflicts,
		collector.WithType(collector.CounterVec),
		collector.WithLabels(labelName),
		collector.WithHelp("Number of conflicts accepted per epoch."),
		collector.WithInitFunc(func() {
			deps.Protocol.Events.Engine.Ledger.ConflictDAG.ConflictAccepted.Hook(func(conflict *conflictdag.Conflict[utxo.TransactionID, utxo.OutputID]) {
				for it := deps.Protocol.Engine().Tangle.Booker.GetAllAttachments(conflict.ID()).Iterator(); it.HasNext(); {
					deps.Collector.Increment(epochNamespace, acceptedConflicts, strconv.Itoa(int(it.Next().ID().Index())))
				}
			}, event.WithWorkerPool(Plugin.WorkerPool))
		}),
	)),
	collector.WithMetric(collector.NewMetric(rejectedConflicts,
		collector.WithType(collector.CounterVec),
		collector.WithLabels(labelName),
		collector.WithHelp("Number of conflicts rejected per epoch."),
		collector.WithInitFunc(func() {
			deps.Protocol.Events.Engine.Ledger.ConflictDAG.ConflictRejected.Hook(func(conflict *conflictdag.Conflict[utxo.TransactionID, utxo.OutputID]) {
				for it := deps.Protocol.Engine().Tangle.Booker.GetAllAttachments(conflict.ID()).Iterator(); it.HasNext(); {
					deps.Collector.Increment(epochNamespace, rejectedConflicts, strconv.Itoa(int(it.Next().ID().Index())))
				}
			}, event.WithWorkerPool(Plugin.WorkerPool))
		}),
	)),
	collector.WithMetric(collector.NewMetric(notConflictingConflicts,
		collector.WithType(collector.CounterVec),
		collector.WithLabels(labelName),
		collector.WithHelp("Number of conflicts rejected per epoch."),
		collector.WithInitFunc(func() {
			deps.Protocol.Events.Engine.Ledger.ConflictDAG.ConflictNotConflicting.Hook(func(conflict *conflictdag.Conflict[utxo.TransactionID, utxo.OutputID]) {
				for it := deps.Protocol.Engine().Tangle.Booker.GetAllAttachments(conflict.ID()).Iterator(); it.HasNext(); {
					deps.Collector.Increment(epochNamespace, notConflictingConflicts, strconv.Itoa(int(it.Next().ID().Index())))
				}
			}, event.WithWorkerPool(Plugin.WorkerPool))
		}),
	)),
)