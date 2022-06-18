package math

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)
type Vec2 struct {
	X,Y float32
}

const epsilon float64 = 10e-6

func Rotate(vec *mgl32.Vec2, angleDeg float32, origin mgl32.Vec2) {
	x := vec.X() - origin.X()
	y := vec.Y() - origin.Y()

	cos := float32(math.Cos(toRadians(angleDeg)))

	sin := float32(math.Sin(toRadians(angleDeg)))

	xPrime := (x * cos) - (y * sin)
	yPrime := (x * sin) + (y * cos)

	xPrime += origin.X()
	yPrime += origin.Y()

	vec[0] = xPrime
	vec[1] = yPrime
}

func toRadians(angleDeg float32) float64 {
	return float64((angleDeg * math.Pi) / 180)
}

func CompareFloat32(x, y float32) bool {
	return math.Abs(float64(x-y)) <= epsilon*math.Max(1, math.Max(math.Abs(float64(x)), math.Abs(float64(y))))
}
func CompareVec2(vec1, vec2 mgl32.Vec2) bool {
	return CompareFloat32(vec1[0], vec2[0]) && CompareFloat32(vec1[1], vec2[1])
}

func MulVec(src, vec mgl32.Vec2) mgl32.Vec2 {
	src[0] *= vec[0] 
	src[1] *= vec[1] 
	return src
}
