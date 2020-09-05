package time

import (
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

func TestAbstractSlotDate_CompareEqualityOfSlotDates_mustReturnTrue(t *testing.T) {
	assert.True(t, NewAbstractSlotDate(2, 15).SameAs(NewAbstractSlotDate(2, 15)),
		"Both slots are the same, and same as must return true.")
}

func TestAbstractSlotDate_CompareEqualityOfSlotDates_mustReturnFalse(t *testing.T) {
	assert.False(t, NewAbstractSlotDate(2, 15).SameAs(NewAbstractSlotDate(3, 15)),
		"Both slots are noz the same, and same as must return false.")
}

func TestAbstractSlotDate_IsSlotDateABeforeSlotDateB_mustReturnTrue(t *testing.T) {
	assert.True(t, NewAbstractSlotDate(2, 15).Before(NewAbstractSlotDate(2, 16)),
		"Slot date 2.15 is before 2.16, so before must return true.")
}

func TestAbstractSlotDate_IsSlotDateABeforeSlotDateB_mustReturnFalse(t *testing.T) {
	assert.False(t, NewAbstractSlotDate(2, 15).Before(NewAbstractSlotDate(2, 15)),
		"Slot date 2.15 is the same as 2.15, so before must return false.")
}

func TestAbstractSlotDate_IsSlotDateAAfterSlotDateB_mustReturnTrue(t *testing.T) {
	assert.True(t, NewAbstractSlotDate(100, 17).After(NewAbstractSlotDate(100, 16)),
		"Slot date 100.7 is after 100.16, so after must return true.")
}

func TestAbstractSlotDate_IsSlotDateAAfterSlotDateB_mustReturnFalse(t *testing.T) {
	assert.False(t, NewAbstractSlotDate(2, 15).After(NewAbstractSlotDate(2, 15)),
		"Slot date 2.15 is the same as 2.15, so after must return false.")
}

func TestAbstractSlotDate_ADiffB_mustReturnPositive(t *testing.T) {
	// setup
	genesisTime, _ := time.Parse(time.RFC3339, "2019-12-13T19:13:37+00:00")
	settings := Settings{
		GenesisBlockDateTime: genesisTime,
		SlotsPerEpoch:        uint64(43200),
		SlotDuration:         time.Duration(2) * time.Second,
	}
	dateA, err := NewConcreteSlotDate(17, 1200, settings)
	if assert.Nil(t, err) {
		dateB, err := NewConcreteSlotDate(16, 35600, settings)
		if assert.Nil(t, err) {
			assert.Equal(t, int64(8800), dateA.Diff(dateB).Int64(),
				"The difference between the two slot dates must be '19600'")
		}
	}
}

func TestAbstractSlotDate_ADiffSubsequentB_mustReturnZero(t *testing.T) {
	genesisTime, _ := time.Parse(time.RFC3339, "2019-12-13T19:13:37+00:00")
	settings := Settings{
		GenesisBlockDateTime: genesisTime,
		SlotsPerEpoch:        uint64(43200),
		SlotDuration:         time.Duration(2) * time.Second,
	}
	dateA, err := NewConcreteSlotDate(18, 43199, settings)
	if assert.Nil(t, err) {
		dateB, err := NewConcreteSlotDate(18, 43199, settings)
		if assert.Nil(t, err) {
			assert.Equal(t, int64(0), dateA.Diff(dateB).Int64(),
				"The difference between the two slot dates must be '19600'")
		}
	}
}

func TestAbstractSlotDate_ADiffPreviousB_mustReturnMinusTwo(t *testing.T) {
	genesisTime, _ := time.Parse(time.RFC3339, "2019-12-13T19:13:37+00:00")
	settings := Settings{
		GenesisBlockDateTime: genesisTime,
		SlotsPerEpoch: uint64(43200),
		SlotDuration: time.Duration(2) * time.Second,
	}
	dateA, err := NewConcreteSlotDate(17, 43199, settings)
	if assert.Nil(t, err) {
		dateB, err := NewConcreteSlotDate(18, 1, settings)
		if assert.Nil(t, err) {
			assert.Equal(t, int64(-2), dateA.Diff(dateB).Int64(),
				"The difference between the two slot dates must be '19600'")
		}
	}
}

func TestAbstractSlotDate_ADiffB_mustReturnNegative(t *testing.T) {
	genesisTime, _ := time.Parse(time.RFC3339, "2019-12-13T19:13:37+00:00")
	settings := Settings{
		GenesisBlockDateTime: genesisTime,
		SlotsPerEpoch: uint64(43200),
		SlotDuration: time.Duration(2) * time.Second,
	}
	dateA, err := NewConcreteSlotDate(8, 40100, settings)
	if assert.Nil(t, err) {
		dateB, err := NewConcreteSlotDate(10, 35600, settings)
		if assert.Nil(t, err) {
			assert.Equal(t, int64(-81900), dateA.Diff(dateB).Int64(),
				"The difference between the two slot dates must be '-81900'.")
		}
	}
}

func TestConcreteSlotDate_GetStartAndEndTime(t *testing.T) {
	// setup
	genesisTime, _ := time.Parse(time.RFC3339, "2019-12-13T19:13:37+00:00")
	settings := Settings{
		GenesisBlockDateTime: genesisTime,
		SlotsPerEpoch: uint64(43200),
		SlotDuration: time.Duration(2) * time.Second,
	}
	date, err := NewConcreteSlotDate(17, 10653, settings)
	if assert.Nil(t, err) {
		// test
		expectedStart, _ := time.Parse(time.RFC3339, "2019-12-31T02:08:43+01:00")
		diff := date.GetStartDateTime().Sub(expectedStart)
		assert.Equal(t, time.Duration(0), diff,
			"The start time must be at '2019-12-31T02:08:43+01:00', but there was a '%s' difference.", diff.String())
	}
}
