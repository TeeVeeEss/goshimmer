package markermanager

import (
	"github.com/iotaledger/hive.go/core/generics/event"

	"github.com/iotaledger/goshimmer/packages/protocol/engine/tangle/booker/markers"
)

type Events struct {
	SequenceEvicted *event.Linkable[markers.SequenceID]

	event.LinkableCollection[Events, *Events]
}

// NewEvents contains the constructor of the Events object (it is generated by a generic factory).
var NewEvents = event.LinkableConstructor(func() (newEvents *Events) {
	return &Events{
		SequenceEvicted: event.NewLinkable[markers.SequenceID](),
	}
})