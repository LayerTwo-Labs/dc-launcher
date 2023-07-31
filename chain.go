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
	Order           int    `json:"order"`
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
	Automine         bool    `json:"automine,omitempty"`
	ChainStateUpdate ChainStateUpdate
}

type State uint

const (
	Unknown State = iota
	Waiting
	Running
)

type ChainStateUpdate struct {
	ID    string `json:"id"`
	timer *time.Ticker
	quit  chan struct{}
}

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

func LaunchChain(cd *ChainData, cs *ChainState, mui *MainUI) {
	if cs.ChainStateUpdate.timer != nil && cs.ChainStateUpdate.quit != nil {
		// TODO: Maybe restart?
	} else {
		csu := ChainStateUpdate{ID: cd.ID}
		cs.ChainStateUpdate = csu
		StartChainStateUpdate(cd, cs, mui)
	}

	//p, err := getChainProcess(cd.BinName)
	//if p != nil && err == nil {
	//println(cd.BinName + " already running...")
	// TODO: Try kill?
	// Seems that running form launcher, when you stop chain with rpc stop method, the process stays
	// alive until you close the launcher.  If you close the launcher, the process dies.
	//return
	//}

	args := []string{"-conf=" + cd.ConfDir + string(os.PathSeparator) + cd.ConfName}
	cmd := exec.Command(cd.BinDir+string(os.PathSeparator)+cd.BinName, args...)

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	cs.State = Waiting
	println(cd.BinName + " Started...")
}

func StopChain(cd *ChainData, cs *ChainState, as *AppState) error {
	if cd.ID == "drivechain" {
		// stop all
		for k := range as.scd {
			scd := as.scd[k]
			scs := as.scs[k]
			StopChain(&scd, &scs, as)
		}
	}

	req, err := MakeRpcRequest(cd, "stop", []interface{}{})
	if err == nil {
		defer req.Body.Close()
	}
	p, err := getChainProcess(cd.BinName)
	if p != nil && err == nil {
		return p.Kill()
	}
	return err
}

func StartChainStateUpdate(cd *ChainData, cs *ChainState, mui *MainUI) {
	cs.ChainStateUpdate.timer = time.NewTicker(1 * time.Second)
	cs.ChainStateUpdate.quit = make(chan struct{})
	go func() {
		for {
			select {
			case <-cs.ChainStateUpdate.timer.C:
				if cd.ID == "drivechain" && cs.Automine {
					DrivechainMine(mui.as, mui)
				}
				updateUI := false
				if GetBlockHeight(cd, cs) && !updateUI {
					mui.as.scs[cd.ID] = *cs
					updateUI = true
				}
				if updateUI {
					mui.Refresh()
				}
			case <-cs.ChainStateUpdate.quit:
				cs.ChainStateUpdate.timer.Stop()
				return
			}
		}
	}()
}

func DrivechainMine(as *AppState, mui *MainUI) {
	_, err := MakeRpcRequest(&as.dcd, "generate", []interface{}{1})
	if err != nil {
		println(err.Error())
	}
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

func NeedsActivation(cd *ChainData, as *AppState) bool {
	ls, err := MakeRpcRequest(&as.dcd, "listactivesidechains", []interface{}{})
	if err != nil {
		println(err.Error())
	} else if ls.StatusCode == 200 {
		defer ls.Body.Close()
		var res RPCCreateSidechainProposalResponse
		err := json.NewDecoder(ls.Body).Decode(&res)
		if err == nil {
			for _, sc := range res.Result {
				if sc.Title == cd.ID {
					return false
				}
			}
			return true
		}
	}
	return true
}

func CreateSidechainProposal(as *AppState, cd *ChainData, cs *ChainState) bool {
	println("Creating sidechain proposal...")
	pr, err := MakeRpcRequest(&as.dcd, "createsidechainproposal", []interface{}{cd.Slot, cd.ID})
	if err != nil {
		println(err.Error())
	} else if pr.StatusCode == 200 {
		_, err := MakeRpcRequest(&as.dcd, "generate", []interface{}{201})
		if err != nil {
			println(err.Error())
		}
	}
	return true
}
