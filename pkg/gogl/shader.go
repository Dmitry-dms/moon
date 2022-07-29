package gogl

import (
	"fmt"
	"os"
	"time"

	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Shader struct {
	ProgramId uint32
	path      string
	beingUsed bool
	modTime   time.Time
}

func NewShader(path string) (*Shader, error) {

	id, err := CreateProgram(path)
	if err != nil {
		return nil, err
	}
	modTime, err := getModifyedTime(path)
	if err != nil {
		return nil, err
	}
	result := Shader{
		ProgramId: id,
		path:      path,
		beingUsed: false,
		modTime:   modTime,
	}

	return &result, nil
}

func (s *Shader) Use() {
	if !s.beingUsed {
		useProgram(s.ProgramId)
		s.beingUsed = true
	}
}
func (s *Shader) Detach() {
	useProgram(0)
	s.beingUsed = false
}
func getModifyedTime(path string) (time.Time, error) {
	fileStat, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	return fileStat.ModTime(), nil
}

func (shader *Shader) CheckShaderForChanges() error {
	modTime, err := getModifyedTime(shader.path)
	if err != nil {
		return err
	}
	if !modTime.Equal(shader.modTime) {
		id, err := CreateProgram(shader.path)
		if err != nil {
			fmt.Println(err)
		} else {
			gl.DeleteProgram(uint32(shader.ProgramId))
			shader.ProgramId = id
			shader.modTime = modTime
		}
	}
	return nil
}

func (s *Shader) UploadFloat(name string, f float32) {
	name_cstr := gl.Str(name + "\x00")
	location := gl.GetUniformLocation(s.ProgramId, name_cstr)
	s.Use()
	gl.Uniform1f(location, f)
}
func (s *Shader) UploadTexture(name string, slot int32) {
	name_cstr := gl.Str(name + "\x00")
	// location := gl.GetUniformLocation(s.ProgramId, name_cstr)
	// s.Use()
	gl.Uniform1i(gl.GetUniformLocation(s.ProgramId, name_cstr), slot)
}

func (s *Shader) UploadInt(name string, slot int32) {
	name_cstr := gl.Str(name + "\x00")
	// location := gl.GetUniformLocation(s.ProgramId, name_cstr)
	// s.Use()
	gl.Uniform1i(gl.GetUniformLocation(s.ProgramId, name_cstr), slot)
}

func (s *Shader) UploadVec2(name string, vec []float32) {
	name_cstr := gl.Str(name + "\x00")
	location := gl.GetUniformLocation(s.ProgramId, name_cstr)
	s.Use()
	gl.Uniform2f(location, vec[0], vec[1])
}
func (s *Shader) UploadVec3(name string, vec mgl32.Vec3) {
	name_cstr := gl.Str(name + "\x00")
	location := gl.GetUniformLocation(s.ProgramId, name_cstr)
	s.Use()
	v3 := [3]float32(vec)
	gl.Uniform3fv(location, 1, &v3[0])
}
func (s *Shader) UploadVec4(name string, vec mgl32.Vec4) {
	name_cstr := gl.Str(name + "\x00")
	location := gl.GetUniformLocation(s.ProgramId, name_cstr)
	s.Use()
	v4 := [4]float32(vec)
	gl.Uniform4fv(location, 1, &v4[0])
}
func (s *Shader) UploadMat4(name string, mat mgl32.Mat4) {
	name_cstr := gl.Str(name + "\x00")
	location := gl.GetUniformLocation(s.ProgramId, name_cstr)
	s.Use()
	m4 := [16]float32(mat)
	gl.UniformMatrix4fv(location, 1, false, &m4[0])
}

func (s *Shader) UploadIntArray(name string, array []int32) {
	name_cstr := gl.Str(name + "\x00")
	location := gl.GetUniformLocation(s.ProgramId, name_cstr)
	s.Use()
	gl.Uniform1iv(location, int32(len(array)), &array[0])
}
