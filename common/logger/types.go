package logger

import (
	"math/big"
	"time"
)

type utctime8601 struct{}

func (utctime8601) MarshalJSON() ([]byte, error) {
	timestr := time.Now().UTC().Format(time.RFC3339Nano)
	// Bounds check
	if len(timestr) > 26 {
		timestr = timestr[:26]
	}
	return []byte(`"` + timestr + `Z"`), nil
}

type JsonLog interface {
	EventName() string
}

type LogEvent struct {
	// Guid string      `json:"guid"`
	Ts utctime8601 `json:"ts"`
	// Level string      `json:"level"`
}

type LogStarting struct {
	ClientString    string `json:"client_impl"`
	ProtocolVersion int    `json:"eth_version"`
	LogEvent
}

func (l *LogStarting) EventName() string {
	return "starting"
}

type P2PConnected struct {
	RemoteId            string `json:"remote_id"`
	RemoteAddress       string `json:"remote_addr"`
	RemoteVersionString string `json:"remote_version_string"`
	NumConnections      int    `json:"num_connections"`
	LogEvent
}

func (l *P2PConnected) EventName() string {
	return "p2p.connected"
}

type P2PDisconnected struct {
	NumConnections int    `json:"num_connections"`
	RemoteId       string `json:"remote_id"`
	LogEvent
}

func (l *P2PDisconnected) EventName() string {
	return "p2p.disconnected"
}

type EthMinerNewBlock struct {
	BlockHash     string   `json:"block_hash"`
	BlockNumber   *big.Int `json:"block_number"`
	ChainHeadHash string   `json:"chain_head_hash"`
	BlockPrevHash string   `json:"block_prev_hash"`
	LogEvent
}

func (l *EthMinerNewBlock) EventName() string {
	return "eth.miner.new_block"
}

type EthChainReceivedNewBlock struct {
	BlockHash     string   `json:"block_hash"`
	BlockNumber   *big.Int `json:"block_number"`
	ChainHeadHash string   `json:"chain_head_hash"`
	BlockPrevHash string   `json:"block_prev_hash"`
	RemoteId      string   `json:"remote_id"`
	LogEvent
}

func (l *EthChainReceivedNewBlock) EventName() string {
	return "eth.chain.received.new_block"
}

type EthChainNewHead struct {
	BlockHash     string   `json:"block_hash"`
	BlockNumber   *big.Int `json:"block_number"`
	ChainHeadHash string   `json:"chain_head_hash"`
	BlockPrevHash string   `json:"block_prev_hash"`
	LogEvent
}

func (l *EthChainNewHead) EventName() string {
	return "eth.chain.new_head"
}

type EthTxReceived struct {
	TxHash   string `json:"tx_hash"`
	RemoteId string `json:"remote_id"`
	LogEvent
}

func (l *EthTxReceived) EventName() string {
	return "eth.tx.received"
}
