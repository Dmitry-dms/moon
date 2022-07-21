package main

import (
	"fmt"

	"github.com/Dmitry-dms/moon/pkg/fonts"
)

// "fmt"

// import "github.com/Dmitry-dms/moon/pkg/ui2/fonts"



func main() {

	f := fonts.NewFont("assets/fonts/Roboto.ttf", 40)
	inf := f.GetCharacter(' ')
	fmt.Println(inf.TexCoords)

}
