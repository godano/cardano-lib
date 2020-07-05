package cardano

import (
    "encoding/json"
    "fmt"
    "math/big"
    "strings"
    "time"
)

type SlotDate interface {
    GetEpoch() *big.Int
    GetSlot() *big.Int
}

// a plain slot date is only determined by its
// epoch and slot. A plain date is missing further
// information about the time settings of the Cardano
// blockchain.
//
// A plain date does not know the number of slots
// in  an epoch as well as the duration of a slot
// in milliseconds. Moreover, the creation time of the
// genesis block is unknown. Hence, a number of methods
// are not available in the plain date, but can be
// used if the plain date is transformed into a full
// date by passing the required time settings of the
// blockchain.
type PlainSlotDate struct {
    epoch *big.Int
    slot  *big.Int
}

// creates the plain slot date from the given epoch and slot number, which
// must both be positive.
func PlainSlotDateFrom(epoch *big.Int, slot *big.Int) (*PlainSlotDate, error) {
    if epoch.Cmp(new(big.Int).SetInt64(0)) < 0 || slot.Cmp(new(big.Int).SetInt64(0)) < 0 {
        return nil, invalidArgument{
            MethodName: "PlainSlotDateFrom",
            Expected:   "The epoch and slot number must be positive.",
        }
    }
    return &PlainSlotDate{epoch: epoch, slot: slot}, nil
}

// creates the plain slot date from the given epoch and slot number.
func PlainSlotDateFromInt(epoch uint64, slot uint64) *PlainSlotDate {
    return &PlainSlotDate{epoch: new(big.Int).SetUint64(epoch), slot: new(big.Int).SetUint64(slot)}
}

// gets the epoch of the date.
func (slotDate *PlainSlotDate) GetEpoch() *big.Int {
    return slotDate.epoch
}

// gets the slot of the date.
func (slotDate *PlainSlotDate) GetSlot() *big.Int {
    return slotDate.slot
}

// returns true, if this date is equal to the given
// other date, otherwise false.
func (slotDate *PlainSlotDate) SameAs(otherDate *PlainSlotDate) bool {
    return (slotDate.GetEpoch().Cmp(otherDate.GetEpoch()) == 0) && (slotDate.GetSlot().Cmp(otherDate.GetSlot()) == 0)
}

// returns true, if this date lies strictly before
// the given other date, otherwise false.
func (slotDate *PlainSlotDate) Before(otherDate *PlainSlotDate) bool {
    if slotDate.GetEpoch().Cmp(otherDate.GetEpoch()) < 0 {
        return true
    } else if slotDate.GetEpoch().Cmp(otherDate.GetEpoch()) == 0 {
        if slotDate.GetSlot().Cmp(otherDate.GetSlot()) < 0 {
            return true
        }
    }
    return false
}

// returns true, if this date lies after the given
// other date.
func (slotDate *PlainSlotDate) After(otherDate *PlainSlotDate) bool {
    if slotDate.GetEpoch().Cmp(otherDate.GetEpoch()) > 0 {
        return true
    } else if slotDate.GetEpoch().Cmp(otherDate.GetEpoch()) == 0 {
        if slotDate.GetSlot().Cmp(otherDate.GetSlot()) > 0 {
            return true
        }
    }
    return false
}

// returns the slot date as a plain string in the
// format <EPOCH>.<SLOT>
func (slotDate *PlainSlotDate) String() string {
    return fmt.Sprintf("%v.%v", slotDate.GetEpoch(), slotDate.GetSlot())
}

// parses plain slot date from the given text. A date must
// be of the format "<EPOCH>.<SLOT>". if parsing fails,
// an error will be returned.
func ParsePlainData(text string) (*PlainSlotDate, error) {
    seps := strings.Split(text, ".")
    if len(seps) == 2 {
        epoch, success := new(big.Int).SetString(seps[0], 10)
        if success {
            if epoch.Cmp(new(big.Int).SetInt64(0)) >= 0 {
                slot, success := new(big.Int).SetString(seps[1], 10)
                if success {
                    if slot.Cmp(new(big.Int).SetInt64(0)) >= 0 {
                        return &PlainSlotDate{epoch: epoch, slot: slot}, nil
                    }
                }
                return nil, parsingError{ParsedText: text, Reason: fmt.Sprintf("Slot must be a positive number, but was '%v'.", seps[1])}
            }
        }
        return nil, parsingError{ParsedText: text, Reason: fmt.Sprintf("Epoch must be a positive number, but was '%v'.", seps[0])}
    } else {
        return nil, parsingError{ParsedText: text, Reason: "The date must be of the format '<EPOCH>.<SLOT>', where epoch and slot are positive numbers."}
    }
}

