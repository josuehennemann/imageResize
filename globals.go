package main

import (
	"flag"
	"github.com/josuehennemann/logger"
	"net/http"
	"runtime/debug"
	"strings"
)

var (
	ServiceName string
	fileConf    = flag.String("fileConf", "", "Endereco do arquivo de configuração")
	config      *Config
	Logger      *logger.Logger
	Access      *logger.Logger
	turnoff     *TurnOffSystem
)

func recoverPanic() {
	if rec := recover(); rec != nil {
		stack := debug.Stack()
		Logger.WritePanic(rec, stack)
		return
	}
}

func getVariablePost(r *http.Request, k string) string {
	return strings.TrimSpace(r.FormValue(k))
}
