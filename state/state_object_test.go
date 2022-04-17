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

package state

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/MonteCarloClub/KBD/common"
)

func BenchmarkCutOriginal(b *testing.B) {
	value := common.HexToHash("0x01")
	for i := 0; i < b.N; i++ {
		bytes.TrimLeft(value[:], "\x00")
	}
}

func BenchmarkCutsetterFn(b *testing.B) {
	value := common.HexToHash("0x01")
	cutSetFn := func(r rune) bool {
		return int32(r) == int32(0)
	}
	for i := 0; i < b.N; i++ {
		bytes.TrimLeftFunc(value[:], cutSetFn)
	}
}

func BenchmarkCutCustomTrim(b *testing.B) {
	value := common.HexToHash("0x01")
	for i := 0; i < b.N; i++ {
		common.TrimLeftZeroes(value[:])
	}
}

func xTestFuzzCutter(t *testing.T) {
	rand.Seed(time.Now().Unix())
	for {
		v := make([]byte, 20)
		zeroes := rand.Intn(21)
		rand.Read(v[zeroes:])
		exp := bytes.TrimLeft(v[:], "\x00")
		got := common.TrimLeftZeroes(v)
		if !bytes.Equal(exp, got) {

			fmt.Printf("Input %x\n", v)
			fmt.Printf("Exp %x\n", exp)
			fmt.Printf("Got %x\n", got)
			t.Fatalf("Error")
		}
		//break
	}
}
