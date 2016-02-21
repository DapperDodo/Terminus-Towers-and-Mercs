package renderer

import (
	"strconv"

	"towmer/api"
	"towmer/ecs"

	"github.com/nsf/termbox-go"
)

func Init() error {
	return termbox.Init()
}

func Size() (int, int) {
	return termbox.Size()
}

func Close() {
	termbox.Close()
}

func PollEvent() api.Key {

	event := termbox.PollEvent()
	if event.Type == termbox.EventKey {
		switch event.Key {
		case termbox.KeyEsc:
			return api.Key_ESC
		case termbox.KeyArrowUp:
			return api.Key_UP
		case termbox.KeyArrowDown:
			return api.Key_DOWN
		case termbox.KeyArrowLeft:
			return api.Key_LEFT
		case termbox.KeyArrowRight:
			return api.Key_RIGHT
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
	}
	return api.Key_UNKNOWN
}

/* convert game coords to viewport coords */
func CalcXY(gx, gy float64) (vx, vy int) {

	// window
	w, h := Size()

	// viewport
	vw := w - 1
	vh := h - 6 - 1

	// cam
	cx, cy := 0.0, 0.0 // offset
	cz := 1.0          // zoom

	// character ratio
	yratio := 2.0

	// viewport offset
	vo := ((float64(vw) - (float64(vh) * yratio)) / 2.0)

	fx := float64(vh) * cz * yratio
	fy := float64(vh) * cz

	x := int((gx * fx) + cx + vo + 0.5)
	y := int((gy * fy) + cy + 2 + 0.5)

	return x, y
}

func Render(pool *ecs.Pool) {

	// window
	w, h := Size()

	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)

	// huds
	for x := 0; x < w; x++ {
		// top hud
		termbox.SetCell(x, 1, '═', termbox.ColorWhite, termbox.ColorBlack)
		// bottom hud
		termbox.SetCell(x, h-4, '═', termbox.ColorWhite, termbox.ColorBlack)
	}

	// game entity layer
	for _, e := range pool.Entities {
		x, y := CalcXY(e.X, e.Y)
		if e.HasAspect(ecs.C_POSITION, ecs.C_TERMINAL) {
			termbox.SetCell(x, y, e.Rune, colorToTerm(e.Color), colorToTerm(e.BgColor))
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

		if e.HasAspect(ecs.C_BASE, ecs.C_TEAM_A, ecs.C_ENERGYSTORE) && len(e.Downstream) == 0 {
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
			x, y := CalcXY(waypoint.X, waypoint.Y)
			termbox.SetCell(x, y, '߉', color, termbox.ColorBlack)
		}
	}
}

func renderEntityHud(e *ecs.Entity, color termbox.Attribute, bold bool) {
	x, y := CalcXY(e.X, e.Y)
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
		renderText(x+1+1, y-1+4, "c:"+strconv.FormatFloat(e.Burden, 'f', 0, 64), color)
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
