package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/soider/d"
	"github.com/soider/schnur/proxy"
	"github.com/soider/schnur/targets"
	"github.com/soider/schnur/targets/manager"
	"github.com/soider/schnur/vnc"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsHandler struct {
	tm  *manager.TargetsManager
	vnc vnc.VncConnector
}

func NewWsHandler(tm *manager.TargetsManager, vnc vnc.VncConnector) *WsHandler {
	return &WsHandler{
		tm:  tm,
		vnc: vnc,
	}
}

// TODO: ping handler
// TODO: stop vnc connector the right way via ctx
func (ws *WsHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	targetName := mux.Vars(req)["target"]
	target, err := ws.tm.Target(targetName)
	if err != nil {
		http.Error(rw, err.Error(), 404)
	}
	h := http.Header{
		"Sec-Websocket-Protocol": {"binary"},
		"Sec-WebSocket-Version":  {"13"},
	}
	wsConn, err := upgrader.Upgrade(rw, req, h)

	if err != nil {
		http.Error(rw, "Can't upgrade ws", 500)
		return
	}
	ws.handle(ctx, wsConn, target)
}

func (ws *WsHandler) handle(ctx context.Context, wsConn *websocket.Conn, target targets.Target) {

	ctx, cancel := context.WithCancel(ctx)
	_ = cancel
	address := target.GetVncAddress()
	d.D(address, target.VncPort, target.VncPassword)
	vncConn, err := ws.vnc.Connect(ctx, address, target.VncPort)
	if err != nil {
		d.D(err)
		wsConn.WriteMessage(websocket.CloseMessage, []byte(err.Error()))
		return
	}
	//defer wsConn.WriteMessage(websocket.CloseMessage, []byte{})
	//defer vncConn.Close()
	proxy := proxy.NewProxyServer(wsConn, vncConn)
	proxy.DoProxy(ctx)
}
