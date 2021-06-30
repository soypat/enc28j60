module github.com/soypat/enc28j60

go 1.16

require (
	github.com/soypat/ether-swtch v0.5.0
	github.com/soypat/net v0.2.0
	tinygo.org/x/drivers v0.16.0
)

replace github.com/soypat/net => ../net

replace github.com/soypat/ether-swtch => ../ether-swtch
