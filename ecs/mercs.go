package ecs

type Guild int

const (
	Guild_UNKNOWN Guild = iota
	Guild_RANGERS
	Guild_WARRIORS
)

type Merc int

const (
	Merc_UNKNOWN Merc = iota

	// rangers
	Merc_ARCHER
	Merc_HUNTER

	// warriors
	Merc_BRAWLER
	Merc_GLADIATOR
)

type Contract struct {
	Guild
	Merc
	Party int     // the number of mercs
	Cost  float64 // payment due each wave
}
