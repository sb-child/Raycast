package utility

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
)

func GetSubscribe(ctx context.Context, cfg *gjson.Json) []*gjson.Json {
	link := cfg.Get("link", "").String()
	file := cfg.Get("file", "").String()
	subName := cfg.Get("name", "Subscribe").String()
	ignoreName := cfg.Get("ignoreName", []string{}).Strings()
	r := make([]*gjson.Json, 0)
	// ignoreAddr := cfg.Get("ignoreAddr", []string{}).Strings()
	var content []byte
	if len(file) != 0 {
		content = gfile.GetBytes(file)
	} else if len(link) != 0 {
		content = g.Client().GetBytes(ctx, link)
	} else {
		return nil
	}
	if content == nil {
		return nil
	}
	clashSub, err := gjson.LoadYaml(content)
	if err != nil {
		g.Log().Errorf(ctx, "GetSubscribe: file type not supported")
		return nil
	}
	proxies := clashSub.GetJsons("proxies")
	if proxies == nil {
		g.Log().Errorf(ctx, "GetSubscribe: no proxies found")
		return nil
	}
	for _, v := range proxies {
		skip := false
		for _, vv := range ignoreName {
			if strings.Contains(v.Get("name").String(), vv) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		switch v.Get("type", "").String() {
		case "trojan":
			n := gjson.New(g.Map{
				"trojan": g.Map{
					"name":    subName + " > " + v.Get("name").String(),
					"through": "0.0.0.0",
					"server": v.Get("server", "").String() + ":" +
						v.Get("port", 0).String(),
					"user": v.Get("password", "").String(),
					"security": g.Map{
						"tls":  v.Get("sni", 0).String(),
						"ver":  "1.2-1.3",
						"alpn": "h2,http/1.1",
					},
				},
			})
			r = append(r, n)
		}
	}
	return r
}
