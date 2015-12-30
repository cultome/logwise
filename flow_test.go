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
    NewSurroundStringTransformation("txId", "\\[", " <-]"),
    NewLineFilter().SetFiles([]string{"logs/invReqRes.log", "logs/invReqRes.log1"}),
  )
  flow.Start()
}

func TestFilterExtractorTransformationFilterWriterFlow(t *testing.T){
  flow := NewFlow(
    NewLineFilter().Set([]string{"logs/nohup.out"}, []string{"invoices - \\[[\\d]+ ->]"}),
    NewPatternExtractor().SetPatterns(map[string]string {"txId": "invoices - \\[([\\d]+) ->]"}),
    NewSurroundStringTransformation("txId", "\\[", " <-]"),
    NewLineFilter().SetFiles([]string{"logs/invReqRes.log", "logs/invReqRes.log1"}),
    NewFileWriter("logs/responses.log", false),
  )
  flow.Start()
}

func TestRealCase(t *testing.T){
  flow := NewFlow(
    NewLineFilter().Set([]string{"logs/orderReqRes.log"}, []string{"<awbNbr>794670441143</awbNbr>"}),
    NewLineFilter().Set([]string{"logs/invReqRes.log"}, []string{"itemnumber=\"794670441143\""}),
    NewFileWriter("tmp.log", false),
    NewPatternExtractor().SetPatterns(map[string]string {"txId": "invoices - \\[([\\d]+) ->]", "serie": "serie=\"([\\w]+)\""}),
    NewSurroundStringTransformation( "txId", "\\[", " <-]"),
    NewFileWriter("tmp.log", true),
    NewLineFilter().SetFiles([]string{"logs/invReqRes.log"}),
    NewFileWriter("tmp.log", true),
    NewPatternExtractor().SetPatterns(map[string]string {"folio": "<folio>([\\d]+)</folio>"}),
    NewFileWriter("tmp.log", true),
    //NewConcatenateTransformation("<InvoiceNumber>", "${serie}", "${folio}", "</InvoiceNumber>"),
    //NewLineFilter().SetFiles([]string{"logs/automaticTask.log"}),
    //NewFileWriter("tmp.log", true),
  )
  flow.Start()
}