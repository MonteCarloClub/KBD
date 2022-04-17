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

package mclock

import (
	"testing"
	"time"
)

var _ Clock = System{}
var _ Clock = new(Simulated)

func TestSimulatedAfter(t *testing.T) {
	const timeout = 30 * time.Minute
	const adv = time.Minute

	var (
		c   Simulated
		end = c.Now().Add(timeout)
		ch  = c.After(timeout)
	)
	for c.Now() < end.Add(-adv) {
		c.Run(adv)
		select {
		case <-ch:
			t.Fatal("Timer fired early")
		default:
		}
	}

	c.Run(adv)
	select {
	case stamp := <-ch:
		want := time.Time{}.Add(timeout)
		if !stamp.Equal(want) {
			t.Errorf("Wrong time sent on timer channel: got %v, want %v", stamp, want)
		}
	default:
		t.Fatal("Timer didn't fire")
	}
}

func TestSimulatedAfterFunc(t *testing.T) {
	var c Simulated

	called1 := false
	timer1 := c.AfterFunc(100*time.Millisecond, func() { called1 = true })
	if c.ActiveTimers() != 1 {
		t.Fatalf("%d active timers, want one", c.ActiveTimers())
	}
	if fired := timer1.Stop(); !fired {
		t.Fatal("Stop returned false even though timer didn't fire")
	}
	if c.ActiveTimers() != 0 {
		t.Fatalf("%d active timers, want zero", c.ActiveTimers())
	}
	if called1 {
		t.Fatal("timer 1 called")
	}
	if fired := timer1.Stop(); fired {
		t.Fatal("Stop returned true after timer was already stopped")
	}

	called2 := false
	timer2 := c.AfterFunc(100*time.Millisecond, func() { called2 = true })
	c.Run(50 * time.Millisecond)
	if called2 {
		t.Fatal("timer 2 called")
	}
	c.Run(51 * time.Millisecond)
	if !called2 {
		t.Fatal("timer 2 not called")
	}
	if fired := timer2.Stop(); fired {
		t.Fatal("Stop returned true after timer has fired")
	}
}

func TestSimulatedSleep(t *testing.T) {
	var (
		c       Simulated
		timeout = 1 * time.Hour
		done    = make(chan AbsTime)
	)
	go func() {
		c.Sleep(timeout)
		done <- c.Now()
	}()

	c.WaitForTimers(1)
	c.Run(2 * timeout)
	select {
	case stamp := <-done:
		want := AbsTime(2 * timeout)
		if stamp != want {
			t.Errorf("Wrong time after sleep: got %v, want %v", stamp, want)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Sleep didn't return in time")
	}
}
