package main

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type MainUI struct {
	headerContainer  *fyne.Container
	contentContainer *fyne.Container
	footerContainer  *fyne.Container
	as               *AppState
	driveChainRow    DrivechainRow
	sideChainRows    []SidechainRow
}

func NewMainUI(as *AppState) *MainUI {
	mui := &MainUI{
		headerContainer:  container.NewVBox(),
		contentContainer: container.NewStack(),
		footerContainer:  container.NewVBox(),
		as:               as,
	}

	menus := fyne.NewMainMenu(&fyne.Menu{
		Label: "File",
		Items: []*fyne.MenuItem{
			{Label: "Reset Everything", Action: func() {
				dialog.NewConfirm("Reset Everything", "This will delete all data and settings for Drivechain and Sidechains.", func(b bool) {
					if b {
						pu := widget.NewModalPopUp(widget.NewLabel("Resetting, please wait..."), mui.as.w.Canvas())
						pu.Show()
						err := ResetEverything(mui.as)
						if err != nil {
							println(err.Error())
						}
						pu.Hide()
						mui.as.w.Content().Refresh()
					}
				}, as.w).Show()
			}},
		},
	})

	as.w.SetMainMenu(menus)

	lv := container.NewVBox()

	mui.driveChainRow = NewDrivechainRow(mui, mui.as.cp["drivechain"], lv)
	for k, cp := range as.cp {
		if k != "drivechain" {
			mui.sideChainRows = append(mui.sideChainRows, NewSidechainRow(mui, cp, lv))
		}
	}

	scrl := container.NewScroll(lv)
	mui.contentContainer.Add(scrl)

	as.w.SetContent(container.NewBorder(mui.headerContainer, mui.footerContainer, nil, nil, mui.contentContainer))
	as.w.Resize(fyne.NewSize(540, 720))
	return mui
}

func (mui *MainUI) Refresh() {
	for _, scr := range mui.sideChainRows {
		scr.Refresh(mui)
	}
	mui.driveChainRow.Refresh(mui)
}

type DrivechainRow struct {
	Title       *widget.RichText
	Desc        *widget.RichText
	Blocks      *widget.RichText
	StartButton *widget.Button
	StopButton  *widget.Button
	MineButton  *widget.Button
}

func NewDrivechainRow(mui *MainUI, cp ChainProvider, c *fyne.Container) DrivechainRow {
	dcr := DrivechainRow{
		Title:  widget.NewRichTextWithText(cp.Name),
		Desc:   widget.NewRichTextWithText(cp.Description),
		Blocks: widget.NewRichTextWithText("Blocks: " + strconv.Itoa(mui.as.dcs.Height)),
		StartButton: widget.NewButtonWithIcon("Launch Chain", mui.as.t.Icon(StartIcon), func() {
			pu := widget.NewModalPopUp(widget.NewLabel("Launching Drivechain..."), mui.as.w.Canvas())
			pu.Show()
			time.AfterFunc(time.Duration(1)*time.Second, func() {
				pu.Hide()
			})
			LaunchChain(&mui.as.dcd, &mui.as.dcs, mui)
		}),
		StopButton: widget.NewButtonWithIcon("Stop Chain", mui.as.t.Icon(StopIcon), func() {
			mui.as.dcs.Automine = false
			pu := widget.NewModalPopUp(widget.NewLabel("Stoping Drivechain..."), mui.as.w.Canvas())
			pu.Show()
			time.AfterFunc(time.Duration(1)*time.Second, func() {
				pu.Hide()
			})
			StopChain(&mui.as.dcd, &mui.as.dcs, mui.as)
		}),
		MineButton: widget.NewButtonWithIcon("Start Mining", mui.as.t.Icon(MineIcon), func() {
			mui.as.dcs.Automine = false
			mui.Refresh()
		}),
	}

	dcr.StartButton.Alignment = widget.ButtonAlignTrailing
	dcr.StartButton.IconPlacement = widget.ButtonIconTrailingText
	dcr.StartButton.Importance = widget.HighImportance

	dcr.StopButton.Alignment = widget.ButtonAlignTrailing
	dcr.StopButton.IconPlacement = widget.ButtonIconTrailingText
	dcr.MineButton.Alignment = widget.ButtonAlignTrailing
	dcr.MineButton.IconPlacement = widget.ButtonIconTrailingText

	dcr.Title.Segments[0].(*widget.TextSegment).Style = widget.RichTextStyle{
		Alignment: fyne.TextAlignLeading,
		SizeName:  theme.SizeNameHeadingText,
		ColorName: theme.ColorNameForeground,
		TextStyle: fyne.TextStyle{Italic: false, Bold: true},
	}

	dcr.Desc.Segments[0].(*widget.TextSegment).Style = widget.RichTextStyle{
		Alignment: fyne.TextAlignLeading,
		SizeName:  theme.SizeNameText,
		ColorName: theme.ColorGray,
		TextStyle: fyne.TextStyle{Italic: false, Bold: false},
	}
	dcr.Desc.Wrapping = fyne.TextWrapWord

	dcr.Blocks.Segments[0].(*widget.TextSegment).Style = widget.RichTextStyle{
		Alignment: fyne.TextAlignLeading,
		SizeName:  theme.SizeNameCaptionText,
		ColorName: theme.ColorGray,
		TextStyle: fyne.TextStyle{Italic: false, Bold: false},
	}

	ftr := container.NewHBox(dcr.Blocks)

	bck := NewThemedRectangle(theme.ColorNameMenuBackground)
	bck.CornerRadius = 8
	bck.Refresh()

	stk := container.NewStack(bck)

	brdr := container.NewBorder(nil, container.NewVBox(&layout.Spacer{FixHorizontal: true, FixVertical: true}, widget.NewSeparator(), ftr), nil,
		container.NewVBox(dcr.StartButton, dcr.StopButton, dcr.MineButton), container.NewVBox(dcr.Title, dcr.Desc))
	stk.Add(container.NewPadded(container.NewPadded(brdr)))
	c.Add(stk)
	return dcr
}

