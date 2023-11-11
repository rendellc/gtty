package serial

import (
	"fmt"
	"log"
	"time"

	"github.com/dsyx/serialport-go"
)

type connection struct {
	device    string
	config    serialport.Config
	txNewline string
	conn      *serialport.SerialPort
	buf       []byte
	sp        *serialport.SerialPort
	rxChan    chan string
	txChan    chan string
	rwChan    chan readWriterCommand
}

type readWriterCommand int

const (
	rwCmdStop readWriterCommand = iota
)

func CreateConnection(config Config) connection {
	parity := serialport.PN
	if config.Parity == "none" {
		parity = serialport.PN
	} else if config.Parity == "odd" {
		parity = serialport.PE
	} else if config.Parity == "even" {
		parity = serialport.PE
	} else if config.Parity == "mark" {
		parity = serialport.PM
	} else if config.Parity == "space" {
		parity = serialport.PS
	}

	serialconfig := serialport.Config{
		BaudRate: config.BaudRate,
		DataBits: config.DataBits,
		StopBits: config.StopBits,
		Parity:   parity,
		Timeout:  config.Timeout,
	}

	txNewline := ""
	if config.TransmitNewline == "auto" {
		txNewline = "\n"
	} else if config.TransmitNewline == "lf" {
		txNewline = "\n"
	} else if config.TransmitNewline == "cr+lf" {
		txNewline = "\r\n"
	}

	return connection{
		device:    config.Device,
		config:    serialconfig,
		txNewline: txNewline,
		conn:      nil,
		buf:       make([]byte, 256),
	}
}

func (c *connection) Start() (Receiver, Transmitter, error) {
	receiver := Receiver{}
	transmitter := Transmitter{}

	if c.rxChan != nil && c.txChan != nil {
		return receiver, transmitter, fmt.Errorf("already running")
	}
	var err error
	c.sp, err = serialport.Open(c.device, c.config)
	if err != nil {
		return receiver, transmitter, err
	}

	c.rxChan = make(chan string)
	c.txChan = make(chan string)
	c.rwChan = make(chan readWriterCommand)
	receiver.channel = c.rxChan
	transmitter.channel = c.txChan

	go c.serialReadWriter()

	return receiver, transmitter, nil
}

func (c *connection) serialReadWriter() {
	rxChanRaw := make(chan string)
	go c.serialReader(rxChanRaw)

	for {
		select {
		case cmd := <-c.txChan:
			err := c.writeSerial(cmd)
			if err != nil {
				log.Printf("Error writing to serial: %s", err.Error())
			}
		case data := <-rxChanRaw:
			c.rxChan <- data
		case rwCmd := <-c.rwChan:
			log.Printf("ReadWriterCommand received: %v", rwCmd)
			break
		}
	}
}

func (h *connection) serialReader(out chan<- string) {
	for {
		readString := h.readSerial()
		if len(readString) == 0 {
			d := 1000 * time.Millisecond
			time.Sleep(d)
			continue
		}
		out <- readString
	}
}

func (c connection) readSerial() string {
	if c.conn == nil {
		log.Fatalf("ReadSerial called, but connection is nil")
	}

	n, _ := c.conn.Read(c.buf)
	return string(c.buf[:n])
}

func (c *connection) writeSerial(cmd string) error {
	if c.sp == nil {
		log.Fatalf("WriteSerial called, but serialport is nil")
	}

	buf := []byte(cmd + c.txNewline)
	n, err := c.sp.Write(buf)
	if err != nil {
		return fmt.Errorf("error writing to serial: %s", err.Error())
	}
	if n < len(buf) {
		return fmt.Errorf("incomplete write, wrote %d out of %d bytes", n, len(buf))
	}

	return nil
}

func (c *connection) Close() {
	if c.sp == nil {
		return
	}

	c.sp.Close()
}

// func mergeSerialDataToLines(serialData <-chan string, serialLine chan<- string) {
// 	partialLines := ""
// 	for {
// 		select {
// 		case s := <-serialData:
// 			partialLines += s
// 		}
// 
// 		lines := strings.Split(partialLines, "\n")
// 		for i, line := range lines {
// 			if i < len(lines)-1 {
// 				line = strings.Trim(line, "\r")
// 				serialLine <- line
// 			}
// 		}
// 		partialLines = lines[len(lines)-1]
// 	}
// }
