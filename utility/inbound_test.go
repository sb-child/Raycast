package utility_test

import (
	"raycast/utility"
	"testing"

	"github.com/gogf/gf/v2/encoding/gjson"
)

func TestHttpProxyInbound_Json(t *testing.T) {
	// Positive test case with one user
	inbound := utility.HttpProxyInbound{
		Listen: "127.0.0.1:8080",
		Users:  []string{"user1:pass1"},
		Tag:    "http-inbound",
	}
	expected := gjson.New(`{"protocol":"http","listen":"127.0.0.1","port":8080,"settings":{"accounts":[{"user":"user1","pass":"pass1"}]},"tag":"http-inbound"}`)
	actual := inbound.Json().String()
	if actual != expected.String() {
		t.Errorf("Json() returned unexpected result:\nexpected: %s\nactual: %s", expected, actual)
	}
	// Positive test case with multiple users
	inbound = utility.HttpProxyInbound{
		Listen: "0.0.0.0:80",
		Users:  []string{"user1:pass1", "user2:pass2"},
		Tag:    "http-inbound",
	}
	expected = gjson.New(`{"protocol":"http","listen":"0.0.0.0","port":80,"settings":{"accounts":[{"user":"user1","pass":"pass1"},{"user":"user2","pass":"pass2"}]},"tag":"http-inbound"}`)
	actual = inbound.Json().String()
	if actual != expected.String() {
		t.Errorf("Json() returned unexpected result:\nexpected: %s\nactual: %s", expected, actual)
	}
	// Negative test case with invalid listen address
	inbound = utility.HttpProxyInbound{
		Listen: "invalid_address",
		Users:  []string{},
		Tag:    "",
	}
	expected = gjson.New(`{"listen":"","port":0,"protocol":"http","settings":{},"tag":""}`)
	actual = inbound.Json().String()
	if actual != expected.String() {
		t.Errorf("Json() returned unexpected result:\nexpected: %s\nactual: %s", expected, actual)
	}
}

func TestSocksProxyInbound_Json(t *testing.T) {
	// Positive test case with one user and UDP enabled
	inbound := utility.SocksProxyInbound{
		Listen: "127.0.0.1:1080",
		Users:  []string{"user1:pass1"},
		Udp:    true,
		Tag:    "socks-inbound",
	}
	expected := gjson.New(`{"protocol":"socks","listen":"127.0.0.1","port":1080,"settings":{"auth":"password","udp":true,"ip":"127.0.0.1","accounts":[{"user":"user1","pass":"pass1"}]},"tag":"socks-inbound"}`)
	actual := inbound.Json().String()
	if actual != expected.String() {
		t.Errorf("Json() returned unexpected result:\nexpected: %s\nactual: %s", expected, actual)
	}
	// Positive test case with multiple users and UDP disabled
	inbound = utility.SocksProxyInbound{
		Listen: "0.0.0.0:1080",
		Users:  []string{"user1:pass1", "user2:pass2"},
		Udp:    false,
		Tag:    "socks-inbound",
	}
	expected = gjson.New(`{"protocol":"socks","listen":"0.0.0.0","port":1080,"settings":{"auth":"password","udp":false,"ip":"0.0.0.0","accounts":[{"user":"user1","pass":"pass1"},{"user":"user2","pass":"pass2"}]},"tag":"socks-inbound"}`)
	actual = inbound.Json().String()
	if actual != expected.String() {
		t.Errorf("Json() returned unexpected result:\nexpected: %s\nactual: %s", expected, actual)
	}
	// Negative test case with invalid listen address
	inbound = utility.SocksProxyInbound{
		Listen: "invalid_address",
		Users:  []string{},
		Udp:    false,
		Tag:    "",
	}
	expected = gjson.New(`{"listen":"","port":0,"protocol":"socks","settings":{"auth":"noauth","udp":false,"ip":""},"tag":""}`)
	actual = inbound.Json().String()
	if actual != expected.String() {
		t.Errorf("Json() returned unexpected result:\nexpected: %s\nactual: %s", expected, actual)
	}
}
