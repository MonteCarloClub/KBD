package chain_manager

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/MonteCarloClub/KBD/common"
	"github.com/MonteCarloClub/KBD/rlp"
	"github.com/MonteCarloClub/KBD/types"
	"github.com/cloudwego/kitex/pkg/klog"
)

// DisabledBadBlockReporting can be set to prevent blocks being reported.
var DisableBadBlockReporting = true

// ReportBlock reports the block to the block reporting tool found at
// badblocks.ethdev.com
func ReportBlock(block *types.Block, err error) {
	if DisableBadBlockReporting {
		return
	}

	const url = "https://badblocks.ethdev.com"

	blockRlp, _ := rlp.EncodeToBytes(block)
	data := map[string]interface{}{
		"block":     common.Bytes2Hex(blockRlp),
		"errortype": err.Error(),
		"hints": map[string]interface{}{
			"receipts": "NYI",
			"vmtrace":  "NYI",
		},
	}
	jsonStr, _ := json.Marshal(map[string]interface{}{"method": "eth_badBlock", "params": []interface{}{data}, "id": "1", "jsonrpc": "2.0"})

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		klog.Error("POST err:%v", err)
		return
	}
	defer resp.Body.Close()
	klog.Debug("response Status:%v", resp.Status)
	klog.Debug("response Headers:%v", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	klog.Debug("response Body:%v", string(body))
}
