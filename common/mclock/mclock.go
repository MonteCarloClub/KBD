/*
Copyright (c) 2022 Zhu Zunxiong <liuzunxiong@qq.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

// Package mclock is a wrapper for a monotonic clock source
package mclock

import (
	"time"

	"github.com/aristanetworks/goarista/monotime"
)

// AbsTime represents absolute monotonic time.
type AbsTime time.Duration

// Now returns the current absolute monotonic time.
func Now() AbsTime {
	return AbsTime(monotime.Now())
}

// Add returns t + d.
func (t AbsTime) Add(d time.Duration) AbsTime {
	return t + AbsTime(d)
}

// The Clock interface makes it possible to replace the monotonic system clock with
// a simulated clock.
type Clock interface {
	Now() AbsTime
	Sleep(time.Duration)
	After(time.Duration) <-chan time.Time
	AfterFunc(d time.Duration, f func()) Timer
}

// Timer represents a cancellable event returned by AfterFunc
type Timer interface {
	Stop() bool
}

// System implements Clock using the system clock.
type System struct{}

// Now returns the current monotonic time.
func (System) Now() AbsTime {
	return AbsTime(monotime.Now())
}

// Sleep blocks for the given duration.
func (System) Sleep(d time.Duration) {
	time.Sleep(d)
}

// After returns a channel which receives the current time after d has elapsed.
func (System) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

// AfterFunc runs f on a new goroutine after the duration has elapsed.
func (System) AfterFunc(d time.Duration, f func()) Timer {
	return time.AfterFunc(d, f)
}
