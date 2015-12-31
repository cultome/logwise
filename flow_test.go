package logwise

import (
  "testing"
  "fmt"
)

func TestFilterExtractorFlow(t *testing.T){
  NewFlow(
    NewLineFilter([]string{"\\[\\+] /invoiceOrder.do"}, "logs/nohup.out"),
    NewPatternExtractor(nil, map[string]string {"path,role,user": "/([\\w]+).do | ([\\d]+) | ([\\d]+)$"}),
  ).Start()
}

func TestFilterExtractorFilterFlow(t *testing.T){
  NewFlow(
    NewLineFilter([]string{"invoices - \\[[\\d]+ ->]"}, "logs/nohup.out"),
    NewPatternExtractor(nil, map[string]string {"txId": "invoices - \\[([\\d]+) ->]"}),
    NewLineFilter(nil, "logs/nohup.out"),
  ).Start()
}

func TestFilterExtractorTransformationFilterFlow(t *testing.T){
  NewFlow(
    NewLineFilter([]string{"invoices - \\[[\\d]+ ->]"}, "logs/nohup.out"),
    NewPatternExtractor(nil, map[string]string {"txId": "invoices - \\[([\\d]+) ->]"}),
    NewSurroundStringTransformation("txId", "\\[", " <-]"),
    NewLineFilter(nil, "logs/invReqRes.log", "logs/invReqRes.log1"),
  ).Start()
}

func TestFilterExtractorTransformationFilterWriterFlow(t *testing.T){
  NewFlow(
    NewLineFilter([]string{"invoices - \\[[\\d]+ ->]"}, "logs/nohup.out"),
    NewPatternExtractor(nil, map[string]string {"txId": "invoices - \\[([\\d]+) ->]"}),
    NewSurroundStringTransformation("txId", "\\[", " <-]"),
    NewLineFilter(nil, "logs/invReqRes.log", "logs/invReqRes.log1"),
    NewFileWriter("logs/responses.log", false, "", ""),
  ).Start()
}

func TestRealCaseTraceOrder(t *testing.T){
  NewFlow(
    NewLineFilter([]string{"<awbNbr>794666000437</awbNbr>"}, "logs/orderReqRes.log"),
    NewLineFilter([]string{"itemnumber=\"794666000437\""}, "logs/invReqRes.log"),
    NewFileWriter("logs/real_case.log", false, "===================== Invoice order and Invoice Request =====================", " "),

    NewPatternExtractor(nil, map[string]string {"txId": "invoices - \\[([\\d]+) ->]"}),
    NewSurroundStringTransformation("txId", "\\[", " <-]"),
    NewLineFilter(nil, "logs/invReqRes.log"),
    NewFileWriter("logs/real_case.log", true, "===================== Invoice Response =====================", " "),

    NewPatternExtractor(nil, map[string]string {"folio": "<folio>([\\d]+)</folio>"}),
    NewCustomTransformation("folio", func(v string) string {
      return fmt.Sprintf("<InvoiceNumber>([A-Z]+%v)</InvoiceNumber>", v)
    }),
    NewLineFilter(nil, "logs/automaticTasks.log"),
    NewLineContext("tasks - \\[\\*] Message", "INFO   "),
    NewFileWriter("logs/real_case.log", true, "===================== LCCS Transaction =====================", " "),
  ).Start()
}

/*
func TestRealCaseCheckRemoteLogs(t *testing.T){
  NewFlow(
    NewServerSession("server:22", "user", "pass"),
    NewGrepFilter([]string{"<awbNbr>794666000437</awbNbr>"}, "$INC_LOG/orderReqRes.log"),
    NewGrepFilter([]string{"itemnumber=\"794666000437\""}, "$INC_LOG/invReqRes.log"),
    NewFileWriter("logs/real_case.log", false, "===================== Invoice order and Invoice Request =====================", " "),

    NewPatternExtractor(nil, map[string]string {"txId": "invoices - \\[([\\d]+) ->]"}),
    NewSurroundStringTransformation("txId", "\\[", " <-]"),
    NewGrepFilter(nil, "$INC_LOG/invReqRes.log"),
    NewFileWriter("logs/real_case.log", true, "===================== Invoice Response =====================", " "),

    NewPatternExtractor(nil, map[string]string {"folio": "<folio>([\\d]+)</folio>"}),
    NewCustomTransformation("folio", func(v string) string {
      return fmt.Sprintf("<InvoiceNumber>([A-Z]+%v)</InvoiceNumber>", v)
    }),
    NewGrepFilter(nil, "$INC_LOG/automaticTasks.log"),
    NewLineContext("tasks - \\[\\*] Message", "INFO   "),
    NewFileWriter("logs/real_case.log", true, "===================== LCCS Transaction =====================", " "),
  ).Start()
}
*/