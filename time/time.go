package time

import (
	"fmt"
	"math"
	"math/big"
	"time"
)

// AbstractSlotDate is determined by its specific epoch and slot number. However, it is missing further
// information about the time settings of the Cardano blockchain. An AbstractSlotDate doesn't know the number of
// slots per epoch as well as the duration of a slot. Moreover, the creation time of the genesis block is unknown.
// Hence, a number of methods are not available in the AbstractSlotDate, but can be used if the AbstractSlotDate
// is transformed into a ConcreteSlotDate by passing the required time settings of the blockchain.
type AbstractSlotDate struct {
	epoch uint64
	slot  uint64
}

// ConcreteSlotDate contains all information of an AbstractSlotDate, but has additionally information about the
// time settings of a concrete blockchain (it exists in).
type ConcreteSlotDate struct {
	AbstractSlotDate
	// time settings of the concrete blockchain
	timeSettings Settings
}

// Settings of a concrete Cardano blockchain relevant for time.
type Settings struct {
	// creation time of the genesis block.
	GenesisBlockDateTime time.Time
	// the number of slots an epoch is divided into.
	SlotsPerEpoch uint64
	// a slot has a fixed duration.
	SlotDuration time.Duration
}

// SlotOutOfBoundsError describes an error, where a given slot number exceeds
// the number of slots per epoch for a concrete blockchain.
type SlotOutOfBoundsError struct {
	Slot          uint64
	SlotsPerEpoch uint64
}

func (e SlotOutOfBoundsError) Error() string {
	return fmt.Sprintf("slot number %d exceeds the number of slots per epoch (%d)", e.Slot, e.SlotsPerEpoch)
}

// NewAbstractSlotDate creates a new abstract slot date from the given epoch and slot number.
func NewAbstractSlotDate(epoch uint64, slot uint64) *AbstractSlotDate {
	return &AbstractSlotDate{epoch: epoch, slot: slot}
}

// NewConcreteSlotDate creates a new concrete slot date from the given epoch and slot number as well as
// the Settings of the concrete blockchain. A SlotOutOfBoundsError will be returned, if the given slot
// number exceeds the number of slots per epoch in the given Settings.
func NewConcreteSlotDate(epoch uint64, slot uint64, settings Settings) (*ConcreteSlotDate, error) {
	if slot >= settings.SlotsPerEpoch {
		return nil, &SlotOutOfBoundsError{
			Slot:          slot,
			SlotsPerEpoch: settings.SlotsPerEpoch,
		}
	}
	return &ConcreteSlotDate{AbstractSlotDate: AbstractSlotDate{epoch: epoch, slot: slot}, timeSettings: settings}, nil
}

// MaterializeSlotDate transforms an AbstractSlotDate into a ConcreteSlotDate by passing the Settings
// of the concrete blockchain. A SlotOutOfBoundsError will be returned, if the slot number of the given
// AbstractSlotDate exceeds the number of slots per epoch in the given Settings.
func MaterializeSlotDate(abstractDate *AbstractSlotDate, settings Settings) (*ConcreteSlotDate, error) {
	return NewConcreteSlotDate(abstractDate.epoch, abstractDate.slot, settings)
}

// GetEpoch gets the epoch of an AbstractSlotDate.
func (abstractDate *AbstractSlotDate) GetEpoch() uint64 {
	return abstractDate.epoch
}

// GetSlot gets the slot of an AbstractSlotDate.
func (abstractDate *AbstractSlotDate) GetSlot() uint64 {
	return abstractDate.slot
}

// SameAs checks whether this AbstractSlotDate and given other AbstractSlotDate are referring to the same slot date.
func (abstractDate *AbstractSlotDate) SameAs(otherDate *AbstractSlotDate) bool {
	return abstractDate.epoch == otherDate.epoch && abstractDate.slot == otherDate.slot
}

