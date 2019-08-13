package swaytui

import (
	"fmt"
	"github.com/ragon000/srandr/pkg/sway"
	tui "github.com/marcusolsson/tui-go"
	"image"
	"math"
)

type MonitorWidget struct {
	tui.WidgetBase
        outputs *sway.OutputsWithSelected

	hasBorder      bool
}

func NewMonitorWidget(opts *sway.OutputsWithSelected) *MonitorWidget {
	return &MonitorWidget{
		outputs:        opts,
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
	maxwh := totalWidthHeight(wid.outputs.Outputs)
	totalwidth := float64(maxwh.X) * widthmodifier
	totalheight := maxwh.Y
	canvaswidth := wid.Size().X - (borderoffset * 2)
	canvasheight := wid.Size().Y - (borderoffset * 2)
	factor := math.Ceil(math.Max(totalwidth/float64(canvaswidth), float64(totalheight)/float64(canvasheight)))
	//log.Printf("factor: %v, tw: %v, th: %v, cw: %v, ch: %v\n",factor,totalwidth,totalheight,canvaswidth,canvasheight)

	for _, o := range wid.outputs.Outputs {

		x := int(float64(o.Rect.X)*widthmodifier/factor) + borderoffset
		y := int(float64(o.Rect.Y)/factor) + borderoffset
		w := int(float64(o.Current_mode.Width) * widthmodifier / factor)
		h := int(float64(o.Current_mode.Height) / factor)
		if o.IsEqualTo(*wid.outputs.SelectedOutput) {
			p.WithStyle("red", func(p *tui.Painter) {
				p.DrawRect(x, y, w, h)
			})
		} else {
			p.DrawRect(x, y, w, h)
		}
		mx, my := middle(x, y, w, h)
		p.DrawText(mx-(len(o.Name)/2), my-1, o.Name)
		whtext := fmt.Sprintf("%dx%d", o.Current_mode.Width, o.Current_mode.Height)
                activetext := fmt.Sprintf("Active: %v",o.Active)
		p.DrawText(mx-(len(whtext)/2), my, whtext)
		p.DrawText(mx-(len(activetext)/2), my+1, activetext)
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
		width = int(math.Max(float64(o.Current_mode.Width+o.Rect.X), float64(width)))
		height = int(math.Max(float64(o.Current_mode.Height+o.Rect.Y), float64(height)))
	}
	return image.Pt(width, height)
}

func middle(x, y, w, h int) (int, int) {
	return x + (w / 2), y + (h / 2)
}
