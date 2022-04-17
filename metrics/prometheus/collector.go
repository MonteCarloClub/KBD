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

package prometheus

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/MonteCarloClub/KBD/metrics"
)

var (
	typeGaugeTpl           = "# TYPE %s gauge\n"
	typeCounterTpl         = "# TYPE %s counter\n"
	typeSummaryTpl         = "# TYPE %s summary\n"
	keyValueTpl            = "%s %v\n\n"
	keyQuantileTagValueTpl = "%s {quantile=\"%s\"} %v\n\n"
)

// collector is a collection of byte buffers that aggregate Prometheus reports
// for different metric types.
type collector struct {
	buff *bytes.Buffer
}

// newCollector createa a new Prometheus metric aggregator.
func newCollector() *collector {
	return &collector{
		buff: &bytes.Buffer{},
	}
}

func (c *collector) addCounter(name string, m metrics.Counter) {
	c.writeGaugeCounter(name, m.Count())
}

func (c *collector) addGauge(name string, m metrics.Gauge) {
	c.writeGaugeCounter(name, m.Value())
}

func (c *collector) addGaugeFloat64(name string, m metrics.GaugeFloat64) {
	c.writeGaugeCounter(name, m.Value())
}

func (c *collector) addHistogram(name string, m metrics.Histogram) {
	pv := []float64{0.5, 0.75, 0.95, 0.99, 0.999, 0.9999}
	ps := m.Percentiles(pv)
	c.writeSummaryCounter(name, m.Count())
	for i := range pv {
		c.writeSummaryPercentile(name, strconv.FormatFloat(pv[i], 'f', -1, 64), ps[i])
	}
}

func (c *collector) addMeter(name string, m metrics.Meter) {
	c.writeGaugeCounter(name, m.Count())
}

func (c *collector) addTimer(name string, m metrics.Timer) {
	pv := []float64{0.5, 0.75, 0.95, 0.99, 0.999, 0.9999}
	ps := m.Percentiles(pv)
	c.writeSummaryCounter(name, m.Count())
	for i := range pv {
		c.writeSummaryPercentile(name, strconv.FormatFloat(pv[i], 'f', -1, 64), ps[i])
	}
}

func (c *collector) addResettingTimer(name string, m metrics.ResettingTimer) {
	if len(m.Values()) <= 0 {
		return
	}
	ps := m.Percentiles([]float64{50, 95, 99})
	val := m.Values()
	c.writeSummaryCounter(name, len(val))
	c.writeSummaryPercentile(name, "0.50", ps[0])
	c.writeSummaryPercentile(name, "0.95", ps[1])
	c.writeSummaryPercentile(name, "0.99", ps[2])
}

func (c *collector) writeGaugeCounter(name string, value interface{}) {
	name = mutateKey(name)
	c.buff.WriteString(fmt.Sprintf(typeGaugeTpl, name))
	c.buff.WriteString(fmt.Sprintf(keyValueTpl, name, value))
}

func (c *collector) writeSummaryCounter(name string, value interface{}) {
	name = mutateKey(name + "_count")
	c.buff.WriteString(fmt.Sprintf(typeCounterTpl, name))
	c.buff.WriteString(fmt.Sprintf(keyValueTpl, name, value))
}

func (c *collector) writeSummaryPercentile(name, p string, value interface{}) {
	name = mutateKey(name)
	c.buff.WriteString(fmt.Sprintf(typeSummaryTpl, name))
	c.buff.WriteString(fmt.Sprintf(keyQuantileTagValueTpl, name, p, value))
}

func mutateKey(key string) string {
	return strings.Replace(key, "/", "_", -1)
}
