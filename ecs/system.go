package ecs

import (
	"math"
	"time"

	"towmer/api"
)

func Control(p *Pool, key api.Key) {

	if key == api.Key_SPACE {
		selected := p.ListAspect(C_SELECTED)
		if len(selected) > 0 {
			selected[0].DelAspect(C_SELECTED, C_BASEBUILDING, C_PATHBUILDING, C_CONTRACTING)
		}
		bases := p.ListAspect(C_TEAM_A, C_BASE)
		for _, base := range bases {
			base.AddAspect(C_TABBABLE, C_BASESELECTION)
		}
		return
	}

	if key == api.Key_BACKSPACE {
		selected := p.ListAspect(C_SELECTED, C_CONTRACTING)
		if len(selected) > 0 {
			if selected[0].Merc != Merc_UNKNOWN {
				selected[0].Merc = Merc_UNKNOWN
			} else if selected[0].Guild != Guild_UNKNOWN {
				selected[0].Guild = Guild_UNKNOWN
			} else {
				selected[0].DelAspect(C_CONTRACTING)
				selected[0].Info = api.InfoBaseMainMenu
			}
		}
		return
	}

	if key == api.Key_A {
		selected := p.ListAspect(C_SELECTED, C_CONTRACTING)
		if len(selected) > 0 {
			selected[0].Merc = Merc_ARCHER
		}
		return
	}

	if key == api.Key_B {
		selected := p.ListAspect(C_SELECTED, C_CONTRACTING)
		if len(selected) > 0 {
			selected[0].Merc = Merc_BRAWLER
			return
		}
		selected = p.ListAspect(C_SELECTED)
		if len(selected) > 0 {
			selected[0].AddAspect(C_BASEBUILDING)
		}
		return
	}

	if key == api.Key_C {
		selected := p.ListAspect(C_SELECTED)
		if len(selected) > 0 {
			selected[0].AddAspect(C_CONTRACTING)
		}
		return
	}

	if key == api.Key_G {
		selected := p.ListAspect(C_SELECTED, C_CONTRACTING)
		if len(selected) > 0 {
			selected[0].Merc = Merc_GLADIATOR
		}
		return
	}

	if key == api.Key_H {
		selected := p.ListAspect(C_SELECTED, C_CONTRACTING)
		if len(selected) > 0 {
			selected[0].Merc = Merc_HUNTER
		}
		return
	}

	if key == api.Key_R {
		selected := p.ListAspect(C_SELECTED, C_CONTRACTING)
		if len(selected) > 0 {
			selected[0].Guild = Guild_RANGERS
		}
		return
	}

	if key == api.Key_P {
		selected := p.ListAspect(C_SELECTED)
		if len(selected) > 0 {
			selected[0].AddAspect(C_PATHBUILDING)
		}
		return
	}

	if key == api.Key_W {
		selected := p.ListAspect(C_SELECTED, C_CONTRACTING)
		if len(selected) > 0 {
			selected[0].Guild = Guild_WARRIORS
		}
		return
	}

	if key == api.Key_TAB {
		tabs := p.ListAspect(C_TABBABLE)
		if len(tabs) > 0 {
			found := false
			for i, tab := range tabs {
				if tab.TabActive {
					found = true
					tab.TabActive = false
					if i == len(tabs)-1 {
						tabs[0].TabActive = true
					} else {
						tabs[i+1].TabActive = true
					}
					break
				}
			}
			if !found && len(tabs) > 0 {
				tabs[0].TabActive = true
			}
		}
		return
	}

	if key == api.Key_ENTER {
		tabs := p.ListAspect(C_TABBABLE)
		if len(tabs) > 0 {
			for _, tab := range tabs {
				if tab.TabActive {
					tab.TabConfirmed = true
				} else {
					tab.TabConfirmed = false
				}
			}
		}
		selected := p.ListAspect(C_SELECTED, C_CONTRACTING)
		if len(selected) > 0 {
			selected[0].Signed = true
		}
		return
	}
}

