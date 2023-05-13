package util

import "time"

type DateTimeFormatter interface {
	ToString(dateTime time.Time) string
	ToTime(dateTime time.Time) (time.Time, error)
	ToTimeFromString(dateTime string) (time.Time, error)
	SetDateTimeFormat(format string)
}

type dateTimeFormatter struct {
	format string
}

func NewDateTimeFormatter(format string) DateTimeFormatter {
	return &dateTimeFormatter{format}
}

func (d *dateTimeFormatter) ToString(dateTime time.Time) string {
	return dateTime.UTC().Format(d.format)
}
func (d *dateTimeFormatter) ToTime(dateTime time.Time) (time.Time, error) {
	return time.Parse(d.format, dateTime.UTC().Format(d.format))
}

func (d *dateTimeFormatter) ToTimeFromString(dateTime string) (time.Time, error) {
	return time.Parse(d.format, dateTime)
}
func (d *dateTimeFormatter) SetDateTimeFormat(format string) {
	d.format = format
}
