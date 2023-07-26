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
	}

	mui = NewMainUI(as)
	mui.Refresh()

	mui.as.w.ShowAndRun()
}
