package swaytui

import (
	"fmt"
	"github.com/ragon000/srandr/pkg/sway"
	"os"
	tui "github.com/marcusolsson/tui-go"
	"image"
	"math"
)

type MonitorWidget struct {
	tui.WidgetBase

	hasBorder      bool
	outputs        *[]sway.Output
	SelectedOutput *sway.Output
}

func NewMonitorWidget(opts *[]sway.Output) *MonitorWidget {
	return &MonitorWidget{
		outputs:        opts,
		SelectedOutput: &(*opts)[0],
	}

}

func (wid *MonitorWidget) Draw(p *tui.Painter) {
	widthmodifier := 2.1 // The Rectangle width gets multiplied by it... fonts are not square but I want to draw a (more or less perfect) square
	var borderoffset int
	if wid.hasBorder {
		borderoffset = 1
		p.DrawRect(0, 0, wid.Size().X, wid.Size().Y)
	} else {
		borderoffset = 0
	}
	maxwh := totalWidthHeight(*wid.outputs)
	totalwidth := float64(maxwh.X) * widthmodifier
	totalheight := maxwh.Y
	canvaswidth := wid.Size().X - (borderoffset * 2)
	canvasheight := wid.Size().Y - (borderoffset * 2)
	factor := math.Ceil(math.Max(totalwidth/float64(canvaswidth), float64(totalheight)/float64(canvasheight)))
	//log.Printf("factor: %v, tw: %v, th: %v, cw: %v, ch: %v\n",factor,totalwidth,totalheight,canvaswidth,canvasheight)

	for _, o := range *wid.outputs {
          fmt.Fprintf(os.Stderr, "wid.SelectedOutput = %v, o = %v", &wid.SelectedOutput.Name, &o.Name)
		x := int(float64(o.Rect.X)*widthmodifier/factor) + borderoffset
		y := int(float64(o.Rect.Y)/factor) + borderoffset
		w := int(float64(o.Rect.Width) * widthmodifier / factor)
		h := int(float64(o.Rect.Height) / factor)
		if &o == wid.SelectedOutput {
			p.WithStyle("red", func(p *tui.Painter) {
				p.DrawRect(x, y, w, h)
			})
		} else {
			p.DrawRect(x, y, w, h)
		}
		mx, my := middle(x, y, w, h)
		p.DrawText(mx-(len(o.Name)/2), my, o.Name)
		whtext := fmt.Sprintf("%dx%d", o.Rect.Width, o.Rect.Height)
		p.DrawText(mx-(len(whtext)/2), my+1, whtext)
	}

}


func (w *MonitorWidget) SizeHint() image.Point {
	return image.Pt(60, 60)

}

func (g *MonitorWidget) SetBorder(enabled bool) {
	g.hasBorder = enabled
}

func totalWidthHeight(outputs []sway.Output) image.Point {
	width := 0
	height := 0
	for _, o := range outputs {
		width = int(math.Max(float64(o.Rect.Width), float64(width)))
		height = int(math.Max(float64(o.Rect.Height), float64(height)))
	}
	return image.Pt(width, height)
}

func middle(x, y, w, h int) (int, int) {
	return x + (w / 2), y + (h / 2)
}
