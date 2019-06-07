package main

import (
	"time"

	"github.com/eaciit/toolkit"
	au "github.com/logrusorgru/aurora"
)

func getLogger(title string) *toolkit.LogEngine {
	logger := toolkit.NewLogEngine(true, false, "", "", "")

	logger.SetStdoutTemplate(func(item toolkit.LogItem) string {
		pattern := au.BgBrightBlue(au.Bold(au.White(" " + title + " "))).String()
		pattern += " " + au.White(time.Now().Format(time.RFC3339)).String() + " "
		switch item.LogType {
		case "ERROR":
			pattern += au.Bold(au.BrightRed("ERR ")).String()
		case "WARNING":
			pattern += au.Bold(au.Yellow("WRN ")).String()
		case "INFO":
			pattern += au.Bold(au.BrightCyan("INF ")).String()
		case "DEBUG":
			pattern += au.Bold(au.Green("DBG ")).String()
		}
		pattern += item.Msg
		return pattern
	})

	return logger
}
