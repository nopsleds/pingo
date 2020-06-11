package timeserie

import (
	"errors"
	"time"
)

type Clock interface {
	Now() time.Time
}

type defaultClock struct{}

func (c *defaultClock) Now() time.Time {
	return time.Now()
}

var (
	DefaultClock = &defaultClock{}
)

type Interval int

const (
	IntervalMonth Interval = iota
	IntervalWeek
	IntervalDay
	IntervalHour
	Interval10m
	IntervalMinute
)

type ValueReducer func(current int, newValue int) int

type Timeserie struct {
	clock        Clock
	ticker       Ticker
	defaultValue int
	interval     Interval
	lastTick     time.Time
	data         []int
	cursor       int
	reduce       ValueReducer
}

type Options struct {
	Clock        Clock
	Ticker       Ticker
	Length       int
	DefaultValue int
	ValueReducer ValueReducer
}

func ReducerLatestValue(current int, newValue int) int {
	return newValue
}

func New(opt Options) (*Timeserie, error) {
	if opt.Clock == nil {
		opt.Clock = DefaultClock
	}
	if opt.Ticker == nil {
		return nil, errors.New("missing ticker")
	}
	if opt.ValueReducer == nil {
		opt.ValueReducer = ReducerLatestValue
	}
	if opt.Length <= 0 {
		return nil, errors.New("can't make timeseries with zero/negative length")
	}

	data := make([]int, opt.Length)
	for i := 0; i < len(data); i++ {
		data[i] = opt.DefaultValue
	}
	return &Timeserie{
		clock:        opt.Clock,
		ticker:       opt.Ticker,
		defaultValue: opt.DefaultValue,
		lastTick:     opt.Ticker(opt.Clock.Now(), 0),
		data:         data,
		cursor:       0,
		reduce:       opt.ValueReducer,
	}, nil
}

func (t *Timeserie) Length() int {
	return len(t.data)
}
func (ts *Timeserie) Insert(value int) {
	now := ts.clock.Now()
	ts.fillUntil(now)
	ts.doInsert(value, now)
}
func (ts *Timeserie) fillUntil(now time.Time) {
	tickMax := ts.ticker(now, 0)
	tickMin := ts.ticker(now, -1*ts.Length())
	if ts.lastTick.After(tickMin) {
		tickMin = ts.ticker(ts.lastTick, 1)
	}
	for t := tickMin; t.Before(tickMax); t = ts.ticker(t, 1) {
		ts.doInsert(ts.defaultValue, t)
	}
}
func (ts *Timeserie) doInsert(value int, now time.Time) {
	tick := ts.ticker(now, 0)
	if tick.After(ts.lastTick) {
		ts.incrCursor()
		ts.lastTick = tick
	}
	ts.data[ts.cursor] = ts.reduce(ts.data[ts.cursor], value)
}
func (ts *Timeserie) incrCursor() {
	ts.cursor = (ts.cursor + 1) % len(ts.data)
}

type Entry struct {
	Time  time.Time
	Value int
}

func (ts *Timeserie) Get() []Entry {
	ts.fillUntil(ts.clock.Now())
	values := append(ts.data[ts.cursor+1:], ts.data[:ts.cursor+1]...)

	nbTicks := len(values)
	entries := make([]Entry, nbTicks)
	for i, value := range values {
		entries[i] = Entry{
			Time:  ts.ticker(ts.lastTick, i-nbTicks),
			Value: value,
		}
	}
	return entries
}