func (dcr *DrivechainRow) Refresh(mui *MainUI) {
	if mui.as.dcs.State == Running {
		dcr.StartButton.Disable()
		dcr.MineButton.Enable()
		dcr.StopButton.Enable()
	} else {
		dcr.StartButton.Enable()
		dcr.MineButton.Disable()
		dcr.StopButton.Disable()
	}
	if mui.as.dcs.Automine {
		dcr.MineButton.Importance = widget.MediumImportance
		dcr.MineButton.SetText("Stop Mining")
		dcr.MineButton.OnTapped = func() {
			mui.as.dcs.Automine = false
			mui.Refresh()
		}
		dcr.MineButton.Refresh()
	} else {
		dcr.MineButton.Importance = widget.HighImportance
		dcr.MineButton.SetText("Start Mining")
		dcr.MineButton.OnTapped = func() {
			mui.as.dcs.Automine = true
			mui.Refresh()
		}
		dcr.MineButton.Refresh()
	}
	mui.driveChainRow.Blocks.Segments[0].(*widget.TextSegment).Text = "Blocks: " + strconv.Itoa(mui.as.dcs.Height)
	mui.driveChainRow.Blocks.Refresh()
	mui.contentContainer.Refresh()
}

type SidechainRow struct {
	Title         *widget.RichText
	Desc          *widget.RichText
	Blocks        *widget.RichText
	StartButton   *widget.Button
	StopButton    *widget.Button
	ChainProivder ChainProvider
}

