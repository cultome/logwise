package logwise

import (
  "bufio"
  "fmt"
  "os"
  "regexp"
)

type Line struct {
  LineNbr int
  Content string
  FileName string
  Regexp string
}

type LineFilter struct {
  files []string
  patterns []string
}

func (line *Line) String() string {
  return fmt.Sprintf("[%6d] %v", line.LineNbr, line.Content)
}

func NewLineFilter() *LineFilter {
  return &LineFilter{}
}

func (filter *LineFilter) SetPatterns(patterns []string) *LineFilter {
  filter.patterns = patterns
  return filter
}

func (filter *LineFilter) SetFiles(files ...string) *LineFilter {
  filter.files = files
  return filter
}

func (filter *LineFilter) Set(files []string, patterns []string) *LineFilter {
  filter.files = files
  filter.patterns = patterns
  return filter
}

func (filter *LineFilter) Filter(files []string, patterns []string) []*Line {
  f, p := filterOperativeParams(filter, files, patterns)

  var linesMatched []*Line
  regexps := make([]*regexp.Regexp, len(p))

  for idx,pattern := range p {
    reg, _ := regexp.Compile(pattern)
    regexps[idx] = reg
  }

  for _,file := range f {
    lines := scanFile(file, regexps)
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

func filterOperativeParams(filter *LineFilter, files []string, patterns []string) ([]string, []string) {
  if (files == nil && filter.files == nil) || (patterns == nil && filter.patterns == nil) {
    panic("Files and Patterns are required for Filter to work!")
  }

  var f []string
  var p []string

  if filter.files == nil {
    f = files
  } else {
    f = filter.files
  }

  if filter.patterns == nil {
    p = patterns
  } else {
    p = filter.patterns
  }

  return f,p
}

func scanFile(filePath string, regexps []*regexp.Regexp) []*Line {
  var lines []*Line

  file,_ := os.Open(filePath)
  defer file.Close()
  scanner := bufio.NewScanner(file)

  lineIdx := 1
  for scanner.Scan() {
    regex := match(scanner.Text(), regexps)
    if regex != "" {
      line := Line{lineIdx, scanner.Text(), filePath, regex}
      lines = append(lines, &line)
    }
    lineIdx++
  }

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
