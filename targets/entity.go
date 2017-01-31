package targets

import (
	"errors"
	"fmt"
)

type Command string
type CommandsList []string

type targetType string
type commandsByTargets map[targetType]CommandsList

var ErrNoCommandsOrType = errors.New("No commands or predifined type in configuration ")

func newUnknownTypeError(t targetType) error {
	return fmt.Errorf("Unknown predifined type %s", t)
}

type TargetList []Target

type Target struct {
	Name       string       `json:"name"`
	Address    string       `json:"address"`
	SshPort    int          `json:"ssh_port"`
	VncAddress string       `json:"vnc_address"`
	VncPort    string       `json:"vnc_port"`
	Username   string       `json:"username"`
	Password   string       `json:"password"`
	KeyPath    string       `json:"key_path"`
	CmdList    CommandsList `json:"cmd_list"` // TODO: fix unparsing from yaml
	Type       targetType   `json:"type"`
}

func (t *Target) GetVncAddress() string {
	if t.VncAddress == "" {
		return t.Address
	}
	return t.VncAddress
}

func (t *Target) GetCmd() (CommandsList, error) {
	if t.CmdList != nil {
		return t.CmdList, nil
	}
	if t.Type == "" {
		return nil, ErrNoCommandsOrType
	}

	commands, found := predifinedTypes[t.Type]
	if !found {
		return nil, newUnknownTypeError(t.Type)
	}
	t.CmdList = commands
	return commands, nil
}