func Update(p *Pool, deltaTime float64) {

	for _, e := range p.Entities {

		if e.HasAspect(C_DYING) {
			die(e, p, deltaTime)
		}

		if e.HasAspect(C_POSITION, C_OBJECTIVES) {
			turn(e)
		}

		if e.HasAspect(C_POSITION, C_ROTATION, C_VELOCITY) {
			move(e)
		}

		if e.HasAspect(C_POSITION, C_VELOCITY) {
			checkBounds(e)
		}

		if e.HasAspect(C_POSITION, C_VELOCITY, C_OBJECTIVES) {
			checkObjective(e)
		}

		if e.HasAspect(C_POSITION, C_SHOOTER, C_COOLDOWN) {
			cooldown(e, deltaTime)
		}

		if e.HasAspect(C_POSITION, C_SHOOTER) {
			fire(e, p)
		}

		if e.HasAspect(C_POSITION, C_BULLET) {
			hit(e)
		}

		if e.HasAspect(C_HEALTH) {
			checkHealth(e)
		}

		if e.HasAspect(C_TABBABLE, C_BASESELECTION) {
			selectBase(e, p)
		}

		if e.HasAspect(C_BASEBUILDING) {
			buildBase(e, p)
		}

		if e.HasAspect(C_BASE, C_RESOURCE, C_ENERGYSTORE) {
			mine(e, deltaTime)
		}

		if e.HasAspect(C_ENERGYNODE, C_ENERGYSTORE) {
			emit(e, p, deltaTime)
		}

		if e.HasAspect(C_POSITION, C_ENERGYSTORE, C_TARGET) {
			transit(e)
		}

		if e.HasAspect(C_PATHBUILDING) {
			buildPath(e, p)
		}

		if e.HasAspect(C_CONTRACTING) {
			contract(e, p)
		}

		if e.HasAspect(C_WAVESTART) {
			wave(e, p, deltaTime)
		}
	}
}

func die(e *Entity, p *Pool, dt float64) {

	e.sickbed += dt
	if e.sickbed >= e.TimeToLive {
		e.DelAspect(C_POSITION, C_TERMINAL, C_DYING)
		p.Delete(e)
	}
}

func turn(e *Entity) {

	if len(e.Objectives.List) > 0 {

		targetX, targetY := e.Objectives.List[0].Point()

		dX := targetX - e.X
		dY := targetY - e.Y
		d := math.Sqrt(math.Pow(dX, 2) + math.Pow(dY, 2))

		e.Add(C_ROTATION)
		e.DX = dX / d
		e.DY = dY / d
	} else {
		e.Del(C_OBJECTIVES)
		e.Del(C_ROTATION)
	}
}

func move(e *Entity) {

	e.X += (e.DX * e.Speed)
	e.Y += (e.DY * e.Speed)
}

func checkBounds(e *Entity) {

	if e.Y < 0 {
		e.Y = 0
		e.Del(C_VELOCITY)
	}
	if e.Y >= 1 {
		e.Y = 1
		e.Del(C_VELOCITY)
	}
	if e.X < 0 {
		e.X = 0
		e.Del(C_VELOCITY)
	}
	if e.X >= 1 {
		e.X = 1
		e.Del(C_VELOCITY)
	}
}

func checkObjective(e *Entity) {

	if len(e.Objectives.List) > 0 {

		targetX, targetY := e.Objectives.List[0].Point()

		if math.Sqrt(math.Pow(e.X-targetX, 2)+math.Pow(e.Y-targetY, 2)) < e.Objectives.List[0].Range {

			last := e.Objectives.List[0]
			e.Objectives.List = e.Objectives.List[1:]

			if len(e.Objectives.List) == 0 {
				e.Add(C_REACHED)
				e.lastReached = last
				e.DelAspect(C_OBJECTIVES, C_VELOCITY, C_ROTATION)
			}
		}
	}
}

