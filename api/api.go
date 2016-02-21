package api

type Key int

const (
	Key_UNKNOWN Key = iota
	Key_ESC
	Key_UP
	Key_DOWN
	Key_LEFT
	Key_RIGHT
	Key_SPACE
	Key_BACKSPACE
	Key_A
	Key_B
	Key_C
	Key_H
	Key_R
	Key_P
	Key_W
	Key_TAB
	Key_ENTER
)

type Color int

const (
	Color_UNKNOWN Color = iota
	Color_BLACK
	Color_WHITE
	Color_GREEN
	Color_RED
	Color_BLUE
	Color_YELLOW
	Color_CYAN
	Color_MAGENTA
)

var (
	InfoBaseMainMenu     = []string{"B: Build Base", "P: Set Path", "C: Contract Merc", "SPACE: Deselect"}
	InfoBaseMainMenuNoB  = []string{"(no patches in range)", "SPACE: Deselect"}
	InfoBaseMainMenuNoP  = []string{"(no waypoints)", "SPACE: Deselect"}
	InfoBaseMainMenuSetP = []string{"(path already set)", "SPACE: Deselect"}
	InfoBaseBuildSelect  = []string{"TAB: select base location", "ENTER: build"}
	InfoPathSelect       = []string{"TAB: select waypoint", "ENTER: next"}
	InfoContractGuilds   = []string{"R: Rangers", "W: Warriors"}
	InfoContractRangers  = []string{"A: Archer", "H: Hunter"}
	InfoContractWarriors = []string{"B: Brawler", "G: Gladiator"}
	InfoContractSign     = []string{"ENTER: sign", "BACKSPACE: cancel", "", "", ""}
)
