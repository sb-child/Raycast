package utility

import "strings"

func SplitUserPwd(x string) (user, pwd string) {
	a := strings.SplitN(x, ":", 2)
	if len(a) == 0 {
		user = ""
		pwd = ""
		return
	} else if len(a) == 1 {
		user = a[0]
		pwd = ""
		return
	} else {
		user = a[0]
		pwd = a[1]
		return
	}
}
