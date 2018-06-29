package cmd

import (
	"bytes"
	"fmt"
	"runtime"
	"time"

	"bitbucket.org/umovme/bootstrapper/helpers"
	"github.com/go-playground/ansi"
	"github.com/go-playground/log"
	"github.com/spf13/viper"
)

var (
	loc              *time.Location
	defaultLogLevels = []log.Level{
		log.InfoLevel,
		log.WarnLevel,
		log.ErrorLevel,
		log.PanicLevel,
		log.AlertLevel,
		log.FatalLevel,
	}
)

const (
	defaultFormat string = "2006-01-02 15:04:05.000ms"
	// defaultTZ     string = "America/Sao_Paulo"
	defaultTZ string = "Brazil/EAST"
)

// TextHandler is a log handdle that print simple text messages
type TextHandler struct{}

// Log accepts log entries to be processed
func (c *TextHandler) Log(e log.Entry) {
	b := new(bytes.Buffer)

	b.Reset()
	b.WriteString(e.Timestamp.In(loc).Format(defaultFormat))
	fmt.Fprint(b, " ")
	fmt.Fprintf(b, "%-6s", e.Level.String())
	fmt.Fprint(b, " ")

	for _, f := range e.Fields {
		fmt.Fprint(b, f.Key)
		fmt.Fprint(b, "=")
		fmt.Fprintf(b, "%-10s", f.Value)
		fmt.Fprint(b, " ")
	}
	b.WriteString(e.Message)
	fmt.Println(b.String())
}

// CollorHandler is a log handdle that print colorfull messages
type CollorHandler struct{}

var defaultColors = [...]ansi.EscSeq{
	log.DebugLevel: ansi.Gray,
	// log.TraceLevel:  ansi.White,
	log.InfoLevel:   ansi.Blue,
	log.NoticeLevel: ansi.LightCyan,
	log.WarnLevel:   ansi.LightYellow,
	log.ErrorLevel:  ansi.LightRed,
	log.PanicLevel:  ansi.Red,
	log.AlertLevel:  ansi.Red + ansi.Underline,
	log.FatalLevel:  ansi.Red + ansi.Underline + ansi.Blink,
}

// Log accepts log entries to be processed
func (c *CollorHandler) Log(e log.Entry) {

	color := defaultColors[e.Level]

	b := new(bytes.Buffer)
	b.Reset()
	b.WriteString(e.Timestamp.In(loc).Format(defaultFormat))
	b.WriteString(" ")
	fmt.Fprintf(b, "%s%-6s%s", ansi.Bold+color, e.Level.String(), ansi.BoldOff+ansi.Reset)

	for _, f := range e.Fields {
		fmt.Fprint(b, ansi.Bold)
		fmt.Fprint(b, f.Key)
		fmt.Fprint(b, ansi.BoldOff)
		fmt.Fprint(b, "=")
		fmt.Fprint(b, ansi.Italics)
		fmt.Fprintf(b, "%-10s", f.Value)
		fmt.Fprint(b, ansi.ItalicsOff)
		fmt.Fprint(b, " ")
	}
	b.WriteString(e.Message)

	fmt.Println(b.String())
}

func setHandler() log.Handler {

	// ugly messages on windows forces me to disable this pretty messages
	if !helpers.IsTerminal() || runtime.GOOS == "windows" {
		return new(TextHandler)
	}

	return new(CollorHandler)
}

func setupLogger() {

	var err error
	loc, err = time.LoadLocation(defaultTZ)
	if err != nil {
		fmt.Println("Default TZ", defaultTZ, " Not found. Falling back to local timezone.")
		loc = time.Local
	}

	if viper.GetBool("debug") {
		log.Info("Enabling DEBUG messages.")
		defaultLogLevels = append(defaultLogLevels, log.DebugLevel)
	}

	log.AddHandler(setHandler(), defaultLogLevels...)

}
