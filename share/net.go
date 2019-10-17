package share

import (
	"fmt"
	"github.com/kpango/glg"
	"github.com/tidwall/evio"
	"net"
	"sync"
)

//communication with all the nodes
func Send(node *Node, msg []byte) error {
	con, err := net.Dial("tcp", node.String())
	if err != nil {
		return err
	}

	defer con.Close()

	_, err = con.Write(msg)
	return err
}

type Server interface {
	Start(port int) error
}

type tcpServer struct {
	events    *evio.Events
	cons      *sync.Map
	callBacks []Handler
}

func NewTcpServer(c []Handler) Server {
	events := &evio.Events{
		NumLoops: -1, //NumProcs
	}
	return &tcpServer{
		events:    events,
		cons:      new(sync.Map),
		callBacks: c,
	}
}

func (s *tcpServer) Start(port int) error {
	addr := fmt.Sprintf("tcp://0.0.0.0:%d?reuseport=true", port)

	s.events.Serving = func(srv evio.Server) (action evio.Action) {
		glg.Infof("TCP Server started on port %d (loops: %d)", srv.Addrs[0], srv.NumLoops)
		return
	}

	s.events.Opened = func(ec evio.Conn) (out []byte, opts evio.Options, action evio.Action) {
		s.cons.Store(ec.RemoteAddr(), struct{}{})
		return
	}

	s.events.Closed = func(ec evio.Conn, err error) (action evio.Action) {
		s.cons.Delete(ec.RemoteAddr())
		return
	}

	s.events.Data = func(c evio.Conn, in []byte) ([]byte, evio.Action) {
		for _, h := range s.callBacks {
			h.Handle(c.RemoteAddr().String(), in)
		}
		return in, evio.Close //close connection //Handle in a goroutine
	}

	return evio.Serve(*s.events, addr)
}

type Handler interface {
	Handle(addr string, data []byte)
}
