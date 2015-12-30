package logwise

import "testing"
import "fmt"

func TestGetLineContext(t *testing.T){
  ctx := NewLineContext("tasks - \\[\\*] Message", "INFO ")
  line := Line{25565, "<InvoiceNumber>FOCCA934</InvoiceNumber>", "logs/automaticTasks.log", "<InvoiceNumber>"}
  lines := ctx.Get(&line)

  for _,l := range lines {
    fmt.Printf("%v\n", l)
  }
}