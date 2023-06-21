package utility

import (
	loggerService "github.com/xtls/xray-core/app/log/command"
	handlerService "github.com/xtls/xray-core/app/proxyman/command"
	routingService "github.com/xtls/xray-core/app/router/command"
	statsService "github.com/xtls/xray-core/app/stats/command"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type XrayController struct {
	HsClient handlerService.HandlerServiceClient
	SsClient statsService.StatsServiceClient
	LsClient loggerService.LoggerServiceClient
	RsClient routingService.RoutingServiceClient
	CmdConn  *grpc.ClientConn
}

func (xrayCtl *XrayController) Init(addr string) (err error) {
	xrayCtl.CmdConn, err = grpc.Dial(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	xrayCtl.HsClient = handlerService.NewHandlerServiceClient(xrayCtl.CmdConn)
	xrayCtl.SsClient = statsService.NewStatsServiceClient(xrayCtl.CmdConn)
	xrayCtl.LsClient = loggerService.NewLoggerServiceClient(xrayCtl.CmdConn)
	xrayCtl.RsClient = routingService.NewRoutingServiceClient(xrayCtl.CmdConn)
	return
}

func (xrayCtl *XrayController) Close() (err error) {
	err = xrayCtl.CmdConn.Close()
	return
}
