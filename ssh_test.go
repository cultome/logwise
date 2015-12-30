package logwise

import "testing"

func TestConnect(t *testing.T){
  session := NewServerSession()
  session.Connect("server:22", "user", "passwd")
}