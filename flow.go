package logwise

import (
  "fmt"
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
    case Extractor:
      prevResult = callExtractor(step, prevStep, prevResult)
    case Filter:
      prevResult = callFilter(step, prevStep, prevResult)
    }

    prevStep = step
  }

  report(prevResult)
}

func report(result interface{}) {
  switch result := result.(type){
  case []*Extraction:
    for _,e := range result {
      fmt.Printf("%v\n", e)
    }
  case []*Line:
    for _,l := range result {
      fmt.Printf("%v\n", l)
    }
  }
}

func callFilter(filter Filter, prevStep interface{}, prevResult interface{}) []*Line {
  var lines []*Line

  if prevStep == nil {
    // nada
    lines = filter.Filter(nil, nil)
  } else if _, ok := prevStep.(Extractor); ok {
    // patrones
    extractions := prevResult.([]*Extraction)
    var valuesForPattern []string

    for _,extraction := range extractions {
      for _,value := range (*extraction.Matches) {
        valuesForPattern = append(valuesForPattern, fmt.Sprintf("\\b%v\\b", value))
      }
    }

    lines = filter.Filter(nil, valuesForPattern)

  } else if prevFilter, ok := prevStep.(Filter); ok {
    // archivos y patrones
    ls := filter.Filter(prevFilter.Files(), prevFilter.Patterns())
    prevLines := prevResult.([]*Line)
    lines = append(prevLines, ls...)
  }

  return lines
}

func callExtractor(extractor Extractor, prevStep interface{}, prevResult interface{}) []*Extraction {
  var extraction []*Extraction
  
  if prevStep == nil {
    extraction = extractor.Extract(nil, nil)
  } else if _, ok := prevStep.(Filter); ok {
    // lineas
    ls := prevResult.([]*Line)
    extraction = extractor.Extract(ls, nil)
  } else if prevExtractor, ok := prevStep.(Extractor); ok {
    // lineas y patrones
    prevExtraction := prevResult.([]*Extraction)
    extra := extractor.Extract(prevExtractor.Lines(), prevExtractor.Patterns())
    extraction = append(prevExtraction, extra...)
  }

  return extraction
}