package main

var (
	as  *AppState
	mui *MainUI
)

func main() {
	as = NewAppState("com.layertwolabs.dclauncher", "Drivechain Launcher")

	err := ConfInit(as)
	if err != nil {
		println(err.Error())
		panic(err)
	}

	mui = NewMainUI(as)
	mui.Refresh()

	StartDrivechainStateUpdate(as, mui)

	// Set intercept to shutdown chains if needed
	mui.as.w.ShowAndRun()
}
