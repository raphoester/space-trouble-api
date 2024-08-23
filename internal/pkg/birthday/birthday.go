package birthday

import "time"

const format = `02/01`

func Parse(date string) (Date, error) {
	t, err := time.Parse(format, date)
	if err != nil {
		return "", err
	}
	return newDay(t.Day(), int(t.Month())), nil
}

func newDay(day, month int) Date {
	date := time.Date(0, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return Date(date.Format(format))
}

type Date string

func (d Date) String() string {
	return string(d)
}
