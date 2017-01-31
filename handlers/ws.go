package handlers

import (
	"github.com/gorilla/websocket"
	"net/http"
	"net"
	"github.com/soider/d"
	"github.com/gorilla/mux"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Proxy(w http.ResponseWriter, r *http.Request) {
	d.D("New connection")
	reqData := mux.Vars(r)
	d.D(reqData)
	h := http.Header{
		"Sec-Websocket-Protocol":{"binary"},
		"Sec-WebSocket-Version": {"13"},
	}
	wsConn, err := upgrader.Upgrade(w, r, h)

	if err != nil {
		d.D(err)
		return
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:5900")
	if err != nil {
		errorMsg := "FAIL(net resolve tcp addr): " + err.Error()
		d.D(errorMsg)
		_ = wsConn.WriteMessage(websocket.CloseMessage, []byte(errorMsg))
		return
	}

	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		errorMsg := "FAIL(net dial tcp): " + err.Error()
		d.D(errorMsg)
		_ = wsConn.WriteMessage(websocket.CloseMessage, []byte(errorMsg))
		return
	}

	proxyserver := NewProxyServer(wsConn, tcpConn)
	go proxyserver.doProxy()
}

type ProxyServer struct {
	wsConn  *websocket.Conn
	tcpConn *net.TCPConn
}

func NewProxyServer(wsConn *websocket.Conn, tcpConn *net.TCPConn) *ProxyServer {
	proxyserver := ProxyServer{wsConn, tcpConn}
	return &proxyserver
}

func (proxyserver *ProxyServer) doProxy() {
	go proxyserver.wsToTcp()
	proxyserver.tcpToWs()
}

func (proxyserver *ProxyServer) tcpToWs() {
	buffer := make([]byte, 1024)
	for {
		n, err := proxyserver.tcpConn.Read(buffer)
		if err != nil {
			proxyserver.tcpConn.Close()
			break
		}
		err = proxyserver.wsConn.WriteMessage(websocket.BinaryMessage, buffer[0:n])
		if err != nil {
			d.D(err.Error())
		}
	}
}

func (proxyserver *ProxyServer) wsToTcp() {
	proxyserver.wsConn.Subprotocol()
	for {
		_, data, err := proxyserver.wsConn.ReadMessage()
		if err != nil {
			break
		}

		_, err = proxyserver.tcpConn.Write(data)
		if err != nil {
			d.D(err.Error())
			break
		}
	}
}