package logwise

import "testing"

func TestFilterOneFileOnePattern(t *testing.T){
  filter := NewLineFilter()
  lines := filter.Filter([]string{"logs/nohup.out"}, []string{"\\[\\+] /invoiceOrder.do"})

  if lines == nil || len(lines) <= 0 {
    t.Error("Lines were not found!")
  }
}
