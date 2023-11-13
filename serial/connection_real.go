package serial

import (
	"fmt"
	"log"
	"time"

	"github.com/dsyx/serialport-go"
)

type connectionReal struct {
	device    string
	config    serialport.Config
	txNewline string
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

func CreateConnection(config Config) Connection {
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

	return connectionReal{
		device:    config.Device,
		config:    serialconfig,
		txNewline: txNewline,
		sp:      nil,
		buf:       make([]byte, 256),
		rxChan: make(chan string),
		txChan: make(chan string),
		rwChan: make(chan readWriterCommand),
	}
}

func (c connectionReal) Start() error {
	if c.rxChan != nil && c.txChan != nil {
		return fmt.Errorf("already running")
	}
	var err error
	c.sp, err = serialport.Open(c.device, c.config)
	if err != nil {
		return err
	}

	go c.serialReadWriter()

	return nil
}

func (c connectionReal) GetReceiver() Receiver {
	receiver := Receiver{
		channel: new(<-chan string),
	}
	*receiver.channel = c.rxChan

	return receiver
}

func (c connectionReal) GetTransmitter() Transmitter {
	transmitter := Transmitter{
		channel: new(chan<- string),
	}
	*transmitter.channel = c.rxChan

	return transmitter
}


func (c *connectionReal) serialReadWriter() {
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

func (h *connectionReal) serialReader(out chan<- string) {
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

func (c connectionReal) readSerial() string {
	if c.sp == nil {
		log.Fatalf("ReadSerial called, but connection is nil")
	}

	n, _ := c.sp.Read(c.buf)
	return string(c.buf[:n])
}

func (c *connectionReal) writeSerial(cmd string) error {
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

func (c connectionReal) Close() {
	if c.sp == nil {
		return
	}

	c.sp.Close()
}
