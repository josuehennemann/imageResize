package main

import (
	"github.com/josuehennemann/logger"
	"sync"
)

const (
	SYSTEM_STATE_ON      = 1 << iota // 1
	SYSTEM_STATE_OFF                 // 2
	SYSTEM_STATE_PAUSING             // 4
)

type TurnOffSystem struct {
	status int            //status do sistema
	http   sync.WaitGroup //wg de controle dos e-mails de return path
	db     sync.WaitGroup
}

func initTurnOff() *TurnOffSystem {
	t := &TurnOffSystem{}
	t.status = SYSTEM_STATE_ON
	t.http = sync.WaitGroup{}
	Logger.Printf(logger.INFO, "Inicializou a estrutura de controle de desligamento do sistema")
	return t
}

// Lista de metodos que adiciona e remove itens dos wg

func (turn *TurnOffSystem) AddHttp() {
	turn.http.Add(1)
}
func (turn *TurnOffSystem) DoneHttp() {
	turn.http.Done()
}

func (turn *TurnOffSystem) AddDatabase() {
	turn.db.Add(1)
}
func (turn *TurnOffSystem) DoneDatabase() {
	turn.db.Done()
}

//Executa o wait dos waitgroups do sistema
func (turn *TurnOffSystem) Wait() {
	turn.http.Wait()
	turn.db.Wait()
}

//altera o status do sistemas
func (turn *TurnOffSystem) SetSystemState(s int) {
	turn.status = s
}

//verifica se o sistema esta em desligamento
func (turn *TurnOffSystem) IsShutdown() bool {
	return turn.checkState(SYSTEM_STATE_OFF | SYSTEM_STATE_PAUSING)
}

//verifica se o sistema já esta up
func (turn *TurnOffSystem) IsAlready() bool {
	return turn.checkState(SYSTEM_STATE_ON)
}

//função interna que valida o status
func (turn *TurnOffSystem) checkState(s int) bool {
	return turn.status&s != 0
}
