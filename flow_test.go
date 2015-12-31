package logwise

import (
  "testing"
  "time"
  "strconv"
  "fmt"
)

func TestFilterExtractorFlow(t *testing.T){
  flow := NewFlow(
    NewLineFilter(NewFileReader("logs/nohup.out"), []string{"\\[\\+] /invoiceOrder.do"}),
    NewPatternExtractor(nil, map[string]string {"path,role,user": "/([\\w]+).do | ([\\d]+) | ([\\d]+)$"}),
  )
  flow.Start()
}

func TestFilterExtractorFilterFlow(t *testing.T){
  flow := NewFlow(
    NewLineFilter(NewFileReader("logs/nohup.out"), []string{"invoices - \\[[\\d]+ ->]"}),
    NewPatternExtractor(nil, map[string]string {"txId": "invoices - \\[([\\d]+) ->]"}),
    NewLineFilter(NewFileReader("logs/nohup.out"), nil),
  )
  flow.Start()
}

func TestFilterExtractorTransformationFilterFlow(t *testing.T){
  flow := NewFlow(
    NewLineFilter(NewFileReader("logs/nohup.out"), []string{"invoices - \\[[\\d]+ ->]"}),
    NewPatternExtractor(nil, map[string]string {"txId": "invoices - \\[([\\d]+) ->]"}),
    NewSurroundStringTransformation("txId", "\\[", " <-]"),
    NewLineFilter(NewFileReader("logs/invReqRes.log", "logs/invReqRes.log1"), nil),
  )
  flow.Start()
}

func TestFilterExtractorTransformationFilterWriterFlow(t *testing.T){
  flow := NewFlow(
    NewLineFilter(NewFileReader("logs/nohup.out"), []string{"invoices - \\[[\\d]+ ->]"}),
    NewPatternExtractor(nil, map[string]string {"txId": "invoices - \\[([\\d]+) ->]"}),
    NewSurroundStringTransformation("txId", "\\[", " <-]"),
    NewLineFilter(NewFileReader("logs/invReqRes.log", "logs/invReqRes.log1"), nil),
    NewFileWriter("logs/responses.log", false, "", ""),
  )
  flow.Start()
}

func TestRealCaseTraceOrder(t *testing.T){
  flow := NewFlow(
    NewLineFilter(NewFileReader("logs/orderReqRes.log"), []string{"<awbNbr>794666000437</awbNbr>"}),
    NewLineFilter(NewFileReader("logs/invReqRes.log"), []string{"itemnumber=\"794666000437\""}),
    NewFileWriter("logs/real_case.log", false, "===================== Invoice order and Invoice Request =====================", " "),

    NewPatternExtractor(nil, map[string]string {"txId": "invoices - \\[([\\d]+) ->]"}),
    NewSurroundStringTransformation("txId", "\\[", " <-]"),
    NewLineFilter(NewFileReader("logs/invReqRes.log"), nil),
    NewFileWriter("logs/real_case.log", true, "===================== Invoice Response =====================", " "),

    NewPatternExtractor(nil, map[string]string {"folio": "<folio>([\\d]+)</folio>"}),
    NewCustomTransformation("folio", func(v string) string {
      return fmt.Sprintf("<InvoiceNumber>([A-Z]+%v)</InvoiceNumber>", v)
    }),
    NewLineFilter(NewFileReader("logs/automaticTasks.log"), nil),
    NewLineContext("tasks - \\[\\*] Message", "INFO   "),
    NewFileWriter("logs/real_case.log", true, "===================== LCCS Transaction =====================", " "),
  )
  flow.Start()
}
