package time

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"
)

// ParsingError describes that an AbstractSlotDate cannot be parsed from a text, due to a
// particular reason.
type ParsingError struct {
	// text that should have been parsed.
	ParsedText string
	// reason why it cannot be parsed.
	Reason string
}

func (err ParsingError) Error() string {
	return fmt.Sprintf("failed to parse '%v': %v", err.ParsedText, err.Reason)
}

// ParseAbstractSlotDate parses an AbstractSlotDate from the given text. An AbstractSlotDate must be of
// the format "<EPOCH>.<SLOT>". If parsing fails, an error will be returned.
func ParseAbstractSlotDate(text string) (*AbstractSlotDate, error) {
	seps := strings.Split(text, ".")
	if len(seps) == 2 {
		epoch, err := strconv.ParseUint(seps[0], 10, 64)
		if err == nil {
			slot, err := strconv.ParseUint(seps[1], 10, 64)
			if err == nil {
				return &AbstractSlotDate{epoch: epoch, slot: slot}, nil
			}
			return nil, ParsingError{
				ParsedText: text,
				Reason:     fmt.Sprintf("slot must be a positive number, but was '%v'", seps[1]),
			}
		}
		return nil, ParsingError{
			ParsedText: text,
			Reason:     fmt.Sprintf("epoch must be a positive number, but was '%v'", seps[0]),
		}
	} else {
		return nil, ParsingError{
			ParsedText: text,
			Reason:     "the date must be of the format '<EPOCH>.<SLOT>', where epoch and slot are positive numbers",
		}
	}
}

// MarshalJSON marshals this AbstractSlotDate into a string field with the format "<EPOCH>.<SLOT>".
func (abstractDate *AbstractSlotDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(abstractDate.String())
}

// UnmarshalJSON unmarshalls an AbstractSlotDate that is expected to be string field of the format
// "<EPOCH>.<SLOT>".
func (abstractDate *AbstractSlotDate) UnmarshalJSON(b []byte) error {
	var slotDateString string
	var err = json.Unmarshal(b, &slotDateString)
	if err == nil {
		parsedSlotDate, err := ParseAbstractSlotDate(slotDateString)
		if err == nil {
			abstractDate.epoch = parsedSlotDate.GetEpoch()
			abstractDate.slot = parsedSlotDate.GetSlot()
		}
	}
	return err
}

// OutOfRangeError describes that a given time is out of range of a concrete blockchain, i.e.
// the time is before the genesis block has been created.
type OutOfRangeError struct {
	// specified time.
	Time time.Time
	// when the genesis block has been created.
	GenesisBlockDateTime time.Time
}

func (err OutOfRangeError) Error() string {
	return fmt.Sprintf("the time %v is before the creation time of the genesis block (%v)",
		err.Time, err.GenesisBlockDateTime)
}

// GetSlotDateFor for the given time and Settings of the concrete blockchain.
func (timeSettings *Settings) GetSlotDateFor(t time.Time) (*ConcreteSlotDate, error) {
	if t.Before(timeSettings.GenesisBlockDateTime) {
		return nil, OutOfRangeError{
			Time:                 t,
			GenesisBlockDateTime: timeSettings.GenesisBlockDateTime,
		}
	}
	if !t.After(timeSettings.GenesisBlockDateTime.Add(math.MaxInt64)) {
		diff := t.Sub(timeSettings.GenesisBlockDateTime)
		slots := uint64(diff / timeSettings.SlotDuration)
		epoch := slots / timeSettings.SlotsPerEpoch
		slotsInEpoch := slots % timeSettings.SlotsPerEpoch
		return NewConcreteSlotDate(epoch, slotsInEpoch, *timeSettings)
	} else {
		diff := new(big.Int).SetInt64(math.MaxInt64)
		cursor := timeSettings.GenesisBlockDateTime.Add(math.MaxInt64)
		for t.After(cursor.Add(math.MaxInt64)) {
			cursor = cursor.Add(math.MaxInt64)
			diff = new(big.Int).Add(diff, new(big.Int).SetInt64(math.MaxInt64))
		}
		diff = new(big.Int).Add(diff, new(big.Int).SetInt64(int64(t.Sub(cursor))))
		slots := new(big.Int).Div(diff, new(big.Int).SetInt64(int64(timeSettings.SlotDuration)))
		slotsPerEpoch := new(big.Int).SetUint64(timeSettings.SlotsPerEpoch)
		epoch := new(big.Int).Div(slots, slotsPerEpoch)
		slotsInEpoch := new(big.Int).Mod(slots, slotsPerEpoch)
		return NewConcreteSlotDate(epoch.Uint64(), slotsInEpoch.Uint64(), *timeSettings)
	}
}
