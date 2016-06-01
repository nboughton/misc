package countdown

import (
	"fmt"
	"math"
	"time"
)

// Countdown contains the target date and a ticker to compare it to
type Countdown struct {
	To time.Time
}

// New creates a new countdown object
func New(to time.Time) (*Countdown, error) {
	if !happened(to) {
		return &Countdown{To: to}, nil
	}
	return &Countdown{}, fmt.Errorf("%v is already in the past", to.Format(time.RFC822))
}

// DHMS returns "<days>d <hours>h <mins>m <sec>s"
func (c *Countdown) DHMS() string {
	return fmt.Sprintf("%vd %vh %vm %vs", c.TotalDays(), c.RemainingHours(), c.RemainingMinutes(), c.RemainingSeconds())
}

// HMS returns "<hours>:<mins>:<sec>"
func (c *Countdown) HMS() string {
	return fmt.Sprintf("%v:%v:%v", c.TotalHours(), c.RemainingMinutes(), c.RemainingSeconds())
}

// MS returns "<mins>:<sec>"
func (c *Countdown) MS() string {
	return fmt.Sprintf("%v:%v", c.TotalMinutes(), c.RemainingSeconds())
}

// TotalDays returns the countdown in days
func (c *Countdown) TotalDays() int {
	if !happened(c.To) {
		return int(math.Floor(c.To.Sub(time.Now()).Seconds() / 3600 / 24))
	}
	return 0
}

// TotalHours return the countdown in hours
func (c *Countdown) TotalHours() int {
	if !happened(c.To) {
		return int(math.Floor(c.To.Sub(time.Now()).Seconds() / 3600))
	}
	return 0
}

// TotalMinutes returns the countdown in minutes
func (c *Countdown) TotalMinutes() int {
	if !happened(c.To) {
		return int(math.Floor(c.To.Sub(time.Now()).Seconds() / 3600 / 60))
	}
	return 0
}

// TotalSeconds returns the conutdown in seconds
func (c *Countdown) TotalSeconds() int {
	if !happened(c.To) {
		return int(c.To.Sub(time.Now()).Seconds())
	}
	return 0
}

// RemainingHours returns the remaining hours
func (c *Countdown) RemainingHours() int {
	if !happened(c.To) {
		return int(math.Floor(math.Mod(c.To.Sub(time.Now()).Seconds()/3600, 24)))
	}
	return 0
}

// RemainingMinutes returns the remaining minutes
func (c *Countdown) RemainingMinutes() int {
	if !happened(c.To) {
		return int(math.Floor(math.Mod(c.To.Sub(time.Now()).Seconds()/60, 60)))
	}
	return 0
}

// RemainingSeconds returns the remaining seconds
func (c *Countdown) RemainingSeconds() int {
	if !happened(c.To) {
		return int(math.Floor(math.Mod(c.To.Sub(time.Now()).Seconds(), 60)))
	}
	return 0
}

func happened(to time.Time) bool {
	if time.Now().Before(to) {
		return false
	}
	return true
}
