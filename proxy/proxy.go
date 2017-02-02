package proxy

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/soider/d"
	"io"
)

type ProxyServer struct {
	wsConn  *websocket.Conn
	tcpConn io.ReadWriteCloser
	errCh   chan error
	done    chan struct{}
}

func NewProxyServer(wsConn *websocket.Conn, tcpConn io.ReadWriteCloser) *ProxyServer {
	proxyserver := ProxyServer{wsConn, tcpConn, make(chan error, 4), make(chan struct{})}
	return &proxyserver
}

func (proxyserver *ProxyServer) DoProxy(ctx context.Context) {
	go proxyserver.WsToTcp(ctx)
	go proxyserver.closer(ctx)
	proxyserver.TcpToWs(ctx)

}

func (ps *ProxyServer) closer(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			ps.tcpConn.Close()
			ps.wsConn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		case err := <-ps.errCh:
			d.D("Have error", err)
		}
	}
}

func (proxyserver *ProxyServer) TcpToWs(ctx context.Context) {
	buffer := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			d.D("Stop tcptows goroutine")
			return
		default:
			n, err := proxyserver.tcpConn.Read(buffer)
			if err != nil {
				//proxyserver.errCh <- err
				return

			}
			err = proxyserver.wsConn.WriteMessage(websocket.BinaryMessage, buffer[0:n])
			if err != nil {
				//proxyserver.errCh <- err
				return

			}
		}
	}
}

func (proxyserver *ProxyServer) WsToTcp(ctx context.Context) {
	proxyserver.wsConn.Subprotocol()
	for {
		select {
		case <-ctx.Done():
			d.D("Stop wstotcp goroutine")
			return
		default:
			_, data, err := proxyserver.wsConn.ReadMessage()
			if err != nil {
				//proxyserver.errCh <- err
				return

			}
			_, err = proxyserver.tcpConn.Write(data)
			if err != nil {
				//proxyserver.errCh <- err
				return
			}
		}

	}
}
