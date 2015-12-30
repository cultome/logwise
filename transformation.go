package logwise

import (
  "fmt"
)

type SurroundStringTransformation struct {
  Key, Prefix, Postfix string
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