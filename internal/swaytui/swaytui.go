package swaytui

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/drawille"
	"github.com/gizak/termui/v3/widgets"
	"github.com/ragon000/srandr/pkg/sway"
	"image"
	"log"
	"math"
	"net"
	"strconv"
	"time"
)

var grid = ui.NewGrid()
var outputtable = widgets.NewTable()
var screenview = ui.NewCanvas()
var outputs []sway.Output

func boolToString(b bool) string {
	if b {
		return "x"
	} else {
		return "-"
	}

}

func createOutputsTable(outputs []sway.Output) [][]string {

	//  log.Fatalf("asdf: %+v\n",outputs)

	outputArray := make([][]string, len(outputs)+1, 4)
	outputArray[0] = []string{"Name", "Active", "Width", "Height"}

	for i, o := range outputs {
		//  log.Fatalf("i: %v, o: %v\n",i,o)
		outputArray[i+1] = []string{o.Name, boolToString(o.Active), strconv.FormatInt((int64)(o.Current_mode.Width), 10), strconv.FormatInt((int64)(o.Current_mode.Height), 10)}
	}
	return outputArray

}

func Start(sock net.Conn) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	outputs = sway.GetOutputs(sock)
	go updateGui(sock)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "r":
			outputs = sway.GetOutputs(sock)
		}
	}
}

func updateGui(sock net.Conn) {
	for {


		outputtable.Rows = createOutputsTable(outputs)
		outputtable.SetRect(0, 0, 60, 10)
		screenview.SetRect(0, 0, 60, 60)
		fillCanvasWithScreens()
		ui.Render(screenview)

		time.Sleep(1006 * time.Millisecond)

	}
}

// Returns the total width and height of the workspace
// (int,int) (width,height)
func totalWidthHeight() (int, int) {
	width := 0
	height := 0
	for _, o := range outputs {
		width = int(math.Max(float64(o.Rect.Width),float64(width)))
		height = int(math.Max(float64(o.Rect.Height),float64(height)))
	}
	return width, height
}

func fillCanvasWithScreens() {
	totalwidth, totalheight := totalWidthHeight()
	canvaswidth := screenview.Bounds().Dx()
	canvasheight := screenview.Bounds().Dy()
	factor := math.Ceil(math.Max(float64(totalwidth)/float64(canvaswidth), float64(totalheight)/float64(canvasheight)))
        //log.Printf("factor: %v, tw: %v, th: %v, cw: %v, ch: %v\n",factor,totalwidth,totalheight,canvaswidth,canvasheight)

	for _, o := range outputs {
		drawrect(image.Pt(int(float64(o.Rect.X) / factor), int(float64(o.Rect.Y) / factor)),
                         image.Pt(int(float64(o.Rect.X+o.Rect.Width) / factor), int(float64(o.Rect.Y+o.Rect.Height) / factor)), []rune(o.Name))
	}

}

func drawrect(p1 image.Point, p2 image.Point, outputname []rune) {
        canvasminx := 0//screenview.GetRect().Min.X
        canvasminy := 0//screenview.GetRect().Min.Y
        ol := image.Pt(p1.X+canvasminx,p1.Y+canvasminy)
        ur := image.Pt(p2.X+canvasminx,p2.Y+canvasminy)
        or := image.Pt(ur.X-1, ol.Y)
        ul := image.Point{ol.X+1, ur.Y}
        //log.Printf("%+v\n",ol,ur)
	screenview.SetLine(ol, ul, ui.ColorWhite)
        screenview.SetLine(ol, or, ui.ColorWhite)
        screenview.SetLine(ul, ur, ui.ColorWhite)
	screenview.SetLine(or, ur, ui.ColorWhite)
        middlepoint := middle(p1,p2)
        for i := 0;i<len(outputname);i++ {
        screenview.CellMap[image.Pt(middlepoint.X-(len(outputname)/2)+i,middlepoint.Y)] = drawille.Cell{outputname[i],drawille.Color(ui.ColorRed)}
        
      }
}

func middle(p1 image.Point, p2 image.Point) image.Point {
  p0 := p2.Sub(p1)
return image.Pt(p1.X+(p0.X/2),p1.Y+(p0.Y/2))
}
