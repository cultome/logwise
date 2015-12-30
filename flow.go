package logwise

import (
)

type Flow struct {
  Steps []interface{}
}

func NewFlow(steps ...interface{}) *Flow {
  return &Flow{steps}
}

func (flow *Flow) Start() {
  var prevStep interface{}
  var prevResult interface{}

  for _,step := range flow.Steps {
    switch step := step.(type) {
    case *PatternExtractor:
      prevResult = callExtractor(step, prevStep, prevResult)
      prevStep = step
    case *LineFilter:
      prevResult = callFilter(step, prevStep, prevResult)
      prevStep = step

    case *CustomTransformation:
      prevResult = callTransformation(step, prevStep, prevResult)
    case *SurroundStringTransformation:
      prevResult = callTransformation(step, prevStep, prevResult)

    case *FileWriter:
      callWriter(step, prevStep, prevResult)
    case *LineContext:
      prevResult = callContext(step, prevStep, prevResult)
    }
  }
}

func callContext(ctx *LineContext, prevStep interface{}, prevResult interface{}) interface{} {
  switch result := prevResult.(type){
  case []*Extraction:
    panic("Unavailable information to get context from!")
  case []*Line:
      return ctx.Get(result...)
  }

  return prevResult
}

func callWriter(writer *FileWriter, prevStep interface{}, prevResult interface{}) interface{} {
  var lines []string

  switch result := prevResult.(type){
  case []*Extraction:
    for _,e := range result {
      lines = append(lines, e.String())
    }
  case []*Line:
    for _,l := range result {
      lines = append(lines, l.String())
    }
  }

  writer.Write(lines)

  return prevResult
}

func callTransformation(trans Transformation, prevStep interface{}, prevResult interface{}) interface{} {
  if prevStep == nil {
    panic("Unavailable information to transform!")
  } else if _, ok := prevStep.(*LineFilter); ok {
    // lineas
    ls := prevResult.([]*Line)
    for _,line := range ls {
      line.Content = trans.Transform("", line.Content)
    }

  } else if _, ok := prevStep.(*PatternExtractor); ok {
    // lineas y patrones
    prevExtraction := prevResult.([]*Extraction)
    for _,extraction := range prevExtraction {
      for k,v := range *extraction.Matches {
        (*extraction.Matches)[k] = trans.Transform(k,v)
      }
    }
  }

  return prevResult
}

func callFilter(filter *LineFilter, prevStep interface{}, prevResult interface{}) []*Line {
  var lines []*Line

  if prevStep == nil {
    // nada
    lines = filter.Filter(nil, nil)
  } else if _, ok := prevStep.(*PatternExtractor); ok {
    // patrones
    extractions := prevResult.([]*Extraction)
    var valuesForPattern []string
    for _,extraction := range extractions {
      for _,value := range (*extraction.Matches) {
        valuesForPattern = append(valuesForPattern, value)
      }
    }

    lines = filter.Filter(nil, valuesForPattern)

  } else if prevFilter, ok := prevStep.(*LineFilter); ok {
    // archivos y patrones
    ls := filter.Filter(prevFilter.Files, prevFilter.Patterns)
    prevLines := prevResult.([]*Line)
    lines = append(prevLines, ls...)
  }

  return lines
}

func callExtractor(extractor *PatternExtractor, prevStep interface{}, prevResult interface{}) []*Extraction {
  var extraction []*Extraction
  
  if prevStep == nil {
    extraction = extractor.Extract(nil, nil)
  } else if _, ok := prevStep.(*LineFilter); ok {
    // lineas
    ls := prevResult.([]*Line)
    extraction = extractor.Extract(ls, nil)
  } else if prevExtractor, ok := prevStep.(*PatternExtractor); ok {
    // lineas y patrones
    prevExtraction := prevResult.([]*Extraction)
    extra := extractor.Extract(prevExtractor.Lines, prevExtractor.Patterns)
    extraction = append(prevExtraction, extra...)
  }

  return extraction
}