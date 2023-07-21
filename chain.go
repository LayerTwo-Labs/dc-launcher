package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/biter777/processex"
)

type ChainProvider struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	ImageURL        string `json:"imageUrl"`
	BinName         string `json:"binName"`
	DefaultDir      string `json:"defaultDir"`
	DefaultConfName string `json:"defaultConfName"`
	DefaultPort     int    `json:"defaultPort"`
	DefaultSlot     int    `json:"defaultSlot,omitempty"`
}

type ChainData struct {
	ID           string `json:"id"`
	IsDrivechain bool   `json:"isdrivechain,omitempty"`
	BinDir       string `json:"bindir,omitempty"`
	BinName      string `json:"binname,omitempty"`
	ConfDir      string `json:"confdir,omitempty"`
	ConfName     string `json:"confname,omitempty"`
	Port         int    `json:"rpcport"`
	RPCUser      string `json:"rpcuser"`
	RPCPass      string `json:"rpcpassword"`
	Slot         int    `json:"slot,omitempty"`       // Only apply to sidechains
	RefreshBMM   bool   `json:"refreshbmm,omitempty"` // Only apply to sidechains
	BMMFee       bool   `json:"bmmfee,omitempty"`     // Only apply to sidechains
}

type ChainState struct {
	ID               string  `json:"id"`
	State            State   `json:"state"`
	AvailableBalance float64 `json:"availablebalance"`
	PendingBalance   float64 `json:"pendingbalance"`
	Height           int     `json:"height,omitempty"`
	Slot             int     `json:"slot,omitempty"` // Only apply to sidechains
}

type State uint

const (
	Unknown State = iota
	Waiting
	Running
)

var (
	drivechainStateUpdate     *time.Ticker
	quitDrivechainStateUpdate chan struct{}
)

func getChainProcess(name string) (*os.Process, error) {
	process, _, err := processex.FindByName(name)
	if err == processex.ErrNotFound {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	if len(process) > 0 {
		return process[0], nil
	}
	return nil, fmt.Errorf("something went wrong finding process")
}

func LaunchChain(cd *ChainData, cs *ChainState) {
	p, err := getChainProcess(cd.BinName)
	if p != nil && err == nil {
		println(cd.BinName + " already running...")
		// TODO: Try kill?
		// Seems that running form launcher, when you stop chain with rpc stop method, the process stays
		// alive until you close the launcher.  If you close the launcher, the process dies.
		return
	}

	args := []string{"-conf=" + cd.ConfDir + string(os.PathSeparator) + cd.ConfName}
	cmd := exec.Command(cd.BinDir+string(os.PathSeparator)+cd.BinName, args...)

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	cs.State = Waiting
	println(cd.BinName + " Started...")
}

func StopChain(cd *ChainData, cs *ChainState) error {
	req, err := MakeRpcRequest(cd, "stop", []interface{}{})
	if err != nil {
		return err
	} else {
		defer req.Body.Close()
		return nil
	}
}

func StartDrivechainStateUpdate(as *AppState, mui *MainUI) {
	drivechainStateUpdate = time.NewTicker(1 * time.Second)
	quitDrivechainStateUpdate = make(chan struct{})
	go func() {
		for {
			select {
			case <-drivechainStateUpdate.C:
				updateUI := false
				if GetBlockHeight(&as.dcd, &as.dcs) && !updateUI {
					updateUI = true
				}
				if updateUI {
					mui.Refresh()
				}
			case <-quitDrivechainStateUpdate:
				drivechainStateUpdate.Stop()
				return
			}
		}
	}()
}

func GetBlockHeight(cd *ChainData, cs *ChainState) bool {
	currentHeight := cs.Height
	currnetState := cs.State
	bcr, err := MakeRpcRequest(cd, "getblockcount", []interface{}{})
	if err != nil {
		println(err.Error())
		cs.State = Unknown
		if currnetState != cs.State {
			return true
		}
	} else {
		defer bcr.Body.Close()
		if bcr.StatusCode == 200 {
			var res RPCGetBlockCountResponse
			err := json.NewDecoder(bcr.Body).Decode(&res)
			if err == nil {
				cs.Height = res.Result
				cs.State = Running
				changed := false
				if currentHeight != cs.Height {
					changed = true
				}
				if currnetState != cs.State {
					changed = true
				}
				return changed
			}
		}
	}
	return false
}

func GetBalance(cd *ChainData, cs *ChainState) bool {
	currentBalance := cs.AvailableBalance
	bcr, err := MakeRpcRequest(cd, "getbalance", []interface{}{})
	if err != nil {
		println(err.Error())
	} else {
		defer bcr.Body.Close()
		if bcr.StatusCode == 200 {
			var res RPCGetBalanceResponse
			err := json.NewDecoder(bcr.Body).Decode(&res)
			if err == nil {
				cs.AvailableBalance = res.Result
				if currentBalance != cs.AvailableBalance {
					return true
				}
			}
		}
	}
	return false
}
