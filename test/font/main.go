package main

import (
	"github.com/Dmitry-dms/moon/pkg/fonts"
	"github.com/Dmitry-dms/moon/pkg/sprite_packer"
	"image"
	"image/png"
	"os"
	"sort"
)

func main() {
	//sheet := sprite_packer.NewSpriteSheet(128, "test")
	//f, d := fonts.NewFont("C:/Windows/Fonts/arial.ttf", 18)
	//CreateImage("fonts-standalone.png", d)
	//ConvertFontToAtlas(f, sheet, d)
	//CreateImage("fonts-in-atlas.png", sheet.Image())
	_, err := sprite_packer.GetSpriteSheetFromFile("atlas.json")
	if err != nil {
		panic(err)
	}
}

func ConvertFontToAtlas(f *fonts.Font, sheet *sprite_packer.SpriteSheet, srcImage *image.RGBA) {

	sort.Slice(f.CharSlice, func(i, j int) bool {
		return f.CharSlice[i].Heigth > f.CharSlice[j].Heigth
	})

	sheet.BeginGroup(f.Filepath, func() []*sprite_packer.SpriteInfo {
		spriteInfo := make([]*sprite_packer.SpriteInfo, len(f.CharSlice))
		for i, info := range f.CharSlice {
			if info.Rune == ' ' || info.Rune == '\u00a0' {
				continue
			}
			ret := srcImage.SubImage(image.Rect(info.SrcX, info.SrcY, info.SrcX+info.Width, info.SrcY-info.Heigth)).(*image.RGBA)
			pixels := sheet.GetData(ret)
			if len(pixels) == 0 {
				continue
			}
			spriteInfo[i] = sheet.AddToSheet(string(info.Rune), pixels)
		}
		return spriteInfo
	})

	//rr, _ := sheet.GetGroup(f.Filepath)

	//for _, info := range rr {
	//	if info != nil {
	//		ll := []rune(info.Id)
	//		char := f.GetCharacter(ll[0])
	//		char.TexCoords = [2]math.Vec2{{info.TextCoords[0], info.TextCoords[1]},
	//			{info.TextCoords[2], info.TextCoords[3]}}
	//	}
	//}
}

func CreateImage(filename string, img image.Image) {
	pngFile, _ := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0664)

	defer pngFile.Close()

	encoder := png.Encoder{
		CompressionLevel: png.BestCompression,
	}
	encoder.Encode(pngFile, img)
}
