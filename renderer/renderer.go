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

func Render(pool *ecs.Pool, w, h float64) {

	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)

	for x := 0; x < int(w+0.5); x++ {
		termbox.SetCell(x, int(h+0.5), '═', termbox.ColorWhite, termbox.ColorBlack)
	}

	// game entity layer
	for _, e := range pool.Entities {
		if e.HasAspect(ecs.C_POSITION, ecs.C_TERMINAL) {
			termbox.SetCell(int(e.X+0.5), int(e.Y+0.5), e.Rune, colorToTerm(e.Color), colorToTerm(e.BgColor))
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
					lx = int(0.333 * w)
					ly -= 3
				case 6, 7, 8:
					lx = int(0.667 * w)
					ly -= 6
				}
				renderText(lx, int(h+0.5)+ly+1, text, termbox.ColorWhite)
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
			termbox.SetCell(int(waypoint.X+0.5), int(waypoint.Y+0.5), '߉', color, termbox.ColorBlack)
		}
	}
}

func renderEntityHud(e *ecs.Entity, color termbox.Attribute, bold bool) {
	if bold {
		termbox.SetCell(int(e.X+0.5-1), int(e.Y+0.5-1), '┏', color, termbox.ColorBlack)
		termbox.SetCell(int(e.X+0.5+1), int(e.Y+0.5-1), '┓', color, termbox.ColorBlack)
		termbox.SetCell(int(e.X+0.5-1), int(e.Y+0.5+1), '┗', color, termbox.ColorBlack)
		termbox.SetCell(int(e.X+0.5+1), int(e.Y+0.5+1), '┛', color, termbox.ColorBlack)
	} else {
		termbox.SetCell(int(e.X+0.5-1), int(e.Y+0.5-1), '┌', color, termbox.ColorBlack)
		termbox.SetCell(int(e.X+0.5+1), int(e.Y+0.5-1), '┐', color, termbox.ColorBlack)
		termbox.SetCell(int(e.X+0.5-1), int(e.Y+0.5+1), '└', color, termbox.ColorBlack)
		termbox.SetCell(int(e.X+0.5+1), int(e.Y+0.5+1), '┘', color, termbox.ColorBlack)
	}
	if e.HasAspect(ecs.C_RESOURCE) {
		renderText(int(e.X+0.5+1)+1, int(e.Y+0.5-1)+3, "r:"+strconv.FormatFloat(e.Resources, 'f', 0, 64), color)
	}
	if e.HasAspect(ecs.C_PAYROLL) {
		renderText(int(e.X+0.5+1)+1, int(e.Y+0.5-1)+4, "c:"+strconv.FormatFloat(e.Burden, 'f', 0, 64), color)
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
