package logwise

import (
  "golang.org/x/crypto/ssh"
  "fmt"
)

type ServerSession struct {
  server string
  config *ssh.ClientConfig
}

func NewServerSession(server, username, password string) *ServerSession {
  config := &ssh.ClientConfig {
    User: username,
    Auth: []ssh.AuthMethod{
      ssh.Password(password),
    },
  }
  return &ServerSession{server, config}
}

func (s *ServerSession) Connect() {
  conn, err := ssh.Dial("tcp", s.server, s.config)
  if err != nil {
    panic("OH NO!! The connection failed!")
  }
  defer conn.Close()

  session, err := conn.NewSession()
  if err != nil {
    panic("OH NO!! The session creation failed!")
  }
  defer session.Close()

  out, err := session.Output("ls -l $INC_LOG")
  if err != nil {
    panic("OH NO!! The command failed!")
  }

  fmt.Printf("> ls -l $INC_LOG\n%v\n", string(out))
}
