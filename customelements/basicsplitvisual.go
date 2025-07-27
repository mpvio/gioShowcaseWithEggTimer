package customelements

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
)

type BasicSplitVisual struct{}

// Test colors.
var (
	background = color.NRGBA{R: 0xC0, G: 0xC0, B: 0xC0, A: 0xFF}
	red        = color.NRGBA{R: 0xC0, G: 0x40, B: 0x40, A: 0xFF}
	green      = color.NRGBA{R: 0x40, G: 0xC0, B: 0x40, A: 0xFF}
	blue       = color.NRGBA{R: 0x40, G: 0x40, B: 0xC0, A: 0xFF}
)

// ColorBox creates a widget with the specified dimensions and color.
func ColorBox(gtx layout.Context, size image.Point, color color.NRGBA) layout.Dimensions {
	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: size}
}

// inputs: graphics context and two components to display together (left & right)
func (s BasicSplitVisual) Layout(gtx layout.Context, left, right layout.Widget) layout.Dimensions {
	// set widget widths
	leftside := gtx.Constraints.Min.X / 2
	rightside := gtx.Constraints.Min.X - leftside

	{
		gtx := gtx
		// define width and height (max height of rendered window)
		gtx.Constraints = layout.Exact(image.Pt(leftside, gtx.Constraints.Max.Y))
		// draw left to gtx using above measurements
		left(gtx)
	}
	{
		gtx := gtx
		// set size of right
		gtx.Constraints = layout.Exact(image.Pt(rightside, gtx.Constraints.Max.Y))
		// set transformation: move this widget so it starts where left ends
		trans := op.Offset(image.Pt(leftside, 0)).Push(gtx.Ops)
		// draw right to gtx
		right(gtx)
		// remove transformation so it doesn't affect future widgets
		trans.Pop()
	}
	return layout.Dimensions{Size: gtx.Constraints.Max}
}

// access the above
func ShowcaseBasicExample(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// returns tuple of two different widgets
	return BasicSplitVisual{}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return FillWithLabel(gtx, th, "Left", red)
	}, func(gtx layout.Context) layout.Dimensions {
		return FillWithLabel(gtx, th, "Right", blue)
	})
}

// helper function
func FillWithLabel(gtx layout.Context, th *material.Theme, text string, backgroundColor color.NRGBA) layout.Dimensions {
	ColorBox(gtx, gtx.Constraints.Max, backgroundColor)
	return layout.Center.Layout(gtx, material.H3(th, text).Layout)
}
