// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
)

type (
	IXrayExec interface {
		Start(ctx context.Context)
		Stop(ctx context.Context)
	}
)

var (
	localXrayExec IXrayExec
)

func XrayExec() IXrayExec {
	if localXrayExec == nil {
		panic("implement not found for interface IXrayExec, forgot register?")
	}
	return localXrayExec
}

func RegisterXrayExec(i IXrayExec) {
	localXrayExec = i
}
