package utility

import (
	"strings"

	"github.com/gogf/gf/v2/util/gconv"
)

func SplitUserPwd(x string) (user, pwd string) {
	a := strings.SplitN(x, ":", 2)
	if len(a) == 0 {
		user = ""
		pwd = ""
		return
	} else if len(a) == 1 {
		user = strings.TrimSpace(a[0])
		pwd = ""
		return
	} else {
		user = strings.TrimSpace(a[0])
		pwd = strings.TrimSpace(a[1])
		return
	}
}

func SplitItems(x string) []string {
	x = strings.TrimSpace(x)
	a := strings.Split(x, ",")
	for i := 0; i < len(a); i++ {
		a[i] = strings.TrimSpace(a[i])
	}
	return a
}

func SplitRange(x string) (r1 string, r2 string) {
	a := strings.SplitN(x, "-", 2)
	if len(a) == 0 {
		r1 = ""
		r2 = ""
		return
	} else if len(a) == 1 {
		r1 = strings.TrimSpace(a[0])
		r2 = strings.TrimSpace(a[0])
		return
	} else {
		r1 = strings.TrimSpace(a[0])
		r2 = strings.TrimSpace(a[1])
		return
	}
}

func SplitFallback(x string) (dest string, xver int) {
	a := strings.SplitN(x, "|", 2)
	if len(a) == 0 {
		dest = ""
		xver = 0
		return
	} else if len(a) == 1 {
		dest = strings.TrimSpace(a[0])
		xver = 0
		return
	} else {
		dest = strings.TrimSpace(a[0])
		xver = gconv.Int(a[1])
		return
	}
}
