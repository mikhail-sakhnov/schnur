package vnc

import (
	"context"
	"fmt"
	"net"
)

type VncConnector struct {
}

type VncConnection struct {
	c *net.TCPConn
}

func (vcc *VncConnection) Close() error {
	return vcc.c.Close()
}

func (vcc *VncConnection) Read(p []byte) (n int, err error) {
	return vcc.c.Read(p)
}

func (vcc *VncConnection) Write(p []byte) (n int, err error) {
	return vcc.c.Write(p)
}

func (vc *VncConnector) Connect(ctx context.Context, host string, port int) (*VncConnection, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}
	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	c := &VncConnection{tcpConn}
	return c, nil
}
