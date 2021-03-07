package main

import (
	"github.com/sbowman/tcod/tcod"
	"github.com/sbowman/tcod/tcod/keys"
)

const (
	Title  = "TCOD Sample"
	Width  = 80
	Height = 50

	SampleScreenWidth = 46
	SampleScreenHeight = 20
)

var Root *tcod.RootConsole

func Setup() {
	var renderer tcod.Renderer
	switch Renderer {
	case "OPENGL":
		renderer = tcod.OpenGL
	case "OPENGL2":
		renderer = tcod.OpenGL2
	case "SDL":
		renderer = tcod.SDL
	case "GLSL":
		renderer = tcod.GLSL
	default:
		renderer = tcod.SDL2
	}

	var fontFlags int
	if Font.InRows {
		fontFlags |= tcod.FontLayoutASCIIInRow
	}
	if Font.Greyscale {
		fontFlags |= tcod.FontTypeGreyscale
	}
	if Font.TCOD {
		fontFlags |= tcod.FontLayoutTCOD
	}

	if FullScreen.Width > 0 {
		tcod.SysForceFullscreenResolution(FullScreen.Width, FullScreen.Height)
	}

	if Font.Filename != "" {
		Root = tcod.NewRootConsoleWithFont(Width, Height, Title, FullScreen.Enabled, Font.Filename, fontFlags, Font.H, Font.V, renderer)
	} else {
		Root = tcod.NewRootConsole(Width, Height, Title, FullScreen.Enabled, renderer)
	}
}

func Run() {
	var creditsEnd bool

	// offscreen := tcod.NewConsole(SampleScreenWidth, SampleScreenHeight)

	for !Root.IsWindowClosed() {
		if !creditsEnd {
			creditsEnd = Root.RenderCredits(50, 40, false)
		}

		Root.SetDefaultBackground(tcod.Black)
		Root.SetDefaultForeground(tcod.Green)

		Root.Print(10, 10, "Hello World")
		Root.Flush()

		var key tcod.Key

		tcod.SysCheckForEvent(tcod.EventKeyPress | tcod.EventMouse, &key, nil)

		// Alt-Enter: toggle fullscreen
		if key.VK == keys.Enter && (key.LAlt || key.RAlt) {
			Root.SetFullscreen(!Root.IsFullscreen())
		}
	}
}
