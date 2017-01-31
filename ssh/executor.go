package ssh

import (
	"golang.org/x/crypto/ssh"
)

type clientBuildFunction func(cfg *ssh.ClientConfig) (*ssh.Client, error)

func buildSshClient(cfg *ssh.ClientConfig) (*ssh.Client, error) {
	client, err := ssh.Dial("tcp", "127.0.0.1:2222", cfg)
	return client, err
}

type Executor struct {
	buildClient clientBuildFunction
	Config      *ssh.ClientConfig
}

func New(cfg *ssh.ClientConfig) *Executor {
	return &Executor{buildClient: buildSshClient, Config: cfg}
}

func NewWithBuilder(f clientBuildFunction) *Executor {
	return &Executor{buildClient: f}
}

func (s *Executor) RunRemoteCMD(cmd string) (string, error) {
	client, err := s.buildClient(s.Config)
	if err != nil {
		return "", err
	}
	session, err := client.NewSession()
	defer session.Close()
	if err != nil {
		return "", err
	}
	d, err := session.CombinedOutput(cmd)
	return string(d), err
}
