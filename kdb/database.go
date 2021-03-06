package kdb

import (
	"sync"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"

	"github.com/MonteCarloClub/KBD/common/logger"
	"github.com/MonteCarloClub/KBD/common/logger/glog"
	"github.com/MonteCarloClub/KBD/compression/rle"
)

type LDBDatabase struct {
	fn string

	mu sync.Mutex
	db *leveldb.DB

	queue map[string][]byte

	quit chan struct{}
}

func NewLDBDatabase(file string) (*LDBDatabase, error) {
	// Open the db
	db, err := leveldb.OpenFile(file, nil)
	if err != nil {
		return nil, err
	}
	database := &LDBDatabase{
		fn:   file,
		db:   db,
		quit: make(chan struct{}),
	}
	database.makeQueue()

	go database.update()

	return database, nil
}

func (self *LDBDatabase) makeQueue() {
	self.queue = make(map[string][]byte)
}

func (self *LDBDatabase) Put(key []byte, value []byte) {
	self.mu.Lock()
	defer self.mu.Unlock()

	self.queue[string(key)] = value
	/*
		value = rle.Compress(value)

		err := self.db.Put(key, value, nil)
		if err != nil {
			fmt.Println("Error put", err)
		}
	*/
}

func (self *LDBDatabase) Get(key []byte) ([]byte, error) {
	self.mu.Lock()
	defer self.mu.Unlock()

	// Check queue first
	if dat, ok := self.queue[string(key)]; ok {
		return dat, nil
	}

	dat, err := self.db.Get(key, nil)
	if err != nil {
		return nil, err
	}

	return rle.Decompress(dat)
}

func (self *LDBDatabase) Delete(key []byte) error {
	self.mu.Lock()
	defer self.mu.Unlock()

	// make sure it's not in the queue
	delete(self.queue, string(key))

	return self.db.Delete(key, nil)
}

func (self *LDBDatabase) LastKnownTD() []byte {
	data, _ := self.Get([]byte("LTD"))

	if len(data) == 0 {
		data = []byte{0x0}
	}

	return data
}

func (self *LDBDatabase) NewIterator() iterator.Iterator {
	return self.db.NewIterator(nil, nil)
}

func (self *LDBDatabase) Flush() error {
	self.mu.Lock()
	defer self.mu.Unlock()

	batch := new(leveldb.Batch)

	for key, value := range self.queue {
		batch.Put([]byte(key), rle.Compress(value))
	}
	self.makeQueue() // reset the queue

	return self.db.Write(batch, nil)
}

func (self *LDBDatabase) Close() {
	self.quit <- struct{}{}
	<-self.quit
	glog.V(logger.Info).Infoln("flushed and closed db:", self.fn)
}

func (self *LDBDatabase) update() {
	ticker := time.NewTicker(1 * time.Minute)
done:
	for {
		select {
		case <-ticker.C:
			if err := self.Flush(); err != nil {
				glog.V(logger.Error).Infof("error: flush '%s': %v\n", self.fn, err)
			}
		case <-self.quit:
			break done
		}
	}

	if err := self.Flush(); err != nil {
		glog.V(logger.Error).Infof("error: flush '%s': %v\n", self.fn, err)
	}

	// Close the leveldb database
	self.db.Close()

	self.quit <- struct{}{}
}

func (self *LDBDatabase) LDB() *leveldb.DB {
	return self.db
}
