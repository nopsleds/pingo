package timeserie

import "time"

type Ticker func(t time.Time, pos int) time.Time

func TickBySeconds(secondsPerTick int) Ticker {
	return func(t time.Time, ticksDelta int) time.Time {
		t = t.Add(time.Duration(ticksDelta*secondsPerTick) * time.Second)
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), int(t.Second()/secondsPerTick)*secondsPerTick, 0, time.Local)
	}
}
func TickByMinutes(minutesPerTick int) Ticker {
	return func(t time.Time, ticksDelta int) time.Time {
		t = t.Add(time.Duration(ticksDelta*minutesPerTick) * time.Minute)
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), int(t.Minute()/minutesPerTick)*minutesPerTick, 0, 0, time.Local)
	}
}
func TickByHours(hoursPerTick int) Ticker {
	return func(t time.Time, ticksDelta int) time.Time {
		t = t.Add(time.Duration(ticksDelta*hoursPerTick) * time.Hour)
		return time.Date(t.Year(), t.Month(), t.Day(), int(t.Hour()/hoursPerTick)*hoursPerTick, 0, 0, 0, time.Local)
	}
}
func TickByDay(t time.Time, ticksDelta int) time.Time {
	t = t.Add(time.Duration(ticksDelta) * time.Hour * 24)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}
