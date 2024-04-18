package control

import (
	//"log"
	"os"
	"strings"

	"math"
	"regexp"
	"runtime/debug"

	"github.com/magiconair/properties"
	"github.com/zaelani23/go-api/agent/agent/active"
	"github.com/zaelani23/go-api/agent/agent/config"
	"github.com/zaelani23/go-api/agent/agent/counter"
	"github.com/zaelani23/go-api/agent/agent/counter/meter"
	"github.com/zaelani23/go-api/agent/agent/data"
	"github.com/zaelani23/go-api/agent/agent/secure"
	"github.com/zaelani23/go-api/agent/agent/topology"

	//"github.com/zaelani23/go-api/agent/dotnet"

	//"github.com/zaelani23/go-api/agent/extension"
	"github.com/whatap/golib/lang"
	"github.com/whatap/golib/lang/pack"
	"github.com/whatap/golib/lang/value"
	"github.com/whatap/golib/util/cmdutil"
	"github.com/zaelani23/go-api/agent/net"
	"github.com/zaelani23/go-api/agent/util/logutil"
)

var start bool = false

func InitControlHandler() {
	if start == false {
		start = true
		go runControl()
	}
}

func runControl() {
	net.InitReceiver()
	for {
		// DEBUG goroutine log
		//logutil.Println("ControlHandler runControl")

		p := <-net.RecvBuffer

		switch p.GetPackType() {
		case pack.PACK_PARAMETER:
			process(p.(*pack.ParamPack))
		default:
		}
	}
}

var secuMaster = secure.GetSecurityMaster()

func process(p *pack.ParamPack) {
	// for 문이 종료 되지 않도록 Recover
	defer func() {
		if r := recover(); r != nil {
			logutil.Println("WA10601", " Recover ", r)
			logutil.Println(string(debug.Stack()))
		}
	}()
	conf := config.GetConfig()

	switch p.Id {
	case net.MODULE_DEPENDENCY:
		if conf.DebugControlEnabled {
			logutil.Infoln("[DEBUG]", "MODULE_DEPENDENCY")
		}
		// PHP 는 command 로 환경 가져옴
		// if conf.AppType == lang.APP_TYPE_PHP || conf.AppType == lang.APP_TYPE_BSM_PHP {
		// 	m := value.NewMapValue()
		// 	moduleinfo := cmdutil.GetPHPModuleInfo()
		// 	for k, v := range moduleinfo {
		// 		m.PutString(k, v)
		// 	}
		// 	p.Put("module", m)
		// } else if conf.AppType == lang.APP_TYPE_GO || conf.AppType == lang.APP_TYPE_BSM_GO {

		// } else {
		// 	extension.SendUdpSession(p.Id, p.Request, []int{})
		// 	return
		// }

	case net.GET_ENV:
		if conf.DebugControlEnabled {
			logutil.Infoln("[DEBUG]", "GET_ENV")
		}
		m := value.NewMapValue()

		for _, e := range os.Environ() {
			pair := strings.Split(e, "=")
			m.PutString(pair[0], pair[1])
		}
		if counter.CounterStaticAgentBootInfo != nil {
			m.PutString("whatap.agent_version", counter.CounterStaticAgentBootInfo.GetString("whatap.agent_version"))
			m.PutString("whatap.agent_buildno", counter.CounterStaticAgentBootInfo.GetString("whatap.agent_buildno"))
		}

		// PHP
		if conf.AppType == lang.APP_TYPE_PHP || conf.AppType == lang.APP_TYPE_BSM_PHP {
			phpinfo := cmdutil.GetPHPInfo()
			for k, v := range phpinfo {
				m.PutString(k, v)
			}

		}

		// .NET
		// if conf.AppType == lang.APP_TYPE_DOTNET || conf.AppType == lang.APP_TYPE_BSM_DOTNET {
		// 	m.PutString("framework.version", dotnet.GetDotnetVersion())
		// 	m.PutString("iis.version", dotnet.GetIISVersion())
		// }

		// Golnag
		if conf.AppType == lang.APP_TYPE_GO || conf.AppType == lang.APP_TYPE_BSM_GO {
		}

		p.Put("env", m)

	case net.CONFIGURE_GET:
		if conf.DebugControlEnabled {
			logutil.Infoln("[DEBUG]", "CONFIGURE_GET")
		}
		path := config.GetConfFile()
		props := properties.MustLoadFile(path, properties.UTF8)
		m := value.NewMapValue()
		for _, key := range props.Keys() {
			match, _ := regexp.MatchString("^\\w", key)
			//if !match || strings.Index(key, "license") > -1 || strings.Index(key, "whatap.server.host") > -1 || strings.Index(key, "OID") > -1 {
			if !match || strings.Index(key, "OID") > -1 {
				continue
			}
			val, _ := props.Get(key)
			m.PutString(key, strings.Replace(val, "\\", "\\\\", -1))
		}
		p.SetMapValue(m)

	case net.SET_CONFIG:
		if conf.DebugControlEnabled {
			logutil.Infoln("[DEBUG]", "SET_CONFIG")
		}
		configmap := p.GetMap("config")
		if configmap != nil {
			keyValues := map[string]string{}
			keyEnumer := configmap.Keys()
			for keyEnumer.HasMoreElements() {
				key := keyEnumer.NextString()
				keyValues[key] = configmap.GetString(key)
			}
			config.SetValues(&keyValues)
		}

	case net.GET_ACTIVE_TRANSACTION_LIST:
		if conf.DebugControlEnabled {
			logutil.Infoln("[DEBUG]", "GET_ACTIVE_TRANSACTION_LIST")
		}
		m := active.GetActiveTxList()
		p.Put("active", m)
	case net.GET_ACTIVE_TRANSACTION_DETAIL:
		if conf.DebugControlEnabled {
			logutil.Infoln("[DEBUG]", "GET_ACTIVE_TRANSACTION_DETAIL")
		}
		//extension.SendUdpSession(p.Id, p.Request, []int{int(p.GetLong("thread_id")), int(p.GetLong("profile"))})
		return

	case net.AGENT_LOG_LIST:
		if conf.DebugControlEnabled {
			logutil.Infoln("[DEBUG]", "AGENT_LOG_LIST")
		}
		m := logutil.GetLogger().GetLogFiles()
		p.Put("files", m)
	case net.AGENT_LOG_READ:
		file := p.GetString("file")
		endpos := p.GetLong("pos")
		length := p.GetLong("length")
		length = int64(math.Min(float64(length), 8000))

		if conf.DebugControlEnabled {
			logutil.Infoln("[DEBUG]", "AGENT_LOG_READ ", file, ", ", endpos, ", ", length)
		}

		logData := logutil.GetLogger().Read(file, endpos, length)
		if logData != nil {
			p.PutLong("before", logData.Before)
			p.PutLong("next", logData.Next)
			p.PutString("text", logData.Text)
		}

	case net.AGENT_STAT:
		if conf.DebugControlEnabled {
			logutil.Infoln("[DEBUG]", "AGENT_STAT ")
		}

		m := meter.GetInstanceMeterSelf().GetMeterSelfStat()
		p.Put("stat", m)

	case net.GET_TOPOLOGY:
		if conf.DebugControlEnabled {
			logutil.Infoln("[DEBUG]", "GET_TOPOLOGY ")
		}
		node := topology.NewStatusDetector().Process()
		if node != nil {
			p.Put("node", value.NewBlobValue(node.ToBytes()))
		}

	}
	if conf.DebugControlEnabled {
		logutil.Infoln("[DEBUG]", "Control Send ")
	}
	data.SendResponse(p.ToResponse())
}
