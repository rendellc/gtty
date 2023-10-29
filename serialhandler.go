package main

import (
	"fmt"

	"github.com/dsyx/serialport-go"
)

type SerialConnection struct {
	Device string
	Config serialport.Config
	handle *serialport.SerialPort
}

func (conn *SerialConnection) Open() error {
	if conn.handle != nil {
		return fmt.Errorf("Connection already has a serial port handle")
	}
	sp, err := serialport.Open(conn.Device, conn.Config)
	if err != nil {
		return err
	}
	defer sp.Close()
	conn.handle = sp

	return nil
}

// conn := SerialConnection{
// 	Device: "COM8",
// 	Config: serialport.Config{
// 		BaudRate: serialport.BR9600,
// 		DataBits: serialport.DB8,
// 		StopBits: serialport.SB1,
// 		Parity:   serialport.PO,
// 		Timeout:  100 * time.Millisecond,
// 	},
// }
// if err := conn.Open(); err != nil {
// 	log.Panic(err)
// }
// log.Println("Connection is open")

// Implement reader interface
func (conn SerialConnection) Read(p []byte) (int, error) {
	if conn.handle == nil {
		return 0, fmt.Errorf("serialconnection: dont have a serialport handle")
	}
	return 0, nil

	//if err != nil {
	//	log.Fatalln(err)
	//}
	//defer sp.Close()

	//fmt.Println("Device is open")

	//buf := make([]byte, 1024)
	//for {
	//	n, _ := sp.Read(buf)
	//	fmt.Printf("%s", string(buf[:n]))
	//}
}
