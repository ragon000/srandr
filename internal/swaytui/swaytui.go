package swaytui

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/ragon000/srandr/pkg/sway"
	"log"
	"net"
)

func createOutputsTable(outputs []sway.Output) string{

  var outputArray [][]string

  for i, o := range outputs {
    [i]outputArray := { o.Name, o.Active, o.Primary, o.Current_mode.Width, o.Current_mode.Height }
  }
  return outputArray

}

func Start(sock net.Conn) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	p := widgets.NewParagraph()
	p.Text = "Hello World!"
	p.SetRect(0, 0, 25, 5)

	ui.Render(p)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		}
	}
}
