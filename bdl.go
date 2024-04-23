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

func main() {
	// Initialize SDL
	sdl.Init(sdl.INIT_VIDEO | sdl.INIT_JOYSTICK)
	defer sdl.Quit()

	// Initialize SDL_ttf
	ttf.Init()
	defer ttf.Quit()

	// Create a window
	window, err := sdl.CreateWindow("File Browser", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		windowWidth, windowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println("Could not create window:", err)
		return
	}
	defer window.Destroy()

	// Create a renderer
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println("Could not create renderer:", err)
		return
	}
	defer renderer.Destroy()

	// Load a font
	font, err := ttf.OpenFont("Weiholmir_regular.ttf", fontSize)
	if err != nil {
		fmt.Println("Could not load font:", err)
		return
	}
	defer font.Close()

	// Define colors
	white := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	yellow := sdl.Color{R: 255, G: 255, B: 0, A: 255}

	var currentSelection int
	path := "testpath"
	files, err := os.ReadDir(path)
	if err != nil {
		// Handle error
	}

	var startPressed, selectPressed bool
	var startPressTime, selectPressTime time.Time

	// Open the first joystick
	joystick := sdl.JoystickOpen(0)
	defer joystick.Close()

	// Main loop
	running := true
	for running {
		if joystick == nil {
			joystick.Close()
			joystick = sdl.JoystickOpen(0)
			if joystick == nil {
				fmt.Println("Reconnecting joystick failed")
			}
		}
		// Inside your main loop
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.KeyboardEvent:
				if e.Type == sdl.KEYDOWN {
					switch e.Keysym.Sym {
					case sdl.K_UP:
						// Move selection up
						if currentSelection > 0 {
							currentSelection--
						}
					case sdl.K_DOWN:
						// Move selection down
						if currentSelection < len(files)-1 {
							currentSelection++
						}
					case sdl.K_RETURN:
						// Enter the selected directory or select the file
						if files[currentSelection].IsDir() {
							// Change the current directory
							path = filepath.Join(path, files[currentSelection].Name())
							// Re-read the directory
							files, _ = os.ReadDir(path)
							currentSelection = 0
						} else {
							// Handle file selection
							fmt.Println("Selected file:", files[currentSelection].Name())
						}
					case sdl.K_BACKSPACE:
						// Go back to the parent directory
						parent := filepath.Dir(path)
						if parent != path {
							path = parent
							// Re-read the directory
							files, _ = os.ReadDir(path)
						}
					}
				}
			case *sdl.JoyAxisEvent:
				// Handle joystick axis motion
				if e.Axis == 1 { // Assuming axis 1 is the vertical axis
					if e.Value < 0 {
						// Move selection up
						if currentSelection > 0 {
							currentSelection--
						}
					} else if e.Value > 0 {
						// Move selection down
						if currentSelection < len(files)-1 {
							currentSelection++
						}
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
							running = false // Exit the application
						}
					}
				} else if e.Button == sdl.CONTROLLER_BUTTON_BACK {
					if e.State == sdl.PRESSED {
						selectPressed = true
						selectPressTime = time.Now()
					} else if e.State == sdl.RELEASED {
						selectPressed = false
						if startPressed && time.Since(startPressTime) >= 3*time.Second {
							running = false // Exit the application
						}
					}
				} else if e.State == sdl.PRESSED {
					switch e.Button {
					case sdl.CONTROLLER_BUTTON_A: // Assuming button 0 is the A button or equivalent
						// Enter the selected directory or select the file
						if files[currentSelection].IsDir() {
							// Change the current directory
							path = filepath.Join(path, files[currentSelection].Name())
							// Re-read the directory
							files, _ = os.ReadDir(path)
							currentSelection = 0
						} else {
							selectedFileName := files[currentSelection].Name()

							// Handle file selection
							// Get the parent directory's name
							parentDirName := filepath.Base(path)
							fmt.Println("Parent directory name:", parentDirName)

							// Get the full path of the selected file
							fullFilePath := filepath.Join(path, selectedFileName)
							fmt.Println("Selected file full path:", fullFilePath)

							// Create a command to execute 'nano example.txt'
							var cmd *exec.Cmd
							if selectedFileName == "self.txt" {
								cmd = exec.Command(parentDirName)
							} else {
								cmd = exec.Command(parentDirName, fullFilePath)
							}

							// Connect the command's stdin, stdout, and stderr to those of the parent process
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
					case sdl.CONTROLLER_BUTTON_B: // Assuming button 1 is the B button or equivalent
						// Go back to the parent directory
						parent := filepath.Dir(path)
						if parent != path {
							path = parent
							// Re-read the directory
							files, _ = os.ReadDir(path)
						}
					}
				}
			}
		}
		// Clear the screen
		renderer.Clear()

		// Read directory
		files, err := os.ReadDir(path)
		if err != nil {
			fmt.Println("Could not read directory:", err)
			return
		}

		// Draw file and folder names
		for i, file := range files {
			text := file.Name()
			if file.IsDir() {
				text += "/"
			}

			// Choose color based on selection
			var color sdl.Color
			if i == currentSelection {
				color = yellow
			} else {
				color = white
			}

			surface, err := font.RenderUTF8Blended(text, color)
			if err != nil {
				fmt.Println("Could not render text:", err)
				return
			}
			defer surface.Free()

			texture, err := renderer.CreateTextureFromSurface(surface)
			if err != nil {
				fmt.Println("Could not create texture:", err)
				return
			}
			defer texture.Destroy()

			renderer.Copy(texture, nil, &sdl.Rect{X: 10, Y: 10 + int32(i*fontSize), W: surface.W, H: surface.H})
		}

		// Update the screen
		renderer.Present()
	}
}
