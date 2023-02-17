package sequencetracker

import (
	"github.com/iotaledger/goshimmer/packages/protocol/engine/tangle/booker/markers"
	"github.com/iotaledger/hive.go/core/generics/event"
	"github.com/iotaledger/hive.go/core/identity"
)

type Events struct {
	VotersUpdated *event.Linkable[*VoterUpdatedEvent]

	event.LinkableCollection[Events, *Events]
}

// NewEvents contains the constructor of the Events object (it is generated by a generic factory).
var NewEvents = event.LinkableConstructor(func() (newEvents *Events) {
	return &Events{
		VotersUpdated: event.NewLinkable[*VoterUpdatedEvent](),
	}
})

type VoterUpdatedEvent struct {
	Voter                 identity.ID
	NewMaxSupportedIndex  markers.Index
	PrevMaxSupportedIndex markers.Index
	SequenceID            markers.SequenceID
}