func fire(e *Entity, p *Pool) {

	if e.Has(C_COOLDOWN) {
		return
	}

	if !e.Has(C_DAMAGER) {
		return
	}

	var targets []*Entity
	if e.Has(C_TEAM_A) {
		targets = p.ListAspect(C_TEAM_B, C_TARGETABLE)
	} else {
		targets = p.ListAspect(C_TEAM_A, C_TARGETABLE)
	}

	var closest_target *Entity
	var closest_range float64 = 9999
	for _, target := range targets {
		d := math.Sqrt(math.Pow(e.X-target.X, 2) + math.Pow(e.Y-target.Y, 2))
		if d < closest_range {
			closest_range = d
			closest_target = target
		}
	}

	if closest_range < e.FireRange {
		bullet, err := p.AddEntity()
		if err != nil {
			panic(err)
		}
		bullet.AddAspect(C_POSITION, C_TERMINAL, C_VELOCITY, C_OBJECTIVES, C_BULLET, C_DAMAGER)
		bullet.Damage = e.Damage
		bullet.X = e.X
		bullet.Y = e.Y
		bullet.Rune = '.'
		bullet.Speed = 0.01
		bullet.List = []*Objective{&Objective{Entity: closest_target, Range: 0.01}}
		if e.Has(C_TEAM_A) {
			bullet.Add(C_TEAM_A)
			bullet.Color = api.Color_GREEN
		} else {
			bullet.Add(C_TEAM_B)
			bullet.Color = api.Color_RED
		}

		e.Add(C_COOLDOWN)
		e.current = 0
	}
}

func cooldown(e *Entity, dt float64) {

	e.current += dt
	if e.current >= e.Cool {
		e.Del(C_COOLDOWN)
	}
}

func hit(e *Entity) {

	if e.Has(C_REACHED) {

		if e.lastReached.Has(C_HEALTH) {
			e.lastReached.Hitpoints -= e.Damage
		}

		e.Rune = '⚡'
		e.Color = api.Color_WHITE
		e.DelAspect(C_ROTATION, C_VELOCITY, C_TEAM_A, C_TEAM_B, C_BASE, C_OBJECTIVES, C_SHOOTER, C_COOLDOWN, C_BULLET, C_DAMAGER, C_REACHED)
		e.Add(C_DYING)
		e.TimeToLive = 0.075
		e.sickbed = 0
	}
}

func checkHealth(e *Entity) {

	if e.Hitpoints <= 0 {
		e.Rune = '☠' //'☨'
		e.Color = api.Color_WHITE
		e.DelAspect(C_ROTATION, C_VELOCITY, C_TEAM_A, C_TEAM_B, C_BASE, C_OBJECTIVES, C_SHOOTER, C_COOLDOWN, C_BULLET, C_DAMAGER, C_HEALTH)
		e.Add(C_DYING)
		e.TimeToLive = 3
		e.sickbed = 0
	}
}

func selectBase(e *Entity, p *Pool) {

	if e.TabConfirmed {
		e.AddAspect(C_SELECTED)
		e.Info = api.InfoBaseMainMenu
		bases := p.ListAspect(C_TABBABLE, C_BASESELECTION)
		for _, base := range bases {
			base.DelAspect(C_TABBABLE, C_BASESELECTION)
		}
	}
}

