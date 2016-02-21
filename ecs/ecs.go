package ecs

import (
	"errors"

	"towmer/api"
)

type Pool struct {
	Size     int
	Entities []*Entity
	Souls    []*Entity
}

func NewPool(maxUnits int) *Pool {
	return &Pool{Size: maxUnits, Entities: make([]*Entity, 0, maxUnits)}
}

var (
	ErrPoolFull = errors.New("pool full")
)

func (p *Pool) AddEntity() (*Entity, error) {

	// reincarnate a soul if possible
	if len(p.Souls) > 0 {
		e := p.Souls[0]
		p.Souls = p.Souls[1:]
		return e, nil
	}

	if len(p.Entities) == p.Size {
		err := ErrPoolFull
		// TODO: log
		return nil, err
	}

	e := &Entity{
		&Position{&Component{}, 0, 0},
		&Rotation{&Component{}, 0, 0},
		&Velocity{&Component{}, 0},
		&Terminal{&Component{}, '?', api.Color_UNKNOWN, api.Color_UNKNOWN},
		&TeamA{&Component{}},
		&TeamB{&Component{}},
		&Base{&Component{}},
		&Objectives{&Component{}, nil, nil},
		&Shooter{&Component{}, 0, 0},
		&Cooldown{&Component{}, 0},
		&Bullet{&Component{}},
		&Targetable{&Component{}},
		&Dying{&Component{}, 0, 0},
		&Baseselection{&Component{}, 0},
		&Selected{&Component{}, nil},
		&Basebuilding{&Component{}},
		&Resource{&Component{}, 0},
		&Tabbable{&Component{}, false, false},
		&Energynode{&Component{}, nil, nil, 0},
		&Energystore{&Component{}, 0},
		&Target{&Component{}, nil},
		&Waypoint{&Component{}, nil},
		&Pathbuilding{&Component{}},
		&Path{&Component{}, nil},
		&Contracting{&Component{}, 0, 0, 0, false, false},
		&Payroll{&Component{}, nil, 0},
	}
	p.Entities = append(p.Entities, e)

	return e, nil
}

func (p *Pool) ListAspect(aspect ...ComponentType) []*Entity {

	list := make([]*Entity, 0)
	for _, e := range p.Entities {
		if e.HasAspect(aspect...) {
			list = append(list, e)
		}
	}
	return list
}

func (p *Pool) Delete(e *Entity) {
	p.Souls = append(p.Souls, e)
}
