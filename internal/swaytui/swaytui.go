package swaytui

import (
	"fmt"
	tui "github.com/marcusolsson/tui-go"
	"github.com/ragon000/srandr/pkg/sway"
	"strconv"
	//	"time"
)


func boolToString(b bool) string {
	if b {
		return "x"
	} else {
		return "-"
	}

}

func Start() {
	swayconn, err := sway.CreateSwayConnection()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", swayconn.Outputs)
	outputs := sway.OutputsWithSelected{Outputs: swayconn.Outputs, SelectedOutput: &swayconn.Outputs[0]}
	defer swayconn.CloseConnection()

	outputtable := tui.NewTable(0, 0)
	outputtable.SetSizePolicy(tui.Maximum, tui.Maximum)
	outputtablebox := tui.NewVBox(outputtable)
	outputtablebox.SetSizePolicy(tui.Maximum, tui.Maximum)
	outputtablebox.SetBorder(true)
	outputtablebox.SetTitle("Outputs")
	redrawOutputTable(outputs, outputtable)

	outputwidget := NewMonitorWidget(&outputs)
	outputwidget.SetSizePolicy(tui.Expanding, tui.Expanding)
	outputwidget.SetBorder(true)
	outputwidgetbox := tui.NewHBox(outputwidget)
	outputwidgetbox.SetBorder(false)
	outputwidgetbox.SetSizePolicy(tui.Expanding, tui.Expanding)
        ttext := "a -> apply, q -> exit, hjkl -> movement, <enter> -> Change Modes, r -> reset"

	tutorialtext := tui.NewLabel(ttext)
	tutorialtextbox := tui.NewVBox(tutorialtext)
	tutorialtextbox.SetBorder(true)
	root := tui.NewVBox(
		outputtablebox,
		outputwidgetbox,
		tutorialtextbox,
	)
	root.SetBorder(true)
	root.SetTitle("srandr")
	ui, err := tui.New(root)
	if err != nil {
		panic(err)
	}
	theme := tui.NewTheme()
	theme.SetStyle("red", tui.Style{Fg: tui.ColorRed})
	theme.SetStyle("selected", tui.Style{Bg: tui.ColorBlack})
	theme.SetStyle("default", tui.Style{Fg: tui.ColorWhite})
	theme.SetStyle("button.focused", tui.Style{Fg: tui.ColorWhite, Bg: tui.ColorBlue})
	theme.SetStyle("button", tui.Style{Fg: tui.ColorWhite, Bg: tui.ColorBlack})
	theme.Style("default")
	ui.SetTheme(theme)

	applyingModes := false
	var focusbuttons []*TextChangeButton
	var activeButton *TextChangeButton
	ui.SetKeybinding("Enter", func() {
		if !applyingModes {
			applyingModes = true
                        tutorialtext.SetText("q -> exit menu, hjkl -> movement, <enter> -> select")

			modeBox := tui.NewVBox()
			modeBox.SetBorder(false)
			focusbuttons = make([]*TextChangeButton, len(outputs.SelectedOutput.Modes))
			for i, m := range outputs.SelectedOutput.Modes {
				if m.IsEqualTo(*outputs.SelectedOutput.Current_mode) {
					focusbuttons[i] = NewTextChangeButton(fmt.Sprintf("> %vx%v@%v", m.Width, m.Height, m.Refresh))
				} else {
					focusbuttons[i] = NewTextChangeButton(fmt.Sprintf("  %vx%v@%v", m.Width, m.Height, m.Refresh))
				}
				focusbuttons[i].SetSizePolicy(tui.Maximum, tui.Maximum)

				focusbuttons[i].OnActivated(func(b *TextChangeButton) {
					for ind, b := range focusbuttons {
						b.SetText(generateModeText(outputs.SelectedOutput.Modes[ind], false))

						if b.IsFocused() {
							outputs.SelectedOutput.Current_mode = &outputs.SelectedOutput.Modes[ind]
							b.SetText(generateModeText(outputs.SelectedOutput.Modes[ind], true))
							redrawOutputTable(outputs, outputtable)
						}
					}
				})
				modeBox.Append(focusbuttons[i])
			}
			if outputs.SelectedOutput.Active {
				activeButton = NewTextChangeButton("Active")
			} else {
				activeButton = NewTextChangeButton("Not Active")
			}
			activeButton.SetSizePolicy(tui.Maximum, tui.Maximum)
			activeButton.OnActivated(func(b *TextChangeButton) {
				if outputs.SelectedOutput.Active {
					outputs.SelectedOutput.Active = false
					activeButton.SetText("Not Active")
					redrawOutputTable(outputs, outputtable)
				} else {
					outputs.SelectedOutput.Active = true
					activeButton.SetText("Active")
					redrawOutputTable(outputs, outputtable)
				}
			})
			spacer := tui.NewPadder(1, 0, tui.NewSpacer())
			selectBox := tui.NewHBox(activeButton, spacer, modeBox)
			selectBox.SetBorder(true)
			selectBox.SetSizePolicy(tui.Preferred, tui.Maximum)
			outputwidgetbox.Prepend(selectBox)
			go activeButton.SetFocused(true)
		}

	})
	ui.SetKeybinding("a", func() {
		//       err := swayconn.ApplyOutputs(outputs.Outputs)
		//       if err != nil {
		//         panic(err)
		//       }

	})
	ui.SetKeybinding("l", func() {
		if !applyingModes {
			outputs.SelectedOutput = sway.RightOf(outputs.SelectedOutput, outputs.Outputs)
					redrawOutputTable(outputs, outputtable)
		} else {
			if activeButton.IsFocused() {
				activeButton.SetFocused(false)
				focusbuttons[0].SetFocused(true)
			}

		}
	})
	ui.SetKeybinding("h", func() {
		if !applyingModes {
			outputs.SelectedOutput = sway.LeftOf(outputs.SelectedOutput, outputs.Outputs)
					redrawOutputTable(outputs, outputtable)
		} else {
			if !activeButton.IsFocused() {
				activeButton.SetFocused(true)
				for _, b := range focusbuttons {
					b.SetFocused(false)
				}
			}

		}
	})
	ui.SetKeybinding("j", func() {
		if !applyingModes {
			outputs.SelectedOutput = sway.DownOf(outputs.SelectedOutput, outputs.Outputs)
					redrawOutputTable(outputs, outputtable)
		} else {
			if !activeButton.IsFocused() {
				length := len(focusbuttons)
				for i, b := range focusbuttons {
					if b.IsFocused() {
						if i != length-1 {
							b.SetFocused(false)
							focusbuttons[i+1].SetFocused(true)
						}
						return
					}
				}
			}
		}
	})
	ui.SetKeybinding("k", func() {
		if !applyingModes {
			outputs.SelectedOutput = sway.UpOf(outputs.SelectedOutput, outputs.Outputs)
					redrawOutputTable(outputs, outputtable)
		} else {
			if !activeButton.IsFocused() {
				for i, b := range focusbuttons {
					if b.IsFocused() {
						if i != 0 {
							b.SetFocused(false)
							focusbuttons[i-1].SetFocused(true)
						}
						return
					}
				}
			}
		}
	})
	ui.SetKeybinding("q", func() {
		if applyingModes == false {
			ui.Quit()
		} else {
			outputwidgetbox.Remove(0)
			applyingModes = false
                        tutorialtext.SetText(ttext)

		}
	})
	ui.SetKeybinding("r", func() {
		if !applyingModes {
			swayconn.GetOutputsFromSocket()
			outputs.Outputs = swayconn.Outputs
			redrawOutputTable(outputs, outputtable)
		}
	})
	if err := ui.Run(); err != nil {
		panic(err)
	}
}

