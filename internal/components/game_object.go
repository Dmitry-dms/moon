package components

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/go-gl/mathgl/mgl32"
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/inkyblackness/imgui-go/v4"
)

type GameObject struct {
	name          string     `json:"name"`
	Transform     *Transform `json:"transform"`
	LastTransform *Transform
	Spr           SpriteRenderer `json:"sprite_renderer"`
	isDirty       bool
	zIndex        int `json:"z_index"`
	uid           int `json:"uid"`
}
type gameObjExported struct {
	Name      string `json:"name"`
	Transform tranformExported `json:"transform"`
	Spr       spriteRendererExported `json:"sprite_renderer"`
	ZIndex    int `json:"z_index"`
	Uid       int `json:"uid"`
}

func NewGameObject(name string, transform *Transform, zIndex int) *GameObject {
	obj := GameObject{
		name:          name,
		Transform:     transform,
		LastTransform: transform.Copy(),
		zIndex:        zIndex,
		uid:           -1,
	}
	obj.generateId()
	return &obj
}

func (g *GameObject) AddSpriteRenderer(spr *SpriteRenderer) {
	g.Spr = *spr
}

func (g *GameObject) generateId() {
	if g.uid == -1 {
		id := ID_COUNTER
		ID_COUNTER++
		g.uid = id
	}
}
func (g *GameObject) GetUid() int {
	return g.uid
}

func (g *GameObject) Update(dt float32) {
	if g.Transform != g.LastTransform {
		g.Transform.CopyTo(g.LastTransform)
		g.isDirty = true
	}
}

func (g *GameObject) GetZIndex() int {
	return g.zIndex
}

func (g *GameObject) SetColor(color mgl.Vec4) {
	if g.Spr.color != color {
		g.isDirty = true
		g.SetColor(color)
	}
}

func (g *GameObject) SetSprite(sprite *gogl.Sprite) {
	g.Spr.SetSprite(sprite)
	g.isDirty = true
}

func (g *GameObject) AddPosition(tr mgl.Vec2) {
	g.LastTransform = g.Transform.Copy()
	g.Transform.position[0] += tr.X()
	g.Transform.position[1] += tr.Y()
	g.isDirty = true
}
func (g *GameObject) SetPosition(tr mgl.Vec2) {
	g.Transform.SetPosition(tr)
	g.isDirty = true
}

func (g *GameObject) IsDirty() bool {
	return g.isDirty
}

func (g *GameObject) SetClean() {
	g.isDirty = false
}

func (g *GameObject) Imgui() {
	colors := [4]float32{g.Spr.color[0], g.Spr.color[1], g.Spr.color[2], g.Spr.color[3]}
	if imgui.ColorPicker4("Color picker", &colors) {
		g.Spr.SetColor(mgl.Vec4{colors[0], colors[1], colors[2], colors[3]})
		g.isDirty = true
	}

}

func reflObj(g *GameObject) {
	v := reflect.ValueOf(*g)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		switch field.Kind() {
		case reflect.Int:
			var val int32 = int32(g.zIndex)
			if imgui.DragInt("zIndex", &val) {
				g.zIndex = int(val)
				g.isDirty = true
			}
		case reflect.ValueOf(mgl.Vec2{}).Kind():
			var fl [2]float32
			if imgui.DragFloat2("floats", &fl) {

			}
		}
	}

}

