package util

import (
	"os/exec"
	"strings"
)

func NewUUID() string {
	if o, e := exec.Command("uuidgen").Output(); e != nil {
		return ""
	} else {
		return strings.Trim(string(o), "\n")
	}
}
