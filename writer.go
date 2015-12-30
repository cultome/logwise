package logwise

import (
  "fmt"
  "os"
)

type FileWriter struct {
  FilePath string
  Append bool
  Prefix string
  Postfix string
}

func NewFileWriter(filePath string, append bool) *FileWriter {
  return &FileWriter{filePath, append, "", ""}
}

func (w *FileWriter) AddPrefix(prefix string) *FileWriter {
  w.Prefix = prefix
  return w
}

func (w *FileWriter) AddPostfix(postfix string) *FileWriter {
  w.Postfix = postfix
  return w
}

func (w *FileWriter) Write(contents []string) {
  var flags int

  if w.Append {
    flags = os.O_WRONLY | os.O_APPEND
  } else {
    flags = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
  }

  file,_ := os.OpenFile(w.FilePath, flags, 0666)
  defer file.Close()
  
  if w.Prefix != ""{
    file.WriteString(fmt.Sprintf("%v\n", w.Prefix))
  }

  for _,line := range contents {
    file.WriteString(fmt.Sprintf("%v\n", line))
  }

  if w.Postfix != ""{
    file.WriteString(fmt.Sprintf("%v\n", w.Postfix))
  }
}