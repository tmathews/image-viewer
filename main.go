package main

import (
	"flag"
	"fmt"
	"github.com/tmathews/gfx-sdl"
	"github.com/veandco/go-sdl2/sdl"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func main() {
	flag.Parse()
	picFilename := flag.Arg(0)
	if picFilename == "" {
		panic("No file to open, closing")
	}

	infos, err := ioutil.ReadDir(filepath.Dir(picFilename))
	index := -1
	if err != nil {
		panic(err)
	}
	files := make([]string, 0)
	for _, info := range infos {
		if info.IsDir() {
			continue
		}
		if ext := filepath.Ext(info.Name()); !IsPhotoFile(ext) {
			continue
		}
		files = append(files, filepath.Join(filepath.Dir(picFilename), info.Name()))
	}
	for i, file := range files {
		if file == picFilename {
			index = i
		}
	}

	var window *sdl.Window
	var renderer *sdl.Renderer

	windowSize := sdl.Rect{0, 0, 800, 600}
	window, err = sdl.CreateWindow("SHELL", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		windowSize.W, windowSize.H, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	if err != nil {
		fmt.Printf("Failed to create window: %s\n", err)
		panic(err)
	}
	defer window.Destroy()

	icon, err := gfx.SurfaceFromBufString(iconPNG)
	if err != nil {
		panic(err)
	}
	defer icon.Free()
	window.SetIcon(icon)

	// don't mess with our compositor, thanks.
	sdl.SetHint(sdl.HINT_VIDEO_X11_NET_WM_BYPASS_COMPOSITOR, "0")
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "2")

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Printf("Failed to create renderer: %s\n", err)
		panic(err)
	}
	defer renderer.Destroy()

	gfx.InitText()

	c := gfx.Container{}
	c.Color.R = 255

	i, err := gfx.NewImage(renderer, picFilename)
	if err != nil {
		panic(err)
	}
	i.Mode = gfx.ImContain
	i.Align = gfx.ImCenter
	i.Width = windowSize.W
	i.Height = windowSize.H

	txt, err := gfx.NewTextFromBufString(filepath.Base(picFilename), 18, ubuntuTTF)
	txt.X = 10
	txt.Y = 10
	txt.Color = sdl.Color{255, 0, 0, 255}
	txt.Render()

LOOP:
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				break LOOP
			case *sdl.WindowEvent:
				w, h := window.GetSize()
				renderer.SetViewport(&sdl.Rect{0, 0, w, h})
				c.Width = int32(float32(w) * 0.5)
				c.Height = int32(float32(h) * 0.5)
				i.Width = w
				i.Height = h
			case *sdl.KeyboardEvent:
				ke := event.(*sdl.KeyboardEvent)
				code := ke.Keysym.Sym
				if ke.State == 0 && (code == sdl.K_DOWN || code == sdl.K_UP) {
					if code == sdl.K_UP {
						index--
						if index < 0 {
							index = 0
						}
					} else if code == sdl.K_DOWN {
						index++
						if index >= len(files) {
							index = len(files) - 1
						}
					}
					//fmt.Println(index, files[index])
					i.Free()
					i.Load(renderer, files[index])
					txt.Content = filepath.Base(files[index])
					txt.Render()
				}
			}
		}
		renderer.Clear()
		i.Draw(renderer)
		//c.Draw(renderer)
		txt.Draw(renderer)
		renderer.Present()
		sdl.Delay(1000 / 60)
	}

	i.Free()
	txt.Free()
}

func IsPhotoFile(ext string) bool {
	ext = strings.ToLower(ext)
	return ext == ".jpg" || ext == ".png" || ext == ".jpeg"
}