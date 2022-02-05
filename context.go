package dilema

import "context"

func (di *dicon) Ctx() context.Context {
	return di.ctx
}