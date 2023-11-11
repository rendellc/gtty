package serial

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

type connectionSim struct {
	period   time.Duration
	rxChan chan string
	txChan chan string
}

func SimulateConnection(period time.Duration) Connection {
	return connectionSim{
		period:  period,
		rxChan: make(chan string),
		txChan: make(chan string),
	}
}

func (c connectionSim) Start() (Receiver, Transmitter, error) {
	receiver := Receiver{
		channel: new(<-chan string),
	}
	transmitter := Transmitter{
		channel: new(chan<- string),
	}

	*receiver.channel = c.rxChan
	*transmitter.channel = c.txChan
	go c.simulateWithLines()

	return receiver, transmitter, nil
}

func (c connectionSim) Close() {}


func (c *connectionSim) simulateWithLines() {
	i := 0 
	data1 := 1.3 
	for {
		select {
		case msg := <-c.txChan:
			c.rxChan <- "echo: " + msg


		default:
			data1 += rand.NormFloat64() * 0.001
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

