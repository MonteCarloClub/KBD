package state

import (
	"fmt"
	checker "gopkg.in/check.v1"
	"math/big"
	"sync"
	"time"

	"github.com/MonteCarloClub/KBD/common"
	"github.com/MonteCarloClub/KBD/kdb"
)

var (
	testTimes = 3000
)

type StateSuite struct {
	state *StateDB
}

var _ = checker.Suite(&StateSuite{})

var toAddr = common.BytesToAddress

func (s *StateSuite) TestDump(c *checker.C) {
	// generate a few entries
	obj1 := s.state.GetOrNewStateObject(toAddr([]byte{0x01}))
	obj1.AddBalance(big.NewInt(22))
	obj2 := s.state.GetOrNewStateObject(toAddr([]byte{0x01, 0x02}))
	obj2.SetCode([]byte{3, 3, 3, 3, 3, 3, 3})
	obj3 := s.state.GetOrNewStateObject(toAddr([]byte{0x02}))
	obj3.SetBalance(big.NewInt(44))

	// write some of them to the trie
	s.state.UpdateStateObject(obj1)
	s.state.UpdateStateObject(obj2)

	// check that dump contains the state objects that are in trie
	got := string(s.state.Dump())
	want := `{
    "root": "61bfdc807b57cc7e7ba22dd01ca7dcd93a94d42ca7c0d329f1bfac2eaef95071",
    "accounts": {
        "0000000000000000000000000000000000000001": {
            "balance": "22",
            "nonce": 0,
            "root": "bc2071a4de846f285702447f2589dd163678e0972a8a1b0d28b04ed5c094547f",
            "codeHash": "a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a",
            "storage": {}
        },
        "0000000000000000000000000000000000000102": {
            "balance": "0",
            "nonce": 0,
            "root": "bc2071a4de846f285702447f2589dd163678e0972a8a1b0d28b04ed5c094547f",
            "codeHash": "cfcb3d7d60805bfaf13773e5c87df8a37d509ebbf8c717f01ed867b3864f77a2",
            "storage": {}
        }
    }
}`
	if got != want {
		c.Errorf("dump mismatch:\ngot: %s\nwant: %s\n", got, want)
	}
}

func (s *StateSuite) SetUpTest(c *checker.C) {
	db, _ := kdb.NewMemDatabase()
	s.state = New(common.Hash{}, db)
}

func (s *StateSuite) TestDB(c *checker.C) {
	db, _ := kdb.NewLDBDatabase("testDB")
	start := time.Now().UnixNano()
	state := New(common.Hash{}, db)
	for i := 0; i < testTimes; i++ {
		cal(state)
	}
	useTime := time.Now().UnixNano() - start
	time := float64(useTime) / 1e9
	fmt.Println(time)
	c.Log(time)
	c.Log()
}

func (s *StateSuite) TestDB_goroutine(c *checker.C) {
	db, _ := kdb.NewLDBDatabase("testDB")
	start := time.Now().UnixNano()
	state := New(common.Hash{}, db)
	wg := &sync.WaitGroup{}
	for i := 0; i < testTimes; i++ {
		ti := i
		fmt.Println(ti)
		wg.Add(1)
		go cal_g(wg, state)
	}
	wg.Wait()
	useTime := time.Now().UnixNano() - start
	time := float64(useTime) / 1e9
	c.Log(time)
	c.Log()
	return
}

func (s *StateSuite) TestMemDB(c *checker.C) {
	db, _ := kdb.NewMemDatabase()
	start := time.Now().UnixNano()
	state := New(common.Hash{}, db)
	for i := 0; i < testTimes; i++ {
		cal(state)
	}
	useTime := time.Now().UnixNano() - start
	time := float64(useTime) / 1e9
	c.Log(time)
}

func (s *StateSuite) TestMemDB_goroutine(c *checker.C) {
	db, _ := kdb.NewMemDatabase()
	start := time.Now().UnixNano()
	state := New(common.Hash{}, db)
	wg := &sync.WaitGroup{}
	for i := 0; i < testTimes; i++ {
		wg.Add(1)
		go cal_g(wg, state)
	}
	wg.Wait()
	useTime := time.Now().UnixNano() - start
	time := float64(useTime) / 1e9
	c.Log(time)
}

func (s *StateSuite) TestSnapshot(c *checker.C) {
	stateobjaddr := toAddr([]byte("aa"))
	storageaddr := common.Big("0")
	data1 := common.NewValue(42)
	data2 := common.NewValue(43)

	// get state object
	stateObject := s.state.GetOrNewStateObject(stateobjaddr)
	// set inital state object value
	stateObject.SetStorage(storageaddr, data1)
	// get snapshot of current state
	snapshot := s.state.Copy()

	// get state object. is this strictly necessary?
	stateObject = s.state.GetStateObject(stateobjaddr)
	// set new state object value
	stateObject.SetStorage(storageaddr, data2)
	// restore snapshot
	s.state.Set(snapshot)

	// get state object
	stateObject = s.state.GetStateObject(stateobjaddr)
	// get state storage value
	res := stateObject.GetStorage(storageaddr)

	c.Assert(data1, checker.DeepEquals, res)
}

func cal_g(wg *sync.WaitGroup, state *StateDB) error {
	defer func() {
		wg.Done()
		if err := recover(); err != nil {
			fmt.Println("error")
		}
	}()
	w := common.NewWallet()
	address := common.StringToAddress(w.NewAddress())
	state.CreateAccount(address)
	//value := common.FromHex("0x823140710bf13990e4500136726d8b55")
	value := make([]byte, 16)
	state.SetState(address, common.Hash{}, value)
	state.Update()
	state.Sync()
	value = state.GetState(address, common.Hash{})
	return nil
}
func cal(state *StateDB) {
	w := common.NewWallet()
	address := common.StringToAddress(w.NewAddress())
	state.CreateAccount(address)
	//value := common.FromHex("0x823140710bf13990e4500136726d8b55")
	value := make([]byte, 16)
	state.SetState(address, common.Hash{}, value)
	state.Update()
	state.Sync()
	value = state.GetState(address, common.Hash{})
}
