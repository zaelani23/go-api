package udp

import (
	"fmt"

	"github.com/whatap/go-api/common/io"
	"github.com/whatap/go-api/common/util/stringutil"
	"github.com/whatap/go-api/common/util/urlutil"
)

type UdpTxHttpcPack struct {
	AbstractPack
	// Pack
	Url          string
	Mcallee      int64
	ErrorType    string
	ErrorMessage string
	Stack        string

	//Processing data
	HttpcURL *urlutil.URL
}

func NewUdpTxHttpcPack() *UdpTxHttpcPack {
	p := new(UdpTxHttpcPack)
	p.Ver = UDP_PACK_VERSION
	p.AbstractPack.Flush = false
	return p
}
func NewUdpTxHttpcPackVer(ver int32) *UdpTxHttpcPack {
	p := new(UdpTxHttpcPack)
	p.Ver = ver
	p.AbstractPack.Flush = false
	return p
}
func (this *UdpTxHttpcPack) GetPackType() uint8 {
	return TX_HTTPC
}

func (this *UdpTxHttpcPack) ToString() string {
	return fmt.Sprint(this.AbstractPack.ToString(), ",Url=", this.Url, ",Callee=", this.Mcallee)
}

func (this *UdpTxHttpcPack) Clear() {
	this.AbstractPack.Clear()
	this.AbstractPack.Flush = false

	this.Url = ""
	this.Mcallee = 0
	this.ErrorType = ""
	this.ErrorMessage = ""
	this.Stack = ""

	this.HttpcURL = nil

}

func (this *UdpTxHttpcPack) Write(dout *io.DataOutputX) {
	this.AbstractPack.Write(dout)
	dout.WriteTextShortLength(this.Url)

	if this.Ver > 40000 {
		// Batch
	} else if this.Ver > 30000 {
		// Dotnet
		dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(this.Mcallee))
	} else if this.Ver > 20000 {
		// Python
		dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(this.Mcallee))
	} else {
		if this.Ver >= 10105 {
			dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(this.Mcallee))
			dout.WriteTextShortLength(this.ErrorType)
			dout.WriteTextShortLength(this.ErrorMessage)
			dout.WriteTextShortLength(this.Stack)
		} else if this.Ver >= 10102 {
			dout.WriteTextShortLength(stringutil.ParseStringZeroToEmpty(this.Mcallee))
		}
	}
}

func (this *UdpTxHttpcPack) Read(din *io.DataInputX) {
	this.AbstractPack.Read(din)

	this.Url = din.ReadTextShortLength()

	if this.Ver > 40000 {
		// Batch
	} else if this.Ver > 30000 {
		// Dotnet
		this.Mcallee = stringutil.ParseInt64(din.ReadTextShortLength())
	} else if this.Ver > 20000 {
		// Python
		this.Mcallee = stringutil.ParseInt64(din.ReadTextShortLength())
	} else {
		// PHP
		if this.Ver >= 10105 {
			this.Mcallee = stringutil.ParseInt64(din.ReadTextShortLength())
			this.ErrorType = din.ReadTextShortLength()
			this.ErrorMessage = din.ReadTextShortLength()
			this.Stack = din.ReadTextShortLength()
		} else if this.Ver >= 10102 {
			this.Mcallee = stringutil.ParseInt64(din.ReadTextShortLength())
		}
	}
}

func (this *UdpTxHttpcPack) Process() {
	this.HttpcURL = urlutil.NewURL(this.Url)
	if this.Ver > 40000 {
		// Batch
	} else if this.Ver > 30000 {
		// Dotnet
	} else if this.Ver > 20000 {
		// Python
	} else {
		// PHP
	}
}
