package game

import (
	"math/rand"
	"time"

	"towmer/api"
	"towmer/ecs"
)

var Hero, Creep, TerminusA, TerminusB, Outpost1, Outpost2, Patch1, Patch2, Patch3 *ecs.Entity

func Spawn(pool *ecs.Pool) {

	rand.Seed(time.Now().UnixNano())

	var err error

	///////////////////////////////////////////////
	// World Corners
	///////////////////////////////////////////////
	LT, err := pool.AddEntity()
	if err != nil {
		panic(err)
	}
	LT.AddAspect(ecs.C_POSITION, ecs.C_TERMINAL)
	LT.X = 0
	LT.Y = 0
	LT.Rune = '⌜'
	LT.Color = api.Color_WHITE

	RT, err := pool.AddEntity()
	if err != nil {
		panic(err)
	}
	RT.AddAspect(ecs.C_POSITION, ecs.C_TERMINAL)
	RT.X = 1
	RT.Y = 0
	RT.Rune = '⌝'
	RT.Color = api.Color_WHITE

	LB, err := pool.AddEntity()
	if err != nil {
		panic(err)
	}
	LB.AddAspect(ecs.C_POSITION, ecs.C_TERMINAL)
	LB.X = 0
	LB.Y = 1
	LB.Rune = '⌞'
	LB.Color = api.Color_WHITE

	RB, err := pool.AddEntity()
	if err != nil {
		panic(err)
	}
	RB.AddAspect(ecs.C_POSITION, ecs.C_TERMINAL)
	RB.X = 1
	RB.Y = 1
	RB.Rune = '⌟'
	RB.Color = api.Color_WHITE

	///////////////////////////////////////////////
	// Main Base Team A
	///////////////////////////////////////////////
	TerminusA, err = pool.AddEntity()
	if err != nil {
		panic(err)
	}
	TerminusA.AddAspect(ecs.C_POSITION, ecs.C_TERMINAL, ecs.C_BASE, ecs.C_MAIN, ecs.C_TEAM_A, ecs.C_TARGETABLE, ecs.C_SELECTED, ecs.C_RESOURCE, ecs.C_ENERGYSTORE, ecs.C_HEALTH)
	TerminusA.X = 0.5
	TerminusA.Y = 1 - 0.1
	TerminusA.Rune = '♛'
	TerminusA.Color = api.Color_GREEN
	TerminusA.Hotkey = api.Key_SPACE
	TerminusA.Info = api.InfoBaseMainMenu
	TerminusA.Resources = 500
	//TerminusA.Tickets = []*Ticket{{Guild: Guild_RANGERS, Merc: Merc_ARCHER, Seniority: time.Now()}}
	TerminusA.Hitpoints = 1000

	///////////////////////////////////////////////
	// Enemy Base
	///////////////////////////////////////////////
	TerminusB, err = pool.AddEntity()
	if err != nil {
		panic(err)
	}
	TerminusB.AddAspect(ecs.C_POSITION, ecs.C_TERMINAL, ecs.C_BASE, ecs.C_MAIN, ecs.C_TEAM_B, ecs.C_TARGETABLE, ecs.C_RESOURCE, ecs.C_ENERGYSTORE, ecs.C_HEALTH)
	TerminusB.X = 0.5
	TerminusB.Y = 0 + 0.1
	TerminusB.Rune = '♛'
	TerminusB.Color = api.Color_RED
	TerminusB.Resources = 500
	TerminusB.Hitpoints = 1000

	///////////////////////////////////////////////
	// Patches
	///////////////////////////////////////////////
	patches := map[string]*struct {
		cx, cy, r float64
		e         *ecs.Entity
	}{
		// "topleftclose":     {0.33, 0.2, 450, nil},
		// "topleftfar":       {0.15, 0.35, 300, nil},
		// "toprightclose":    {0.66, 0.2, 450, nil},
		// "toprightfar":      {0.85, 0.35, 300, nil},
		"bottomleftclose":  {0.33, 0.8, 450, nil},
		"bottomleftfar":    {0.15, 0.65, 300, nil},
		"bottomrightclose": {0.66, 0.8, 450, nil},
		"bottomrightfar":   {0.85, 0.65, 300, nil},
	}
	for _, patch := range patches {

		p, err := pool.AddEntity()
		if err != nil {
			panic(err)
		}
		p.AddAspect(ecs.C_POSITION, ecs.C_TERMINAL, ecs.C_RESOURCE)
		p.X = patch.cx
		p.Y = patch.cy
		p.Rune = '፨'
		p.Color = api.Color_BLUE
		p.Resources = patch.r
		patch.e = p
	}

	///////////////////////////////////////////////
	// Waypoints
	///////////////////////////////////////////////

	waypoints := map[string]*struct {
		cx, cy float64
		e      *ecs.Entity
	}{
		"topmidclose":    {0.5, 0.3, nil},
		"topleftfar":     {0.375, 0.4, nil},
		"toprightfar":    {0.625, 0.4, nil},
		"midleft":        {0.25, 0.5, nil},
		"midright":       {0.75, 0.5, nil},
		"bottommidclose": {0.5, 0.7, nil},
		"bottomleftfar":  {0.375, 0.6, nil},
		"bottomrightfar": {0.625, 0.6, nil},
	}
	for _, waypoint := range waypoints {

		p, err := pool.AddEntity()
		if err != nil {
			panic(err)
		}
		p.AddAspect(ecs.C_POSITION, ecs.C_TERMINAL, ecs.C_WAYPOINT)
		p.X = waypoint.cx
		p.Y = waypoint.cy
		//p.Rune = 'ߜ'
		//p.Rune = '.'
		p.Rune = ' '
		p.Color = api.Color_BLUE
		waypoint.e = p
	}

	///////////////////////////////////////////////
	// Waypoint Graph
	///////////////////////////////////////////////

	// onramps (from node -base- to edge -waypoint-)
	TerminusA.AddAspect(ecs.C_WAYPOINT)
	TerminusA.Outward = []*ecs.Entity{waypoints["bottommidclose"].e}
	patches["bottomleftclose"].e.AddAspect(ecs.C_WAYPOINT)
	patches["bottomleftclose"].e.Outward = []*ecs.Entity{waypoints["bottomleftfar"].e}
	patches["bottomrightclose"].e.AddAspect(ecs.C_WAYPOINT)
	patches["bottomrightclose"].e.Outward = []*ecs.Entity{waypoints["bottomrightfar"].e}
	patches["bottomleftfar"].e.AddAspect(ecs.C_WAYPOINT)
	patches["bottomleftfar"].e.Outward = []*ecs.Entity{waypoints["midleft"].e}
	patches["bottomrightfar"].e.AddAspect(ecs.C_WAYPOINT)
	patches["bottomrightfar"].e.Outward = []*ecs.Entity{waypoints["midright"].e}

	// path options
	waypoints["bottommidclose"].e.Outward = []*ecs.Entity{waypoints["bottomleftfar"].e, waypoints["bottomrightfar"].e}
	waypoints["bottomleftfar"].e.Outward = []*ecs.Entity{waypoints["topleftfar"].e, waypoints["midleft"].e}
	waypoints["bottomrightfar"].e.Outward = []*ecs.Entity{waypoints["toprightfar"].e, waypoints["midright"].e}
	waypoints["topleftfar"].e.Outward = []*ecs.Entity{waypoints["topmidclose"].e}
	waypoints["toprightfar"].e.Outward = []*ecs.Entity{waypoints["topmidclose"].e}
	waypoints["midleft"].e.Outward = []*ecs.Entity{waypoints["topleftfar"].e}
	waypoints["midright"].e.Outward = []*ecs.Entity{waypoints["toprightfar"].e}

	// offramps (from edge -waypoint- to node -base-)
	waypoints["topmidclose"].e.Outward = []*ecs.Entity{TerminusB}
	TerminusB.AddAspect(ecs.C_WAYPOINT)

	///////////////////////////////////////////////
	// Enemy
	///////////////////////////////////////////////

	Creep, err = pool.AddEntity()
	if err != nil {
		panic(err)
	}
	Creep.AddAspect(ecs.C_POSITION, ecs.C_TERMINAL, ecs.C_ROTATION, ecs.C_VELOCITY, ecs.C_OBJECTIVES, ecs.C_DAMAGER, ecs.C_SHOOTER, ecs.C_TEAM_B, ecs.C_TARGETABLE, ecs.C_HEALTH)
	Creep.X = 0.5
	Creep.Y = 0 + 0.06
	Creep.Rune = '☠'
	Creep.Color = api.Color_RED
	Creep.Speed = 0.000225
	Creep.List = []*ecs.Objective{
		&ecs.Objective{Entity: TerminusA, Range: 0.14},
	}
	Creep.Damage = 20
	Creep.Cool = 1
	Creep.FireRange = 0.15
	Creep.Hitpoints = 100

	///////////////////////////////////////////////
	// Test Units
	///////////////////////////////////////////////

	// W0, err := pool.AddEntity()
	// if err != nil {
	// 	panic(err)
	// }
	// W0.AddAspect(ecs.C_POSITION, ecs.C_TERMINAL)
	// W0.X = 0
	// W0.Y = 0
	// W0.Rune = 'T'
	// W0.Color = api.Color_MAGENTA

	// W1, err := pool.AddEntity()
	// if err != nil {
	// 	panic(err)
	// }
	// W1.AddAspect(ecs.C_POSITION, ecs.C_TERMINAL, ecs.C_ROTATION, ecs.C_VELOCITY, ecs.C_OBJECTIVES)
	// W1.X = 0
	// W1.Y = 0
	// W1.Rune = 'T'
	// W1.Color = api.Color_MAGENTA
	// W1.Speed = 0.002
	// W1.List = []*ecs.Objective{
	// 	&ecs.Objective{X: 1, Y: 0, Range: 0},
	// }

	// W2, err := pool.AddEntity()
	// if err != nil {
	// 	panic(err)
	// }
	// W2.AddAspect(ecs.C_POSITION, ecs.C_TERMINAL, ecs.C_ROTATION, ecs.C_VELOCITY, ecs.C_OBJECTIVES)
	// W2.X = 0
	// W2.Y = 0
	// W2.Rune = 'T'
	// W2.Color = api.Color_MAGENTA
	// W2.Speed = 0.002
	// W2.List = []*ecs.Objective{
	// 	&ecs.Objective{X: 0, Y: 1, Range: 0},
	// }

	// W3, err := pool.AddEntity()
	// if err != nil {
	// 	panic(err)
	// }
	// W3.AddAspect(ecs.C_POSITION, ecs.C_TERMINAL, ecs.C_ROTATION, ecs.C_VELOCITY, ecs.C_OBJECTIVES)
	// W3.X = 0
	// W3.Y = 0
	// W3.Rune = 'T'
	// W3.Color = api.Color_MAGENTA
	// W3.Speed = 0.002
	// W3.List = []*ecs.Objective{
	// 	&ecs.Objective{X: 1, Y: 1, Range: 0},
	// }
}
