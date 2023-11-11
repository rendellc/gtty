package serial

import (
	"time"

	"github.com/dsyx/serialport-go"
)

var validBaudrates []int = []int{
	serialport.BR110,
	serialport.BR300,
	serialport.BR600,
	serialport.BR1200,
	serialport.BR2400,
	serialport.BR4800,
	serialport.BR9600,
	serialport.BR14400,
	serialport.BR19200,
	serialport.BR38400,
	serialport.BR57600,
	serialport.BR115200,
	serialport.BR128000,
	serialport.BR256000,
}

var validDatabits []int = []int{
	serialport.DB5,
	serialport.DB6,
	serialport.DB7,
	serialport.DB8,
}

var validStopbits []int = []int{
	serialport.SB1,
	serialport.SB1_5,
	serialport.SB2,
}

var validParities []string = []string{
	"none",
	"odd",
	"even",
	"mark",
	"space",
}

var validLineendings []string = []string{
	"auto",
	"lf",
	"cr+lf",
}

type Config struct {
	Device       string
	BaudRate         int
	DataBits         int
	StopBits         int
	Parity           string
	Timeout          time.Duration
	TransmitNewline  string
	SimulationEnable bool
}
