package connect

import "fmt"

type ClientInfo struct {
	Name    string
	Version string
	Build   string
}

func (ci *ClientInfo) String() string {
	return fmt.Sprintf("%s/%s (build %s)", name, version, build)
}
