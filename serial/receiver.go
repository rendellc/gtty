package serial

type Receiver struct {
	channel <-chan string
}

func (r Receiver) Get() string {
	return <-r.channel
}
