package targets

import "golang.org/x/crypto/ssh"

func (t Target) SshConfig() *ssh.ClientConfig {
	cfg := &ssh.ClientConfig{
		User: t.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(t.Password),

		},
	}
	return cfg
}
