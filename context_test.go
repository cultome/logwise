package logwise

import "testing"

func TestGetLineContext(t *testing.T){
  ctx := NewLineContext("tasks - \\[\\*] Message", "INFO ")
  line := Line{25565, "<InvoiceNumber>FOCCA934</InvoiceNumber>", "logs/automaticTasks.log", "<InvoiceNumber>"}
  lines := ctx.Get(&line)

  if len(lines) <= 0 {
    t.Error("Did not extract context from file")
  }
}