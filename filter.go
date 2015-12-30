package logwise

import (
  "bufio"
  "fmt"
  "os"
  "regexp"
)

type Filter interface {
  Filter(reader *LineReader, patterns []string) []*Line
}

type LineReader interface {
  Read() *Line
  Source() string
}

type Line struct {
  LineNbr int
  Content string
  FileName string
  Regexp string
}

type LineFilter struct {
  Reader *LineReader
  Patterns []string
}

func (line *Line) String() string {
  return fmt.Sprintf("[%6d] %v", line.LineNbr, line.Content)
}

func NewLineFilter(reader *LineReader, patterns []string) Filter {
  return &LineFilter{files, patterns}
}

/*
func (filter *LineFilter) SetPatterns(patterns []string) *LineFilter {
  filter.Patterns = patterns
  return filter
}

func (filter *LineFilter) SetFiles(files ...string) *LineFilter {
  filter.Files = files
  return filter
}

func (filter *LineFilter) Set(files []string, patterns []string) *LineFilter {
  filter.Files = files
  filter.Patterns = patterns
  return filter
}
*/

func (filter *LineFilter) Filter(reader *LineReader, patterns []string) []*Line {
  f, p := filterOperativeParams(filter, reader, patterns)

  var linesMatched []*Line
  regexps := make([]*regexp.Regexp, len(p))

  for idx,pattern := range p {
    reg, _ := regexp.Compile(pattern)
    regexps[idx] = reg
  }

  //for _,file := range f {
    lines := scanFile(reader, regexps)
    if len(lines) > 0{
      newLines := filterExisting(linesMatched, lines)
      linesMatched = append(linesMatched, newLines...)
    }
  //}
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

func filterOperativeParams(filter *LineFilter, reader *LineReader, patterns []string) (*LineReader, []string) {
  if (reader == nil && filter.Reader == nil) || (patterns == nil && filter.Patterns == nil) {
    panic("Reader and Patterns are required for Filter to work!")
  }

  var f []string
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

func scanFile(reader *LineReader, regexps []*regexp.Regexp) []*Line {
  var lines []*Line

  //file,_ := os.Open(filePath)
  //defer file.Close()
  //scanner := bufio.NewScanner(file)

  lineIdx := 1
  //for scanner.Scan() {
    regex := match(reader.Read(), regexps)
    if regex != "" {
      line := Line{lineIdx, scanner.Text(), reader.Source(), regex}
      lines = append(lines, &line)
    }
    lineIdx++
  //}

  return lines
}

func match(line string, regexps []*regexp.Regexp) string {
  for _, reg :=  range regexps {
    if reg.Match([]byte(line)) {
      return reg.String()
    }
  }
  return ""
}