func NewSidechainRow(mui *MainUI, cp ChainProvider, c *fyne.Container) SidechainRow {
	scr := SidechainRow{
		Title:  widget.NewRichTextWithText(cp.Name),
		Desc:   widget.NewRichTextWithText(cp.Description),
		Blocks: widget.NewRichTextWithText("Blocks: " + strconv.Itoa(mui.as.scs[cp.ID].Height)),
		StartButton: widget.NewButtonWithIcon("Launch Chain", mui.as.t.Icon(StartIcon), func() {
			cd := mui.as.scd[cp.ID]
			cs := mui.as.scs[cp.ID]
			if NeedsActivation(&cd, mui.as) {
				CreateSidechainProposal(mui.as, &cd, &cs)
				ap := widget.NewModalPopUp(widget.NewLabel(fmt.Sprintf("Activating %s...", cp.Name)), mui.as.w.Canvas())
				ap.Show()
				time.AfterFunc(time.Duration(2)*time.Second, func() {
					ap.Hide()

					pu := widget.NewModalPopUp(widget.NewLabel(fmt.Sprintf("Launching %s...", cp.Name)), mui.as.w.Canvas())
					pu.Show()
					time.AfterFunc(time.Duration(1)*time.Second, func() {
						pu.Hide()
					})
					LaunchChain(&cd, &cs, mui)
				})
			} else {
				pu := widget.NewModalPopUp(widget.NewLabel(fmt.Sprintf("Launching %s...", cp.Name)), mui.as.w.Canvas())
				pu.Show()
				time.AfterFunc(time.Duration(1)*time.Second, func() {
					pu.Hide()
				})
				LaunchChain(&cd, &cs, mui)
			}
		}),
		StopButton: widget.NewButtonWithIcon("Stop Chain", mui.as.t.Icon(StopIcon), func() {
			pu := widget.NewModalPopUp(widget.NewLabel(fmt.Sprintf("Stoping %s...", cp.Name)), mui.as.w.Canvas())
			pu.Show()
			time.AfterFunc(time.Duration(1)*time.Second, func() {
				pu.Hide()
			})
			cd := mui.as.scd[cp.ID]
			cs := mui.as.scs[cp.ID]
			StopChain(&cd, &cs, mui.as)
		}),
		ChainProivder: cp,
	}

	scr.StartButton.Alignment = widget.ButtonAlignTrailing
	scr.StartButton.IconPlacement = widget.ButtonIconTrailingText
	scr.StartButton.Importance = widget.HighImportance
	scr.StopButton.Alignment = widget.ButtonAlignTrailing
	scr.StopButton.IconPlacement = widget.ButtonIconTrailingText

	scr.Title.Segments[0].(*widget.TextSegment).Style = widget.RichTextStyle{
		Alignment: fyne.TextAlignLeading,
		SizeName:  theme.SizeNameHeadingText,
		ColorName: theme.ColorNameForeground,
		TextStyle: fyne.TextStyle{Italic: false, Bold: true},
	}

	scr.Desc.Segments[0].(*widget.TextSegment).Style = widget.RichTextStyle{
		Alignment: fyne.TextAlignLeading,
		SizeName:  theme.SizeNameText,
		ColorName: theme.ColorGray,
		TextStyle: fyne.TextStyle{Italic: false, Bold: false},
	}
	scr.Desc.Wrapping = fyne.TextWrapWord

	scr.Blocks.Segments[0].(*widget.TextSegment).Style = widget.RichTextStyle{
		Alignment: fyne.TextAlignLeading,
		SizeName:  theme.SizeNameCaptionText,
		ColorName: theme.ColorGray,
		TextStyle: fyne.TextStyle{Italic: false, Bold: false},
	}

	ftr := container.NewHBox(scr.Blocks)

	bck := NewThemedRectangle(theme.ColorNameMenuBackground)
	bck.CornerRadius = 8
	bck.Refresh()

	stk := container.NewStack(bck)

	brdr := container.NewBorder(nil, container.NewVBox(&layout.Spacer{FixHorizontal: true, FixVertical: true}, widget.NewSeparator(), ftr), nil, container.NewVBox(scr.StartButton, scr.StopButton), container.NewVBox(scr.Title, scr.Desc))
	stk.Add(container.NewPadded(container.NewPadded(brdr)))
	c.Add(stk)
	return scr
}

func (scr *SidechainRow) Refresh(mui *MainUI) {
	if mui.as.dcs.State != Running {
		scr.StartButton.Disable()
		scr.StopButton.Disable()
		return
	}
	if mui.as.scs[scr.ChainProivder.ID].State == Running {
		scr.StartButton.Disable()
		scr.StopButton.Enable()
	} else {
		scr.StartButton.Enable()
		scr.StopButton.Disable()
	}
	scr.Blocks.Segments[0].(*widget.TextSegment).Text = "Blocks: " + strconv.Itoa(mui.as.scs[scr.ChainProivder.ID].Height)
	scr.Blocks.Refresh()
	mui.contentContainer.Refresh()
}
