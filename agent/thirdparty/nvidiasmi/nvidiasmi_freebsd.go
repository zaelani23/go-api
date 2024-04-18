// +build freebsd

package nvidiasmi

import (
	"github.com/zaelani23/go-api/agent/agent/config"
	"os"
)

const (
	DOCLINELIMIT = 50000
)

var (
	nvidiaexe              = "nvidia-smi"
	conf                   = config.GetConfig()
	NvidiaEnabled     bool = false
	NvidiaInitialized bool = false
	proc              *os.Process
)

func InitNvidiaWatch() {

}