func redrawOutputTable(outputs sway.OutputsWithSelected, outputtable *tui.Table) {
	outputtable.RemoveRows()
	outputstrings := generateOutputTableRows(outputs)
	for _, o := range outputstrings {
		outputtable.AppendRow(o...)
	}

}

func generateOutputTableRows(outputs sway.OutputsWithSelected) [][]tui.Widget {
	table := make([][]tui.Widget, len(outputs.Outputs)+1, 4)
	table[0] = []tui.Widget{
		tui.NewLabel("  Name"),
		tui.NewLabel("Make"),
		tui.NewLabel("Model"),
		tui.NewLabel("Active"),
		tui.NewLabel("Width"),
		tui.NewLabel("Height"),
		tui.NewLabel("X"),
		tui.NewLabel("Y"),
	}

	for i, o := range outputs.Outputs {
		table[i+1] = []tui.Widget{
			tui.NewLabel(selectedSign(o, outputs) + o.Name),
			tui.NewLabel(o.Make),
			tui.NewLabel(o.Model),
			tui.NewLabel(strconv.FormatBool(o.Active)),
			tui.NewLabel(strconv.FormatInt(int64(o.Current_mode.Width), 10) + "px"),
			tui.NewLabel(strconv.FormatInt(int64(o.Current_mode.Height), 10) + "px"),
			tui.NewLabel(strconv.FormatInt(int64(o.Rect.X), 10)),
			tui.NewLabel(strconv.FormatInt(int64(o.Rect.Y), 10)),
		}
	}
	return table
}

func selectedSign(o sway.Output, outputs sway.OutputsWithSelected) string {
	if o.IsEqualTo(*outputs.SelectedOutput) {
		return "> "
	}
	return "  "
}

func generateModeText(m sway.Mode, a bool) string {
	asdf := " "
	if a {
		asdf = ">"
	}
	return fmt.Sprintf("%v %vx%v@%v", asdf, m.Width, m.Height, m.Refresh)

}
