// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	IXrayApi interface {
		AddSystemInbound(ctx context.Context, addr string, tag string) (err error)
		DelSystemInbound(ctx context.Context, tag string) (err error)
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
