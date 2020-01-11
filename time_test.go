package cardano

import (
    "encoding/json"
    "github.com/stretchr/testify/assert"
    "math/big"
    "testing"
    "time"
)

func TestPlainSlotDate_ParseCorrectSlotDateString(t *testing.T) {
    plainSlotDate, err := ParsePlainData("4.15")
    if assert.Nil(t, err, "Parsing must be successful, but throw error.") {
        if assert.NotNil(t, plainSlotDate, "Result of the parsing must not be nil.") {
            assert.Equal(t, plainSlotDate.GetEpoch().Uint64(), uint64(4))
            assert.Equal(t, plainSlotDate.GetSlot().Uint64(), uint64(15))
        }
    }
}

func TestPlainSlotDate_ParseSlotDateStringInWrongFormat_mustThrowError(t *testing.T) {
    _, err := ParsePlainData("1-17")
    assert.NotNil(t, err, "'1-17' is in the wrong format and parser must return an error.")

    _, err = ParsePlainData("1.A")
    assert.NotNil(t, err, "'1.A' is in the wrong format and parser must return an error.")

    _, err = ParsePlainData("A.444")
    assert.NotNil(t, err, "'A.444' is in the wrong format and parser must return an error.")

    _, err = ParsePlainData("-1.444")
    assert.NotNil(t, err, "'-1.444' is in the wrong format and parser must return an error.")

    _, err = ParsePlainData("2.-666")
    assert.NotNil(t, err, "'2.-666' is in the wrong format and parser must return an error.")
}

func TestPlainSlotDate_MarshallPlainSlotDate(t *testing.T) {
    var slotDate, _ = PlainSlotDateFrom(new(big.Int).SetInt64(2), new(big.Int).SetInt64(15))
    data, err := json.Marshal(slotDate)
    if assert.Nil(t, err, "Marshalling the slot date failed.") {
        if assert.NotNil(t, data, "Marshaled data must not be nil.") {
            assert.Equal(t, string(data), "\"2.15\"")
        }
    }
}

type data struct {
    BlockHeight uint64         `json:"block-height"`
    SlotDate    *PlainSlotDate `json:"slot-date"`
}

func TestPlainSlotDate_MarshallStruct(t *testing.T) {
    d := data{BlockHeight: 12, SlotDate: PlainSlotDateFromInt(2, 15)}
    data, err := json.Marshal(d)
    if assert.Nil(t, err, "Marshalling the slot date failed.") {
        if assert.NotNil(t, data, "Marshaled data must not be nil.") {
            assert.Contains(t, string(data), "\"slot-date\":\"2.15\"")
        }
    }
}

func TestPlainSlotDate_UnmarshalStruct(t *testing.T) {
    dataString := "{\"block-height\":12,\"slot-date\":\"2.15\"}"
    var dataStruct data
    err := json.Unmarshal([]byte(dataString), &dataStruct)
    if assert.Nil(t, err, "Unmarshal of the struct with plain slot date failed.") {
        if assert.NotNil(t, dataStruct.SlotDate, "The plain slot date in the struct must not be nil.") {
            assert.Equal(t, dataStruct.SlotDate.GetEpoch().Uint64(), uint64(2))
            assert.Equal(t, dataStruct.SlotDate.GetSlot().Uint64(), uint64(15))
        }
        assert.Equal(t, dataStruct.BlockHeight, uint64(12))
    }
}

func TestPlainSlotDate_CompareEqualityOfSlotDates_mustReturnTrue(t *testing.T) {
    assert.True(t, PlainSlotDateFromInt(2, 15).SameAs(PlainSlotDateFromInt(2, 15)),
        "Both slots are the same, and same as must return true.")
}

func TestPlainSlotDate_CompareEqualityOfSlotDates_mustReturnFalse(t *testing.T) {
    assert.False(t, PlainSlotDateFromInt(2, 15).SameAs(PlainSlotDateFromInt(3, 15)),
        "Both slots are noz the same, and same as must return false.")
}

func TestPlainSlotDate_IsSlotDateABeforeSlotDateB_mustReturnTrue(t *testing.T) {
    assert.True(t, PlainSlotDateFromInt(2, 15).Before(PlainSlotDateFromInt(2, 16)),
        "Slot date 2.15 is before 2.16, so before must return true.")
}

