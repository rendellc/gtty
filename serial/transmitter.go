package serial

type Transmitter struct {
	channel *chan<- string
}

func (t Transmitter) Send(data string) {
	*t.channel <- data
}
