package main

import (
	"machine"
	"time"

	"github.com/soypat/net"

	"github.com/soypat/enc28j60"
	swtch "github.com/soypat/ether-swtch"
)

func main() {
	var (
		// SPI Chip select pin. Can be any Digital pin
		spiCS = machine.D53
		MAC   = net.HardwareAddr{0xDE, 0xAD, 0xBE, 0xEF, 0xFE, 0xFF}
		MyIP  = net.IP{192, 168, 1, 5} //static setup is the only one available
	)

	// Configure writer/reader integrated circuit.
	dev := enc28j60.New(spiCS, machine.SPI0)

	err := dev.Init(MAC)
	if err != nil {
		println(err.Error())
	}
	const okHeader = "HTTP/1.0 200 OK\r\nContent-Type: text/html\r\nPragma: no-cache\r\n\r\n"
	timeout := time.Second * 1
	// Spin up HTTP server which responds with "Hello world!"
	swtch.HTTPListenAndServe(dev, MAC, MyIP, timeout, func(URL []byte) (response []byte) {
		return []byte(okHeader + "Hello world!")
	}, printNonNilErr)
}

func printNonNilErr(err error) {
	if err != nil {
		println(err.Error())
	}
}

/*
Arduino Uno SPI pins:
sck:  PB5, // is D13
sdo:  PB3, // MOSI is D11
sdi:  PB4, // MISO is D12
cs:   PB2} // CS  is D10

Arduino MEGA 2560 SPI pins as taken from tinygo library (online documentation seems to be wrong at times)
SCK: PB1 == D52
MOSI(sdo): PB2 == D51
MISO(sdi): PB3 == D50
CS: PB0 == D53
*/
