package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type Colors struct {
	white  sdl.Color
	yellow sdl.Color
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
