package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type AppState struct {
	a   fyne.App
	w   fyne.Window
	t   CustomTheme
	dcd ChainData
	dcs ChainState
	scd map[string]ChainData
	scs map[string]ChainState
	cp  map[string]ChainProvider
}

func NewAppState(id string, title string) *AppState {
	a := app.NewWithID(id)
	w := a.NewWindow(title)
	t := NewCustomTheme()
	a.Settings().SetTheme(t)

	return &AppState{
		a:   a,
		w:   w,
		t:   *t,
		scd: make(map[string]ChainData),
		scs: make(map[string]ChainState),
	}
}
