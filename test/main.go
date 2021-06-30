package main

import (
	"machine"
	"unsafe"

	"github.com/soypat/net"

	"github.com/soypat/enc28j60"
	swtch "github.com/soypat/ether-swtch"
)

/* Arduino Uno SPI pins:
sck:  PB5, // is D13
sdo:  PB3, // MOSI is D11
sdi:  PB4, // MISO is D12
cs:   PB2} // CS  is D10
*/

/* Arduino MEGA 2560 SPI pins as taken from tinygo library (online documentation seems to be wrong at times)
SCK: PB1 == D52
MOSI(sdo): PB2 == D51
MISO(sdi): PB3 == D50
CS: PB0 == D53
*/

// Arduino uno CS Pin
// var spiCS = machine.D10 // on Arduino Uno

func main() {
	println("start")
	// SPI Chip select pin. Can be any Digital pin
	var spiCS = machine.D53
	// Inline declarations so not used as RAM
	var (
		MAC = net.HardwareAddr{0xDE, 0xAD, 0xBE, 0xEF, 0xFE, 0xFF}
		// MyIP = net.IP{192, 168, 1, 5} //static setup is the only one available
	)

	// 8MHz SPI clk for older than Rev 6 boards (See Rev. B4 Silicon Errata)
	machine.SPI0.Configure(machine.SPIConfig{Frequency: 8e6})

	e := enc28j60.New(spiCS, machine.SPI0)
	enc28j60.SDB = true
	err := e.Init(MAC)
	if err != nil {
		println(err.Error())
	}
	if e == nil { // WITHOUT THIS NIL CHECK THE PROGRAM GC's e??
		panic("nil device")
	}
	swtch.SDB = true
	c := swtch.NewTCPConn(e, nil, MAC)
	err = c.Decode()
	if err != nil {
		println(err.Error())
	}
	println(c.ARPv4.String())
	err = c.SendResponse()
	if err != nil {
		println(err.Error())
	}

	err = c.Decode()
	if err != nil {
		println(err.Error())
	}
	println(c.TCP.String())
}

func codeFromErrorUnsafeArduino(err error) uint8 {
	if err != nil {
		type eface struct {
			typ uintptr
			val uint8
		}
		ptr := unsafe.Pointer(&err)
		return (*eface)(ptr).val
	}
	return 0
}
