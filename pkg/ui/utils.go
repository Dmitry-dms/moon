package ui

type stack[T any] struct {
	Push   func(T)
	Pop    func() T
	GetTop func() T
	Length func() int
}

func Stack[T any]() stack[T] {
	slice := make([]T, 0)
	return stack[T]{
		Push: func(i T) {
			slice = append(slice, i)
		},
		Pop: func() T {
			res := slice[len(slice)-1]
			slice = slice[:len(slice)-1]
			return res
		},
		GetTop: func() T {
			return slice[len(slice)-1]
		},
		Length: func() int {
			return len(slice)
		},
	}
}

type Vec2 struct {
	X, Y float32
}

func (v Vec2) Sub(x Vec2) Vec2 {
	return Vec2{v.X - x.X, v.Y - x.Y}
}

func (v Vec2) LengthSqr() float32 {
	return v.X*v.X + v.Y*v.Y
}

type Vec4 struct {
	X, Y, W, H float32
}

func Max[T int | float32](x1, x2 T) T {
	if x1 > x2 {
		return x1
	} else {
		return x2
	}
}

type Rect struct {
	Min, Max Vec2
}

func NewRect(x, y, w, h float32) Rect {
	min := Vec2{x, y}
	max := Vec2{x + w, y + h}
	return Rect{min, max}
}

func (r Rect) Contains(p Vec2) bool {
	return p.X >= r.Min.X && p.Y >= r.Min.Y && p.X < r.Max.X && p.Y < r.Max.Y
}

func PointInRect(point Vec2, box Rect) bool {
	min := box.Min
	max := box.Max

	return point.X <= max.X && min.X <= point.X &&
		point.Y <= max.Y && min.Y <= point.Y
}