// JSON Marshaling of plain slot date.
func (slotDate *PlainSlotDate) MarshalJSON() ([]byte, error) {
    return json.Marshal(slotDate.String())
}

// JSON Unmarshaling of plain slot date.
func (slotDate *PlainSlotDate) UnmarshalJSON(b []byte) error {
    var slotDateString string
    var err = json.Unmarshal(b, &slotDateString)
    if err == nil {
        var parsedSlotDate, err = ParsePlainData(slotDateString)
        if err == nil {
            slotDate.epoch = parsedSlotDate.GetEpoch()
            slotDate.slot = parsedSlotDate.GetSlot()
            return nil
        }
    }
    return err
}

// a full slot data contains information of the plain slot
// date including the epoch and slot number together with
// the time settings of the blockchain.
type FullSlotDate struct {
    PlainSlotDate
    timeSettings TimeSettings
}

// settings of the Cardano blockchain relevant for time.
type TimeSettings struct {
    // creation time of the genesis block.
    GenesisBlockDateTime time.Time
    // the number of slots an epoch is divided into.
    SlotsPerEpoch *big.Int
    // a slot has fixed duration, and thus has a start and end date.
    SlotDuration time.Duration
}

// returns the number of slots that are between this
// date and the other date.
func (date *FullSlotDate) Diff(otherDate SlotDate) *big.Int {
    a := new(big.Int).Add(new(big.Int).Mul(date.GetEpoch(), date.timeSettings.SlotsPerEpoch), date.GetSlot())
    b := new(big.Int).Add(new(big.Int).Mul(otherDate.GetEpoch(), date.timeSettings.SlotsPerEpoch), otherDate.GetSlot())
    return new(big.Int).Sub(a, b)
}

// gets an instant in time with nanosecond precision of
// the start of the slot.
func (date *FullSlotDate) GetStartDateTime() time.Time {
    slots := new(big.Int).Add(new(big.Int).Mul(date.GetEpoch(), date.timeSettings.SlotsPerEpoch), date.GetSlot()).Uint64()
    return date.timeSettings.GenesisBlockDateTime.Add(time.Duration(slots) * date.timeSettings.SlotDuration)
}

// gets an instant in time with nanosecond precision of
// the end of the slot.
func (date *FullSlotDate) GetEndDateTime() time.Time {
    return date.GetStartDateTime().Add(date.timeSettings.SlotDuration)
}

func FullSlotDateFrom(epoch *big.Int, slot *big.Int, settings TimeSettings) (*FullSlotDate, error) {
    if epoch.Cmp(new(big.Int).SetInt64(0)) < 0 || slot.Cmp(new(big.Int).SetInt64(0)) < 0 {
        return nil, invalidArgument{
            MethodName: "PlainSlotDateFrom",
            Expected:   "The epoch and slot number must be positive.",
        }
    }
    return &FullSlotDate{PlainSlotDate: PlainSlotDate{epoch: epoch, slot: slot}, timeSettings: settings}, nil
}

func FullSlotDateFromInt(epoch uint64, slot uint64, settings TimeSettings) *FullSlotDate {
    return &FullSlotDate{PlainSlotDate: PlainSlotDate{epoch: new(big.Int).SetUint64(epoch), slot: new(big.Int).SetUint64(slot)}, timeSettings: settings}
}

func MakeFullSlotDate(plainDate *PlainSlotDate, settings TimeSettings) *FullSlotDate {
    return &FullSlotDate{PlainSlotDate: PlainSlotDate{epoch: plainDate.GetEpoch(), slot: plainDate.GetSlot()}, timeSettings: settings}
}

// gets the full slot date for the given time.
func (timeSettings *TimeSettings) GetSlotDateFor(t time.Time) (*FullSlotDate, error) {
    if t.Before(timeSettings.GenesisBlockDateTime) {
        return nil, invalidArgument{
            MethodName: "GetSlotDateFor",
            Expected:   fmt.Sprintf("The given time \"%v\" must not be before the creation time of the genesis block %v.", t, timeSettings.GenesisBlockDateTime),
        }
    }
    diff := t.Sub(timeSettings.GenesisBlockDateTime)
    totalSlots := new(big.Int).SetInt64(int64(diff / timeSettings.SlotDuration))
    epoch := new(big.Int).Div(totalSlots, timeSettings.SlotsPerEpoch)
    slotsInEpoch := new(big.Int).Mod(totalSlots, timeSettings.SlotsPerEpoch)
    return FullSlotDateFrom(epoch, slotsInEpoch, *timeSettings)
}
