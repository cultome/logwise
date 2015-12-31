package logwise

import (
  "fmt"
  "regexp"
  "bufio"
  "os"
)

type Filter interface {
  Filter(reader LineReader, patterns []string) []*Line
}

type LineReader interface {
  Read() (*Line, error)
  Source() string
}

/*
 * Line
 */
type Line struct {
  LineNbr int
  Content string
  FileName string
  Regexp string
}

func (line *Line) String() string {
  return fmt.Sprintf("[%6d] %v", line.LineNbr, line.Content)
}

/*
 * FileReader
 */
type FileReader struct {
  Files []string
  currentFile int
  currentLineNbr int
  fileInScanner *os.File
  scanner *bufio.Scanner
}

func NewFileReader(filePaths ...string) LineReader {
  return &FileReader{filePaths, -1, -1, nil, nil}
}

func (r *FileReader) Read() (*Line, error) {
  var err error
  if r.currentFile < 0 {
    r.currentFile = 0
    r.currentLineNbr = 0
    r.scanner, r.fileInScanner, err = scanner(r.Files[r.currentFile])
  } else if r.scanner == nil {
    r.currentFile += 1
    r.fileInScanner.Close()
    if r.currentFile >= len(r.Files) {
      return nil, fmt.Errorf("End of files")
    }
    r.scanner, r.fileInScanner, err = scanner(r.Files[r.currentFile])
  }

  if err != nil {
    return nil, err
  }

  if r.scanner.Scan() {
    r.currentLineNbr += 1
    return &Line{r.currentLineNbr, r.scanner.Text(), r.Files[r.currentFile], ""}, nil
  }

  r.scanner = nil
  return r.Read()
}

func scanner(filePath string) (*bufio.Scanner, *os.File, error) {
  file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
  if err != nil {
    return nil, nil, fmt.Errorf("File [%v] cannot be open: %v", filePath, err)
  }
  scanner := bufio.NewScanner(file)
  return scanner, file, nil
}

func (r *FileReader) Source() string {
  return r.Files[r.currentFile]
}

/*
 * LineFilter
 */
type LineFilter struct {
  Reader LineReader
  Patterns []string
}

func NewLineFilter(patterns []string, files ...string) Filter {
  reader := NewFileReader(files...)
  return &LineFilter{reader, patterns}
}

func (filter *LineFilter) Filter(reader LineReader, patterns []string) []*Line {
  f, p := filterOperativeParams(filter, reader, patterns)

  var linesMatched []*Line
  regexps := make([]*regexp.Regexp, len(p))

  for idx,pattern := range p {
    reg, _ := regexp.Compile(pattern)
    regexps[idx] = reg
  }

  lines, err := scanFile(f, regexps)
  if err == nil {
    if len(lines) > 0{
      newLines := filterExisting(linesMatched, lines)
      linesMatched = append(linesMatched, newLines...)
    }
  }

  return linesMatched
}

func filterExisting(existing, news []*Line) []*Line {
  var newLines []*Line
  exist := false

  for _,n := range news {
    exist = false
    for _,e := range existing {
      if n.LineNbr == e.LineNbr && e.FileName == n.FileName {
        exist = true
        break
      }
    }

    if !exist {
      newLines = append(newLines, n)
    }
  }

  return newLines
}

func filterOperativeParams(filter *LineFilter, reader LineReader, patterns []string) (LineReader, []string) {
  if (reader == nil && filter.Reader == nil) || (patterns == nil && filter.Patterns == nil) {
    panic("Reader and Patterns are required for Filter to work!")
  }

  var f LineReader
  var p []string

  if filter.Reader == nil {
    f = reader
  } else {
    f = filter.Reader
  }

  if filter.Patterns == nil {
    p = patterns
  } else {
    p = filter.Patterns
  }

  return f,p
}

func scanFile(reader LineReader, regexps []*regexp.Regexp) ([]*Line, error) {
  var lines []*Line
  lineIdx := 1
  var err error

  for line, err := reader.Read(); err == nil; line, err = reader.Read() {
    regex := match(line.Content, regexps)
    if regex != "" {
      line := Line{lineIdx, line.Content, reader.Source(), regex}
      lines = append(lines, &line)

      if len(lines) > 200000 {
        panic("Too many lines detected. Please refine the search pattern")
      }
    }
    lineIdx++
  }

  if err != nil && err.Error() != "End of files" {
    return nil, err
  }

  return lines, nil
}

func match(line string, regexps []*regexp.Regexp) string {
  for _, reg :=  range regexps {
    if reg.Match([]byte(line)) {
      return reg.String()
    }
  }
  return ""
}
