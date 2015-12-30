package logwise

import "testing"
//import "fmt"

func TestExtractOnePattern(t *testing.T){
  line := Line{50426, "INFO   2015-12-24 14:28:16,788 resourceAccess - [+] /invoiceOrder.do | 30 | 910744", "logs/nohup.out", "/([\\w]+).do"}
  extractor := NewPatternExtractor([]*Line{&line}, map[string]string {"path": "/([\\w]+).do"})

  extractions := extractor.Extract(nil, nil)
  if len(extractions) != 1 {
    t.Errorf("Wrong extractions. Expecting [1], got [%v]\n", len(extractions))
  }

  extraction := extractions[0]
  if _, ok := (*extraction.Matches)["path"]; !ok {
    t.Error("Wrong extractions. Expecting match \"path\" not found")
  }
}

func TestExtractMultiplePatterns(t *testing.T){
  line := Line{50426, "INFO   2015-12-24 14:28:16,788 resourceAccess - [+] /invoiceOrder.do | 30 | 910744", "logs/nohup.out", "/([\\w]+).do"}
  extractor := NewPatternExtractor([]*Line{&line}, map[string]string {"path": "/([\\w]+).do", "role": "\\| ([\\d]+) \\|", "user": "([\\d]+)$"})

  extractions := extractor.Extract(nil, nil)
  if len(extractions) <= 0 {
    t.Error("No extractions were made! Expecting [1]")
  }

  extraction := extractions[0]
  _, ok1 := (*extraction.Matches)["path"];
  _, ok2 := (*extraction.Matches)["role"];
  _, ok3 := (*extraction.Matches)["user"];

  if !ok1 || !ok2 || !ok3 {
    t.Errorf("Wrong extractions. Expecting matches on \"path\" [%v], \"role\" [%v] and \"user\" [%v]\n", ok1, ok2, ok3)
  }
}

func TestExtractSinglePatternMultipleGroups(t *testing.T){
  line := Line{50426, "INFO   2015-12-24 14:28:16,788 resourceAccess - [+] /invoiceOrder.do | 30 | 910744", "logs/nohup.out", "/([\\w]+).do"}
  extractor := NewPatternExtractor([]*Line{&line}, map[string]string {"path,role,user": "/([\\w]+).do | ([\\d]+) | ([\\d]+)$"})

  extractions := extractor.Extract(nil, nil)
  if len(extractions) <= 0 {
    t.Error("No extractions were made! Expecting [1]")
  }

  extraction := extractions[0]

  _, ok1 := (*extraction.Matches)["path"];
  _, ok2 := (*extraction.Matches)["role"];
  _, ok3 := (*extraction.Matches)["user"];

  if !ok1 || !ok2 || !ok3 {
    t.Errorf("Wrong extractions. Expecting matches on \"path\" [%v], \"role\" [%v] and \"user\" [%v]\n", ok1, ok2, ok3)
  }
}