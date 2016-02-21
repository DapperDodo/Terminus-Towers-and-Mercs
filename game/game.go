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
	TerminusA.AddAspect(ecs.C_POSITION, ecs.C_TERMINAL, ecs.C_BASE, ecs.C_TEAM_A, ecs.C_TARGETABLE, ecs.C_SELECTED, ecs.C_RESOURCE, ecs.C_ENERGYSTORE)
	TerminusA.X = 0.5
	TerminusA.Y = 1 - 0.1
	TerminusA.Rune = '♛'
	TerminusA.Color = api.Color_GREEN
	TerminusA.Hotkey = api.Key_SPACE
	TerminusA.Info = api.InfoBaseMainMenu
	TerminusA.Resources = 1000

	///////////////////////////////////////////////
	// Enemy Base
	///////////////////////////////////////////////
	TerminusB, err = pool.AddEntity()
	if err != nil {
		panic(err)
	}
	TerminusB.AddAspect(ecs.C_POSITION, ecs.C_TERMINAL, ecs.C_BASE, ecs.C_TEAM_B, ecs.C_TARGETABLE, ecs.C_RESOURCE, ecs.C_ENERGYSTORE)
	TerminusB.X = 0.5
	TerminusB.Y = 0 + 0.1
	TerminusB.Rune = '♛'
	TerminusB.Color = api.Color_RED
	TerminusB.Resources = 1000

	///////////////////////////////////////////////
	// Patches
	///////////////////////////////////////////////
	// for p := 0; p < 6; p++ {

	// 	Patch, err := pool.AddEntity()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	Patch.AddAspect(ecs.C_POSITION, ecs.C_TERMINAL, ecs.C_RESOURCE)
	// 	Patch.X = rand.Float64()
	// 	Patch.Y = rand.Float64()/2
	// 	Patch.Rune = '⚛'
	// 	Patch.Color = api.Color_BLUE
	// 	Patch.Resources = float64(100 * (p + 1))

	// 	Mirror, err := pool.AddEntity()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	Mirror.AddAspect(ecs.C_POSITION, ecs.C_TERMINAL, ecs.C_RESOURCE)
	// 	Mirror.X = 1 - Patch.X
	// 	Mirror.Y = 1 - Patch.Y
	// 	Mirror.Rune = '⚛'
	// 	Mirror.Color = api.Color_BLUE
	// 	Mirror.Resources = Patch.Resources
	// }

	patches := map[string]*struct {
		cx, cy, r float64
		e         *ecs.Entity
	}{
		// "topleftclose":     {0.33, 0.2, 900, nil},
		// "topleftfar":       {0.15, 0.35, 600, nil},
		// "toprightclose":    {0.66, 0.2, 900, nil},
		// "toprightfar":      {0.85, 0.35, 600, nil},
		"bottomleftclose":  {0.33, 0.8, 900, nil},
		"bottomleftfar":    {0.15, 0.65, 600, nil},
		"bottomrightclose": {0.66, 0.8, 900, nil},
		"bottomrightfar":   {0.85, 0.65, 600, nil},
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

	//TODO: user interface for creating path from terminus A

	// ///////////////////////////////////////////////
	// // Outpost 1
	// ///////////////////////////////////////////////
	// Outpost1, err = pool.AddEntity()
	// if err != nil {
	// 	panic(err)
	// }
	// Outpost1.AddAspect(ecs.C_POSITION, ecs.C_TERMINAL, ecs.C_TEAM_A, ecs.C_TARGETABLE)
	// Outpost1.X = w * 0.50
	// Outpost1.Y = h * 0.75
	// Outpost1.Rune = '♜'
	// Outpost1.Color = api.Color_GREEN

	// ///////////////////////////////////////////////
	// // Outpost2
	// ///////////////////////////////////////////////
	// Outpost2, err = pool.AddEntity()
	// if err != nil {
	// 	panic(err)
	// }
	// Outpost2.AddAspect(ecs.C_POSITION, ecs.C_TERMINAL, ecs.C_TEAM_A, ecs.C_TARGETABLE)
	// Outpost2.X = w * 0.25
	// Outpost2.Y = h * 0.25
	// Outpost2.Rune = '♜'
	// Outpost2.Color = api.Color_GREEN

	// ///////////////////////////////////////////////
	// // Hero
	// ///////////////////////////////////////////////
	// Hero, err = pool.AddEntity()
	// if err != nil {
	// 	panic(err)
	// }
	// Hero.AddAspect(ecs.C_POSITION, ecs.C_TERMINAL, ecs.C_CONTROLLABLE, ecs.C_SHOOTER, ecs.C_TEAM_A, ecs.C_TARGETABLE)
	// Hero.X = w / 2.0
	// Hero.Y = h / 2.0
	// Hero.Rune = '☃'
	// Hero.Color = api.Color_GREEN
	// Hero.Cool = 0.66
	// Hero.FireRange = 15

	///////////////////////////////////////////////
	// Enemy
	///////////////////////////////////////////////

	Creep, err = pool.AddEntity()
	if err != nil {
		panic(err)
	}
	Creep.AddAspect(ecs.C_POSITION, ecs.C_TERMINAL, ecs.C_ROTATION, ecs.C_VELOCITY, ecs.C_OBJECTIVES, ecs.C_SHOOTER, ecs.C_TEAM_B, ecs.C_TARGETABLE)
	Creep.X = 0.5
	Creep.Y = 0 + 0.06
	Creep.Rune = '☠'
	Creep.Color = api.Color_RED
	Creep.DX = 0
	Creep.DY = 0.1
	Creep.Speed = 0.00125
	Creep.List = []*ecs.Objective{
		&ecs.Objective{Entity: TerminusA, Range: 0.14},
	}
	Creep.Cool = 1
	Creep.FireRange = 0.15

}
