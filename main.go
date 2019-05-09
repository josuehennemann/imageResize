package main

import (
	"flag"
	"fmt"
	"github.com/josuehennemann/logger"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	//faz o parse da variaveis que sao passadas para o binario

	flag.Parse()
	initService()
	//trava a execução da main
	select {}
}

//verifica se deu erro ou nao, em caso de erro, printa o erro e mata o script
func checkErrorAndKillMe(e error) {
	if e == nil {
		return
	}
	fmt.Println(e.Error())
	os.Exit(2)
}

func initService() {
	ServiceName = "image-resize"

	var err error

	//inicia o processo para subir o binario como serviço no linux
	//	procDaemon, err = daemon.Daemonize(*pidPath)
	checkErrorAndKillMe(err)

	//cria a goroutine que sabe tratar o desligamento do serviço
	go killMeSignal()

	_init()

	//inicia o servidor http
	go startHttpServer()

}

//inicializa as variaveis que podem ser utilizadas caso rode sem ser serviço
func _init() {
	//carrega para a memoria o arquivo de inicialização
	err := initConfig()
	checkErrorAndKillMe(err)
	setOutput()
	//inicia o arquivo de log
	Logger, err = logger.New(config.LogPath+ServiceName+".log", logger.LEVEL_PRODUCTION, true)
	checkErrorAndKillMe(err)
	turnoff = initTurnOff()
	Access, err = logger.New(config.LogPath+ServiceName+"-Access.log", logger.LEVEL_PRODUCTION, true)
	checkErrorAndKillMe(err)
}

func setOutput() {
	if config.IsDev {
		return
	}

	dirLog := filepath.Dir(config.LogPath)
	if dirLog == "." {
		dirLog = ""
	} else {
		dirLog += "/"
	}

	//IMPORTANTE ISSO NAO FUNCIONA NO WINDOWS
	daemon, err := os.OpenFile(dirLog+ServiceName+"-Daemon.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
	checkErrorAndKillMe(err)
	syscall.Dup2(int(daemon.Fd()), 1)
	syscall.Dup2(int(daemon.Fd()), 2)
}

func killMeSignal() {
	c := make(chan os.Signal)
	signal.Notify(c)
	for {
		s := <-c
		if s == syscall.SIGTERM || s == syscall.SIGINT {
			Logger.Printf(logger.INFO, "Sinal de desligamento recebido")
			turnoff.SetSystemState(SYSTEM_STATE_OFF)
			turnoff.Wait() //espera os wait groups terminarem

			//Fecha o arquivo de log
			Logger.Close()
			Logger.Printf(logger.INFO, "Pronto para morrer. Signal[%v]", s)
			Access.Close()
			os.Exit(0)
		}
	}
	return
}

func PrepareClose(w http.ResponseWriter, r *http.Request) {
	r.Close = true
	r.Body.Close()
	return
}

func startHttpServer() {
	// HTTP LISTEN

	http.HandleFunc("/resize", HttpResizeImage)

	http.HandleFunc("/", http.NotFound)
	if tr, ok := http.DefaultTransport.(*http.Transport); ok {
		tr.DisableKeepAlives = true
		tr.MaxIdleConnsPerHost = 1
		tr.CloseIdleConnections()
	}
	Logger.Printf(logger.INFO, "Iniciando serviço [%s] ...", config.HttpAddress)

	//serviço na porta definida
	server := &http.Server{Addr: config.HttpAddress, ReadTimeout: 2 * time.Second, WriteTimeout: 5 * time.Second}
	err := server.ListenAndServe()
	checkErrorAndKillMe(err)
}
