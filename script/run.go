package script

import (
	"github.com/go-vgo/robotgo"
)

func Run() {
	robotgo.ScrollMouse(10, "up")
	robotgo.MouseClick("left", true)
	robotgo.MoveMouseSmooth(100, 200, 1.0, 100.0)
}
