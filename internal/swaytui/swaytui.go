package swaytui

import (
	tui "github.com/marcusolsson/tui-go"
	"github.com/ragon000/srandr/pkg/sway"
	//	"log"
	//	"math"
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
	redrawOutputTable(swayconn.Outputs, outputtable)

	outputwidget := NewMonitorWidget(&swayconn.Outputs)
	outputwidget.SetSizePolicy(tui.Expanding, tui.Expanding)
	outputwidget.SetBorder(true)

	tutorialtext := tui.NewLabel("q -> exit, hjkl -> movement, <enter> -> select, r -> reset")
	tutorialtextbox := tui.NewVBox(tutorialtext)
	tutorialtextbox.SetBorder(true)
	root := tui.NewVBox(
		outputtablebox,
		outputwidget,
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
	theme.SetStyle("default", tui.Style{Fg: tui.ColorWhite})
	theme.Style("default")
	ui.SetTheme(theme)
  ui.SetKeybinding("l", func() {
    outputwidget.SelectedOutput = sway.RightOf(outputwidget.SelectedOutput, &swayconn.Outputs)
  })
  ui.SetKeybinding("h", func() {
    outputwidget.SelectedOutput = sway.LeftOf(outputwidget.SelectedOutput, &swayconn.Outputs)
  })
  ui.SetKeybinding("j", func() {
    outputwidget.SelectedOutput = sway.DownOf(outputwidget.SelectedOutput, &swayconn.Outputs)
  })
  ui.SetKeybinding("k", func() {
    outputwidget.SelectedOutput = sway.UpOf(outputwidget.SelectedOutput, &swayconn.Outputs)
  })
	ui.SetKeybinding("q", func() { ui.Quit() })
	ui.SetKeybinding("r", func() {
		swayconn.GetOutputsFromSocket()
		redrawOutputTable(swayconn.Outputs, outputtable)
	})
	if err := ui.Run(); err != nil {
		panic(err)
	}
}

func redrawOutputTable(outputs []sway.Output, outputtable *tui.Table) {
	outputtable.RemoveRows()
	for _, o := range generateOutputTableRows(outputs) {
		outputtable.AppendRow(o...)

	}

}

func generateOutputTableRows(outputs []sway.Output) [][]tui.Widget {

	table := make([][]tui.Widget, len(outputs)+1, 4)
	table[0] = []tui.Widget{tui.NewLabel("Name"),
		tui.NewLabel("Make"),
		tui.NewLabel("Model"),
		tui.NewLabel("Active"),
		tui.NewLabel("Width"),
		tui.NewLabel("Height"),
		tui.NewLabel("X"),
		tui.NewLabel("Y"),
	}

	for i, o := range outputs {
		table[i+1] = []tui.Widget{
			tui.NewLabel(o.Name),
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
