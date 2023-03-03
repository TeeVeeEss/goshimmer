package tangletime

import (
	"time"

	"github.com/iotaledger/goshimmer/packages/core/module"
	"github.com/iotaledger/goshimmer/packages/protocol/engine/clock"
)

// Clock is a clock that is used to derive some Time parameters from the Tangle.
type Clock struct {
	events        *clock.Events
	lastAccepted  *timeUpdate
	lastConfirmed *timeUpdate

	module.Module
}

// New creates a new Clock with the given genesisTime.
func New() *Clock {
	return &Clock{
		events:        clock.NewEvents(),
		lastAccepted:  &timeUpdate{},
		lastConfirmed: &timeUpdate{},
	}
}

func (c *Clock) Events() *clock.Events {
	return c.events
}

// AcceptedTime returns the Time of the last accepted Block.
func (c *Clock) AcceptedTime() (acceptedTime time.Time) {
	return c.lastAccepted.Time()
}

// SetAcceptedTime sets the Time of the last accepted Block.
func (c *Clock) SetAcceptedTime(acceptedTime time.Time) (updated bool) {
	now := time.Now()
	if updated = c.lastAccepted.Update(now, acceptedTime); updated {
		c.events.AcceptanceTimeUpdated.Trigger(acceptedTime, now)
	}

	return
}

// RelativeAcceptedTime returns the real-Time adjusted version of the Time of the last accepted Block.
func (c *Clock) RelativeAcceptedTime() (relativeAcceptedTime time.Time) {
	return c.lastAccepted.RelativeTime()
}

// ConfirmedTime returns the Time of the last confirmed Block.
func (c *Clock) ConfirmedTime() (confirmedTime time.Time) {
	return c.lastConfirmed.Time()
}

// SetConfirmedTime sets the Time of the last confirmed Block.
func (c *Clock) SetConfirmedTime(confirmedTime time.Time) (updated bool) {
	now := time.Now()
	if updated = c.lastConfirmed.Update(now, confirmedTime); updated {
		c.events.ConfirmedTimeUpdated.Trigger(confirmedTime, now)
	}

	return
}

// RelativeConfirmedTime returns the real-Time adjusted version of the Time of the last confirmed Block.
func (c *Clock) RelativeConfirmedTime() (relativeConfirmedTime time.Time) {
	return c.lastConfirmed.RelativeTime()
}

var _ clock.Clock = &Clock{}
