package primitives

import (
	"math"

	"github.com/Dmitry-dms/moon/pkg/gogl"
	gom "github.com/Dmitry-dms/moon/pkg/math"
	mgl "github.com/go-gl/mathgl/mgl32"
)

func PointOnLine(point mgl.Vec2, line gogl.Line2D) bool {
	dy := line.To().Y() - line.From().Y()
	dx := line.To().X() - line.From().X()
	if dx == 0 {
		return gom.CompareFloat32(point.X(), line.From().X())
	}
	m := dy / dx

	b := line.To().Y() - (m * line.To().X())

	// Check the line equation
	return point.Y() == m*point.X()+b
}

func PointInCircle(point mgl.Vec2, circle Circle) bool {
	circleCenter := circle.GetCenter()
	centerToPoint := point.Sub(circleCenter)

	return centerToPoint.LenSqr() <= circle.GetRadius()*circle.GetRadius()
}

func PointInAABB(point mgl.Vec2, box AABB) bool {
	min := box.GetMin()
	max := box.GetMax()

	return point.X() <= max.X() && min.X() <= point.X() &&
		point.Y() <= max.Y() && min.Y() <= point.Y()
}

func PointInBox2D(point mgl.Vec2, box *Box2D) bool {
	// Translate the point into local space
	pointLocalBoxSpace := point
	gom.Rotate(&pointLocalBoxSpace, box.GetRigidbody().GetRotation(),
		box.GetRigidbody().GetPosition())

	min := box.GetMin()
	max := box.GetMax()

	return pointLocalBoxSpace.X() <= max.X() && min.X() <= pointLocalBoxSpace.X() &&
		pointLocalBoxSpace.Y() <= max.Y() && min.Y() <= pointLocalBoxSpace.Y()
}

// Line and Circle
func LineAndCircle(line gogl.Line2D, circle Circle) bool {
	if PointInCircle(line.From(), circle) || PointInCircle(line.To(), circle) {
		return true
	}

	ab := line.To().Sub(line.From())

	// Project point (circle position) onto ab (line segment)
	// parameterized position t
	circleCenter := circle.GetCenter()
	centerToLineStart := circleCenter.Sub(line.From())
	t := centerToLineStart.Dot(ab) / ab.Dot(ab)

	if t < 0.0 || t > 1.0 {
		return false
	}

	// Find the closest point to the line segment
	closestPoint := line.From().Add(ab.Mul(t))

	return PointInCircle(closestPoint, circle)
}

func LineAndAABB(line gogl.Line2D, box AABB) bool {
	if PointInAABB(line.From(), box) || PointInAABB(line.To(), box) {
		return true
	}

	unitVector := line.To().Sub(line.From())
	unitVector.Normalize()
	if unitVector.X() != 0 {
		unitVector[0] = 1
	} else {
		unitVector[0] = 0
	}
	if unitVector.Y() != 0 {
		unitVector[1] = 1
	} else {
		unitVector[1] = 0
	}

	min := box.GetMin()
	min = min.Sub(line.From())
	min = gom.MulVec(min, unitVector)
	max := box.GetMax()
	max = max.Sub(line.From())
	max = gom.MulVec(max, unitVector)

	tmin := float32(math.Max(math.Min(float64(min.X()), float64(max.X())), math.Min(float64(min.Y()), float64(max.Y()))))
	tmax := float32(math.Min(math.Max(float64(min.X()), float64(max.X())), math.Max(float64(min.Y()), float64(max.Y()))))
	if tmax < 0 || tmin > tmax {
		return false
	}
	var t float32
	if tmin < 0 {
		t = tmax
	} else {
		t = tmin
	}

	return t > 0 && t*t < line.LenSqr()
}

func LineAndBox2D(line gogl.Line2D, box Box2D) bool {
	theta := -box.GetRigidbody().GetRotation()
	center := box.GetRigidbody().GetPosition()
	localStart := line.From()
	localEnd := line.To()
	gom.Rotate(&localStart, theta, center)
	gom.Rotate(&localEnd, theta, center)

	localLine := gogl.NewLine2Df(localStart, localEnd)
	aabb := NewAABB(box.GetMin(), box.GetMax())

	return LineAndAABB(*localLine, *aabb)
}

