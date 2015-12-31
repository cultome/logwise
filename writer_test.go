package logwise

import (
  "testing"
  "io/ioutil"
)

func TestWriteFile(t *testing.T) {
  writer := NewFileWriter("logs/result.log", false, "", "")
  writer.Write([]string{"uno", "dos", "tres"})

  content, err := ioutil.ReadFile("logs/result.log")

  if err != nil {
    t.Errorf("Error reading file: %v\n", err)
  }

  if content == nil {
    t.Errorf("File is empty!\n")
  }

  c := string(content)
  if c != "uno\ndos\ntres\n" {
    t.Errorf("Writing was not successful!\n")
  }
}