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
	IXrayApi interface {
		AddOutbound(ctx context.Context, json *gjson.Json) error
		DelOutbound(ctx context.Context, tag string) (err error)
		AddSystemInbound(ctx context.Context, addr string, tag string) (err error)
		DelInbound(ctx context.Context, tag string) (err error)
		Stat(ctx context.Context, inbound bool, tag string, down bool) (val int64, err error)
		Start(ctx context.Context)
		Stop(ctx context.Context)
	}
)

var (
	localXrayApi IXrayApi
)

func XrayApi() IXrayApi {
	if localXrayApi == nil {
		panic("implement not found for interface IXrayApi, forgot register?")
	}
	return localXrayApi
}

func RegisterXrayApi(i IXrayApi) {
	localXrayApi = i
}
