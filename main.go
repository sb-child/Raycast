package main

import (
	_ "raycast/internal/logic"
	_ "raycast/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"raycast/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
