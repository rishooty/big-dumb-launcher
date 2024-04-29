package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

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

type Colors struct {
	white  sdl.Color
	yellow sdl.Color
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

var startPressed, selectPressed bool
var startPressTime, selectPressTime time.Time

func (app *App) pollInputs() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.KeyboardEvent:
			if e.Type == sdl.KEYDOWN {
				switch e.Keysym.Sym {
				case sdl.K_UP:
					app.moveSelectUp()
				case sdl.K_DOWN:
					app.moveSelectDown()
				case sdl.K_RETURN:
					app.selectOrLaunch()
				case sdl.K_BACKSPACE:
					app.moveSelectBack()
				}
			}
		case *sdl.JoyAxisEvent:
			if e.Axis == 1 {
				if e.Value < 0 {
					app.moveSelectDown()
				} else if e.Value > 0 {
					app.moveSelectUp()
				}
			}
		case *sdl.JoyButtonEvent:
			fmt.Println(e.Button)
			if e.Button == sdl.CONTROLLER_BUTTON_START {
				if e.State == sdl.PRESSED {
					startPressed = true
					startPressTime = time.Now()
				} else if e.State == sdl.RELEASED {
					startPressed = false
					if selectPressed && time.Since(selectPressTime) >= 3*time.Second {
						app.running = false
					}
				}
			} else if e.Button == sdl.CONTROLLER_BUTTON_BACK {
				if e.State == sdl.PRESSED {
					selectPressed = true
					selectPressTime = time.Now()
				} else if e.State == sdl.RELEASED {
					selectPressed = false
					if startPressed && time.Since(startPressTime) >= 3*time.Second {
						app.running = false
					}
				}
			} else if e.State == sdl.PRESSED {
				switch e.Button {
				case sdl.CONTROLLER_BUTTON_A:
					app.selectOrLaunch()
				case sdl.CONTROLLER_BUTTON_B:
					app.moveSelectBack()
				}
			}
		}
	}
}

func (app *App) draw() {
	for i, file := range app.files {
		text := file.Name()
		if file.IsDir() {
			text += "/"
		}

		// Choose color based on selection
		var color sdl.Color
		if i == app.currentSelection {
			color = app.colors.yellow
		} else {
			color = app.colors.white
		}

		surface, err := app.font.RenderUTF8Blended(text, color)
		if err != nil {
			fmt.Println("Could not render text:", err)
			return
		}
		defer surface.Free()

		texture, err := app.renderer.CreateTextureFromSurface(surface)
		if err != nil {
			fmt.Println("Could not create texture:", err)
			return
		}
		defer texture.Destroy()

		app.renderer.Copy(texture, nil, &sdl.Rect{X: 10, Y: 10 + int32(i*fontSize), W: surface.W, H: surface.H})
	}
}

func (app *App) moveSelectUp() {
	if app.currentSelection > 0 {
		app.currentSelection--
	}
}

func (app *App) moveSelectDown() {
	if app.currentSelection < len(app.files)-1 {
		app.currentSelection++
	}
}

func (app *App) moveSelectBack() {
	parent := filepath.Dir(app.path)
	if parent != app.path {
		app.path = parent
		app.files, _ = os.ReadDir(app.path)
	}
}

func (app *App) selectOrLaunch() {
	if app.files[app.currentSelection].IsDir() {
		app.path = filepath.Join(app.path, app.files[app.currentSelection].Name())
		app.files, _ = os.ReadDir(app.path)
		app.currentSelection = 0
	} else {
		selectedFileName := app.files[app.currentSelection].Name()

		parentDirName := filepath.Base(app.path)
		// fmt.Println("Parent directory name:", parentDirName)

		fullFilePath := filepath.Join(app.path, selectedFileName)
		// fmt.Println("Selected file full path:", fullFilePath)

		var cmd *exec.Cmd
		if selectedFileName == "self.txt" {
			cmd = exec.Command(parentDirName)
		} else {
			cmd = exec.Command(parentDirName, fullFilePath)
		}

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Run the command
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error executing command:", err)
			return
		}
	}
}

func (app *App) Run() {
	for app.running {
		app.pollInputs()
		app.renderer.Clear()
		app.draw()
		app.renderer.Present()
	}
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
