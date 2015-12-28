package logwise

import "testing"

func TestBasicFlow(t *testing.T){
  flow := NewFlow(
    NewLineFilter().Set([]string{"nohup.out"}, []string{"\\[\\+] /invoiceOrder.do"}),
    NewPatternExtractor().SetPatterns(map[string]string {"path,role,user": "/([\\w]+).do | ([\\d]+) | ([\\d]+)$"}),
  )
  flow.Start()
}

func TestBasicFlow2(t *testing.T){
  flow := NewFlow(
    NewLineFilter().Set([]string{"nohup.out"}, []string{"invoices - \\[[\\d]+ ->]"}),
    NewPatternExtractor().SetPatterns(map[string]string {"txId": "invoices - \\[([\\d]+) ->]"}),
    NewLineFilter().SetFiles([]string{"nohup.out"}),
  )
  flow.Start()
}