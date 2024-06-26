package conn

type ConnType string

const (
	TCP ConnType = "TCP"
	TUN ConnType = "TUN"
	UDP ConnType = "UDP"
)

func (d ConnType) String() string {
	return string(d)
}

const BufferSize = 64 * 1024 // 64 KB buffer size for io.CopyBuffer
