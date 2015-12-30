package logwise

import (
  "golang.org/x/crypto/ssh"
  "fmt"
)

type ServerSession struct {

}

func NewServerSession() *ServerSession {
  return &ServerSession{}
}

func (s *ServerSession) Connect(server, username, password string) {
  config := &ssh.ClientConfig{
    User: username,
    Auth: []ssh.AuthMethod{
      ssh.Password(password),
    },
  }

  conn, err := ssh.Dial("tcp", server, config)
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
