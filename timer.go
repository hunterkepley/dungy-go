package main

import (
	"fmt"
	"time"
)

// TimerHandler holds/handles all timers in the game
type TimerHandler struct {
	timers    []Timer
	numTimers int
}

// Timer is a timer that is X seconds long
type Timer struct {
	timer *time.Timer
	id    int
	fired bool
}

// Creates a new timer for X seconds
func (t *TimerHandler) createTimer(seconds int) {
	t.timers = append(t.timers, Timer{time.NewTimer(time.Duration(seconds) * time.Second), t.numTimers, false})
	t.numTimers++
}

// This runs the timer then removes it after it fires
func (t *TimerHandler) runTimer(id int) {
	go func() {
		<-t.timers[id].timer.C
		t.timers[id].fired = true
		fmt.Println("timer with ID ", id, " fired")
		t.timers = t.remove(id)
		t.numTimers--
	}()
}

// Removes a timer from the timers slice
func (t *TimerHandler) remove(id int) []Timer {
	return append(t.timers[:id], t.timers[id+1:]...)
}
