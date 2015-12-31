package logwise

import (
  "testing"
  "time"
  "strconv"
//  "fmt"
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
    NewFileWriter("logs/responses.log", false),
  )
  flow.Start()
}

func TestRealCaseTraceOrder(t *testing.T){
  flow := NewFlow(
    NewLineFilter(NewFileReader("logs/orderReqRes.log"), []string{"<awbNbr>794666000437</awbNbr>"}),
    NewLineFilter(NewFileReader("logs/invReqRes.log"), []string{"itemnumber=\"794666000437\""}),
    NewFileWriter("logs/real_case.log", false), //.AddPrefix("===================== Invoice order and Invoice Request ====================="),

    NewPatternExtractor(nil, map[string]string {"txId": "invoices - \\[([\\d]+) ->]"}),
    NewSurroundStringTransformation("txId", "\\[", " <-]"),
    NewLineFilter(NewFileReader("logs/invReqRes.log"), nil),
    NewFileWriter("logs/real_case.log", true), //.AddPrefix("\n\n===================== Invoice Response ====================="),

    NewPatternExtractor(nil, map[string]string {"folio": "<folio>([\\d]+)</folio>"}),
    NewSurroundStringTransformation( "folio", "", "</InvoiceNumber>"),
    NewLineFilter(NewFileReader("logs/automaticTasks.log"), nil),
    NewLineContext("tasks - \\[\\*] Message", "INFO   "),
    NewFileWriter("logs/real_case.log", true), //.AddPrefix("\n\n===================== LCCS Transaction ====================="),
  )
  flow.Start()
}

func tenSecondsLater(value string) string {
  // 2015-12-24 13:07:00,241
  year,_ := strconv.Atoi(value[:4])
  month,_ := strconv.Atoi(value[5:7])
  day,_ := strconv.Atoi(value[8:10])

  hour,_ := strconv.Atoi(value[11:13])
  minute,_ := strconv.Atoi(value[14:16])
  second,_ := strconv.Atoi(value[17:19])

  location,_ := time.LoadLocation("America/Mexico_City")

  date := time.Date(year, time.Month(month), day, hour, minute, second, 0, location)

  after := date.Add(10 * time.Second)

  return after.String()[:18]
}