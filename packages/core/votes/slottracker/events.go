package slottracker

import (
	"github.com/iotaledger/hive.go/core/slot"
	"github.com/iotaledger/hive.go/crypto/identity"
	"github.com/iotaledger/hive.go/runtime/event"
)

type Events struct {
	VotersUpdated *event.Event1[*VoterUpdatedEvent]

	event.Group[Events, *Events]
}

// NewEvents contains the constructor of the Events object (it is generated by a generic factory).
var NewEvents = event.CreateGroupConstructor(func() (newEvents *Events) {
	return &Events{
		VotersUpdated: event.New1[*VoterUpdatedEvent](),
	}
})

type VoterUpdatedEvent struct {
	Voter               identity.ID
	NewLatestSlotIndex  slot.Index
	PrevLatestSlotIndex slot.Index
}