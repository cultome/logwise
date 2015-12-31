# Logwise
Utility to automate information extraction from logs.

Its common for me to be in front of a console grep-ing log files to extract pieces of information to feed another grep to get another piece of information, an so and so. This utility allows me to create commands to automate this workflow.

## Usage

To automate a work flow, you create a ```NewFlow``` and add a pipeline-like of commands. Basically it does 3 things:
  * Filter lines
  * Extract information from those lines
  * Convert the extracted information

For example, the following script does:
  * Look for the pattern ```<awbNbr>794666000437</awbNbr>``` in the file ```logs/orderReqRes.log```.
  * Look for the pattern ```itemnumber="794666000437"``` in the file ```logs/invReqRes.log```.
  * The resulting lines (of both filters) are writen in the file ```logs/real_case.log``` (deleting its contents if any using the flag ```append``` in false).
  * In the resulting lines, match the regexp ```invoices - \\[([\\d]+) ->]```, storing the capture group within the name ```txId```.
  * Apply a transformation to the ```txId```, in this case only surrounds the value of ```txId``` with prefix ```[``` and sufix ``` <-]```.
  * With the result of this transformation, filter the file ```logs/invReqRes.log```.
  * Write the resulting lines in the file ```logs/real_case.log``` (this time appending to the existing content using the flag ```append``` in true).
  * Using result of this transformation, match the pattern ```<folio>([\\d]+)</folio>``` and put the results within the name ```folio```.
  * Apply a custom transformation to the filter's results, to convert ```folio``` into ```<InvoiceNumber>([A-Z]+%v)</InvoiceNumber>```
  * Using this transformed pattern, filter the file ```logs/automaticTasks.log```
  * With the resulting lines, extract the content of these lines, that mean, extract the lines between and initial pattern and end pattern that contains the line. You can think of this like a ```grep file.log -A10 -B10``` but instead of line numers we use patterns.
  * Finally append this "context extraction" into the ```logs/real_case.log```

```go
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
```

At the end of the process you have all the information in the file ```logs/real_case.log```, so you can analize it easily.

You can find a couple more examples in ```flow_test.go```

## TODO
 
 * Process files in remote servers using ssh
