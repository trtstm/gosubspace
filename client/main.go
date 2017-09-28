package main

import (
	"flag"
	"os"
	"path"
	"runtime"
	"strconv"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/trtstm/gosubspace/client/renderer"
	"github.com/trtstm/gosubspace/helpers"
	"github.com/trtstm/gosubspace/log"
	"github.com/trtstm/gosubspace/protocol"
	//	_ "net/http/pprof"
)

type ClientSettings struct {
	ContinuumPath string
	DataPath      string
	ZonesPath     string
	ResX          uint
	ResY          uint
	Server        string
	ServerHTTP    string
}

var clientSettings = ClientSettings{
	ResX: 1280,
	ResY: 800,
}

func init() {
	runtime.LockOSThread()

	flag.StringVar(&clientSettings.DataPath, "data", "data", "The path to the data folder.")
	flag.StringVar(&clientSettings.Server, "server", "127.0.0.1", "The server address.")
	flag.StringVar(&clientSettings.ContinuumPath, "continuum_path", "", "The path to the continuum folder(required).")
}

var defaultShader *renderer.ShaderProgram
var serverConnection *ServerConnection

func main() {
	flag.Parse()
	if !flag.Parsed() || clientSettings.ContinuumPath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	clientSettings.ZonesPath = path.Join(clientSettings.DataPath, "zones")
	clientSettings.ServerHTTP = "http://" + clientSettings.Server + ":" + strconv.Itoa(protocol.ServerHTTPPort)

	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	// These should be compatible with our renderer.
	glfw.WindowHint(glfw.ClientAPI, glfw.OpenGLESAPI)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)

	window, err := glfw.CreateWindow(int(clientSettings.ResX), int(clientSettings.ResY), "GoSubspace", nil, nil)
	helpers.AssertNoError(err)

	window.MakeContextCurrent()

	renderer.Init()
	defaultShader, err := renderer.NewProgram(
		path.Join(clientSettings.DataPath, "shaders", "default.vs"),
		path.Join(clientSettings.DataPath, "shaders", "default.fs"),
	)
	helpers.AssertNoError(err)
	_ = defaultShader

	var currentState GameState = &LoadingState{}

	for !window.ShouldClose() {
		renderer.ClearScreen()

		if currentState.Run() {
			prevState := currentState
			currentState = currentState.NextState()
			if currentState == nil {
				panic("No current gamestate.")
			}
			log.Debugf("Gamestate change %s -> %s.", prevState.Name(), currentState.Name())
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
