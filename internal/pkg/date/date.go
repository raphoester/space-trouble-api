package date

import "time"

const format = "02/01/2006"

type Date time.Time

func (d Date) String() string {
	return time.Time(d).Format(format)
}

func Parse(date string) (Date, error) {
	t, err := time.Parse(format, date)
	if err != nil {
		return Date{}, err
	}
	return Date(t), nil
}
