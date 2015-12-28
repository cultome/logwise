package logwise

import (
  "fmt"
)

type SurroundStringTransformation struct {
  Prefix, Postfix string
}

func NewSurroundStringTransformation(prefix, postfix string) *SurroundStringTransformation {
  return &SurroundStringTransformation{prefix, postfix}
}

func (trans *SurroundStringTransformation) Transform(value string) string {
  return fmt.Sprintf("%v%v%v", trans.Prefix, value, trans.Postfix)
}