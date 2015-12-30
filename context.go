package logwise

import(
  "os"
  "bufio"
)

type LineContext struct {
  BeforeMatch, AfterMatch string
}

type lineRange struct {
  begin, end int
  fileName string
}

func NewLineContext(beforeMatch, afterMatch string) *LineContext {
  return &LineContext{beforeMatch, afterMatch}
}

func (ctx *LineContext) Get(lines ...*Line) []*Line {
  var result []*Line
  var existingRanges []lineRange

  for _,line := range lines {
    rangeLines, begin, end := getLines(ctx, line)
    if !rangeExists(&existingRanges, &lineRange{begin.LineNbr, end.LineNbr, line.FileName}) {
      result = append(result, rangeLines...)
    }
  }
  return result
}

func rangeExists(existingRanges *[]lineRange, lineRange *lineRange) bool {
  if existingRanges != nil {
    for _,r := range *existingRanges {
      if r.begin == lineRange.begin && r.end == lineRange.end && r.fileName == lineRange.fileName {
        return true
      }
    }
  }

  *existingRanges = append(*existingRanges, *lineRange)
  return false
}

func getLines(ctx *LineContext, line *Line) ([]*Line, *Line, *Line) {
  filter := NewLineFilter()
  filter.Set([]string{line.FileName}, []string{ctx.BeforeMatch, ctx.AfterMatch})
  lines := filter.Filter(nil, nil)

  beginLine := findBegin(lines, line, ctx.BeforeMatch)
  endLine := findEnd(lines, line, ctx.AfterMatch)

  lineCtx := extractLineRange(line.FileName, beginLine.LineNbr, endLine.LineNbr)

  return lineCtx, beginLine, endLine
}

func extractLineRange(filePath string, beginLine, endLine int) []*Line {
  var lines []*Line
  file,_ := os.Open(filePath)
  defer file.Close()
  scanner := bufio.NewScanner(file)

  lineIdx := 0
  for ; lineIdx < beginLine; lineIdx++ {    
    scanner.Scan() 
  }

  for ; lineIdx < endLine; lineIdx++ {
    lines = append(lines, &Line{lineIdx, scanner.Text(), filePath, ""})
    scanner.Scan()
  }

  return lines
}

func findBegin(lines []*Line, line *Line, regexp string) *Line {
  var begin *Line
  for _,l := range lines {
    if l.Regexp == regexp && l.LineNbr < line.LineNbr && (begin == nil || begin.LineNbr < l.LineNbr) {
      begin = l
    }
  }
  return begin
}

func findEnd(lines []*Line, line *Line, regexp string) *Line {
  var end *Line
  for _,l := range lines {
    if l.Regexp == regexp && l.LineNbr > line.LineNbr && (end == nil || end.LineNbr > l.LineNbr) {
      end = l
    }
  }
  return end
}
