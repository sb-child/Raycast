package utility

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

func GetCertificates(name string) []*gjson.Json {
	names := SplitItems(name)
	rr := make([]*gjson.Json, 0)
	for i := 0; i < len(names); i++ {
		cert := gjson.New(g.Config().MustGet(context.TODO(), "certificates."+names[i], ""))
		r := gjson.New(g.Map{
			"ocspStapling":    cert.Get("refresh", 3600).Int(),
			"oneTimeLoading":  !cert.Contains("refresh"),
			"usage":           "encipherment",
			"certificateFile": cert.Get("cert", "").String(),
			"keyFile":         cert.Get("key", "").String(),
		})
		rr = append(rr, r)
	}
	return rr
}

// return nil if not detected
func AutoSecurityJson(c *gjson.Json, inbound bool) *gjson.Json {
	if c == nil {
		return nil
	}
	switch {
	case c.Contains("reality"):
		x := RealitySecurity{}
		return x.FromCfg(c, inbound).Json()
	case c.Contains("tls"):
		x := TlsSecurity{}
		return x.FromCfg(c, inbound).Json()
	}
	return nil
}

type RealitySecurity struct {
	Domains      string
	PriKey       string
	PubKey       string
	Secrets      []string
	FingerPrint  string
	Fallback     string
	VersionRange string
	TimeDiff     int
	Spider       string
	ForInbound   bool
}

func (x *RealitySecurity) FromCfg(c *gjson.Json, inbound bool) *RealitySecurity {
	x.Domains = c.Get("reality", "").String()
	x.PriKey = c.Get("priKey", "").String()
	x.PubKey = c.Get("pubKey", "").String()
	x.Secrets = c.Get("secret", "").Strings()
	x.FingerPrint = c.Get("fingerprint", "").String()
	x.Fallback = c.Get("fallback", "").String()
	x.VersionRange = c.Get("ver", "").String()
	x.TimeDiff = c.Get("timeDiff", 0).Int()
	x.Spider = c.Get("spider", "").String()
	x.ForInbound = inbound
	return x
}

func (x *RealitySecurity) Json() *gjson.Json {
	fbDest, fbVer := SplitFallback(x.Fallback)
	minVer, maxVer := SplitRange(x.VersionRange)
	r := gjson.New(g.Map{
		"network":  "tcp",
		"security": "reality",
		"realitySettings": g.Map{
			"show":         false,
			"dest":         fbDest,
			"xver":         fbVer,
			"privateKey":   x.PriKey,
			"publicKey":    x.PubKey,
			"minClientVer": minVer,
			"maxClientVer": maxVer,
			"maxTimeDiff":  x.TimeDiff,
			"fingerprint":  x.FingerPrint,
			"spiderX":      x.Spider,
		},
	})
	if x.ForInbound {
		r.Set("realitySettings.serverNames", SplitItems(x.Domains))
		r.Set("realitySettings.shortIds", x.Secrets)
	} else {
		r.Remove("realitySettings.dest")
		r.Remove("realitySettings.xver")
		r.Remove("realitySettings.minClientVer")
		r.Remove("realitySettings.maxClientVer")
		r.Set("realitySettings.serverName", x.Domains)
		if len(x.Secrets) == 0 {
			r.Set("realitySettings.shortId", "")
		} else {
			r.Set("realitySettings.shortId", x.Secrets[0])
		}
	}
	return r
}

type TlsSecurity struct {
	Domain             string
	VersionRange       string
	SniCheck           bool
	Alpn               string
	CipherSuites       string
	Fingerprint        string
	Certificates       string
	CertificatesObject *gjson.Json
	ForInbound         bool
}

func (x *TlsSecurity) FromCfg(c *gjson.Json, inbound bool) *TlsSecurity {
	x.Domain = c.Get("tls", "").String()
	x.VersionRange = c.Get("ver", "1.2-1.3").String()
	x.SniCheck = c.Get("sniCheck", true).Bool()
	x.Alpn = c.Get("alpn", "").String()
	x.CipherSuites = c.Get("cipherSuites", "").String()
	x.Fingerprint = c.Get("fingerprint", "").String()
	x.Certificates = c.Get("certificates", "").String()
	x.ForInbound = inbound
	return x
}

func (x *TlsSecurity) Json() *gjson.Json {
	minVer, maxVer := SplitRange(x.VersionRange)
	r := gjson.New(g.Map{
		"network":  "tcp",
		"security": "tls",
		"tlsSettings": g.Map{
			"serverName":                       x.Domain,
			"rejectUnknownSni":                 x.SniCheck,
			"allowInsecure":                    false,
			"alpn":                             SplitItems(x.Alpn),
			"minVersion":                       minVer,
			"maxVersion":                       maxVer,
			"cipherSuites":                     x.CipherSuites,
			"certificates":                     GetCertificates(x.Certificates),
			"disableSystemRoot":                false,
			"enableSessionResumption":          false,
			"fingerprint":                      x.Fingerprint,
			"pinnedPeerCertificateChainSha256": nil,
		},
	})
	if len(x.CipherSuites) == 0 {
		r.Remove("tlsSettings.cipherSuites")
	}
	if x.ForInbound {
		// r.Set("realitySettings.serverNames", SplitItems(x.Domains))
		// r.Set("realitySettings.shortIds", x.Secrets)
	} else {
		r.Remove("tlsSettings.certificates")
	}
	return r
}
