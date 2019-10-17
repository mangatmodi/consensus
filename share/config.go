package share

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

//Contains the config of other nodes and Port

const DEFAULT_PORT = uint16(6800)
var Port uint16
type Node struct {
	Host string
	Port uint16
}

func (n *Node) String() string {
	return fmt.Sprintf("%s:%d", n.Host, n.Port)
}

func NewNode(host string, port uint16) (*Node, error) {
	if host == "" || port < 0 {
		return nil, errors.New("invalid Host and Port")
	}

	if port == 0 {
		port = DEFAULT_PORT
	}
	return &Node{
		Host: strings.TrimSpace(host),
		Port: port,
	}, nil
}

//Parse takes comma separated Host:Port string and returns an array of Node
func Parse(str string) ([]*Node, error) {
	arr := strings.Split(str, ",")
	nodes := make([]*Node, len(arr))

	for i, el := range arr {
		temp := strings.SplitN(el, ":", 2)
		host := temp[0]

		var err error
		var port uint16

		if len(temp) == 1 {
			port = DEFAULT_PORT
		} else {
			p, err := strconv.Atoi(temp[1])
			if err != nil {
				return nil, err
			}
			port = uint16(p)
		}

		nodes[i], err = NewNode(host, port)
		if err != nil {
			return nil, err
		}
	}

	return nodes, nil
}
