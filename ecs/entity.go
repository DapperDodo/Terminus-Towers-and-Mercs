package ecs

import (
	"towmer/api"
)

type Entity struct {
	*Position
	*Rotation
	*Velocity
	*Terminal
	*TeamA
	*TeamB
	*Base
	*Objectives
	*Shooter
	*Cooldown
	*Bullet
	*Targetable
	*Dying
	*Baseselection
	*Selected
	*Basebuilding
	*Resource
	*Tabbable
	*Energynode
	*Energystore
	*Target
	*Waypoint
	*Pathbuilding
	*Path
	*Contracting
	*Payroll
}

func (e *Entity) AddAspect(aspect ...ComponentType) {
	for _, t := range aspect {
		e.Add(t)
	}
}

func (e *Entity) HasAspect(aspect ...ComponentType) bool {
	for _, t := range aspect {
		if !e.Has(t) {
			return false
		}
	}
	return true
}

func (e *Entity) DelAspect(aspect ...ComponentType) {
	for _, t := range aspect {
		if e.Has(t) {
			e.Del(t)
		}
	}
}

func (e *Entity) Add(t ComponentType) {
	switch t {
	case C_POSITION:
		e.Position.Active = true
	case C_ROTATION:
		e.Rotation.Active = true
	case C_VELOCITY:
		e.Velocity.Active = true
	case C_TERMINAL:
		e.Terminal.Active = true
	case C_TEAM_A:
		e.TeamA.Active = true
	case C_TEAM_B:
		e.TeamB.Active = true
	case C_BASE:
		e.Base.Active = true
	case C_OBJECTIVES:
		e.Objectives.Active = true
	case C_SHOOTER:
		e.Shooter.Active = true
	case C_COOLDOWN:
		e.Cooldown.Active = true
	case C_BULLET:
		e.Bullet.Active = true
	case C_TARGETABLE:
		e.Targetable.Active = true
	case C_DYING:
		e.Dying.Active = true
	case C_BASESELECTION:
		e.Baseselection.Active = true
	case C_SELECTED:
		e.Selected.Active = true
	case C_BASEBUILDING:
		e.Basebuilding.Active = true
	case C_RESOURCE:
		e.Resource.Active = true
	case C_TABBABLE:
		e.Tabbable.Active = true
	case C_ENERGYNODE:
		e.Energynode.Active = true
	case C_ENERGYSTORE:
		e.Energystore.Active = true
	case C_TARGET:
		e.Target.Active = true
	case C_WAYPOINT:
		e.Waypoint.Active = true
	case C_PATHBUILDING:
		e.Pathbuilding.Active = true
	case C_PATH:
		e.Path.Active = true
	case C_CONTRACTING:
		e.Contracting.Active = true
	case C_PAYROLL:
		e.Payroll.Active = true
	}
}

func (e *Entity) Del(t ComponentType) {
	switch t {
	case C_POSITION:
		e.Position.Active = false
		e.X = 0
		e.Y = 0
	case C_ROTATION:
		e.Rotation.Active = false
		e.DX = 0
		e.DY = 0
	case C_VELOCITY:
		e.Velocity.Active = false
		e.Speed = 0
	case C_TERMINAL:
		e.Terminal.Active = false
		e.Rune = '?'
		e.Color = api.Color_UNKNOWN
		e.BgColor = api.Color_UNKNOWN
	case C_TEAM_A:
		e.TeamA.Active = false
	case C_TEAM_B:
		e.TeamB.Active = false
	case C_BASE:
		e.Base.Active = false
	case C_OBJECTIVES:
		e.Objectives.Active = false
		e.List = nil
		e.lastReached = nil
	case C_SHOOTER:
		e.Shooter.Active = false
		e.Cool = 0
		e.FireRange = 0
	case C_COOLDOWN:
		e.Cooldown.Active = false
		e.current = 0
	case C_BULLET:
		e.Bullet.Active = false
	case C_TARGETABLE:
		e.Targetable.Active = false
	case C_DYING:
		e.Dying.Active = false
		e.TimeToLive = 0
		e.sickbed = 0
	case C_BASESELECTION:
		e.Baseselection.Active = false
		e.Hotkey = api.Key_UNKNOWN
	case C_SELECTED:
		e.Selected.Active = false
		e.Info = nil
	case C_BASEBUILDING:
		e.Basebuilding.Active = false
	case C_RESOURCE:
		e.Resource.Active = false
		e.Resources = 0
	case C_TABBABLE:
		e.Tabbable.Active = false
		e.TabConfirmed = false
		e.TabActive = false
	case C_ENERGYNODE:
		e.Energynode.Active = false
		e.Upstream = nil
		e.Downstream = nil
		e.timeSinceLastEmission = 0
	case C_ENERGYSTORE:
		e.Energystore.Active = false
		e.Energy = 0
	case C_TARGET:
		e.Target.Active = false
		e.TargetEntity = nil
	case C_WAYPOINT:
		e.Waypoint.Active = false
		e.Outward = nil
	case C_PATHBUILDING:
		e.Pathbuilding.Active = false
	case C_PATH:
		e.Path.Active = false
		e.Waypoints = nil
	case C_CONTRACTING:
		e.Contracting.Active = false
		e.Guild = 0
		e.Merc = 0
		e.Party = 0
		e.Tail = false
		e.Signed = false
	case C_PAYROLL:
		e.Payroll.Active = false
		e.Contracts = nil
		e.Burden = 0
	}
}

func (e *Entity) Has(t ComponentType) bool {
	switch t {
	case C_POSITION:
		return e.Position.Active
	case C_ROTATION:
		return e.Rotation.Active
	case C_VELOCITY:
		return e.Velocity.Active
	case C_TERMINAL:
		return e.Terminal.Active
	case C_TEAM_A:
		return e.TeamA.Active
	case C_TEAM_B:
		return e.TeamB.Active
	case C_BASE:
		return e.Base.Active
	case C_OBJECTIVES:
		return e.Objectives.Active
	case C_SHOOTER:
		return e.Shooter.Active
	case C_COOLDOWN:
		return e.Cooldown.Active
	case C_BULLET:
		return e.Bullet.Active
	case C_TARGETABLE:
		return e.Targetable.Active
	case C_DYING:
		return e.Dying.Active
	case C_BASESELECTION:
		return e.Baseselection.Active
	case C_SELECTED:
		return e.Selected.Active
	case C_BASEBUILDING:
		return e.Basebuilding.Active
	case C_RESOURCE:
		return e.Resource.Active
	case C_TABBABLE:
		return e.Tabbable.Active
	case C_ENERGYNODE:
		return e.Energynode.Active
	case C_ENERGYSTORE:
		return e.Energystore.Active
	case C_TARGET:
		return e.Target.Active
	case C_WAYPOINT:
		return e.Waypoint.Active
	case C_PATHBUILDING:
		return e.Pathbuilding.Active
	case C_PATH:
		return e.Path.Active
	case C_CONTRACTING:
		return e.Contracting.Active
	case C_PAYROLL:
		return e.Payroll.Active
	}
	return false
}
