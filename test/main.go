package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
)

func main() {
	g := 5
	th := Third{
		address: "address third",
	}
	sec := Second{
		name: "second string",
		th:   &th,
	}
	o := Object{
		Pos: mgl32.Vec2{1, 2},
		sec: sec,
		age: &g,
	}

	str := reflectJson(reflect.ValueOf(o), reflect.TypeOf(o))
	fmt.Println(str)
	//reflectFromJson(str)
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
		//запятую ставим вначале, т.к. так проще проверить на нулевые указатели, структуры
		if i != 0 && jsonTag != "" {
			jsonBuilder.WriteString(",")
		}
		if jsonTag == "" {
			continue
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
					jsonBuilder.WriteString(fmt.Sprintf("\"%v\"",reflect.Indirect(field)))
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
		case reflect.Int, reflect.Float32, reflect.Int32:
			jsonBuilder.WriteString(fmt.Sprintf("%v", field))
		case reflect.String:
			jsonBuilder.WriteString(fmt.Sprintf("\"%v\"", field))
		default:
			//jsonBuilder.WriteString(fmt.Sprintf("%v", field))

		}
	}
	jsonBuilder.WriteString("}")
	return jsonBuilder.String()
}

func reflectFromJson(data string) {
	var obj Object
	err := json.Unmarshal([]byte(data), &obj)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(obj)
	}
}

type Object struct {
	Pos mgl32.Vec2 `json:"pos"`
	sec Second     `json:"second"`
	age *int       `json:"age"`
}

type Second struct {
	name string `json:"name"`
	th   *Third `json:"third"`
}

type Third struct {
	address string `json:"address"`
}
