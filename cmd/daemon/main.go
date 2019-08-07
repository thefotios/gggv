package main

import (
	"bufio"
	"fmt"
	_ "image/png"
	"log"
	"os"
	"runtime"

	"github.com/dianelooney/gvd/pkg/gvdl"

	"github.com/dianelooney/gvd/pkg/daemon"
	"github.com/fsnotify/fsnotify"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

var dmn *daemon.D

func main() {
	dmn = daemon.New()

	// startup, type:
	//
	// add source default "sample2.mp4"
	//
	// into the console

	if err := dmn.Scene.LoadProgram("shaders/vert/default.glsl", "shaders/frag/default.glsl"); err != nil {
		panic(err)
	}

	dmn.Scene.BindBuffers()
	dmn.Scene.TextureInit("default")

	go watchShaders()

	go listenForInput()

	dmn.DrawLoop()
}

func listenForInput() {
	r := bufio.NewReader(os.Stdin)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		err = gvdl.Exec(line, dmn)
		if err != nil {
			fmt.Println("Error executing line:", err)
		}
	}
}

func watchShaders() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Panic(err)
	}
	watcher.Add("shaders/frag/default.glsl")
	watcher.Add("shaders/vert/default.glsl")

	go func() {
		for {
			if (<-watcher.Events).Op != fsnotify.Write {
				continue
			}
			dmn.ReloadShaders()
		}
	}()
}
