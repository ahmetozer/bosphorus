package flags

type ArrFlag []string

func (i *ArrFlag) String() string {
	return ""
}

func (i *ArrFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var TcpFlag ArrFlag

//var TcpRawFlag ArrFlag

var TunFlag ArrFlag

type connType string
