package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func Eggtimer() {
	go func() {
		// create new window
		w := new(app.Window)
		w.Option(app.Title("Timer"))
		// width, height
		w.Option(app.Size(unit.Dp(400), unit.Dp(600)))

		// ops = operations from the ui
		var ops op.Ops
		// define start button as a clickable widget
		var startButton widget.Clickable
		// th = font, color, etc.
		th := material.NewTheme()

		// listen for events in window
		for {
			event := w.Event()
			// handle types of events
			switch typ := event.(type) {
			// FrameEvent = app should re-render display
			case app.FrameEvent:
				gtx := app.NewContext(&ops, typ)
				// actual button, but points to button widget
				btn := material.Button(th, &startButton, "Start")
				// button applies itself to the context (has click animation by default)
				btn.Layout(gtx)
				// send operation to the frame
				typ.Frame(gtx.Ops)
			// DestroyEvent = handle app being closed
			case app.DestroyEvent:
				os.Exit(0)
			}
		}
	}()
	// start by handing control of thread to OS
	app.Main()
}

func EggtimerWithLayout() {
	go func() {
		// create new window
		w := new(app.Window)
		w.Option(app.Title("Timer"))
		// width, height
		w.Option(app.Size(unit.Dp(400), unit.Dp(600)))

		// ops = operations from the ui
		var ops op.Ops
		// define start button as a clickable widget
		var startButton widget.Clickable
		// th = font, color, etc.
		th := material.NewTheme()

		// listen for events in window
		for {
			event := w.Event()
			// handle types of events
			switch typ := event.(type) {
			// FrameEvent = app should re-render display
			case app.FrameEvent:
				gtx := app.NewContext(&ops, typ)

				layout.Flex{
					Axis:    layout.Vertical,   // items are placed top > bottom
					Spacing: layout.SpaceStart, // empty space is left at start (vertical so top)
				}.Layout(gtx,
					// rigid accepts a widget (anything returning dimensions)
					layout.Rigid(
						func(gtx layout.Context) layout.Dimensions {
							// button is defined same as before; only returned
							btn := material.Button(th, &startButton, "Start")
							return btn.Layout(gtx)
						},
					),
					layout.Rigid(
						// spacer struct's layout function returns dimensions
						// a bit of space at bottom of window (beneath button)
						layout.Spacer{Height: unit.Dp(25)}.Layout,
					),
				)

				// send operation to the frame
				typ.Frame(gtx.Ops)
			// DestroyEvent = handle app being closed
			case app.DestroyEvent:
				os.Exit(0)
			}
		}
	}()
	// start by handing control of thread to OS
	app.Main()
}

// shorthand
type C = layout.Context
type D = layout.Dimensions

// for progressbar
var progress float32
var progressIncrementer chan float32