func TestPlainSlotDate_IsSlotDateABeforeSlotDateB_mustReturnFalse(t *testing.T) {
    assert.False(t, PlainSlotDateFromInt(2, 15).Before(PlainSlotDateFromInt(2, 15)),
        "Slot date 2.15 is the same as 2.15, so before must return false.")
}

func TestPlainSlotDate_IsSlotDateAAfterSlotDateB_mustReturnTrue(t *testing.T) {
    assert.True(t, PlainSlotDateFromInt(100, 17).After(PlainSlotDateFromInt(100, 16)),
        "Slot date 100.7 is after 100.16, so after must return true.")
}

func TestPlainSlotDate_IsSlotDateAAfterSlotDateB_mustReturnFalse(t *testing.T) {
    assert.False(t, PlainSlotDateFromInt(2, 15).After(PlainSlotDateFromInt(2, 15)),
        "Slot date 2.15 is the same as 2.15, so after must return false.")
}

func TestPlainSlotDate_ADiffB_mustReturnPositive(t *testing.T) {
    // setup
    genesisTime, _ := time.Parse(time.RFC3339, "2019-12-13T19:13:37+00:00")
    settings := TimeSettings{GenesisBlockDateTime: genesisTime, SlotsPerEpoch: new(big.Int).SetInt64(43200), SlotDuration: time.Duration(2) * time.Second}
    dateA := FullSlotDateFromInt(17, 1200, settings)
    dateB := FullSlotDateFromInt(16, 35600, settings)
    assert.Equal(t, dateA.Diff(dateB).Int64(), int64(8800), "The difference between the two slot dates must be '19600'")
}

func TestPlainSlotDate_ADiffSubsequentB_mustReturnNegative(t *testing.T) {
    genesisTime, _ := time.Parse(time.RFC3339, "2019-12-13T19:13:37+00:00")
    settings := TimeSettings{GenesisBlockDateTime: genesisTime, SlotsPerEpoch: new(big.Int).SetInt64(43200), SlotDuration: time.Duration(2) * time.Second}
    dateA := FullSlotDateFromInt(18, 43200, settings)
    dateB := FullSlotDateFromInt(18, 43200, settings)
    assert.Equal(t, dateA.Diff(dateB).Int64(), int64(0), "The difference between the two slot dates must be '19600'")
}

func TestPlainSlotDate_ADiffPreviousB_mustReturn0(t *testing.T) {
    genesisTime, _ := time.Parse(time.RFC3339, "2019-12-13T19:13:37+00:00")
    settings := TimeSettings{GenesisBlockDateTime: genesisTime, SlotsPerEpoch: new(big.Int).SetInt64(43200), SlotDuration: time.Duration(2) * time.Second}
    dateA := FullSlotDateFromInt(17, 43200, settings)
    dateB := FullSlotDateFromInt(18, 1, settings)
    assert.Equal(t, dateA.Diff(dateB).Int64(), int64(-1), "The difference between the two slot dates must be '19600'")
}

func TestPlainSlotDate_ADiffB_mustReturnNegative(t *testing.T) {
    genesisTime, _ := time.Parse(time.RFC3339, "2019-12-13T19:13:37+00:00")
    settings := TimeSettings{GenesisBlockDateTime: genesisTime, SlotsPerEpoch: new(big.Int).SetInt64(43200), SlotDuration: time.Duration(2) * time.Second}
    dateA := FullSlotDateFromInt(8, 40100, settings)
    dateB := FullSlotDateFromInt(10, 35600, settings)
    assert.Equal(t, dateA.Diff(dateB).Int64(), int64(-81900), "The difference between the two slot dates must be '-81900'.")
}

func TestFullSlotDate_GetStartAndEndTime(t *testing.T) {
    // setup
    genesisTime, _ := time.Parse(time.RFC3339, "2019-12-13T19:13:37+00:00")
    settings := TimeSettings{GenesisBlockDateTime: genesisTime, SlotsPerEpoch: new(big.Int).SetInt64(43200), SlotDuration: time.Duration(2) * time.Second}
    date := FullSlotDateFromInt(17, 10653, settings)
    // test
    expectedStart, _ := time.Parse(time.RFC3339, "2019-12-31T02:08:43+01:00")
    diff := date.GetStartDateTime().Sub(expectedStart)
    assert.Equal(t, diff, time.Duration(0), "The start time must be at '2019-12-31T02:08:43+01:00', but there was a '%s' difference.", diff.String())
}
