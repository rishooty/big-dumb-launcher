package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

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
