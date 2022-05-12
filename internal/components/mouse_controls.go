package components

import (


	"github.com/Dmitry-dms/moon/internal/listeners"
	"github.com/Dmitry-dms/moon/internal/utils"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type MouseControls struct {
	holdingObject *GameObject
	addToScene    func(g *GameObject)
}

func NewMouseControls(addToScene func(g *GameObject)) *MouseControls {
	return &MouseControls{
		addToScene: addToScene,
	}
}

func (c *MouseControls) PickupObject(obj *GameObject) {
	c.holdingObject = obj
	c.addToScene(obj)
}

func (c *MouseControls) place() {
	c.holdingObject = nil
}

func (c *MouseControls) Update(dt float32) {
	if c.holdingObject != nil {

		c.holdingObject.SetPosition(mgl32.Vec2{float32(listeners.GetOrthoX()) - 16, float32(listeners.GetOrthoY()) - 16})
		//snap to grid
		x := int((c.holdingObject.Transform.position.X() / utils.GRID_WIDTH)) * int(utils.GRID_WIDTH)
		y := int((c.holdingObject.Transform.position.Y()  / utils.GRID_HEIGHT)) * int(utils.GRID_HEIGHT)
		c.holdingObject.SetPosition(mgl32.Vec2{float32(x), float32(y)})
	//	fmt.Println(x,y)
		//
		if listeners.MouseButtonDown(glfw.MouseButtonLeft) {

			c.place()
		}
	}
}
