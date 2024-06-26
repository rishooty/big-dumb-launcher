package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	windowWidth  = 800
	windowHeight = 600
	fontSize     = 24
)

type App struct {
	window           *sdl.Window
	renderer         *sdl.Renderer
	joystick         *sdl.Joystick
	font             *ttf.Font
	currentSelection int
	files            []os.DirEntry
	path             string
	running          bool
	colors           Colors
}

func main() {
	app := &App{}
	err := app.Init()
	if err != nil {
		fmt.Println("Failed to initialize the application:", err)
		return
	}
	defer app.Cleanup()

	app.Run()
}

func (app *App) Init() error {
	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_JOYSTICK); err != nil {
		return err
	}
	if err := ttf.Init(); err != nil {
		sdl.Quit()
		return err
	}

	window, err := sdl.CreateWindow("File Browser", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, windowWidth, windowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		ttf.Quit()
		sdl.Quit()
		return err
	}
	app.window = window

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		window.Destroy()
		ttf.Quit()
		sdl.Quit()
		return err
	}
	app.renderer = renderer

	font, err := ttf.OpenFont("Weiholmir_regular.ttf", fontSize)
	if err != nil {
		renderer.Destroy()
		window.Destroy()
		ttf.Quit()
		sdl.Quit()
		return err
	}
	app.font = font

	app.joystick = sdl.JoystickOpen(0)
	app.running = true
	app.colors = Colors{
		white:  sdl.Color{R: 255, G: 255, B: 255, A: 255},
		yellow: sdl.Color{R: 255, G: 255, B: 0, A: 255},
	}
	app.path = "testpath"
	app.files, err = os.ReadDir(app.path)
	if err != nil {
		return err
	}
	return nil
}

func (app *App) Run() {
	for app.running {
		app.pollInputs()
		app.renderer.Clear()
		app.draw()
		app.renderer.Present()
	}
}

func (app *App) Cleanup() {
	if app.font != nil {
		app.font.Close()
	}
	if app.renderer != nil {
		app.renderer.Destroy()
	}
	if app.window != nil {
		app.window.Destroy()
	}
	if app.joystick != nil {
		app.joystick.Close()
	}
	ttf.Quit()
	sdl.Quit()
}