func (g *GameObject) MarshalJSON() ([]byte, error) {
	// var transfExp tranformExported
	// var spriteExported *gogl.SpriteExported
	// var sprExported spriteRendererExported
	// var texExported *gogl.TextureExported
	// transfExp = tranformExported{
	// 	Position: g.Transform.position,
	// 	Scale:    g.Transform.scale,
	// }
	// if g.Spr.GetTexture() != nil {
	// 	texExported = &gogl.TextureExported{
	// 		Filepath:  g.Spr.GetTexture().GetFilepath(),
	// 		TextureId: g.Spr.GetTexture().GetId(),
	// 		Width:     g.Spr.GetTexture().GetWidth(),
	// 		Height:    g.Spr.GetTexture().GetHeight(),
	// 	}
	// 	spriteExported = &gogl.SpriteExported{
	// 		Texture:   texExported,
	// 		TexCoords: g.Spr.GetTextureCoords(),
	// 		Width:     g.Spr.sprite.GetWidth(),
	// 		Height:    g.Spr.sprite.GetHeight(),
	// 	}
	// } else {
	// 	spriteExported = &gogl.SpriteExported{
	// 		Texture:   nil,
	// 		TexCoords: g.Spr.GetTextureCoords(),
	// 	}
	// }

	// sprExported = spriteRendererExported{
	// 	Color:  g.Spr.GetColor(),
	// 	Sprite: spriteExported,
	// }

	// rt := gameObjExported{
	// 	Name:      g.name,
	// 	Transform: transfExp,
	// 	Spr:       sprExported,
	// 	ZIndex:    g.zIndex,
	// 	Uid:       g.GetUid(),
	// }
	// return json.Marshal(rt)
	data := reflectJson(reflect.Indirect(reflect.ValueOf(g)), reflect.Indirect(reflect.ValueOf(g)).Type())
	return []byte(data), nil
}
func reflectJson(t reflect.Value, ty reflect.Type) string {
	num := t.NumField()
	jsonBuilder := strings.Builder{}

	jsonBuilder.WriteString("{")
	for i := 0; i < num; i++ {
		field := t.Field(i)
		//если ссылка пуста, пропускаем
		if field.Kind() == reflect.Pointer {
			if field.IsNil() {
				continue
			}
			//если структура пуста, пропускаем
		} else if field.Kind() == reflect.Struct {
			if field.IsZero() {
				continue
			}
		}
		jsonTag := ty.Field(i).Tag.Get("json")
		if jsonTag == "" {
			continue
		}
		//запятую ставим вначале, т.к. так проще проверить на нулевые указатели, структуры
		if i != 0 && jsonTag != "" {
			jsonBuilder.WriteString(",")
		}

		jsonBuilder.WriteString(fmt.Sprintf("\"%s\"", jsonTag))
		jsonBuilder.WriteString(":")
		switch field.Kind() {
		case reflect.Pointer:
			//проверяем, что ссылка указывает на структуру
			if reflect.Indirect(field).Kind() == reflect.Struct {
				jsonBuilder.WriteString(reflectJson(reflect.Indirect(field), reflect.Indirect(field).Type()))
			} else {
				switch reflect.Indirect(field).Kind() {
				case reflect.Int, reflect.Float32, reflect.Int32:
					jsonBuilder.WriteString(fmt.Sprintf("%v", reflect.Indirect(field)))
				case reflect.String:
					jsonBuilder.WriteString(fmt.Sprintf("\"%v\"", reflect.Indirect(field)))
				}
				//jsonBuilder.WriteString(fmt.Sprintf("\"%v\"", reflect.Indirect(field)))
			}
		case reflect.Struct:
			jsonBuilder.WriteString(reflectJson(field, field.Type()))
		case reflect.Array:
			l := field.Len()
			jsonBuilder.WriteString("[")
			for i := 0; i < l; i++ {

				d := field.Index(i)
				jsonBuilder.WriteString(fmt.Sprintf("%v", d))
				if i != l-1 {
					jsonBuilder.WriteString(",")
				}
			}
			jsonBuilder.WriteString("]")
		case reflect.TypeOf([]mgl32.Vec2{}).Kind():
			l := field.Len()
			jsonBuilder.WriteString("[")
			for i := 0; i < l; i++ {
				jsonBuilder.WriteString("[")
				num := field.Index(i).Len()
				for j := 0; j < num; j++ {
					d := field.Index(i).Index(j)
					jsonBuilder.WriteString(fmt.Sprintf("%v", d))
					if j != num-1 {
						jsonBuilder.WriteString(",")
					}
				}
				jsonBuilder.WriteString("]")
				if i != l-1 {
					jsonBuilder.WriteString(",")
				}
			}
			jsonBuilder.WriteString("]")
		case reflect.Int, reflect.Float32, reflect.Int32, reflect.Uint32:
			jsonBuilder.WriteString(fmt.Sprintf("%v", field))
		case reflect.String:
			jsonBuilder.WriteString(fmt.Sprintf("\"%v\"", field))
		default:
			fmt.Println("DEFAULT = ", field)
			jsonBuilder.WriteString(fmt.Sprintf("%v", field))

		}
	}
	jsonBuilder.WriteString("}")
	return jsonBuilder.String()
}
func (g *GameObject) UnmarshalJSON(data []byte) error {
	var tr gameObjExported
	err := json.Unmarshal(data, &tr)
	if err != nil {
		return err
	}
	g.isDirty = true
	g.name = tr.Name
	g.zIndex = tr.ZIndex
	var newTransform Transform
	newTransform.position = tr.Transform.Position
	newTransform.scale = tr.Transform.Scale
	g.Transform = &newTransform
	g.LastTransform = &newTransform
	var spr SpriteRenderer
	spr.color = tr.Spr.Color
	sprite := gogl.DefSprite()

	if tr.Spr.Sprite.Texture != nil {
		tex := gogl.CreateTexture(tr.Spr.Sprite.Texture.Filepath, tr.Spr.Sprite.Texture.TextureId,
			tr.Spr.Sprite.Texture.Width, tr.Spr.Sprite.Texture.Height)
		sprite.SetTexture(tex)
		sprite.SetTexCoords(tr.Spr.Sprite.TexCoords)
		sprite.SetWidth(tr.Spr.Sprite.Width)
		sprite.SetHeight(tr.Spr.Sprite.Height)
	}

	spr.sprite = sprite
	g.Spr = spr
	g.uid = tr.Uid
	return err
	return nil
}

// func (g *GameObject) Start() {
// 	for _, c := range g.components {
// 		c.Start()
// 	}
// }
