package logwise

import "testing"
//import "fmt"

func TestFilterOneFileOnePattern(t *testing.T){
  filter := NewLineFilter()
  lines := filter.Filter([]string{"nohup.out"}, []string{"\\[\\+] /invoiceOrder.do"})
  //lines := filter.Filter([]string{"nohup.out"}, []string{"invoices - \\[[\\d]+ ->]"})

  if lines == nil || len(lines) <= 0 {
    t.Error("Lines were not found!")
  }
}
