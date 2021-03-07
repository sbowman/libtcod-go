package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var Font struct {
	Filename                string
	H, V                    int
	InRows, TCOD, Greyscale bool
}

var FullScreen struct {
	Enabled       bool
	Width, Height int
}

var Renderer string

var rootCmd = &cobra.Command{
	Use:   "sample",
	Short: "Run the tcod package sample app",

	Run: func(cmd *cobra.Command, args []string) {
		Setup()
		Run()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	runtime.LockOSThread()

	rootCmd.Flags().StringVar(&Font.Filename, "font", "data/fonts/dejavu10x10_gs_tc.png", "use a custom font")
	rootCmd.Flags().IntVar(&Font.H, "font.h", 0, "number of characters horizontally in the font sheet")
	rootCmd.Flags().IntVar(&Font.V, "font.v", 0, "number of characters vertically in the font sheet")
	rootCmd.Flags().BoolVar(&Font.InRows, "font.rows", false, "the font layout is in rows instead of columns")
	rootCmd.Flags().BoolVar(&Font.TCOD, "font.tcod", false, "the font uses the TCOD layout instead of ASCII")
	rootCmd.Flags().BoolVar(&Font.Greyscale, "font.greyscale", false, "antialiased font using greyscale bitmap")
	rootCmd.Flags().BoolVar(&FullScreen.Enabled, "fullscreen", false, "start in fullscreen")
	rootCmd.Flags().IntVar(&FullScreen.Width, "fullscreen.width", 0, "force fullscreen width")
	rootCmd.Flags().IntVar(&FullScreen.Height, "fullscreen.height", 0, "force fullscreen height")
	rootCmd.Flags().StringVar(&Renderer, "renderer", "GLSL", "renderer type: GLSL, OPENGL, SDL")
}
