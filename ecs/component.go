package ecs

import (
	"towmer/api"
)

type Component struct {
	Active bool
}

type ComponentType int

const (
	C_UNKNOWN ComponentType = iota
	C_POSITION
	C_ROTATION
	C_VELOCITY
	C_TERMINAL
	C_TEAM_A
	C_TEAM_B
	C_CONTROLLABLE
	C_BASE
	C_OBJECTIVES
	C_SHOOTER
	C_COOLDOWN
	C_BULLET
	C_TARGETABLE
	C_DYING
	C_BASESELECTION
	C_SELECTED
	C_BASEBUILDING
	C_RESOURCE
	C_TABBABLE
	C_ENERGYNODE
	C_ENERGYSTORE
	C_TARGET
	C_WAYPOINT
	C_PATHBUILDING
	C_PATH
	C_CONTRACTING
	C_PAYROLL
)

type Position struct {
	*Component
	X, Y float64
}

type Rotation struct {
	*Component
	DX, DY float64
}

type Velocity struct {
	*Component
	Speed float64
}

type Terminal struct {
	*Component
	Rune           rune
	Color, BgColor api.Color
}

type TeamA struct {
	*Component
}

type TeamB struct {
	*Component
}

type Controllable struct {
	*Component
}

type Base struct {
	*Component
}

type Objectives struct {
	*Component
	List        []*Objective
	lastReached *Objective
}

type Objective struct {
	*Entity
	X, Y, Range float64
}

func (o Objective) Point() (x, y float64) {

	if o.Entity != nil {
		return o.Entity.X, o.Entity.Y
	}
	return o.X, o.Y
}

type Shooter struct {
	*Component
	Cool, FireRange float64
}

type Cooldown struct {
	*Component
	current float64
}

type Bullet struct {
	*Component
}

type Targetable struct {
	*Component
}

type Dying struct {
	*Component
	TimeToLive, sickbed float64
}

type Baseselection struct {
	*Component
	Hotkey api.Key
}

type Selected struct {
	*Component
	Info []string
}

type Basebuilding struct {
	*Component
}

type Resource struct {
	*Component
	Resources float64
}

type Tabbable struct {
	*Component
	TabActive, TabConfirmed bool
}

type Energynode struct {
	*Component
	Upstream, Downstream  []*Entity
	timeSinceLastEmission float64
}

type Energystore struct {
	*Component
	Energy float64
}

type Target struct {
	*Component
	TargetEntity *Entity
}

type Waypoint struct {
	*Component
	Outward []*Entity // when constructing a path outward (towards the enemy base), this contains the options for the next path node
}

type Pathbuilding struct {
	*Component
}

type Path struct {
	*Component
	Waypoints []*Entity
}

type Contracting struct {
	*Component
	Guild
	Merc
	Party  int
	Tail   bool
	Signed bool
}

type Payroll struct {
	*Component
	Contracts []*Contract
	Burden    float64
}
