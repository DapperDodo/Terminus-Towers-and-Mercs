package renderer

import (
	"strconv"

	"github.com/DapperDodo/Terminus-Towers-and-Mercs/api"
	"github.com/DapperDodo/Terminus-Towers-and-Mercs/ecs"

	"github.com/nsf/termbox-go"
)

func Init() error {
	err := termbox.Init()
	if err != nil {
		return err
	}
	_ = termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	return nil
}

func Size() (int, int) {
	return termbox.Size()
}

func Close() {
	termbox.Close()
}

func PollEvent() api.Key {

	event := termbox.PollEvent()
	switch event.Type {
	case termbox.EventKey:
		switch event.Key {
		case termbox.KeyEsc:
			return api.Key_ESC
		case termbox.KeyArrowUp:
			return api.Key_PANUP
		case termbox.KeyArrowDown:
			return api.Key_PANDOWN
		case termbox.KeyArrowLeft:
			return api.Key_PANLEFT
		case termbox.KeyArrowRight:
			return api.Key_PANRIGHT
		case termbox.KeySpace:
			return api.Key_SPACE
		case termbox.KeyBackspace, termbox.KeyBackspace2:
			return api.Key_BACKSPACE
		case termbox.KeyTab:
			return api.Key_TAB
		case termbox.KeyEnter:
			return api.Key_ENTER
		}
		switch event.Ch {
		case 'a', 'A':
			return api.Key_A
		case 'b', 'B':
			return api.Key_B
		case 'c', 'C':
			return api.Key_C
		case 'h', 'H':
			return api.Key_H
		case 'p', 'P':
			return api.Key_P
		case 'r', 'R':
			return api.Key_R
		case 'w', 'W':
			return api.Key_W
		}
	case termbox.EventMouse:
		switch event.Key {
		case termbox.MouseWheelUp:
			return api.Key_ZOOMIN
		case termbox.MouseWheelDown:
			return api.Key_ZOOMOUT
		}
	}
	return api.Key_UNKNOWN
}

func Render(pool *ecs.Pool) {

	// window
	w, h := Size()

	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)

	// huds
	for x := 0; x < w; x++ {
		termbox.SetCell(x, 1, '═', termbox.ColorWhite, termbox.ColorBlack)
		termbox.SetCell(x, h-4, '═', termbox.ColorWhite, termbox.ColorBlack)
	}

	// game entity layer
	for _, e := range pool.Entities {
		if e.HasAspect(ecs.C_POSITION, ecs.C_TERMINAL) {
			x, y, cull := CalcXY(e.X, e.Y)
			if !cull {
				termbox.SetCell(x, y, e.Rune, colorToTerm(e.Color), colorToTerm(e.BgColor))
			}
		}
	}
	// overlays
	for _, e := range pool.Entities {

		if e.HasAspect(ecs.C_TABBABLE) {

			if e.TabActive {
				renderEntityHud(e, termbox.ColorYellow, false)
				if e.HasAspect(ecs.C_PATH) {
					renderPath(e, termbox.ColorYellow)
				}
			} else {
				renderEntityHud(e, termbox.ColorMagenta, false)
			}
		}

		if e.HasAspect(ecs.C_BASE, ecs.C_TEAM_A, ecs.C_SELECTED) {
			renderEntityHud(e, termbox.ColorWhite, false)
			for line, text := range e.Info {
				lx := 0
				ly := line
				switch line {
				case 3, 4, 5:
					lx = int(0.333 * float64(w))
					ly -= 3
				case 6, 7, 8:
					lx = int(0.667 * float64(w))
					ly -= 6
				}
				renderText(lx, h-3+ly, text, termbox.ColorWhite)
			}
			if e.HasAspect(ecs.C_PATH) {
				renderPath(e, termbox.ColorWhite)
			}
		}

		if e.HasAspect(ecs.C_MAIN, ecs.C_TEAM_A, ecs.C_ENERGYSTORE) {
			renderText(0, 0, "e:"+strconv.FormatFloat(e.Energy, 'f', 0, 64), termbox.ColorWhite)
		}
	}
	termbox.Flush()
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// PRIVATE
///////////////////////////////////////////////////////////////////////////////////////////////////////////////

