package scenes

// import (
// 	"fmt"

// 	"github.com/Dmitry-dms/moon/internal/components"
// 	"github.com/Dmitry-dms/moon/internal/listeners"
// 	"github.com/go-gl/glfw/v3.3/glfw"
// )

// type LevelScene struct {
// 	showDemoWindow bool

// 	changeSceneCallback func(scene int)

// 	changingScene     bool
// 	timeToChangeScene float32
// }

// func NewLevelScene(changeSceneCallback func(scene int)) *LevelScene {
// 	edtrScene := LevelScene{
// 		showDemoWindow:      true,
// 		timeToChangeScene:   2,
// 		changeSceneCallback: changeSceneCallback,
// 	}

// 	return &edtrScene
// }

// func (e *LevelScene) Init() {
// 	fmt.Println("level scene init")
// }
// func (e *LevelScene) Start() {

// }

// func (e *LevelScene) AddGameObjectToScene(obj *components.GameObject) {
	
// }
// func (e *LevelScene) Destroy() {

// }
// func (e *LevelScene) Update(dt float32) {
// 	if listeners.IsKeyPressed(glfw.KeySpace) && !e.changingScene {
// 		e.changingScene = true
// 	}
// 	if e.changingScene && e.timeToChangeScene > 0 {
// 		e.timeToChangeScene -= dt
// 	} else if e.changingScene {
// 		e.changeSceneCallback(0)
// 	}

// }