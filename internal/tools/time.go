package tools

import "time"

func Now() time.Time {
	t := time.Now()
	t2 := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
	return t2
}

func EndOfTheDay(t time.Time) time.Time {
	t2 := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
	return t2
}
