package conn

type ConnType string

const (
	TCP ConnType = "TCP"
	TUN ConnType = "TUN"
)

func (d ConnType) String() string {
	return string(d)
}
