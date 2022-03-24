package main

import (
	"fmt"
	"reflect"

	"github.com/go-gl/mathgl/mgl32"
)

func main() {
	o := Object{
		Pos: mgl32.Vec2{1,2},
		color: mgl32.Vec4{1,1,1,1},
	}

	val := reflect.ValueOf(&o)

	p := reflect.ValueOf(mgl32.Vec2{0,0})
	val.Elem().Field(0).Set(p)
	

	fmt.Println(o)
}

type Object struct {
	Pos   mgl32.Vec2
	color mgl32.Vec4
}
