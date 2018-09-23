package common

import (
	"os/exec"
	"strings"
)

func UUID() string {
	if o, e := exec.Command("uuidgen").Output(); e != nil {
		return ""
	} else {
		return strings.Trim(string(o), "\n")
	}
}

func Shorten(uuid string) string {
	return "#" + uuid[:4]
}

func Btoi(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}
