package renderer

import (
	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/trtstm/gosubspace/helpers"
	"github.com/trtstm/gosubspace/log"
)

func Init() {
	log.Info("Initializing rendering engine.")
	helpers.AssertNoError(gl.Init())

	log.Infof("GL_VENDOR: %s", gl.GoStr(gl.GetString(gl.VENDOR)))
	log.Infof("GL_RENDERER: %s", gl.GoStr(gl.GetString(gl.RENDERER)))
	log.Infof("GL_VERSION: %s", gl.GoStr(gl.GetString(gl.VERSION)))
	log.Infof("GL_SHADING_LANGUAGE_VERSION: %s", gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)))
	log.Info("Initializing rendering engine complete.")
}

func ClearScreen() {
	gl.ClearColor(0, 0, 0, 0)
}
