package gogl

import (
	"github.com/go-gl/gl/v4.2-core/gl"
	"log"
	// "github.com/go-gl-legacy/glu"
)

func OpenGLSentinel() func() {
	check := func() {
		e := gl.GetError()
		var msg string
		if e != gl.NO_ERROR {
			switch e {
			case gl.INVALID_ENUM:
				msg = "INVALID_ENUM"
			case gl.INVALID_VALUE:
				msg = "INVALID_VALUE"
			case gl.INVALID_OPERATION:
				msg = "INVALID_OPERATION"
			case gl.STACK_OVERFLOW:
				msg = "STACK_OVERFLOW"
			case gl.STACK_UNDERFLOW:
				msg = "STACK_UNDERFLOW"
			case gl.OUT_OF_MEMORY:
				msg = "OUT_OF_MEMORY"
			case gl.INVALID_FRAMEBUFFER_OPERATION:
				msg = "INVALID_FRAMEBUFFER_OPERATION"
			}
		}
		log.Panic("Encountered GLError: ", e, " = ", msg)
	}
	// check()
	return check
}

