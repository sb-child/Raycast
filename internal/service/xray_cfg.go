// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
)

type (
	IXrayCfg interface {
		Generate(ctx context.Context)
		GetUserOutboundList(ctx context.Context) []string
		GetOutboundSetting(ctx context.Context, n int, tag string) *gjson.Json
		GetOutboundName(ctx context.Context, n int) string
		Start(ctx context.Context)
		Stop(ctx context.Context)
	}
)

var (
	localXrayCfg IXrayCfg
)

func XrayCfg() IXrayCfg {
	if localXrayCfg == nil {
		panic("implement not found for interface IXrayCfg, forgot register?")
	}
	return localXrayCfg
}

func RegisterXrayCfg(i IXrayCfg) {
	localXrayCfg = i
}
