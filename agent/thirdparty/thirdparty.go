package thirdparty

import (
	"github.com/zaelani23/go-api/agent/thirdparty/nvidiasmi"
)

func StartAll() {
	nvidiasmi.InitNvidiaWatch()
}
