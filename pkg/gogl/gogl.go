package gogl

import (
	_ "bufio"
	"errors"
	"fmt"
	"unsafe"

	//	"regexp"

	"time"

	"io/ioutil"

	"strings"

	"github.com/go-gl/gl/v4.2-core/gl"
)

type shaderInfo struct {
	pathVert string
	pathFrag string
	modified time.Time
}

func GetVersion() string {
	return gl.GoStr(gl.GetString(gl.VERSION))
}

func LoadShader(path string, shaderType uint32) (uint32, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	shader := string(file)
	shId, err := CreateShader(shader, shaderType)
	if err != nil {
		return 0, err
	}
	return shId, nil
}

func LoadShaders(path string) (uint32, uint32, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, 0, err
	}
	shaders := string(file)
	spl := strings.Split(shaders, "#type")
	var vertexSource, fragmentSource string

	for i := 1; i < len(spl); i++ {
		tmp := strings.Split(spl[i], "\r\n")
		shaderTypeStr := strings.TrimSpace(tmp[0])
		switch shaderTypeStr {
		case "vertex":
			vertexSource = spl[i][len(tmp[0]):]
		case "fragment":
			fragmentSource = spl[i][len(tmp[0]):]
		}
	}
	vertId, err := CreateShader(vertexSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, 0, err
	}
	fragId, err := CreateShader(fragmentSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, 0, err
	}
	return vertId, fragId, nil
}

func CreateShader(source string, shaderType uint32) (uint32, error) {
	shaderId := gl.CreateShader(shaderType)
	vsource, free := gl.Strs(source, "\x00")
	gl.ShaderSource(shaderId, 1, vsource, nil)
	free()
	gl.CompileShader(shaderId)
	var status int32
	gl.GetShaderiv(shaderId, gl.COMPILE_STATUS, &status) //logging
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shaderId, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength)+1)
		gl.GetShaderInfoLog(shaderId, logLength, nil, gl.Str(log))
		return 0, errors.New(log)
	}
	return shaderId, nil
}

func Str(src string) *uint8 {
	return gl.Str(src + "\x00")
}

func CreateProgram(path string) (uint32, error) {

	vert, frag, err := LoadShaders(path)
	if err != nil {
		return 0, err
	}
	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vert)
	gl.AttachShader(shaderProgram, frag)
	gl.LinkProgram(shaderProgram)
	var status int32
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &status) //logging
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength)+1)
		gl.GetProgramInfoLog(shaderProgram, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to link program: %s \n", log)
	}

	gl.DeleteShader(vert)
	gl.DeleteShader(frag)

	return shaderProgram, nil
}

func GenBindBuffer(target uint32) uint32 {
	var buffer uint32
	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(target, buffer)
	return buffer
}
func GenBindVAO() uint32 {
	var VAO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)
	return VAO
}

func GenEBO() uint32 {
	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	return ebo
}

func BufferData[T float32 | int32 | int | uint](target uint32, data []T, usage uint32) {
	size := int(unsafe.Sizeof(data[0]))
	gl.BufferData(target, size*len(data), gl.Ptr(data), usage)
}

func SetVertexAttribPointer(index uint32, size int32, xtype uint32, stride, offset int) {
	var memSize int = 0
	switch xtype {
	case gl.INT:
		fallthrough
	case gl.FLOAT:
		memSize = 4
	}
	gl.VertexAttribPointer(index, size, xtype, false, int32(stride), gl.PtrOffset(offset*memSize))
	gl.EnableVertexAttribArray(index)
}

func useProgram(progId uint32) {
	gl.UseProgram(progId)
}

func BindVertexArray(vaoId uint32) {
	gl.BindVertexArray(vaoId)
}

func InitGLdebug() {
	var flags int32
	gl.GetIntegerv(gl.CONTEXT_FLAGS, &flags)
	gl.Enable(gl.DEBUG_OUTPUT)
	gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS)
	gl.DebugMessageCallback(glDebug, nil)
	gl.DebugMessageControl(gl.DONT_CARE, gl.DONT_CARE, gl.DONT_CARE, 0, nil, true)

	var sizeTex int32
	gl.GetIntegerv(gl.MAX_COMBINED_TEXTURE_IMAGE_UNITS, &sizeTex)
	fmt.Println(sizeTex)
}

func glDebug(source uint32, gltype uint32, id uint32, severity uint32, length int32, message string, userParam unsafe.Pointer) {
	// ignore non-significant error/warning codes
	if id == 131169 || id == 131185 || id == 131218 || id == 131204 {
		return
	}
	fmt.Printf("Debug message (%d): %s \n", id, message)
	switch source {
	case gl.DEBUG_SOURCE_API:
		fmt.Println("Source: API")
	case gl.DEBUG_SOURCE_WINDOW_SYSTEM:
		fmt.Println("Source: Window System")
	case gl.DEBUG_SOURCE_SHADER_COMPILER:
		fmt.Println("Source: Shader Compiler")
	case gl.DEBUG_SOURCE_THIRD_PARTY:
		fmt.Println("Source: Third Party")
	case gl.DEBUG_SOURCE_APPLICATION:
		fmt.Println("Source: Application")
	case gl.DEBUG_SOURCE_OTHER:
		fmt.Println("Source: Other")
	}

	switch gltype {
	case gl.DEBUG_TYPE_ERROR:
		fmt.Println("Type: Error")
	case gl.DEBUG_TYPE_DEPRECATED_BEHAVIOR:
		fmt.Println("Type: Deprecated Behaviour")
	case gl.DEBUG_TYPE_UNDEFINED_BEHAVIOR:
		fmt.Println("Type: Undefined Behaviour")
	case gl.DEBUG_TYPE_PORTABILITY:
		fmt.Println("Type: Portability")
	case gl.DEBUG_TYPE_PERFORMANCE:
		fmt.Println("Type: Performance")
	case gl.DEBUG_TYPE_MARKER:
		fmt.Println("Type: Marker")
	case gl.DEBUG_TYPE_PUSH_GROUP:
		fmt.Println("Type: Push Groups")
	case gl.DEBUG_TYPE_POP_GROUP:
		fmt.Println("Type: Pop Groups")
	case gl.DEBUG_TYPE_OTHER:
		fmt.Println("Type: Other")
	}
	switch severity {
	case gl.DEBUG_SEVERITY_HIGH:
		fmt.Println("Severity: high")
	case gl.DEBUG_SEVERITY_MEDIUM:
		fmt.Println("Severity: medium")
	case gl.DEBUG_SEVERITY_LOW:
		fmt.Println("Severity: low")
	case gl.DEBUG_SEVERITY_NOTIFICATION:
		fmt.Println("Severity: notification")
	}
}
