package logwise

import (
  "regexp"
  "fmt"
  "strings"
  "sort"
)

type Extractor interface {
  Extract(lines []*Line, patterns map[string]string) []*Extraction
  Lines() []*Line
  Patterns() map[string]string
}

type PatternExtractor struct {
  lines []*Line
  patterns map[string]string
}

type Extraction struct {
  Line *Line
  Matches *map[string]string
}

func (extractor *Extraction) String() string {
  keys := make([]string, len(*extractor.Matches))
  for k,_ := range *extractor.Matches {
    keys = append(keys, k)
  }
  sort.Strings(keys)

  str := ""
  separator := ""

  for _, k := range keys {
    if k != "" {
      str += fmt.Sprintf("%v%v: %q", separator, k, (*extractor.Matches)[k])
      separator = ", "
    }
  }
  
  return str
}

func NewPatternExtractor() *PatternExtractor {
  return &PatternExtractor{}
}

func (extractor *PatternExtractor) SetLines(lines []*Line) *PatternExtractor {
  extractor.lines = lines
  return extractor
}

func (extractor *PatternExtractor) SetPatterns(patterns map[string]string) *PatternExtractor {
  extractor.patterns = patterns
  return extractor
}

func (extractor *PatternExtractor) Set(lines []*Line, patterns map[string]string) *PatternExtractor {
  extractor.lines = lines
  extractor.patterns = patterns
  return extractor
}

func (extractor *PatternExtractor) Lines() []*Line {
  return extractor.lines
}

func (extractor *PatternExtractor) Patterns() map[string]string {
  return extractor.patterns
}

func (extractor *PatternExtractor) Extract(lines []*Line, patterns map[string]string) []*Extraction {
  l,p := extractorOperativeParams(extractor, lines, patterns)

  extractions := make([]*Extraction, 0)
  regexps := make(map[string]*regexp.Regexp)

  for attr,pattern := range p {
    reg := regexp.MustCompile(pattern)
    regexps[attr] = reg
  }

  for _,line := range l {
    extract,err := extract(line, regexps)
    if err != nil {
      fmt.Printf("Error extracting information: %v\n", err)
    }

    if extract != nil {
      extractions = append(extractions, extract)
    }
  }
  return extractions
}

func extractorOperativeParams(extractor *PatternExtractor, lines []*Line, patterns map[string]string) ([]*Line, map[string]string) {
  if (lines == nil && extractor.lines == nil) || (patterns == nil && extractor.patterns == nil) {
    panic("Lines and Patterns are required for Extractor to work!")
  }

  var l []*Line
  var p map[string]string

  if extractor.lines == nil {
    l = lines
  } else {
    l = extractor.lines
  }

  if extractor.patterns == nil {
    p = patterns
  } else {
    p = extractor.patterns
  }

  return l,p
}

func extract(line *Line, regexps map[string]*regexp.Regexp) (*Extraction, error) {
  matches := make(map[string]string)

  for attr,regex := range regexps {
    match := regex.FindAllStringSubmatch(line.Content,-1)
    if len(match) > 0 {      
      if split := strings.Split(attr, ","); len(split) == 1{
        matches[attr] = match[0][1]

      } else if len(match) != len(split) {
        return nil,  fmt.Errorf(fmt.Sprintf("Pattern dont match required groups. Expecting [%v], got [%v]\n", len(split), len(match)))

      } else {
        for idx,a := range split {
          matches[a] = match[idx][idx+1]
        }
      }
    }
  }

  if len(matches) > 0 {
    extraction := Extraction{line, &matches}
    return &extraction, nil
  }

  return nil, fmt.Errorf("Could not extract any information [%v] ~= [%v]", line.Content, regexps)
}