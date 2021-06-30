package enc28j60

import (
	"io"

	swtch "github.com/soypat/ether-swtch"
)

type Packet struct {
	ic     *Dev
	cursor uint16
	end    uint16
}

func (d *Dev) NextPacket() (swtch.Reader, error) {
	dbp("NextPacket")
	var err error
	p := &Packet{ic: d}
	for d.read(EPKTCNT) == 0 { // loop until a packet is received.
	}
	// Set the read pointer to the start of the next packet
	d.write16(ERDPTL, d.nextPacketPtr)
	p.cursor = d.nextPacketPtr // Packet reader

	d.readBuffer(d.dummy[:])
	d.nextPacketPtr = uint16(d.dummy[0]) + uint16(d.dummy[1])<<8
	// read the packet length (see datasheet page 43)
	plen := uint16(d.dummy[2]) + uint16(d.dummy[3])<<8 - 4 //remove the CRC count (minus 4)
	p.end = p.cursor + plen
	// read the receive status (see datasheet page 43)
	rxstat := uint16(d.dummy[4]) + uint16(d.dummy[5])<<8
	// check CRC and symbol errors (see datasheet page 44, table 7-3):
	// The ERXFCON.CRCEN is set by default. Normally we should not
	// need to check this.
	if (rxstat & 0x80) == 0 {
		err = ErrCRC
	}
	return p, err
}
func (p *Packet) Discard() error {
	dbp("DiscardPacket")
	if p.cursor != p.end {
		p.cursor = p.end
		p.ic.writeOp(BIT_FIELD_SET, ECON2, ECON2_PKTDEC)
	}
	return nil
}
func (p *Packet) Read(buff []byte) (n uint16, err error) {
	dbp("ReadPacket")
	// total remaining packet length
	plen := p.end - p.cursor
	if plen == 0 {
		return 0, io.EOF
	}
	if len(buff) == 0 {
		return 0, nil
	}
	// Limit retreive length if Total packet length is greater than buffer length
	if plen > uint16(len(buff)) {
		plen = uint16(len(buff))
	}
	println(p.ic)
	// copy the packet from the receive buffer
	p.ic.readBuffer(buff[:plen])
	dbp("ReadPacket2")
	p.cursor += plen
	// Move the RX read pointer to where we ended reading
	p.ic.write16(ERXRDPTL, p.cursor)
	if p.cursor == p.end { // minus CRC length
		// decrement packet counter to indicate we are done with it.
		p.ic.writeOp(BIT_FIELD_SET, ECON2, ECON2_PKTDEC)
		err = io.EOF
	}
	return plen, err
}

func (d *Dev) Write(buff []byte) (uint16, error) {
	plen := uint16(len(buff))
	if plen+d.tcursor > MAX_FRAMELEN {
		d.tcursor = 0
		return 0, ErrBufferSize
	}

	d.write16(EWRPTL, TXSTART_INIT+d.tcursor)
	if d.tcursor == 0 {
		// write per-packet control byte (0x00 means use macon3 settings)
		d.writeOp(WRITE_BUF_MEM, 0, 0x00)
		d.tcursor++ // WBM spurious increment
	}
	// copy the packet into the transmit buffer
	d.writeBuffer(buff)
	d.tcursor += plen
	return plen, nil
}

func (d *Dev) Flush() error {
	dbp("send response")
	d.write16(ETXNDL, TXSTART_INIT+d.tcursor-1) // subtract WBM spurious increment
	// send the contents of the transmit buffer onto the network
	d.writeOp(BIT_FIELD_SET, ECON1, ECON1_TXRTS)
	// Reset the transmit logic problem. See Rev. B4 Silicon Errata point 12.
	if d.read(EIR)&EIR_TXERIF != 0 {
		d.writeOp(BIT_FIELD_CLR, ECON1, ECON1_TXRTS)
	}
	d.tcursor = 0
	return nil
}
