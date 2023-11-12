package serial

type Connection interface {
	Close()
	Start() error
	GetReceiver() Receiver
	GetTransmitter() Transmitter
}

