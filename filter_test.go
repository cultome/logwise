package logwise

import "testing"

func TestFilterOneFileOnePattern(t *testing.T){
  filter := NewLineFilter([]string{"logs/nohup.out"}, []string{"\\[\\+] /invoiceOrder.do"})
  lines := filter.Filter(nil, nil)

  if lines == nil || len(lines) <= 0 {
    t.Error("Lines were not found!")
  }
}
