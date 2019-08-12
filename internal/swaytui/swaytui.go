package swaytui

import (
	tui "github.com/marcusolsson/tui-go"
	"github.com/ragon000/srandr/pkg/sway"
        "fmt"
	//	"math"
	"strconv"
	//	"time"
)

type OutputsWithSelected struct {
	SelectedOutput sway.Output
	Outputs        []sway.Output
}

func boolToString(b bool) string {
	if b {
		return "x"
	} else {
		return "-"
	}

}

func Start() {
	swayconn, err := sway.CreateSwayConnection()
	outputs := OutputsWithSelected{Outputs: swayconn.Outputs, SelectedOutput: swayconn.Outputs[0]}
	defer swayconn.CloseConnection()
	if err != nil {
		panic(err)
	}

	outputtable := tui.NewTable(0, 0)
	outputtable.SetSizePolicy(tui.Maximum, tui.Maximum)
	outputtablebox := tui.NewVBox(outputtable)
	outputtablebox.SetSizePolicy(tui.Maximum, tui.Maximum)
	outputtablebox.SetBorder(true)
	outputtablebox.SetTitle("Outputs")
	redrawOutputTable(outputs, outputtable)

	outputwidget := NewMonitorWidget(swayconn.Outputs)
	outputwidget.SetSizePolicy(tui.Expanding, tui.Expanding)
	outputwidget.SetBorder(false)
        outputwidgetbox := tui.NewVBox(outputwidget)
        outputwidgetbox.SetBorder(true)
        outputwidgetbox.SetSizePolicy(tui.Expanding, tui.Expanding)

	tutorialtext := tui.NewLabel("a -> apply, q -> exit, hjkl -> movement, <enter> -> Change Modes, r -> reset")
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
	theme.Style("default")
	ui.SetTheme(theme)

        applyingModes := false

        ui.SetKeybinding("Enter", func() {
          applyingModes = true

        modeBox := tui.NewVBox()
        for _, m := range outputs.SelectedOutput.Modes {
          button := tui.NewButton(fmt.Sprintf("%vx%v@%v", m.Width, m.Height, m.Refresh))

          button.OnActivated(func(b *tui.Button){
            outputs.SelectedOutput.Current_mode = m
          })
          modeBox.Append(button)
        }
        activeButton := tui.NewButton("Active")
        selectBox := tui.NewHBox(activeButton, modeBox)
        outputwidgetbox.Remove(0)
        outputwidgetbox.Append(selectBox)

        })
        ui.SetKeybinding("a", func() {
          err := swayconn.ApplyOutputs(outputs.Outputs)
          if err != nil {
            panic(err)
          }

        })
	ui.SetKeybinding("l", func() {
          if !applyingModes {
		outputwidget.SelectedOutput = sway.RightOf(outputs.SelectedOutput, outputs.Outputs)
              }
	})
	ui.SetKeybinding("h", func() {
          if !applyingModes {
		outputwidget.SelectedOutput = sway.LeftOf(outputs.SelectedOutput, outputs.Outputs)
              }
	})
	ui.SetKeybinding("j", func() {
          if !applyingModes {
		outputwidget.SelectedOutput = sway.DownOf(outputs.SelectedOutput, outputs.Outputs)
              }
	})
	ui.SetKeybinding("k", func() {
          if !applyingModes {
		outputwidget.SelectedOutput = sway.UpOf(outputs.SelectedOutput, outputs.Outputs)
              }
	})
	ui.SetKeybinding("q", func() {
          if !applyingModes {
          ui.Quit() 
              } else {
            outputwidgetbox.Remove(0)
            outputwidgetbox.Append(outputwidget)

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




func redrawOutputTable(outputs OutputsWithSelected, outputtable *tui.Table) {
	outputtable.RemoveRows()
	outputstrings := generateOutputTableRows(outputs)
	for _, o := range outputstrings {
		outputtable.AppendRow(o...)
	}

}

func generateOutputTableRows(outputs OutputsWithSelected) [][]tui.Widget {
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
			tui.NewLabel(strconv.FormatInt(int64(o.Rect.Width), 10) + "px"),
			tui.NewLabel(strconv.FormatInt(int64(o.Rect.Height), 10) + "px"),
			tui.NewLabel(strconv.FormatInt(int64(o.Rect.X), 10)),
			tui.NewLabel(strconv.FormatInt(int64(o.Rect.Y), 10)),
		}
	}
	return table
}

func selectedSign(o sway.Output, outputs OutputsWithSelected) string{
    if o.IsEqualTo(outputs.SelectedOutput) {
      return "> "
    }
    return "  "
}
