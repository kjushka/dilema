package dilema

import "context"

type aliasType string

func (di *dicon) Ctx() context.Context {
	return di.ctx
}

func (di *dicon) SetCtx(ctx context.Context) {
	di.ctx = ctx
}

func (di *dicon) AddToCtx(alias string, value interface{}) {
	di.ctx = context.WithValue(di.ctx, aliasType(alias), value)
}

func (di *dicon) GetFromCtx(alias string) interface{} {
	return di.ctx.Value(aliasType(alias))
}
