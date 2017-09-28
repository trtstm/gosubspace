package renderer

import (
	"fmt"
	"io/ioutil"
	"strings"

	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/trtstm/gosubspace/log"
)

type ShaderProgram struct {
	id uint32
}

func NewProgram(vsPath, fsPath string) (*ShaderProgram, error) {
	vsBytes, err := ioutil.ReadFile(vsPath)
	if err != nil {
		return nil, err
	}
	fsBytes, err := ioutil.ReadFile(fsPath)
	if err != nil {
		return nil, err
	}

	log.Debugf("Compiling %s.", vsPath)
	vs, err := compileShader(string(vsBytes)+"\x00", gl.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}

	log.Debugf("Compiling %s.", fsPath)
	fs, err := compileShader(string(fsBytes)+"\x00", gl.FRAGMENT_SHADER)
	if err != nil {
		gl.DeleteShader(vs)
		return nil, err
	}

	program := &ShaderProgram{
		id: gl.CreateProgram(),
	}

	gl.AttachShader(program.id, vs)
	gl.AttachShader(program.id, fs)
	log.Debugf("Linking %s, %s.", vsPath, fsPath)
	gl.LinkProgram(program.id)

	var status int32
	gl.GetProgramiv(program.id, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program.id, gl.INFO_LOG_LENGTH, &logLength)

		linkLog := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program.id, logLength, nil, gl.Str(linkLog))

		gl.DeleteShader(vs)
		gl.DeleteShader(fs)
		return nil, fmt.Errorf("linking error: %v", linkLog)
	}

	gl.DeleteShader(vs)
	gl.DeleteShader(fs)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		compileLog := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(compileLog))
		return 0, fmt.Errorf("compilation error: %v", compileLog)
	}

	return shader, nil
}
