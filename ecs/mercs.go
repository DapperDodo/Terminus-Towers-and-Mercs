package ecs

import (
	"sort"
	"time"
)

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
	Cost      float64 // payment due each wave
	Base      *Entity
	Seniority time.Time
}

type bySeniorityAsc []*Contract

func (s bySeniorityAsc) Len() int {
	return len(s)
}

func (s bySeniorityAsc) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s bySeniorityAsc) Less(i, j int) bool {
	return s[i].Seniority.Before(s[j].Seniority)
}

type Ticket struct {
	Guild
	Merc
	Seniority time.Time
	WaitForIt float64
}

type bySeniorityDesc []*Ticket

func (s bySeniorityDesc) Len() int {
	return len(s)
}

func (s bySeniorityDesc) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s bySeniorityDesc) Less(i, j int) bool {
	return s[j].Seniority.Before(s[i].Seniority)
}

func WaveStart(pool *Pool) {

	// get the terminus holding the necessary funds
	terminii := pool.ListAspect(C_MAIN, C_TEAM_A)
	if len(terminii) != 1 {
		panic("team A terminal not found")
	}
	terminus := terminii[0]

	// get and merge all contracts from this team
	contracts := []*Contract{}
	bases := pool.ListAspect(C_PAYROLL, C_TEAM_A)
	for _, base := range bases {

		contracts = append(contracts, base.Contracts...)

		// (re)set wave start tickets
		base.Tickets = []*Ticket{}
	}

	// sort contracts by seniority DESC (oldest contractor gets paid first)
	sort.Sort(bySeniorityAsc(contracts))

	for _, contract := range contracts {

		// enough funds?
		if contract.Cost <= terminus.Energy {

			// pay merc and add wave ticket
			terminus.Energy -= contract.Cost
			contract.Base.Tickets = append(contract.Base.Tickets, &Ticket{Guild: contract.Guild, Merc: contract.Merc, Seniority: contract.Seniority})

		} else {
			// // no:
			// // // TODO: notify guild management
		}
	}

	for _, base := range bases {

		if len(base.Tickets) > 0 {

			// order by seniority ASC (last contractor starts the wave)
			sort.Sort(bySeniorityDesc(base.Tickets))
			base.AddAspect(C_WAVESTART)

		} else {

			base.Tickets = nil
		}
	}
}
