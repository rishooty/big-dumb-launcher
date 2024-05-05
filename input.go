package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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
					app.navigateToDirectory()
				case sdl.K_SPACE:
					app.launchFileOrExecuteCommand()
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
					app.navigateToDirectory()
				case sdl.CONTROLLER_BUTTON_X:
					app.launchFileOrExecuteCommand()
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
	currentNode := app.getCurrentNode()
	if app.currentSelection < len(currentNode.Children)-1 {
		app.currentSelection++
	}
}

func (app *App) moveSelectBack() {
	currentNode := app.getCurrentNode()
	if currentNode.Path != app.fileTree.Path {
		app.currentSelection = app.getParentIndex(currentNode)
	}
}

func (app *App) navigateToDirectory() {
	currentNode := app.getCurrentNode()
	if currentNode.IsDir {
		app.currentSelection = 0
	}
}

func (app *App) launchFileOrExecuteCommand() {
	currentNode := app.getCurrentNode()
	if !currentNode.IsDir {
		app.executeCommand(currentNode.Command)
	} else {
		app.executeCommand(currentNode.Command) //trim off the second segment
	}
}

func (app *App) buildCommandArgs(filePath string) []string {
	pathSegments := strings.Split(filePath, "/")[1:] // Skip the 'testpath'
	var args []string

	if len(pathSegments) > 0 {
		args = app.handleCommandSegment(pathSegments[0])
	}

	for segIndex, segment := range pathSegments[1:] {
		args = app.handleSegment(segment, segIndex, args)
	}

	args = append(args, filePath)
	return args
}

func (app *App) handleCommandSegment(segment string) []string {
	var args []string
	if dotIndex := strings.Index(segment, "."); dotIndex != -1 {
		command := segment[:dotIndex]
		args = append(args, command)
		flag := "--" + segment[dotIndex+1:]
		args = append(args, flag)
	} else {
		args = append(args, segment)
	}
	return args
}

func (app *App) handleSegment(segment string, segIndex int, args []string) []string {
	if dotIndex := strings.Index(segment, "."); dotIndex != -1 {
		flag := ""
		value := ""
		if segIndex != 0 {
			flag = "--" + segment[:dotIndex]
		} else {
			flag = segment[:dotIndex]
			value = "--" + segment[dotIndex+1:]
		}
		args = append(args, flag, value) // else if a matching file was found (above), swap in its full filepath
	} else {
		args = append(args, segment)
	}
	return args
}

func (app *App) getCurrentNode() *FileNode {
	currentNode := app.fileTree
	for _, index := range app.getSelectionIndexes() {
		currentNode = currentNode.Children[index]
	}
	return currentNode
}

func (app *App) getSelectionIndexes() []int {
	indexes := []int{}
	node := app.fileTree
	index := app.currentSelection

	for {
		if len(node.Children) == 0 {
			break
		}

		indexes = append(indexes, index)
		node = node.Children[index]
		index = 0
	}

	return indexes
}

func (app *App) getParentIndex(node *FileNode) int {
	parentNode := app.fileTree
	for _, index := range app.getSelectionIndexes()[:len(app.getSelectionIndexes())-1] {
		parentNode = parentNode.Children[index]
	}

	for i, child := range parentNode.Children {
		if child == node {
			return i
		}
	}

	return 0
}

func (app *App) executeCommand(args []string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Executing command:", cmd.String())
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error executing command:", err)
		app.retryWithShortFormOptions(args)
	}
}

func (app *App) retryWithShortFormOptions(args []string) {
	for i, arg := range args[1:] {
		if strings.HasPrefix(arg, "--") {
			args[i+1] = "-" + arg[2:] // Convert to short-form
		}
	}
	fmt.Println("Retrying with short-form options...")
	cmd := exec.Command(args[0], args[1:]...)
	fmt.Println("Executing command:", cmd.String())
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error executing command with short-form options:", err)
	}
}