func renderPath(e *ecs.Entity, color termbox.Attribute) {
	if len(e.Waypoints) > 0 {
		for _, waypoint := range e.Waypoints {
			x, y, cull := CalcXY(waypoint.X, waypoint.Y)
			if !cull {
				termbox.SetCell(x, y, '߉', color, termbox.ColorBlack)
			}
		}
	}
}

func renderEntityHud(e *ecs.Entity, color termbox.Attribute, bold bool) {
	x, y, cull := CalcXY(e.X, e.Y)
	if !cull {
		if bold {
			termbox.SetCell(x-1, y-1, '┏', color, termbox.ColorBlack)
			termbox.SetCell(x+1, y-1, '┓', color, termbox.ColorBlack)
			termbox.SetCell(x-1, y+1, '┗', color, termbox.ColorBlack)
			termbox.SetCell(x+1, y+1, '┛', color, termbox.ColorBlack)
		} else {
			termbox.SetCell(x-1, y-1, '┌', color, termbox.ColorBlack)
			termbox.SetCell(x+1, y-1, '┐', color, termbox.ColorBlack)
			termbox.SetCell(x-1, y+1, '└', color, termbox.ColorBlack)
			termbox.SetCell(x+1, y+1, '┘', color, termbox.ColorBlack)
		}
		if e.HasAspect(ecs.C_RESOURCE) {
			renderText(x+1+1, y-1+3, "r:"+strconv.FormatFloat(e.Resources, 'f', 0, 64), color)
		}
		if e.HasAspect(ecs.C_PAYROLL) {
			renderText(x+1+1, y-1+4, "c:"+strconv.Itoa(len(e.Contracts)), color)
		}
		if e.HasAspect(ecs.C_WAVESTART) {
			renderText(x+1+1, y-1+5, "w:"+strconv.Itoa(len(e.Tickets)), color)
		}
	}
}

func renderText(x, y int, text string, color termbox.Attribute) {
	for idx, char := range text {
		termbox.SetCell(x+idx, y, rune(char), color, termbox.ColorBlack)
	}
}

func colorToTerm(c api.Color) termbox.Attribute {
	switch c {
	case api.Color_BLACK:
		return termbox.ColorBlack
	case api.Color_WHITE:
		return termbox.ColorWhite
	case api.Color_GREEN:
		return termbox.ColorGreen
	case api.Color_RED:
		return termbox.ColorRed
	case api.Color_BLUE:
		return termbox.ColorBlue
	case api.Color_YELLOW:
		return termbox.ColorYellow
	case api.Color_CYAN:
		return termbox.ColorCyan
	case api.Color_MAGENTA:
		return termbox.ColorMagenta
	default:
		return termbox.ColorBlack
	}
}

var (
	cx, cy float64 = 0.0, 0.0 // cam offset
	cz     float64 = 1.0      // cam zoom
)

func CamAction(key api.Key) {

	camspeed := 0.2

	switch key {
	case api.Key_PANUP:
		cy += camspeed
	case api.Key_PANDOWN:
		cy -= camspeed
	case api.Key_PANLEFT:
		cx += camspeed
	case api.Key_PANRIGHT:
		cx -= camspeed
	case api.Key_ZOOMIN:
		cz += 0.1
	case api.Key_ZOOMOUT:
		cz -= 0.1
	}
}

// convert game coords to viewport coords
// TODO: parts of this only need to be calculated when terminal size changes
func CalcXY(gx, gy float64) (vx, vy int, cullable bool) {

	// window
	w, h := Size()

	// viewport
	vw := w - 1
	vh := h - 6 - 1

	// character ratio
	yratio := 2.0

	// viewport offset
	vo := ((float64(vw) - (float64(vh) * yratio)) / 2.0)

	fx := float64(vh) * cz * yratio
	fy := float64(vh) * cz

	x := int((gx * fx) + cx + vo + 0.5)
	y := int((gy * fy) + cy + 2 + 0.5)

	if x < 0 || x > vw || y < 2 || y > (vh+2) {
		cullable = true
	}

	return x, y, cullable
}
