package datetime

import (
	"fmt"
	"strings"
	"time"
)

type DateTime struct {
	Format   string
	Timezone *time.Location
}

func Init() (*DateTime, error) {
	location, err := time.LoadLocation("UTC")
	if err != nil {
		return &DateTime{}, err
	}
	return &DateTime{
		Timezone: location,
		Format:   "2006-01-02 15:04:05",
	}, nil
}

func (d *DateTime) Now() string {
	now := time.Now().In(d.Timezone)
	return now.Format(d.Format)
}

func (d *DateTime) UnixNano() string {
	now := time.Now().In(d.Timezone).UnixNano()
	return fmt.Sprintf("%v", now)
}

func (d *DateTime) Unix() int64 {
	return time.Now().In(d.Timezone).Unix()
}

func (d *DateTime) StringToUnix(date string) int64 {
	var timeLocal time.Time
	var err error
	timeLocal, err = time.ParseInLocation(d.Format, date, d.Timezone)
	if err != nil {
		return 0
	}
	return timeLocal.Unix()
}

func (d *DateTime) Parse(date string) (string, error) {
	var timeLocal time.Time
	var err error
	timeLocal, err = time.ParseInLocation(d.Format, date, d.Timezone)
	if err != nil {
		return "", err
	}
	return timeLocal.Format(d.Format), nil
}

func (d *DateTime) ParseDateOfDefaultLayout(date string) (string, error) {
	t, err := time.ParseInLocation(d.Format, date, d.Timezone)
	if err != nil {
		return "", err
	}
	return t.Format(d.Format), nil
}

func (d *DateTime) ParseDateOfLayout(layout, date string) (string, error) {
	t, err := time.ParseInLocation(layout, date, d.Timezone)
	if err != nil {
		return "", err
	}
	return t.Format(d.Format), nil
}

func (d *DateTime) AddSeconds(date string, seconds int) (string, error) {
	stringNew := strings.ReplaceAll(date, "T", " ")
	t, err := time.ParseInLocation(d.Format, stringNew, d.Timezone)
	if err != nil {
		return "", err
	}
	t.Add(time.Duration(seconds) * time.Second)
	return t.Format(d.Format), nil
}

func (d *DateTime) SubtractFromTime(date string, count time.Duration) (string, error) {
	t, errParse := time.Parse(d.Format, date)
	if errParse != nil {
		return "", errParse
	}
	timeSub := t.Add(-count * time.Second)
	return timeSub.Format(d.Format), nil
}

func (d *DateTime) ParseDateAndSubtractHour(date string) (string, error) {
	t, err := time.ParseInLocation(d.Format, date, d.Timezone)
	if err != nil {
		return "", err
	}
	timeSub := t.Add(time.Duration(-3) * time.Hour)
	return timeSub.Format(d.Format), nil
}

// func (d *DateTime) ParseDateWithTimeZoneAndLayout(date, timezone, layout string) (string, error) {
// 	location, err := time.LoadLocation(timezone)
// 	if err != nil {
// 		return "", err
// 	}
// 	t, errParse := time.ParseInLocation(layout, date, location)
// 	if errParse != nil {
// 		return "", errParse
// 	}
// 	fmt.Println(t.Zone())
// 	t.In(d.Timezone)
// 	return t.Format(d.Format), nil
// }
