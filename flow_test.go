package logwise

import "testing"

func TestFilterExtractorFlow(t *testing.T){
  flow := NewFlow(
    NewLineFilter().Set([]string{"logs/nohup.out"}, []string{"\\[\\+] /invoiceOrder.do"}),
    NewPatternExtractor().SetPatterns(map[string]string {"path,role,user": "/([\\w]+).do | ([\\d]+) | ([\\d]+)$"}),
  )
  flow.Start()
}

func TestFilterExtractorFilterFlow(t *testing.T){
  flow := NewFlow(
    NewLineFilter().Set([]string{"logs/nohup.out"}, []string{"invoices - \\[[\\d]+ ->]"}),
    NewPatternExtractor().SetPatterns(map[string]string {"txId": "invoices - \\[([\\d]+) ->]"}),
    NewLineFilter().SetFiles([]string{"logs/nohup.out"}),
  )
  flow.Start()
}

func TestFilterExtractorTransformationFilterFlow(t *testing.T){
  flow := NewFlow(
    NewLineFilter().Set([]string{"logs/nohup.out"}, []string{"invoices - \\[[\\d]+ ->]"}),
    NewPatternExtractor().SetPatterns(map[string]string {"txId": "invoices - \\[([\\d]+) ->]"}),
    NewSurroundStringTransformation("\\[", " <-]"),
    NewLineFilter().SetFiles([]string{"logs/invReqRes.log", "logs/invReqRes.log1"}),
  )
  flow.Start()
}

func TestFilterExtractorTransformationFilterWriterFlow(t *testing.T){
  flow := NewFlow(
    NewLineFilter().Set([]string{"logs/nohup.out"}, []string{"invoices - \\[[\\d]+ ->]"}),
    NewPatternExtractor().SetPatterns(map[string]string {"txId": "invoices - \\[([\\d]+) ->]"}),
    NewSurroundStringTransformation("\\[", " <-]"),
    NewLineFilter().SetFiles([]string{"logs/invReqRes.log", "logs/invReqRes.log1"}),
    NewFileWriter("logs/responses.log"),
  )
  flow.Start()
}