func EggtimerComplete() {
	// setup channel to generate progress
	// (progress is TRACKED in draw)
	progressIncrementer = make(chan float32)
	go func() {
		for {
			time.Sleep(time.Second / 25)
			progressIncrementer <- 0.004
		}
	}()
	// gio channel
	go func() {
		// create new window
		w := new(app.Window)
		w.Option(app.Title("Timer"))
		// width, height
		w.Option(app.Size(unit.Dp(400), unit.Dp(600)))
		// call draw function and handle potential error inline using ';'
		if err := draw(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

// helper function, returns an error (or nil) so main func knows if it worked
func draw(w *app.Window) error {

	// ops = operations from the ui
	var ops op.Ops
	// define start button as a clickable widget
	var startButton widget.Clickable
	// th = font, color, etc.
	th := material.NewTheme()

	// for timer text field
	var durationInput widget.Editor
	var active bool
	var duration float32

	// listen for events in incrementer channel
	go func() {
		for range progressIncrementer {
			// if active and progress isn't complete, add predefined value
			// to progress (0.004)
			if active && progress < 1 {
				progress += 1.0 / 25.0 / duration
				if progress >= 1 {
					progress = 1
				}
				// invalidate = tells window it's out of date so should refresh
				//w.Invalidate()
			}
		}
	}()

	for {
		// simplified version of checking type
		switch e := w.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// check if button has been clicked and flip "active" status:
			if startButton.Clicked(gtx) {
				active = !active

				// reset progress if complete
				if progress >= 1 {
					progress = 0
				}

				// get max time from input widget
				inputString := durationInput.Text()
				inputString = strings.TrimSpace(inputString)
				inputFloat, err := strconv.ParseFloat(inputString, 32)
				if err == nil {
					// if converting str to float was successful, set duration
					duration = float32(inputFloat)
					duration = duration / (1 - progress)
				}
			}

			layout.Flex{
				Axis:    layout.Vertical,   // items are placed top > bottom
				Spacing: layout.SpaceStart, // empty space is left at start (vertical so top)
			}.Layout(gtx,
				// rigid accepts a widget (anything returning dimensions)
				layout.Rigid(
					// draw circle/ egg
					func(gtx C) D {
						return makeEggShape(gtx)
					},
				),
				layout.Rigid(
					// render input box
					func(gtx C) D {
						// convert editor to text field
						editor := material.Editor(th, &durationInput, "sec(s)")
						// set properties
						durationInput.SingleLine = true
						durationInput.Alignment = text.Middle
						// text-based countdown (show remaining time in input box)
						if active && progress < 1 {
							remaining := (1 - progress) * duration
							oneDecimal := fmt.Sprintf("%.1f", math.Round(float64(remaining)*10)/10)
							durationInput.SetText(oneDecimal)
						}
						// define margin/ inset and border
						margins := layout.Inset{
							Top:    unit.Dp(0),
							Right:  unit.Dp(170),
							Bottom: unit.Dp(40),
							Left:   unit.Dp(170),
						}
						border := widget.Border{
							Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
							CornerRadius: unit.Dp(3),
							Width:        unit.Dp(2),
						}
						// return margin. within it is returned border containing editor
						return margins.Layout(gtx,
							func(gtx C) D {
								return border.Layout(gtx, editor.Layout)
							},
						)
					},
				),
				layout.Rigid(
					func(gtx C) D {
						// progress is value between 0 & 1 (used as %)
						bar := material.ProgressBar(th, progress)
						// use more efficient re-rendering
						if active && progress < 1 {
							// if active and in progress: refresh in 1/25 seconds
							inv := op.InvalidateCmd{At: gtx.Now.Add(time.Second / 25)}
							gtx.Execute(inv)
						}
						return bar.Layout(gtx)
					},
				),
				layout.Rigid(
					func(gtx C) D {
						// add margin around components
						// uniform = all four margins are the same
						margins := layout.UniformInset(unit.Dp(25))
						// place other components inside the margins
						return margins.Layout(gtx,
							// margins is returned by a function returning the button function
							func(gtx C) D {
								var text string
								if !active {
									text = "Start"
								} else if active && progress < 1 {
									text = "Stop"
								} else {
									text = "Reset"
								}
								// button is defined same as before; only returned
								btn := material.Button(th, &startButton, text)
								return btn.Layout(gtx)
							},
						)
					},
				),
				layout.Rigid(
					// spacer struct's layout function returns dimensions
					// a bit of space at bottom of window (beneath button)
					layout.Spacer{Height: unit.Dp(25)}.Layout,
				),
			)

			// send operation to the frame
			e.Frame(gtx.Ops)
		case app.DestroyEvent:
			// returns nil if normal closure, otherwise error
			return e.Err
		}
	}
}

func makeCircle(gtx C) D {
	circle := clip.Ellipse{
		// using constraints means position of circle changes
		// based on size of window
		Min: image.Pt(gtx.Constraints.Max.X/2-120, 0),
		Max: image.Pt(gtx.Constraints.Max.X/2+120, 240),
	}.Op(gtx.Ops)
	// define color (red, fully opaque) and apply it to circle
	color := color.NRGBA{R: 200, A: 255}
	paint.FillShape(gtx.Ops, color, circle)
	// return dimensions based on height of circle
	d := image.Point{Y: 400}
	return D{Size: d}
}

func makeEggShape(gtx C) D {
	// Draw a custom path, shaped like an egg
	var eggPath clip.Path
	op.Offset(image.Pt(gtx.Dp(200), gtx.Dp(125))).Add(gtx.Ops)
	eggPath.Begin(gtx.Ops)
	// Rotate from 0 to 360 degrees
	for deg := 0.0; deg <= 360; deg++ {
		// Source for egg drawing function: https://observablehq.com/@toja/egg-curve
		// Convert degrees to radians
		rad := deg * math.Pi / 180
		// Trig gives the distance in X and Y direction
		cosT := math.Cos(rad)
		sinT := math.Sin(rad)
		// Constants to define the eggshape
		a := 110.0
		b := 150.0
		d := 20.0
		// The x/y coordinates
		x := a * cosT
		y := -(math.Sqrt(b*b-d*d*cosT*cosT) + d*sinT) * sinT
		// Finally the point on the outline
		p := f32.Pt(float32(x), float32(y))
		// Draw the line to this point
		eggPath.LineTo(p)
	}
	// Close the path
	eggPath.Close()

	// Get hold of the actual clip
	eggArea := clip.Outline{Path: eggPath.End()}.Op()

	// Fill the shape
	// color changes to reflect progress (yellow -> red)
	color := color.NRGBA{R: 255, G: uint8(239 * (1 - progress)), B: uint8(174 * (1 - progress)), A: 255}
	paint.FillShape(gtx.Ops, color, eggArea)

	d := image.Point{Y: 375}
	return layout.Dimensions{Size: d}
}
