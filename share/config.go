package share

import (
	"errors"
	"strconv"
	"strings"
)

//Contains the config of other nodes and Port

const DEFAULT_PORT = int16(6800)

type Node struct {
	Host string
	Port int16
}

func NewNode(host string, port int16) (*Node, error) {
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

		var port int16
		if len(temp) == 1 {
			port = DEFAULT_PORT
		} else {
			p, err := strconv.Atoi(temp[1])
			if err != nil {
				return nil, err
			}

			port = int16(p)
		}

		var err error

		nodes[i], err = NewNode(host, port)
		if err != nil {
			return nil, err
		}
	}

	return nodes, nil
}
