package logwise

import "testing"

func TestWriteFile(t *testing.T) {
  writer := NewFileWriter("logs/result.log", false)
  writer.Write([]string{"uno", "dos", "tres"})
}