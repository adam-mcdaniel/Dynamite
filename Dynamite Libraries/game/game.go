package game

import (
	. "../virtualmachine"
	// "math"
	// "fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var size_x float64
var size_y float64
var title string

func Start() {
	cfg := pixelgl.WindowConfig{
		Title:  title,
		Bounds: pixel.R(0, 0, size_x, size_y),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.Clear(colornames.Skyblue)

	for !win.Closed() {
		win.Update()
	}
}

func MakeGame(vm *VirtualMachine) {
	// fmt.Println("Making game...")
	title = vm.Pop().(string)
	size_x = vm.Pop().(float64)
	size_y = vm.Pop().(float64)
	pixelgl.Run(Start)
}

func InstallLibrary(vm *VirtualMachine) {
	vm.Library["game.make"] = MakeGame
	vm.Push("{ game.make % }")
	vm.Push("make")
	vm.Op_store(false)
}
