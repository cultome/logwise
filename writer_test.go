package logwise

import "testing"

func TestWriteFile(t *testing.T) {
  writer := NewFileWriter("logs/result.log")
  writer.Write([]string{"uno", "dos", "tres"})
}