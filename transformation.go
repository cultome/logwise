package logwise

import (
  "fmt"
)

type Transformation interface {
  Transform(key,value string) string
}

type SurroundStringTransformation struct {
  Key, Prefix, Postfix string
}

type CustomTransformation struct {
  Key string
  Fnc func(string) string
}

func NewSurroundStringTransformation(key, prefix, postfix string) *SurroundStringTransformation {
  return &SurroundStringTransformation{key, prefix, postfix}
}

func (trans *SurroundStringTransformation) Transform(key,value string) string {
  if trans.Key == key {
    return fmt.Sprintf("%v%v%v", trans.Prefix, value, trans.Postfix)
  }
  return value
}

func NewCustomTransformation(key string, fnc func(string) string) *CustomTransformation {
  return &CustomTransformation{key, fnc}
}

func (trans *CustomTransformation) Transform(key,value string) string {
  if key == trans.Key {
    return trans.Fnc(value)
  }
  return value
}