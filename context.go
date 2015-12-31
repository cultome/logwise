package logwise

import(
  "regexp"
)

type Contexter interface {
  Get(lines ...*Line) []*Line
}

type LineContext struct {
  BeforeMatch, AfterMatch string
}

type lineRange struct {
  begin, end int
  fileName string
}

func NewLineContext(beforeMatch, afterMatch string) Contexter {
  return &LineContext{beforeMatch, afterMatch}
}

func (ctx *LineContext) Get(lines ...*Line) []*Line {
  var result []*Line
  var existingRanges []lineRange

  for _,line := range lines {
    rangeLines := getLines(ctx, line)
    if !rangeExists(&existingRanges, &rangeLines) {
      result = append(result, rangeLines...)
    }
  }
  return result
}

func rangeExists(existingRanges *[]lineRange, lrange *[]*Line) bool {
  begin := (*lrange)[0]
  end := (*lrange)[len(*lrange)-1]

  if existingRanges != nil {
    for _,r := range *existingRanges {
      if r.begin == begin.LineNbr && r.end == end.LineNbr && r.fileName == begin.FileName {
        return true
      }
    }
  }

  r := lineRange{begin.LineNbr, end.LineNbr, begin.FileName}
  *existingRanges = append(*existingRanges, r)
  return false
}

func getLines(ctx *LineContext, line *Line) []*Line {
  reader := NewFileReader(line.FileName)

  topRange := topRange(reader, ctx.BeforeMatch, line.LineNbr)
  bottomRange := bottomRange(reader, ctx.AfterMatch)
  return append(topRange, bottomRange...)
}

func bottomRange(reader LineReader, expression string) []*Line {
  var lines []*Line
  reg := regexp.MustCompile(expression)
  
  for l,err := reader.Read(); err == nil; l,err = reader.Read() {
    if reg.MatchString(l.Content) {
      break
    } else {
      lines = append(lines, l)
    }
  }
  return lines
}

func topRange(reader LineReader, expression string, refLineNbr int) []*Line {
  var lines []*Line
  capturing := false

  reg := regexp.MustCompile(expression)

  for l,err := reader.Read(); err == nil; l,err = reader.Read() {
    if reg.MatchString(l.Content) {
      capturing = true
      lines = nil
    }

    if l.LineNbr > refLineNbr {
      lines = append(lines, l)
      break
    }

    if capturing {
      lines = append(lines, l)
    }
  }

  return lines
}
