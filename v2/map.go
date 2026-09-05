package views

import (
	"iter"
)

type MapWrapper[UK comparable, UV, EK, EV any] interface {
	WrapKey(UK) EK
	WrapValue(UV) EV
	UnwrapKey(EK) (UK, bool)
}

type Map[UK comparable, UV, EK, EV any, C MapWrapper[UK, UV, EK, EV]] map[UK]UV

func (m Map[UK, UV, EK, EV, C]) Len() int {
	return len(m)
}

func (m Map[UK, UV, EK, EV, C]) Get(ek EK) (EV, bool) {
	var c C
	if uk, ok := c.UnwrapKey(ek); ok {
		if uv, ok := m[uk]; ok {
			ev := c.WrapValue(uv)
			return ev, true
		}
	}
	var zero EV
	return zero, false
}

func (m Map[UK, UV, EK, EV, C]) All() iter.Seq2[EK, EV] {
	var c C
	return func(yield func(EK, EV) bool) {
		for uk, uv := range m {
			ek := c.WrapKey(uk)
			ev := c.WrapValue(uv)
			if !yield(ek, ev) {
				return
			}
		}
	}
}

func (m Map[UK, UV, EK, EV, C]) Keys() iter.Seq[EK] {
	var c C
	return func(yield func(EK) bool) {
		for uk := range (m) {
			ek := c.WrapKey(uk)
			if !yield(ek) {
				return
			}
		}
	}
}

func (m Map[UK, UV, EK, EV, C]) Values() iter.Seq[EV] {
	var c C
	return func(yield func(EV) bool) {
		for _, uv := range (m) {
			ev := c.WrapValue(uv)
			if !yield(ev) {
				return
			}
		}
	}
}