// Before checks whether this AbstractSlotDate is strictly before the given other AbstractSlotDate.
func (abstractDate *AbstractSlotDate) Before(otherDate *AbstractSlotDate) bool {
	return abstractDate.epoch < otherDate.epoch ||
		(abstractDate.epoch == otherDate.epoch && abstractDate.slot < otherDate.slot)
}

// After checks whether this AbstractSlotDate is strictly after the given other AbstractSlotDate.
func (abstractDate *AbstractSlotDate) After(otherDate *AbstractSlotDate) bool {
	return abstractDate.epoch > otherDate.epoch ||
		(abstractDate.epoch == otherDate.epoch && abstractDate.slot > otherDate.slot)
}

// String returns the AbstractSlotDate as a plain string in the format <EPOCH>.<SLOT>
func (abstractDate *AbstractSlotDate) String() string {
	return fmt.Sprintf("%v.%v", abstractDate.GetEpoch(), abstractDate.GetSlot())
}

// Same checks whether this Settings is the same as the given other Settings
func (timeSettings *Settings) Same(otherSettings *Settings) bool {
	return (timeSettings.SlotsPerEpoch == otherSettings.SlotsPerEpoch) &&
		(timeSettings.SlotDuration == otherSettings.SlotDuration) &&
		(timeSettings.GenesisBlockDateTime.Sub(otherSettings.GenesisBlockDateTime) == 0)
}

// GetStartDateTime gets the start time of the given ConcreteSlotDate.
func (concreteDate *ConcreteSlotDate) GetStartDateTime() time.Time {
	thisEpoch := new(big.Int).SetUint64(concreteDate.epoch)
	thisSlot := new(big.Int).SetUint64(concreteDate.slot)
	slotDuration := new(big.Int).SetInt64(int64(concreteDate.timeSettings.SlotDuration))
	maxDuration := new(big.Int).SetUint64(math.MaxInt64)
	timeSettings := new(big.Int).SetUint64(concreteDate.timeSettings.SlotsPerEpoch)

	duration := new(big.Int).Mul(new(big.Int).Add(new(big.Int).Mul(thisEpoch, timeSettings), thisSlot), slotDuration)
	if duration.Cmp(maxDuration) <= 0 {
		return concreteDate.timeSettings.GenesisBlockDateTime.Add(time.Duration(duration.Uint64()))
	}
	//todo: efficient handling of big slot dates
	startTime := concreteDate.timeSettings.GenesisBlockDateTime
	n := new(big.Int).Div(duration, maxDuration).Uint64()
	for i := uint64(0); i < n; i++ {
		startTime = startTime.Add(time.Duration(math.MaxInt64))
	}
	mod := new(big.Int).Mod(duration, maxDuration).Int64()
	return startTime.Add(time.Duration(mod))
}

// GetEndDateTime gets the end time of the given ConcreteSlotDate.
func (concreteDate *ConcreteSlotDate) GetEndDateTime() time.Time {
	return concreteDate.GetStartDateTime().Add(concreteDate.timeSettings.SlotDuration)
}

// Diff computes the difference in slots between this ConcreteSlotDate and the given other ConcreteSlotDate.
func (concreteDate *ConcreteSlotDate) Diff(otherDate *ConcreteSlotDate) *big.Int {
	thisSettings := new(big.Int).SetUint64(concreteDate.timeSettings.SlotsPerEpoch)
	thisEpoch := new(big.Int).SetUint64(concreteDate.epoch)
	otherEpoch := new(big.Int).SetUint64(otherDate.epoch)
	thisSlot := new(big.Int).SetUint64(concreteDate.slot)
	otherSlot := new(big.Int).SetUint64(otherDate.slot)
	a := new(big.Int).Add(new(big.Int).Mul(thisEpoch, thisSettings), thisSlot)
	b := new(big.Int).Add(new(big.Int).Mul(otherEpoch, thisSettings), otherSlot)
	return new(big.Int).Sub(a, b)
}
