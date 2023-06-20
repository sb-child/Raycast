package utility

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

type CfgFramework struct {
	CfgLog       *gjson.Json   `json:"log"`
	CfgApi       *gjson.Json   `json:"api"`
	CfgDns       *gjson.Json   `json:"dns"`
	CfgRouting   *gjson.Json   `json:"routing"`
	CfgPolicy    *gjson.Json   `json:"policy"`
	CfgInbounds  []*gjson.Json `json:"inbounds"`
	CfgOutbounds []*gjson.Json `json:"outbounds"`
	CfgTransport *gjson.Json   `json:"transport"`
	CfgStats     *gjson.Json   `json:"stats"`
	CfgReverse   *gjson.Json   `json:"reverse"`
	CfgFakeDns   *gjson.Json   `json:"fakedns"`
}

func (x *CfgFramework) Init() {
	x.CfgLog = gjson.New(g.Map{})
	x.CfgApi = nil
	x.CfgDns = gjson.New(g.Map{})
	x.CfgRouting = gjson.New(g.Map{})
	x.CfgPolicy = gjson.New(g.Map{})
	x.CfgInbounds = make([]*gjson.Json, 0)
	x.CfgOutbounds = make([]*gjson.Json, 0)
	x.CfgTransport = nil
	x.CfgStats = gjson.New(g.Map{})
	x.CfgReverse = gjson.New(g.Map{})
	x.CfgFakeDns = nil
}

func (x *CfgFramework) Json() *gjson.Json {
	return gjson.New(x)
}

func (x *CfgFramework) Api(b bool) {
	if b {
		x.CfgApi = gjson.New(g.Map{})
		x.CfgApi.Set("tag", "api")
		x.CfgApi.Set("services", g.SliceStr{"HandlerService", "LoggerService", "StatsService"})
	} else {
		x.CfgApi = nil
	}
}

func (x *CfgFramework) Transport(b bool) {
	if b {
		x.CfgTransport = gjson.New(g.Map{})
	} else {
		x.CfgTransport = nil
	}
}

func (x *CfgFramework) Inbounds(a ...*gjson.Json) {
	if x.CfgInbounds == nil {
		x.CfgInbounds = []*gjson.Json{}
	}
	x.CfgInbounds = append(x.CfgInbounds, a...)
}

func (x *CfgFramework) Outbounds(a ...*gjson.Json) {
	if x.CfgOutbounds == nil {
		x.CfgOutbounds = []*gjson.Json{}
	}
	x.CfgOutbounds = append(x.CfgOutbounds, a...)
}
