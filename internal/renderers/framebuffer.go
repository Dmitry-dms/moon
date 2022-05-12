package renderers

import (
	"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/pkg/errors"
)

type Framebuffer struct {
	width, height int32
	fboId         uint32
	texture       *gogl.Texture
}

func NewFramebuffer(width, height int32) (*Framebuffer, error) {

	frameBuffer := Framebuffer{
		width:  width,
		height: height,
	}
	var fboid uint32
	//generate framebuufer
	gl.GenFramebuffers(1, &fboid)
	frameBuffer.fboId = fboid
	gl.BindFramebuffer(gl.FRAMEBUFFER, fboid)

	//создание текстуры, в которую рендерить данные и привязка к фрэйм буферу
	frameBuffer.texture = gogl.NewTextureFramebuffer(width, height)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0,
		gl.TEXTURE_2D, frameBuffer.texture.GetId(), 0)
	//create render buffer to store depth info
	var rboId uint32
	gl.GenRenderbuffers(1, &rboId)
	gl.BindRenderbuffer(gl.RENDERBUFFER, rboId)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT32, width, height)
	//attach to framebuffer
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, rboId)

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		return nil, errors.New("Attachment on framebuffer has failed")
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	return &frameBuffer, nil
}
func (f *Framebuffer) GetTextureId() uint32 {
	return f.texture.GetId()
}
func (f *Framebuffer) GetFboId() uint32 {
	return f.fboId
}

func (f *Framebuffer) Bind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, f.fboId)
}
func (f *Framebuffer) Unbind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}