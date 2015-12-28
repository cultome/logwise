package logwise

import (
  "fmt"
  "bufio"
  "os"
)

type FileWriter struct {
  FilePath string
}

func NewFileWriter(filePath string) *FileWriter {
  return &FileWriter{filePath}
}

func (w *FileWriter) Write(contents []string) {
  file,_ := os.Create(w.FilePath)
  defer file.Close()
  
  writer := bufio.NewWriter(file)

  for _,line := range contents {
    writer.Write([]byte(fmt.Sprintf("%v\n", line)))
  }

  writer.Flush()
}