package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDateTimeFormatter_ToString(t *testing.T) {
	formatter := NewDateTimeFormatter(time.DateTime)
	dateTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	exp := "2023-01-01 00:00:00"
	act := formatter.ToString(dateTime)

	assert.Equal(t, exp, act)
}

func TestDateTimeFormatter_ToString2(t *testing.T) {
	formatter := NewDateTimeFormatter(time.DateOnly)
	dateTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	exp := "2023-01-01"
	act := formatter.ToString(dateTime)

	assert.Equal(t, exp, act)
}

func TestDateTimeFormatter_ToTime(t *testing.T) {
	formatter := NewDateTimeFormatter(time.DateOnly)
	exp := time.Date(2023, 1, 1, 1, 1, 1, 1, time.UTC)
	formatExp := exp.Format(time.DateOnly)
	act, _ := formatter.ToTime(exp)
	fromString, err := formatter.ToTimeFromString(formatExp)
	assert.NoError(t, err)
	assert.Equal(t, fromString, act)

}

func TestDateTimeFormatter_ToTimeFromString(t *testing.T) {
	formatter := NewDateTimeFormatter(time.DateOnly)
	exp := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	format := exp.Format(time.DateOnly)
	act, err := formatter.ToTimeFromString(format)
	assert.NoError(t, err)
	assert.Equal(t, exp, act)
}

func TestDateTimeFormatter_ToTimeFromString_WrongStr(t *testing.T) {
	formatter := NewDateTimeFormatter(time.DateOnly)
	exp := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	format := "3"
	act, err := formatter.ToTimeFromString(format)
	assert.Error(t, err)
	assert.NotEqual(t, exp, act)
}

func TestDateTimeFormatter_SetDateTimeFormat(t *testing.T) {
	formatter := NewDateTimeFormatter(time.DateOnly)
	format := time.DateTime
	formatter.SetDateTimeFormat(format)
	exp := time.Date(2023, 1, 1, 1, 1, 1, 0, time.UTC)
	act, err := formatter.ToTimeFromString("2023-01-01 01:01:01")
	assert.NoError(t, err)
	assert.Equal(t, exp, act)
}

func FuzzDateTimeFormatter_ToString(f *testing.F) {
	formatter := NewDateTimeFormatter(time.DateTime)
	f.Fuzz(func(t *testing.T, inputSec int64) {
		format := time.Unix(inputSec, 0).Format(time.DateTime)
		_, err := formatter.ToTimeFromString(format)
		assert.NoError(t, err)
	})
}
