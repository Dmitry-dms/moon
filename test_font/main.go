package main

import (
	"fmt"

	"github.com/Dmitry-dms/moon/pkg/ui2/fonts"
)

// "fmt"

// import "github.com/Dmitry-dms/moon/pkg/ui2/fonts"



func main() {

	f := fonts.NewFont("assets/fonts/mono.ttf", 30)
	inf := f.GetCharacter(' ')
	fmt.Println(inf.TexCoords)

}
