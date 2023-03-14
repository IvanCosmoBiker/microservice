package datetime

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseTimeFromKass(t *testing.T) {
	date, errTime := Init()
	if errTime != nil {
		t.Error(errTime)
	}
	var timeKass = "2022-10-03T16:59:00"
	timeStaring1, err1 := date.ParseDateOfLayout("2006-01-02T15:04:05", timeKass)
	if err1 != nil {
		t.Error(err1)
	}
	timeUtc, errSub := date.SubtractFromTime(timeStaring1, 10800)
	if errSub != nil {
		t.Error(errSub)
	}
	assert.Equal(t, "2022-10-03 13:59:00", timeUtc, "they should be equal")
	timeStaring, err := date.ParseDateOfLayout("2006-01-02T15:04:05", timeKass)
	if err != nil {
		t.Error(err)
	}
	timeUtc, errSub = date.SubtractFromTime(timeStaring, 10800)
	if errSub != nil {
		t.Error(errSub)
	}
	assert.Equal(t, "2022-10-03 13:59:00", timeUtc, "they should be equal")
}

func TestParseTimeOfKass(t *testing.T) {
	date, errTime := Init()
	if errTime != nil {
		t.Error(errTime)
	}
	stringTime := "2023-03-01T13:34:00"
	timeStaring1, err1 := date.ParseDateOfLayout("2006-01-02T15:04:05", stringTime)
	if err1 != nil {
		t.Error(err1)
	}
	assert.Equal(t, "2023-03-01 13:34:00", timeStaring1, "they should be equal")
}
