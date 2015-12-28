package logwise

import (
  "bufio"
  "fmt"
  "os"
  "regexp"
)

type Filter interface {
  Filter(files []string, pattern []string) []*Line
  Files() []string
  Patterns() []string
}

type Line struct {
  LineNbr int
  Content string
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

func (filter *LineFilter) SetFiles(files []string) *LineFilter {
  filter.files = files
  return filter
}

func (filter *LineFilter) Set(files []string, patterns []string) *LineFilter {
  filter.files = files
  filter.patterns = patterns
  return filter
}

func (filter *LineFilter) Files() []string {
  return filter.files
}

func (filter *LineFilter) Patterns() []string {
  return filter.patterns
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
      linesMatched = append(linesMatched, lines...)
    }
  }
  return linesMatched
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
  scanner := bufio.NewScanner(file)

  lineIdx := 1
  for scanner.Scan() {
    if match(scanner.Text(), regexps) {
      line := Line{lineIdx, scanner.Text()}
      lines = append(lines, &line)
    }
    lineIdx++
  }

  return lines
}

func match(line string, regexps []*regexp.Regexp) bool {
  for _, reg :=  range regexps {
    if reg.Match([]byte(line)) {
      return true
    }
  }
  return false
}
