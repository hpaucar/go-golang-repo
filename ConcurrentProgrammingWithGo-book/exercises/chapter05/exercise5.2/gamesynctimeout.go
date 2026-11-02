/*
 * Archivo:     gamesynctimeout.go
 *
 * Propósito:   Solución del ejercicio 5.2 del libro *Learn Concurrent
 *              Programming with Go*, que pide agregar un tiempo límite a la
 *              sincronización de jugadores: si no se conectan todos en 10
 *              segundos, la goroutine `timeout` activa `cancel` y despierta a
 *              todos con `cond.Broadcast()` para cancelar la partida.
 *
 * Compilación: go build gamesynctimeout.go
 * Ejecución:   ./gamesynctimeout
 *              (Alternativamente: go run gamesynctimeout.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Se conectan 4 de los 5 jugadores esperados, uno por segundo. A
 *              los 10 segundos cada jugador imprime "Game cancelled". El
 *              programa termina a los 60 segundos.
 *
 * Notas:
 *    1. `main` solo lanza 4 jugadores para que la partida nunca se complete y
 *       se vea la cancelación.
 *    2. La bandera `cancel` se lee y escribe siempre con el mutex de la
 *       variable de condición, así no hay carreras de datos.
 *    3. La espera usa `for ... && !*cancel`, de modo que un jugador sale tanto
 *       si llegan todos como si se cancela.
 */

package main

import (
	"fmt"
	"sync"
	"time"
)

func playerHandler(cond *sync.Cond, playersRemaining *int,
	playerId int, cancel *bool) {
	cond.L.Lock()
	fmt.Println(playerId, ": Connected")
	*playersRemaining--
	if *playersRemaining == 0 {
		cond.Broadcast()
	}
	for *playersRemaining > 0 && !*cancel {
		fmt.Println(playerId, ": Waiting for more players")
		cond.Wait()
	}
	cond.L.Unlock()
	if *cancel {
		fmt.Println(playerId, ": Game cancelled")
	} else {
		fmt.Println("All players connected. Ready player", playerId)
	}
}

func timeout(cond *sync.Cond, cancel *bool) {
	time.Sleep(10 * time.Second)
	cond.L.Lock()
	*cancel = true
	cond.Broadcast()
	cond.L.Unlock()
}

func main() {
	cond := sync.NewCond(&sync.Mutex{})
	cancel := false
	go timeout(cond, &cancel)
	playersInGame := 5
	for i := 0; i < 4; i++ {
		go playerHandler(cond, &playersInGame, i, &cancel)
		time.Sleep(1 * time.Second)
	}
	time.Sleep(60 * time.Second)
}
