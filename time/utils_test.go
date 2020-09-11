package time

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAbstractSlotDate_ParseCorrectSlotDateString(t *testing.T) {
	AbstractSlotDate, err := ParseAbstractSlotDate("4.15")
	if assert.Nil(t, err, "Parsing must be successful, but returns error.") {
		if assert.NotNil(t, AbstractSlotDate, "Result of the parsing must not be nil.") {
			assert.Equal(t, AbstractSlotDate.GetEpoch(), uint64(4))
			assert.Equal(t, AbstractSlotDate.GetSlot(), uint64(15))
		}
	}
}

func TestAbstractSlotDate_ParseSlotDateStringInWrongFormat_mustThrowError(t *testing.T) {
	_, err := ParseAbstractSlotDate("1-17")
	assert.NotNil(t, err, "'1-17' is in the wrong format and parser must return an error.")

	_, err = ParseAbstractSlotDate("1.A")
	assert.NotNil(t, err, "'1.A' is in the wrong format and parser must return an error.")

	_, err = ParseAbstractSlotDate("A.444")
	assert.NotNil(t, err, "'A.444' is in the wrong format and parser must return an error.")

	_, err = ParseAbstractSlotDate("-1.444")
	assert.NotNil(t, err, "'-1.444' is in the wrong format and parser must return an error.")

	_, err = ParseAbstractSlotDate("2.-666")
	assert.NotNil(t, err, "'2.-666' is in the wrong format and parser must return an error.")
}

func TestAbstractSlotDate_MarshallSlotDate(t *testing.T) {
	slotDate := NewAbstractSlotDate(2, 15)
	data, err := json.Marshal(slotDate)
	if assert.Nil(t, err, "Marshalling the slot date failed.") {
		if assert.NotNil(t, data, "Marshaled data must not be nil.") {
			assert.Equal(t, string(data), "\"2.15\"")
		}
	}
}

type data struct {
	BlockHeight uint64            `json:"block-height"`
	SlotDate    *AbstractSlotDate `json:"slot-date"`
}

func TestAbstractSlotDate_MarshallStruct(t *testing.T) {
	d := data{BlockHeight: 12, SlotDate: NewAbstractSlotDate(2, 15)}
	data, err := json.Marshal(d)
	if assert.Nil(t, err, "Marshalling the slot date failed.") {
		if assert.NotNil(t, data, "Marshaled data must not be nil.") {
			assert.Contains(t, string(data), "\"slot-date\":\"2.15\"")
		}
	}
}

func TestAbstractSlotDate_UnmarshalStruct(t *testing.T) {
	dataString := "{\"block-height\":12,\"slot-date\":\"2.15\"}"
	var dataStruct data
	err := json.Unmarshal([]byte(dataString), &dataStruct)
	if assert.Nil(t, err, "Unmarshal of the struct with plain slot date failed.") {
		if assert.NotNil(t, dataStruct.SlotDate, "The plain slot date in the struct must not be nil.") {
			assert.Equal(t, dataStruct.SlotDate.GetEpoch(), uint64(2))
			assert.Equal(t, dataStruct.SlotDate.GetSlot(), uint64(15))
		}
		assert.Equal(t, dataStruct.BlockHeight, uint64(12))
	}
}

func TestTimeSettings_GetConcreteSlotDateForGenesisCreationTime(t *testing.T) {
	// setup
	genesisTime, _ := time.Parse(time.RFC3339, "2019-12-13T19:13:37+00:00")
	settings := Settings{
		GenesisBlockDateTime: genesisTime,
		SlotsPerEpoch: uint64(43200),
		SlotDuration: time.Duration(2) * time.Second,
	}

	data, err := settings.GetSlotDateFor(genesisTime)
	if assert.Nil(t, err, "Method must not return an error.") {
		assert.Equal(t, uint64(0), data.GetEpoch(), "The epoch must be 0.")
		assert.Equal(t, uint64(0), data.GetSlot(), "The epoch must be 0.")
	}
}

func TestTimeSettings_GetConcreteSlotDateForArbitraryTime(t *testing.T) {
	// setup
	genesisTime, _ := time.Parse(time.RFC3339, "2019-12-13T19:13:37+00:00")
	randomTime, _ := time.Parse(time.RFC3339, "2020-01-14T16:51:37+01:00")
	settings := Settings{
		GenesisBlockDateTime: genesisTime,
		SlotsPerEpoch: uint64(43200),
		SlotDuration: time.Duration(2) * time.Second,
	}

	data, err := settings.GetSlotDateFor(randomTime)
	if assert.Nil(t, err, "Method must not return an error.") {
		assert.Equal(t, uint64(31), data.GetEpoch(), "The epoch must be 31.")
		assert.Equal(t, uint64(37140), data.GetSlot(), "The epoch must be 37140.")
	}
}

func TestTimeSettings_GetConcreteSlotDateOfTimeBeforeGenesis_MustReturnError(t *testing.T) {
	// setup
	genesisTime, _ := time.Parse(time.RFC3339, "2019-12-13T19:13:37+00:00")
	randomTime, _ := time.Parse(time.RFC3339, "2019-01-14T16:51:37+01:00")
	settings := Settings{
		GenesisBlockDateTime: genesisTime,
		SlotsPerEpoch: uint64(43200),
		SlotDuration: time.Duration(2) * time.Second,
	}

	_, err := settings.GetSlotDateFor(randomTime)
	assert.NotNil(t, err, "Method must return an error, because time is before the genesis block.")
}
