package logwise

import "testing"

func TestFilterOneFileOnePattern(t *testing.T){
  filter := NewLineFilter(NewFileReader("logs/nohup.out"), []string{"\\[\\+] /invoiceOrder.do"})
  lines := filter.Filter(nil, nil)

  if lines == nil || len(lines) <= 0 {
    t.Error("Lines were not found!")
  }
}

func TestFileReader(t *testing.T){
  reader := NewFileReader("logs/one.log","logs/two.log", "logs/three.log")

  if one,_ := reader.Read(); one.Content != "one" {
    t.Errorf("Error reading files. Expecting %q, got %q\n", "one", one.Content)
  }
  
  if two,_ := reader.Read(); two.Content != "two" {
    t.Errorf("Error reading files. Expecting %q, got %q\n", "two", two.Content)
  }
  
  if three,_ := reader.Read(); three.Content != "three" {
    t.Errorf("Error reading files. Expecting %q, got %q\n", "three", three.Content)
  }
  
  if four,_ := reader.Read(); four.Content != "four" {
    t.Errorf("Error reading files. Expecting %q, got %q\n", "four", four.Content)
  }
}
