package dilema

import "reflect"

type Constructor interface {
	Creator() reflect.Value
}

type constructor struct {
	creator reflect.Value
}

func (c *constructor) Creator() reflect.Value {
	return c.creator
}
