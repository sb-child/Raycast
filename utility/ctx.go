package utility

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func CheckCancel(ctx context.Context) bool {
	if ctx.Err() == context.Canceled || gerror.Is(ctx.Err(), context.Canceled) {
		g.Log().Error(ctx, "Context canceled")
		// panic(ctx.Err())
		return true
	}
	return false
}