func buildBase(e *Entity, p *Pool) {

	// set tabbable selection
	patches := p.ListAspect(C_RESOURCE)

	var closest_patch *Entity
	var closest_range float64 = 9999
	for _, patch := range patches {
		if patch.HasAspect(C_TABBABLE) || patch.HasAspect(C_TEAM_B) || patch.HasAspect(C_BASE) {
			continue
		}
		d := math.Sqrt(math.Pow(e.X-patch.X, 2) + math.Pow(e.Y-patch.Y, 2))
		if d < 0.25 {
			if d < closest_range {
				closest_range = d
				closest_patch = patch
			}
			patch.AddAspect(C_TABBABLE)
		}
	}
	if closest_patch != nil {
		closest_patch.TabActive = true
	}

	// return to main base menu if there are no patches in range
	patches = p.ListAspect(C_TABBABLE)
	if len(patches) == 0 {
		e.DelAspect(C_BASEBUILDING)
		e.Info = api.InfoBaseMainMenuNoB
		return
	}

	// build on confirmation
	e.Info = api.InfoBaseBuildSelect
	for _, patch := range patches {
		if patch.TabConfirmed {

			// remove tabs
			tabs := p.ListAspect(C_TABBABLE)
			for _, tab := range tabs {
				tab.DelAspect(C_TABBABLE)
			}

			// build base on confirmed location, mark this base as energynode downstream and select it
			patch.AddAspect(C_BASE, C_TEAM_A, C_SELECTED, C_ENERGYNODE, C_ENERGYSTORE)
			patch.Rune = '♛'
			patch.Color = api.Color_GREEN
			patch.Info = api.InfoBaseMainMenu
			patch.Downstream = []*Entity{e}
			patch.timeSinceLastEmission = 9999

			// mark new base as energynode upstream
			e.AddAspect(C_ENERGYNODE)
			if e.Upstream == nil {
				e.Upstream = []*Entity{patch}
			} else {
				e.Upstream = append(e.Upstream, patch)
			}

			// exit builder mode and deselect this base
			e.DelAspect(C_BASEBUILDING, C_SELECTED)
		}
	}
}

func mine(e *Entity, dt float64) {

	var yield float64
	if e.Resources > 0 {
		yield = 1 * dt
	}
	if yield > e.Resources {
		yield = e.Resources
		e.Resources = 0
	} else {
		e.Resources -= yield
	}

	e.Energy += yield
}

func emit(e *Entity, p *Pool, dt float64) {

	// bail on no energy to emit
	if e.Energy <= 0 {
		return
	}

	// bail on no downstream available
	// this will effectively buffer the energy until a downstream becomes available
	if len(e.Downstream) == 0 {
		return
	}

	// bail on emission not yet due
	if e.timeSinceLastEmission < 0.333 {
		e.timeSinceLastEmission += dt
		return
	}

	// emit energy packet
	e.timeSinceLastEmission = 0
	packet, err := p.AddEntity()
	if err != nil {
		panic(err)
	}
	packet.AddAspect(C_POSITION, C_TERMINAL, C_VELOCITY, C_OBJECTIVES, C_ENERGYSTORE, C_TARGET)
	packet.X = e.X
	packet.Y = e.Y
	packet.Rune = '+'
	packet.Speed = 0.0125
	packet.List = []*Objective{&Objective{Entity: e.Downstream[0], Range: packet.Speed / 2}}
	packet.TargetEntity = e.Downstream[0]
	if e.Has(C_TEAM_A) {
		packet.Add(C_TEAM_A)
		packet.Color = api.Color_GREEN
	} else {
		packet.Add(C_TEAM_B)
		packet.Color = api.Color_RED
	}

	// transfer energy from base to packet
	packet.Energy = e.Energy
	e.Energy = 0
}

func transit(e *Entity) {

	// bail on packet still travelling
	if len(e.Objectives.List) > 0 {
		return
	}

	// transfer energy to base
	if e.TargetEntity != nil {
		if e.TargetEntity.Has(C_ENERGYSTORE) {
			e.TargetEntity.Energy += e.Energy
			e.Energy = 0
		}
	}

	// kill packet
	e.Add(C_DYING)
	e.sickbed = 0
	if len(e.TargetEntity.Downstream) == 0 {
		e.Rune = '$' // when packet reaches Terminus, celebrate a little...
		e.TimeToLive = 0.05
	} else {
		e.Rune = ' '
		e.TimeToLive = 0
		e.DelAspect(C_POSITION)
	}
	e.Color = api.Color_GREEN
	e.DelAspect(C_ROTATION, C_VELOCITY, C_TEAM_A, C_TEAM_B, C_OBJECTIVES, C_ENERGYSTORE, C_TARGET)
}

