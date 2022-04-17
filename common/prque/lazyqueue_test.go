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

package prque

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/MonteCarloClub/KBD/common/mclock"
)

const (
	testItems        = 1000
	testPriorityStep = 100
	testSteps        = 1000000
	testStepPeriod   = time.Millisecond
	testQueueRefresh = time.Second
	testAvgRate      = float64(testPriorityStep) / float64(testItems) / float64(testStepPeriod)
)

type lazyItem struct {
	p, maxp int64
	last    mclock.AbsTime
	index   int
}

func testPriority(a interface{}, now mclock.AbsTime) int64 {
	return a.(*lazyItem).p
}

func testMaxPriority(a interface{}, until mclock.AbsTime) int64 {
	i := a.(*lazyItem)
	dt := until - i.last
	i.maxp = i.p + int64(float64(dt)*testAvgRate)
	return i.maxp
}

func testSetIndex(a interface{}, i int) {
	a.(*lazyItem).index = i
}

func TestLazyQueue(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	clock := &mclock.Simulated{}
	q := NewLazyQueue(testSetIndex, testPriority, testMaxPriority, clock, testQueueRefresh)

	var (
		items  [testItems]lazyItem
		maxPri int64
	)

	for i := range items[:] {
		items[i].p = rand.Int63n(testPriorityStep * 10)
		if items[i].p > maxPri {
			maxPri = items[i].p
		}
		items[i].index = -1
		q.Push(&items[i])
	}

	var lock sync.Mutex
	stopCh := make(chan chan struct{})
	go func() {
		for {
			select {
			case <-clock.After(testQueueRefresh):
				lock.Lock()
				q.Refresh()
				lock.Unlock()
			case stop := <-stopCh:
				close(stop)
				return
			}
		}
	}()

	for c := 0; c < testSteps; c++ {
		i := rand.Intn(testItems)
		lock.Lock()
		items[i].p += rand.Int63n(testPriorityStep*2-1) + 1
		if items[i].p > maxPri {
			maxPri = items[i].p
		}
		items[i].last = clock.Now()
		if items[i].p > items[i].maxp {
			q.Update(items[i].index)
		}
		if rand.Intn(100) == 0 {
			p := q.PopItem().(*lazyItem)
			if p.p != maxPri {
				t.Fatalf("incorrect item (best known priority %d, popped %d)", maxPri, p.p)
			}
			q.Push(p)
		}
		lock.Unlock()
		clock.Run(testStepPeriod)
		clock.WaitForTimers(1)
	}

	stop := make(chan struct{})
	stopCh <- stop
	<-stop
}
