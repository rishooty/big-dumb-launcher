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
	app.drawFileTree(app.fileTree, 0, 0)
}

func (app *App) drawFileTree(node *FileNode, depth, index int) int {
	text := node.Name
	if node.IsDir {
		text += "/"
	}

	// Choose color based on selection
	var color sdl.Color
	if depth == 0 && index == app.currentSelection {
		color = app.colors.yellow
	} else {
		color = app.colors.white
	}

	surface, err := app.font.RenderUTF8Blended(text, color)
	if err != nil {
		fmt.Println("Could not render text:", err)
		return index
	}
	defer surface.Free()

	texture, err := app.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		fmt.Println("Could not create texture:", err)
		return index
	}
	defer texture.Destroy()

	app.renderer.Copy(texture, nil, &sdl.Rect{X: 10 + int32(depth*20), Y: 10 + int32(index*fontSize), W: surface.W, H: surface.H})

	if node.IsDir {
		for _, child := range node.Children {
			index = app.drawFileTree(child, depth+1, index+1)
		}
	}

	return index
}