func buildPath(e *Entity, p *Pool) {

	// bail on path already set
	if len(e.Waypoints) > 0 && e.Waypoints[len(e.Waypoints)-1].HasAspect(C_TERMINAL, C_TEAM_B) {
		e.DelAspect(C_PATHBUILDING)
		e.Info = api.InfoBaseMainMenuSetP
		return
	}

	// bail on missing onramp
	if len(e.Outward) == 0 {
		e.DelAspect(C_PATHBUILDING)
		e.Info = api.InfoBaseMainMenuNoP
		return
	}

	e.Info = api.InfoPathSelect

	// start the path by adding the onramp
	if !e.HasAspect(C_PATH) {
		e.AddAspect(C_PATH)
		e.Waypoints = []*Entity{e.Outward[0]}
	}

	// get the path tip
	tip := e.Waypoints[len(e.Waypoints)-1]

	// bail on no (further) waypoints
	if len(tip.Outward) == 0 {
		e.DelAspect(C_PATHBUILDING, C_PATH)
		e.Info = api.InfoBaseMainMenuNoP
		return
	}

	// immediately add the waypoint if there is only one
	if len(tip.Outward) == 1 {
		e.Waypoints = append(e.Waypoints, tip.Outward[0])

		// finish the path once the enemy terminus has been reached
		if tip.Outward[0].HasAspect(C_TERMINAL, C_TEAM_B) {

			// exit builder mode
			e.DelAspect(C_PATHBUILDING)
			e.Info = api.InfoBaseMainMenu
		}
		return
	}

	// what options do we have?
	active := false
	for _, waypoint := range tip.Outward {
		if !waypoint.HasAspect(C_TABBABLE) {
			waypoint.AddAspect(C_TABBABLE)
			if !active {
				active = true
				tip.Outward[0].TabActive = true
			}
			waypoint.Rune = '߉'
		}
	}

	// check selection result
	for _, waypoint := range tip.Outward {

		if waypoint.TabConfirmed {

			// remove tabs
			tabs := p.ListAspect(C_TABBABLE)
			for _, tab := range tabs {
				tab.DelAspect(C_TABBABLE)
				tab.Rune = ' '
			}

			// add selected waypoint to path
			e.Waypoints = append(e.Waypoints, waypoint)

			return
		}
	}
}

func contract(e *Entity, p *Pool) {

	// main contracting menu
	e.Info = api.InfoContractGuilds

	// guild already selected?
	switch e.Guild {
	case Guild_RANGERS:
		e.Info = api.InfoContractRangers
	case Guild_WARRIORS:
		e.Info = api.InfoContractWarriors
	default:
		e.Signed = false
		return
	}

	// merc already selected?
	switch e.Merc {
	case Merc_UNKNOWN:
		e.Signed = false
		return
	default:
		e.Info = api.InfoContractSign
	}

	if e.Signed {

		c := &Contract{e.Guild, e.Merc, 5.0, e, time.Now()}

		e.Burden += c.Cost

		if !e.HasAspect(C_PAYROLL) {
			e.AddAspect(C_PAYROLL)
			e.Contracts = []*Contract{c}
		} else {
			e.Contracts = append(e.Contracts, c)
		}

		e.DelAspect(C_CONTRACTING)
		e.Info = api.InfoBaseMainMenu
	}
}

func wave(e *Entity, p *Pool, dt float64) {

	if len(e.Tickets) > 0 {
		if e.Tickets[0].WaitForIt < 1.0 {

			e.Tickets[0].WaitForIt += dt

		} else {

			// spawn merc
			merc, err := p.CloneEntity(MercPrefab(e.Tickets[0].Merc))
			if err != nil {
				panic(err)
			}
			if merc == nil {
				panic("merc is nil")
			}
			merc.X = e.X
			merc.Y = e.Y
			merc.List = make([]*Objective, len(e.Waypoints))
			for idx, waypoint := range e.Waypoints {
				merc.List[idx] = &Objective{Entity: waypoint, Range: 0.0025}
			}
			if e.Has(C_TEAM_A) {
				merc.Add(C_TEAM_A)
				merc.Color = api.Color_GREEN
			} else {
				merc.Add(C_TEAM_B)
				merc.Color = api.Color_RED
			}

			// started, so remove ticket
			e.Tickets = e.Tickets[1:]
		}
	} else {
		e.DelAspect(C_WAVESTART)
	}
}
