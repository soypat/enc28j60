package enc28j60

type ErrorCode uint8

const (
	errUndefined ErrorCode = iota
	ErrOutOfBound
	// buff size not in 64..1500
	ErrBufferSize
	// got rev=0. is dev connected?
	ErrBadRev
	// mac addr len not 6
	ErrBadMac
	// invalid IP address
	ErrBadIP
	// unable to resolve ARP
	ErrUnableToResolveARP
	// ARP protocol violation
	ErrARPViolation
	// internet protocol procedure not implemented
	ErrIPNotImplemented
	// I/O
	ErrIO
	// read deadline exceeded
	ErrRXDeadlineExceeded
	// CRC checksum fail
	ErrCRC
	// EOF
	ErrEOF
)

func (err ErrorCode) Error() string {
	switch err {
	case ErrOutOfBound:
		return "out of buff bounds"
	case ErrBufferSize:
		return "buff size not in 64..1500"
	case ErrBadRev:
		return "got rev=0. is dev connected?"
	case ErrBadIP:
		return "invalid IP address"
	case ErrBadMac:
		return "mac addr len not 6"
	case ErrUnableToResolveARP:
		return "unable to resolve ARP"
	case ErrARPViolation:
		return "ARP protocol violation"
	case ErrIO:
		return "I/O"
	case ErrRXDeadlineExceeded:
		return "rx deadline exceeded"
	case ErrCRC:
		return "CRC error"
	case ErrEOF:
		return "encEOF"
	case ErrIPNotImplemented:
		return "internet protocol procedure not implemented"
	}
	return "undefined"
}
