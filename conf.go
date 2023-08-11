package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const (
	defaultChainProvidersConfName = "chains.json"
)

//go:embed binaries/linux/drivechain-qt-linux
var drivechainLinux []byte

//go:embed binaries/linux/testchain-qt-linux
var testchainLinux []byte

//go:embed binaries/linux/bitassets-qt-linux
var bitassetsLinux []byte

//go:embed binaries/linux/thunder-linux
var thunderLinux []byte

//go:embed binaries/linux/bitcoin-qt-linux
var latestCoreLinux []byte

//go:embed binaries/linux/bitnames.zip
var bitnamesZipLinux []byte

//go:embed chain.conf
var chainConfBytes []byte

//go:embed chains.json
var chainsBytes []byte

func ResetEverything(as *AppState) error {
	// Stop all chains
	// Stoping Drivechain will also stop sidechains
	err := StopChain(&as.dcd, &as.dcs, as)
	if err != nil {
		println(err.Error())
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		println(err.Error())
	}

	err = os.RemoveAll(homeDir + string(os.PathSeparator) + ".dclauncher")
	if err != nil {
		println(err.Error())
	}

	err = os.RemoveAll(homeDir + string(os.PathSeparator) + ".drivechain")
	if err != nil {
		println(err.Error())
	}

	for _, chainData := range as.scd {
		err = os.RemoveAll(chainData.ConfDir)
		if err != nil {
			println(err.Error())
		}
	}

	go func() {
		as.dcs.ChainStateUpdate.quit <- struct{}{}
		for _, chainState := range as.scs {
			chainState.ChainStateUpdate.quit <- struct{}{}
		}
	}()

	return ConfInit(as)
}

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
	chains, err := os.ReadFile(defaultChainProvidersConf)
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
			var confBytes []byte
			if k == "thunder" {
				confBytes = append(confBytes, fmt.Sprintf("\nrpcport=%v", chainProvider.DefaultPort)...)
				confBytes = append(confBytes, fmt.Sprintf("\nslot=%v", chainProvider.DefaultSlot)...)
			} else if k == "latestcore" {
				confBytes = append(confBytes, "chain=regtest"...)
				confBytes = append(confBytes, "\nserver=1"...)
				confBytes = append(confBytes, "\nsplash=0"...)
				confBytes = append(confBytes, fmt.Sprintf("\nslot=%v", chainProvider.DefaultSlot)...)
				confBytes = append(confBytes, "\ndatadir="+confDir...)
				confBytes = append(confBytes, "\n[regtest]"...)
				confBytes = append(confBytes, "\nrpcuser=user"...)
				confBytes = append(confBytes, "\nrpcpassword=password"...)
				confBytes = append(confBytes, fmt.Sprintf("\nrpcport=%v", chainProvider.DefaultPort)...)
			} else {
				confBytes = chainConfBytes
				confBytes = append(confBytes, "\ndatadir="+confDir...)
				confBytes = append(confBytes, fmt.Sprintf("\nrpcport=%v", chainProvider.DefaultPort)...)
				if k != "drivechain" {
					confBytes = append(confBytes, fmt.Sprintf("\nslot=%v", chainProvider.DefaultSlot)...)
				}
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
		if k == "bitnames" {
			chainData.BinDir = confDir + string(os.PathSeparator) + "usr" + string(os.PathSeparator) + "bin"
		} else {
			chainData.BinDir = confDir
		}
		chainData.ConfDir = confDir

		err = loadConf(&chainData)
		if err != nil {
			println(err.Error())
			return err
		}

		chainData.Port = chainProvider.DefaultPort
		if k != "drivechain" {
			chainData.Slot = chainProvider.DefaultSlot
		}

		if k == "drivechain" {
			as.dcd = chainData
			as.dcs = ChainState{ID: k}
		} else {
			as.scd[k] = chainData
			as.scs[k] = ChainState{ID: k, Slot: chainData.Slot}
		}

		// Write chain binary
		if k == "drivechain" || k == "testchain" || k == "bitassets" || k == "thunder" || k == "latestcore" || k == "bitnames" {
			err = writeBinary(&chainData)
			if err != nil {
				println(err.Error())
				return err
			}
		}

	}

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
			// println(k + " = " + v)
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

func writeBinary(cd *ChainData) error {
	var binBytes []byte
	binDir := cd.BinDir + string(os.PathSeparator) + cd.BinName
	target := runtime.GOOS
	switch target {
	case "linux":
		switch cd.ID {
		case "drivechain":
			binBytes = drivechainLinux
		case "testchain":
			binBytes = testchainLinux
		case "bitassets":
			binBytes = bitassetsLinux
		case "thunder":
			binBytes = thunderLinux
		case "latestcore":
			binBytes = latestCoreLinux
		case "bitnames":
			return WriteBitnamesZipContents(cd)
		}
	}
	if len(binBytes) > 0 {
		err := os.WriteFile(binDir, binBytes, 0o755)
		if err != nil {
			return err
		}
	}

	return nil
}

func IsDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// read in ONLY one file
	_, err = f.Readdir(1)

	// and if the file is EOF... well, the dir is empty.
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func WriteBitnamesZipContents(cd *ChainData) error {
	dst := cd.ConfDir
	reader := bytes.NewReader(bitnamesZipLinux)
	zipReader, err := zip.NewReader(reader, int64(len(bitnamesZipLinux)))
	if err != nil {
		return err
	}
	for _, f := range zipReader.File {
		filePath := filepath.Join(dst, f.Name)
		// fmt.Println("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			fmt.Println("invalid file path")
			return err
		}
		if f.FileInfo().IsDir() {
			// fmt.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}
	return nil
}
