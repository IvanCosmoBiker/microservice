package datetime

import (
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

func (d *DateTime) Parse(date string) (string, error) {
	var timeLocal time.Time
	var err error
	timeLocal, err = time.ParseInLocation(d.Format, date, d.Timezone)
	if err != nil {
		return "", err
	}
	timeLocal.Format(d.Format)
	return timeLocal.Format(d.Format), nil
}
