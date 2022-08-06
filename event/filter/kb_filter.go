package filter

// TODO make use of the generic filtering system

import (
	"github.com/MonteCarloClub/KBD/manager"
	"sync"

	"github.com/MonteCarloClub/KBD/event"
	"github.com/MonteCarloClub/KBD/state"
)

type FilterManager struct {
	eventMux *event.TypeMux

	filterMu sync.RWMutex
	filterId int
	filters  map[int]*manager.Filter

	quit chan struct{}
}

func NewFilterManager(mux *event.TypeMux) *FilterManager {
	return &FilterManager{
		eventMux: mux,
		filters:  make(map[int]*manager.Filter),
	}
}

func (self *FilterManager) Start() {
	go self.filterLoop()
}

func (self *FilterManager) Stop() {
	close(self.quit)
}

func (self *FilterManager) InstallFilter(filter *manager.Filter) (id int) {
	self.filterMu.Lock()
	defer self.filterMu.Unlock()
	id = self.filterId
	self.filters[id] = filter
	self.filterId++

	return id
}

func (self *FilterManager) UninstallFilter(id int) {
	self.filterMu.Lock()
	defer self.filterMu.Unlock()
	if _, ok := self.filters[id]; ok {
		delete(self.filters, id)
	}
}

// GetFilter retrieves a filter installed using InstallFilter.
// The filter may not be modified.
func (self *FilterManager) GetFilter(id int) *manager.Filter {
	self.filterMu.RLock()
	defer self.filterMu.RUnlock()
	return self.filters[id]
}

func (self *FilterManager) filterLoop() {
	// Subscribe to events
	events := self.eventMux.Subscribe(
		//core.PendingBlockEvent{},
		event.ChainEvent{},
		event.TxPreEvent{},
		state.Logs(nil))

out:
	for {
		select {
		case <-self.quit:
			break out
		case e := <-events.Chan():
			switch e := e.(type) {
			case event.ChainEvent:
				self.filterMu.RLock()
				for _, filter := range self.filters {
					if filter.BlockCallback != nil {
						filter.BlockCallback(e.Block, e.Logs)
					}
				}
				self.filterMu.RUnlock()

			case event.TxPreEvent:
				self.filterMu.RLock()
				for _, filter := range self.filters {
					if filter.TransactionCallback != nil {
						filter.TransactionCallback(e.Tx)
					}
				}
				self.filterMu.RUnlock()

			case state.Logs:
				self.filterMu.RLock()
				for _, filter := range self.filters {
					if filter.LogsCallback != nil {
						msgs := filter.FilterLogs(e)
						if len(msgs) > 0 {
							filter.LogsCallback(msgs)
						}
					}
				}
				self.filterMu.RUnlock()
			}
		}
	}
}
