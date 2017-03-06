package logger

import (
	"os"

	. "github.com/tendermint/go-common"
	"github.com/tendermint/log15"
	"path/filepath"
	"fmt"
	"github.com/pkg/errors"
)

var mainHandler log15.Handler
var bypassHandler log15.Handler

func init() {
	ResetWithLogLevel("debug", "")
}

// TODO: srm - either enhance this with logDir or throw it out.
func SetLogLevel(logLevel string) {
	ResetWithLogLevel(logLevel, "")
}

func ResetWithLogLevel(logLevel string, logDir string) (error) {
	// main handler
	handlers := []log15.Handler{}
	mainHandler = log15.LvlFilterHandler(
		getLevel(logLevel),
		log15.StreamHandler(os.Stdout, log15.TerminalFormat()),
	)
	handlers = append(handlers, mainHandler)

	// only provide a logfile if a path has been given. log to terminal anyway
	if (logDir != "") {

		var logFile = "node.log"

		var absoluteLogPath, dirErr = filepath.Abs(logDir)
		var absoluteLogFile, fileErr = filepath.Abs(logFile)

		if fileErr != nil || dirErr != nil {
			var err = Fmt("Given log_dir cannot be parsed [log_dir=%s]\n", logFile)
			fmt.Println(err)
			return errors.New(err)
		}

		if _, err := os.Stat(absoluteLogPath); err != nil {
			var err = Fmt("The given log_dir does not exist [log_dir=%s]\n", absoluteLogPath)
			fmt.Println(err)
			return errors.New(err)
		}

		handlers = append(handlers, log15.Must.FileHandler(absoluteLogFile, log15.TerminalFormat()))
	}

	// bypass handler for not filtering on global logLevel.
	bypassHandler = log15.StreamHandler(os.Stdout, log15.TerminalFormat())
	//handlers = append(handlers, bypassHandler)

	// By setting handlers on the root, we handle events from all loggers.
	log15.Root().SetHandler(handlers)

	return nil
}

// See go-wire/log for an example of usage.
func MainHandler() log15.Handler {
	return mainHandler
}

func BypassHandler() log15.Handler {
	return bypassHandler
}

func New(ctx ...interface{}) log15.Logger {
	return NewMain(ctx...)
}

func NewMain(ctx ...interface{}) log15.Logger {
	return log15.Root().New(ctx...)
}

func NewBypass(ctx ...interface{}) log15.Logger {
	bypass := log15.New(ctx...)
	bypass.SetHandler(bypassHandler)
	return bypass
}

func getLevel(lvlString string) log15.Lvl {
	lvl, err := log15.LvlFromString(lvlString)
	if err != nil {
		Exit(Fmt("Invalid log level %v: %v", lvlString, err))
	}
	return lvl
}

//----------------------------------------
// Exported from log15

var LvlFilterHandler = log15.LvlFilterHandler
var LvlDebug = log15.LvlDebug
var LvlInfo = log15.LvlInfo
var LvlNotice = log15.LvlNotice
var LvlWarn = log15.LvlWarn
var LvlError = log15.LvlError
