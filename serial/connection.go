package serial

type Connection interface {
	Close()
	Start() (Receiver, Transmitter, error)
}

