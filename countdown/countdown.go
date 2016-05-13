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
	if time.Now().Before(to) {
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
	return int(math.Floor(c.To.Sub(time.Now()).Seconds() / 3600 / 24))
}

// TotalHours return the countdown in hours
func (c *Countdown) TotalHours() int {
	return int(math.Floor(c.To.Sub(time.Now()).Seconds() / 3600))
}

// TotalMinutes returns the countdown in minutes
func (c *Countdown) TotalMinutes() int {
	return int(math.Floor(c.To.Sub(time.Now()).Seconds() / 3600 / 60))
}

// TotalSeconds returns the conutdown in seconds
func (c *Countdown) TotalSeconds() int {
	return int(c.To.Sub(time.Now()).Seconds())
}

// RemainingHours returns the remaining hours
func (c *Countdown) RemainingHours() int {
	return int(math.Floor(math.Mod(c.To.Sub(time.Now()).Seconds()/3600, 24)))
}

// RemainingMinutes returns the remaining minutes
func (c *Countdown) RemainingMinutes() int {
	return int(math.Floor(math.Mod(c.To.Sub(time.Now()).Seconds()/60, 60)))
}

// RemainingSeconds returns the remaining seconds
func (c *Countdown) RemainingSeconds() int {
	return int(math.Floor(math.Mod(c.To.Sub(time.Now()).Seconds(), 60)))
}
