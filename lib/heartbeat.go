package lib

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/kpango/glg"
	"github.com/mangatmodi/consensus/share"
	"time"
)

//Send bytes as `echo -n -e \\x0f\\x0f\\xc8\\x00| nc 127.0.0.1 300` for port 200
//A heartbeat message is header, followed by port-number(uint16)
var beat = []byte{0x0f, 0x0f}

const (
	BEAT_SIZE = 2
	TTLms     = 5000 //5 seconds
)

func IsHeartBeat(data []byte) bool {
	if len(data) < BEAT_SIZE {
		return false
	}

	return bytes.Equal(data[0:BEAT_SIZE], beat)
}

//Get and send heartbeats
type HeartBeatHandler interface {
	Handle(addr string, data []byte)
	Send(node *share.Node) error
	SendHeartBeats(nodes []*share.Node)
}

type heartBeat struct {
	states []*syncState
}

type syncState struct {
	Node *share.Node
	At   time.Time
}

func (s *syncState) String() string {
	return fmt.Sprintf("Node:%s, at:%v", s.Node, s.At)
}

func (h *heartBeat) Handle(addr string, data []byte) {
	if !IsHeartBeat(data) {
		return //not an heartbeat
	}
	n := extractNode(addr, data)
	h.sync(n)
}

func msg() []byte {
	msg := make([]byte, BEAT_SIZE+2)
	for i, _ := range beat {
		msg[i] = beat[i]
	}

	p := make([]byte, 2)
	binary.LittleEndian.PutUint16(p, share.Port)
	msg[2] = p[0]
	msg[3] = p[1]
	return msg
}

func (h *heartBeat) Send(node *share.Node) error {
	return share.Send(node, msg())
}

func extractNode(addr string, data []byte) *share.Node {
	n, _ := share.Parse(addr)
	p := binary.LittleEndian.Uint16(data[BEAT_SIZE:])
	n[0].Port = p
	return n[0]
}

func (h *heartBeat) exists(node *share.Node) bool {
	for _, v := range h.states {
		if (v.Node.Host == node.Host) &&
			(v.Node.Port == node.Port) && //Adding only 1 node
			(time.Now().Sub(v.At).Milliseconds() < TTLms) {
			v.At = time.Now()
			return true
		}
	}
	return false
}

func (h *heartBeat) outSync(node *share.Node) {
	for _, v := range h.states {
		if v.Node.Host == node.Host {
			v.At = time.Unix(0, 0)
		}
	}
}

func (h *heartBeat) sync(node *share.Node) {
	s := &syncState{
		Node: node,
		At:   time.Now(),
	}
	if !h.exists(node) { //Only if not synced
		h.states = append(h.states, s)
	}
}

func (h *heartBeat) syncedNodes() []*syncState {
	var n = []*syncState{}
	for _, v := range h.states {
		if time.Now().Sub(v.At).Milliseconds() < TTLms {
			n = append(n, v)
		}
	}
	return n

}

func NewHeartbeatHandler() HeartBeatHandler {
	h := &heartBeat{states: []*syncState{}}

	go func(h *heartBeat) {
		for {
			j, _ := json.Marshal(h.syncedNodes())
			glg.Debugf("All Nodes:%v\n", string(j))
			time.Sleep(1 * time.Second)
		}
	}(h)
	return h
}

func (h *heartBeat) SendHeartBeats(nodes []*share.Node) {
	go func(h *heartBeat, nodes []*share.Node) {
		for {
			for _, n := range nodes {
				err := h.Send(n)
				if err != nil {
					glg.Errorf("Unable to send heartbeat to %s, because %v", n, err)
				}
			}
			time.Sleep(500 * time.Millisecond)
		}
	}(h, nodes)
}
