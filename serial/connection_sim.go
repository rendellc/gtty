package serial

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

type connectionSim struct {
	period  time.Duration
	rxChan  chan string
	txChan  chan string
}

func SimulateConnection(period time.Duration) Connection {
	c := connectionSim{
		period:  period,
		rxChan:  make(chan string),
		txChan:  make(chan string),
	}

	return c
}

func (c connectionSim) Start() error {
	log.Printf("Starting connection")
	go c.simulateWithLines()
	return nil
}

func (c connectionSim) Close() {}

func (c connectionSim) GetReceiver() Receiver {
	receiver := Receiver{
		channel: new(<-chan string),
	}
	*receiver.channel = c.rxChan

	return receiver
}

func (c connectionSim) GetTransmitter() Transmitter {
	transmitter := Transmitter{
		channel: new(chan<- string),
	}
	*transmitter.channel = c.rxChan

	return transmitter
}

func (c *connectionSim) simulateWithLines() {
	i := 0
	data1 := 1.3
	for {
		select {
		case msg := <-c.txChan:
			c.rxChan <- "echo: " + msg

		default:
			data1 += rand.NormFloat64() * 0.01
			line := fmt.Sprintf("%d, %.3f", i, data1)

			log.Printf("Sending '%v' to rxChan", line)
			c.rxChan <- line

			rateMillis := 1 / float64(c.period.Milliseconds())
			sleepMillis := rand.ExpFloat64() / rateMillis
			d := time.Duration(int64(sleepMillis) * int64(time.Millisecond))
			time.Sleep(d)
			i += 1
		}
	}
}
