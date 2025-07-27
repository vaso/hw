package bucket

import "time"

type Timer interface {
	Now() time.Time
}

type timer struct{}

func NewTimer() Timer {
	return &timer{}
}

func (t timer) Now() time.Time {
	return time.Now()
}
