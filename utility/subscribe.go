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
	ignoreName := cfg.Get("ignoreName", []string{}).Strings()
	// ignoreAddr := cfg.Get("ignoreAddr", []string{}).Strings()
	content := make([]byte, 0)
	if len(file) != 0 {
		content = gfile.GetBytes(file)
	} else if len(link) != 0 {
		content = g.Client().GetBytes(ctx, link)
	} else {
		return nil
	}
	if content == nil || len(content) == 0 {
		return nil
	}
	clashSub, err := gjson.LoadYaml(content)
	if err != nil {
		g.Log().Errorf(ctx, "GetSubscribe: file type not supported")
		return nil
	}
	proxies := clashSub.GetJsons("proxies")
	if proxies == nil || len(proxies) == 0 {
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

		}
	}
}