func RaycastCircle(circle Circle, ray Ray2D, result *RaycastResult) bool {
	reset(result)
	originToCircle := circle.GetCenter().Sub(ray.GetOrigin())
	radiusSquared := circle.GetRadius() * circle.GetRadius()
	originToCircleLengthSquared := originToCircle.LenSqr()

	// Project the vector from the ray origin onto the direction of the ray
	a := originToCircle.Dot(ray.GetDirection())
	bSq := originToCircleLengthSquared - (a * a)
	if radiusSquared-bSq < 0 {
		return false
	}

	f := float32(math.Sqrt(float64(radiusSquared - bSq)))
	var t float32 = 0
	if originToCircleLengthSquared < radiusSquared {
		// Ray starts inside the circle
		t = a + f
	} else {
		t = a - f
	}

	if result != nil {
		point := ray.GetOrigin().Add(
			ray.GetDirection().Mul(t))
		normal := point.Sub(circle.GetCenter())
		normal.Normalize()

		result.Init(point, normal, t, true)
	}
	return true
}

func RaycastAABB(box AABB, ray Ray2D, result *RaycastResult) bool {
	reset(result)

	unitVector := ray.GetDirection()
	unitVector.Normalize()
	if unitVector.X() != 0 {
		unitVector[0] = 1 / unitVector[0]
	} else {
		unitVector[0] = 0
	}
	if unitVector.Y() != 0 {
		unitVector[1] = 1 / unitVector[1]
	} else {
		unitVector[1] = 0
	}

	min := box.GetMin()
	min = min.Sub(ray.GetOrigin())
	min = gom.MulVec(min, unitVector)
	max := box.GetMax()
	max = max.Sub(ray.GetOrigin())
	max = gom.MulVec(max, unitVector)

	tmin := float32(math.Max(math.Min(float64(min.X()), float64(max.X())), math.Min(float64(min.Y()), float64(max.Y()))))
	tmax := float32(math.Min(math.Max(float64(min.X()), float64(max.X())), math.Max(float64(min.Y()), float64(max.Y()))))
	if tmax < 0 || tmin > tmax {
		return false
	}
	var t float32
	if tmin < 0 {
		t = tmax
	} else {
		t = tmin
	}
	hit := t > 0
	if !hit {
		return false
	}

	if result != nil {
		point := ray.GetOrigin().Add(
			ray.GetDirection().Mul(t))
		normal := ray.GetOrigin().Sub(point)
		normal.Normalize()

		result.Init(point, normal, t, true)
	}

	return true
}

func RaycastBox2D(box Box2D, ray Ray2D, result *RaycastResult) bool {
	reset(result)

	size := box.GetHalfSize()
	xAxis := mgl.Vec2{1, 0}
	yAxis := mgl.Vec2{0, 1}
	gom.Rotate(&xAxis, -box.GetRigidbody().GetRotation(), mgl.Vec2{0, 0})
	gom.Rotate(&yAxis, -box.GetRigidbody().GetRotation(), mgl.Vec2{0, 0})

	p := box.GetRigidbody().GetPosition().Sub(ray.GetOrigin())
	// Project the direction of the ray onto each axis of the box
	f := mgl.Vec2{xAxis.Dot(ray.GetDirection()), yAxis.Dot(ray.GetDirection())}

	// Next, project p onto every axis of the box
	e := mgl.Vec2{xAxis.Dot(p), yAxis.Dot(p)}

	tArr := []float32{0, 0, 0, 0}
	for i := 0; i < 2; i++ {
		if gom.CompareFloat32(f[i], 0) {
			// If the ray is parallel to the current axis, and the origin of the
			// ray is not inside, we have no hit
			if -e[i]-size[i] > 0 || -e[i]+size[i] < 0 {
				return false
			}
			f[i] = 0.00001 // Set it to small value, to avoid divide by zero
		}
		tArr[i*2+0] = (e[i] + size[i]) / f[i] // tmax for this axis
		tArr[i*2+1] = (e[i] - size[i]) / f[i] // tmin for this axis
	}

	tmin := float32(math.Max(math.Min(float64(tArr[0]), float64(tArr[1])), math.Min(float64(tArr[2]), float64(tArr[3]))))
	tmax := float32(math.Min(math.Max(float64(tArr[0]), float64(tArr[1])), math.Max(float64(tArr[2]), float64(tArr[3]))))

	var t float32
	if tmin < 0 {
		t = tmax
	} else {
		t = tmin
	}
	hit := t > 0
	if !hit {
		return false
	}

	if result != nil {
		point := ray.GetOrigin().Add(
			ray.GetDirection().Mul(t))
		normal := ray.GetOrigin().Sub(point)
		normal.Normalize()

		result.Init(point, normal, t, true)
	}

	return true
}
