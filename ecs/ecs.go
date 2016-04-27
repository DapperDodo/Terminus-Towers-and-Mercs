package ecs

import (
	"errors"
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

	e := NewEntity()

	p.Entities = append(p.Entities, e)

	return e, nil
}

func (p *Pool) CloneEntity(prefab *Entity) (*Entity, error) {

	e, err := p.AddEntity()
	if err != nil {
		return nil, err
	}

	e.Clone(prefab)

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
