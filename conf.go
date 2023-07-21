package main

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	defaultChainProvidersConfName = "chains.json"
)

//go:embed binaries/drivechain-qt-linux
var drivechainLinux []byte

//go:embed chain.conf
var chainConfBytes []byte

//go:embed chains.json
var chainsBytes []byte

func ConfInit(as *AppState) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// Setup dc launcher directory and write if not found
	defaultLauncherDir := homeDir + string(os.PathSeparator) + ".dclauncher"
	if _, err := os.Stat(defaultLauncherDir); os.IsNotExist(err) {
		println("Creating " + defaultLauncherDir)
		err = os.Mkdir(defaultLauncherDir, 0o755)
		if err != nil {
			println(err.Error())
			return err
		}
	}

	// Look for chains.json write if not found
	// TODO: This should be pulled from a remote source
	defaultChainProvidersConf := defaultLauncherDir + string(os.PathSeparator) + defaultChainProvidersConfName
	if _, err := os.Stat(defaultChainProvidersConf); os.IsNotExist(err) {
		println("Creating " + defaultChainProvidersConf)
		err = os.WriteFile(defaultChainProvidersConf, chainsBytes, 0o755)
		if err != nil {
			println(err.Error())
			return err
		}
	}

	// Now read in the chains.json file
	chains, err := ioutil.ReadFile(defaultChainProvidersConf)
	if err != nil {
		println(err.Error())
		return err
	}

	var chainProviders map[string]ChainProvider
	if err := json.Unmarshal(chains, &chainProviders); err != nil {
		println(err.Error())
		return err
	}
	as.cp = chainProviders

	for k, chainProvider := range chainProviders {

		confDir := homeDir + string(os.PathSeparator) + chainProvider.DefaultDir
		if _, err := os.Stat(confDir); os.IsNotExist(err) {
			println("Creating " + confDir)
			err = os.Mkdir(confDir, 0o755)
			if err != nil {
				println(err.Error())
				return err
			}
		}

		conf := confDir + string(os.PathSeparator) + chainProvider.DefaultConfName
		if _, err := os.Stat(conf); os.IsNotExist(err) {
			confBytes := chainConfBytes
			confBytes = append(confBytes, "\ndatadir="+confDir...)
			confBytes = append(confBytes, fmt.Sprintf("\nrpcport=%v", chainProvider.DefaultPort)...)
			if k != "drivechain" {
				confBytes = append(confBytes, fmt.Sprintf("\nslot=%v", chainProvider.DefaultSlot)...)
			}
			err := os.WriteFile(conf, confBytes, 0o755)
			println("Writing " + conf)
			if err != nil {
				println(err.Error())
				return err
			}
		}

		// Read back in the confs
		chainData := ChainData{}
		chainData.ID = k
		chainData.IsDrivechain = k == "drivechain"
		chainData.BinName = chainProvider.BinName
		chainData.ConfName = chainProvider.DefaultConfName
		chainData.BinDir = confDir
		chainData.ConfDir = confDir
		chainData.Port = chainProvider.DefaultPort
		if k != "drivechain" {
			chainData.Slot = chainProvider.DefaultSlot
		}
		err = loadConf(&chainData)
		if err != nil {
			println(err.Error())
			return err
		}

		// fmt.Print("%+v\n", chainData)

		if k == "drivechain" {
			as.dcd = chainData
			as.dcs = ChainState{ID: k}
		} else {
			as.scd[k] = chainData
			as.scs[k] = ChainState{ID: k}
		}

		// Write chain binary if not found

	}

	// No

	return nil
}

func loadConf(chainData *ChainData) error {
	readFile, err := os.Open(chainData.ConfDir + string(os.PathSeparator) + chainData.ConfName)
	if err != nil {
		println(err.Error())
		return err
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	readFile.Close()

	confMap := make(map[string]interface{})

	for _, line := range fileLines {
		a := strings.Split(line, "=")
		if len(a) == 2 {
			k := strings.TrimSpace(a[0])
			v := strings.TrimSpace(a[1])
			println(k + " = " + v)
			if k != "" {
				iv, err := (strconv.ParseInt(v, 0, 64))
				if err != nil {
					confMap[k] = v
				} else {
					confMap[k] = int(iv)
				}
			}
		}
	}

	jsonData, _ := json.Marshal(confMap)
	err = json.Unmarshal(jsonData, &chainData)
	if err != nil {
		println(err.Error())
		return err
	}
	return nil
}
