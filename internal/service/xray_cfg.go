// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	IXrayCfg interface {
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
