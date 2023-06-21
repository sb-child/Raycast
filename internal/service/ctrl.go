// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	ICtrl interface {
		Start(ctx context.Context)
		Stop(ctx context.Context)
	}
)

var (
	localCtrl ICtrl
)

func Ctrl() ICtrl {
	if localCtrl == nil {
		panic("implement not found for interface ICtrl, forgot register?")
	}
	return localCtrl
}

func RegisterCtrl(i ICtrl) {
	localCtrl = i
}
