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
	sheet := sprite_packer.NewSpriteSheet(128, "test")
	f, d := fonts.NewFont("C:/Windows/Fonts/arial.ttf", 28, 127, 33, 126)
	f2, d2 := fonts.NewFont("C:/Windows/Fonts/times.ttf", 36, 127, 33, 126)
	//CreateImage("fonts-standalone.png", d)
	ConvertFontToAtlas(f2, sheet, d2)
	ConvertFontToAtlas(f, sheet, d)
	CreateImage("fonts-in-atlas.png", sheet.Image())
	//_, err := sprite_packer.GetSpriteSheetFromFile("atlas.json")
	//if err != nil {
	//	panic(err)
	//}
}

func ConvertFontToAtlas(f *fonts.Font, sheet *sprite_packer.SpriteSheet, srcImage *image.RGBA) {

	sort.Slice(f.CharSlice, func(i, j int) bool {
		return f.CharSlice[i].Height > f.CharSlice[j].Height
	})

	sheet.BeginGroup(f.Filepath, func() map[string]*sprite_packer.SpriteInfo {
		spriteInfo := make(map[string]*sprite_packer.SpriteInfo, len(f.CharSlice))
		//gr := sprite_packer.NewGroup(f.Filepath)
		for _, info := range f.CharSlice {
			if info.Rune == ' ' || info.Rune == '\u00a0' {
				continue
			}
			ret := srcImage.SubImage(image.Rect(info.SrcX, info.SrcY, info.SrcX+info.Width, info.SrcY-info.Height)).(*image.RGBA)
			pixels := sheet.GetData(ret)
			if len(pixels) == 0 {
				continue
			}
			spriteInfo[string(info.Rune)] = sheet.AddToSheet(string(info.Rune), pixels)
		}
		return spriteInfo
	})

	//rr, _ := sheet.GetGroup(f.Filepath)
	//
	//for _, info := range rr.Contents {
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
