package main

import (
	"fmt"
	"os"
	"path/filepath"

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
	sdl.Init(sdl.INIT_VIDEO)
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

	// Main loop
	running := true
	for running {
		// Handle events
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
