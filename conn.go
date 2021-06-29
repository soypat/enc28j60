package enc28j60

import (
	"io"

	swtch "github.com/soypat/ether-swtch"
)

type Packet struct {
	*Dev
	cursor uint16
	end    uint16
}

func (e *Dev) NextPacket() (swtch.Reader, error) {
	var err error
	p := &Packet{Dev: e}
	for e.read(EPKTCNT) == 0 { // loop until a packet is received.
	}
	// Set the read pointer to the start of the next packet
	e.write16(ERDPTL, e.nextPacketPtr)
	p.cursor = e.nextPacketPtr // Packet reader

	e.readBuffer(e.dummy[:])
	e.nextPacketPtr = uint16(e.dummy[0]) + uint16(e.dummy[1])<<8
	// read the packet length (see datasheet page 43)
	plen := uint16(e.dummy[2]) + uint16(e.dummy[3])<<8 - 4 //remove the CRC count (minus 4)
	p.end = p.cursor + plen
	// read the receive status (see datasheet page 43)
	rxstat := uint16(e.dummy[4]) + uint16(e.dummy[5])<<8
	// check CRC and symbol errors (see datasheet page 44, table 7-3):
	// The ERXFCON.CRCEN is set by default. Normally we should not
	// need to check this.
	if (rxstat & 0x80) == 0 {
		err = ErrCRC
	}
	return p, err
}
func (p *Packet) Discard() error {
	if p.cursor != p.end {
		p.cursor = p.end
		p.writeOp(BIT_FIELD_SET, ECON2, ECON2_PKTDEC)
	}
	return nil
}
func (p *Packet) Read(buff []byte) (n uint16, err error) {
	// total remaining packet length
	plen := p.end - p.cursor
	if plen == 0 {
		return 0, io.EOF
	}
	// Limit retreive length if Total packet length is greater than buffer length
	if plen > uint16(len(buff)) {
		plen = uint16(len(buff))
	}
	// copy the packet from the receive buffer
	p.readBuffer(buff[:plen])
	p.cursor += plen
	// Move the RX read pointer to where we ended reading
	p.write16(ERXRDPTL, p.cursor)

	if p.cursor == p.end { // minus CRC length
		// decrement packet counter to indicate we are done with it.
		p.writeOp(BIT_FIELD_SET, ECON2, ECON2_PKTDEC)
		err = io.EOF
	}
	return plen, err
}

func (e *Dev) Write(buff []byte) (uint16, error) {
	plen := uint16(len(buff))
	if plen+e.tcursor > MAX_FRAMELEN {
		e.tcursor = 0
		return 0, ErrBufferSize
	}

	e.write16(EWRPTL, TXSTART_INIT+e.tcursor)
	if e.tcursor == 0 {
		// write per-packet control byte (0x00 means use macon3 settings)
		e.writeOp(WRITE_BUF_MEM, 0, 0x00)
		e.tcursor++ // WBM spurious increment
	}
	// copy the packet into the transmit buffer
	e.writeBuffer(buff)
	e.tcursor += plen
	return plen, nil
}

func (e *Dev) Flush() error {
	dbp("send response")
	e.write16(ETXNDL, TXSTART_INIT+e.tcursor-1) // subtract WBM spurious increment
	// send the contents of the transmit buffer onto the network
	e.writeOp(BIT_FIELD_SET, ECON1, ECON1_TXRTS)
	// Reset the transmit logic problem. See Rev. B4 Silicon Errata point 12.
	if e.read(EIR)&EIR_TXERIF != 0 {
		e.writeOp(BIT_FIELD_CLR, ECON1, ECON1_TXRTS)
	}
	e.tcursor = 0
	return nil
}
