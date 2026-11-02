/*
 * Archivo:     gamesync.go
 *
 * Propósito:   Usar `cond.Broadcast()` para sincronizar jugadores: se conectan
 *              cuatro jugadores, uno por segundo, y cada uno espera con
 *              `cond.Wait()` hasta que estén todos. El último jugador despierta
 *              a todos los que esperan.
 *
 * Compilación: go build gamesync.go
 * Ejecución:   ./gamesync
 *              (Alternativamente: go run gamesync.go)
 *
 * Entrada:     Ninguna.
 * Salida:      "N : Connected" y "N : Waiting for more players" para los tres
 *              primeros jugadores; al conectarse el cuarto, cuatro líneas "All
 *              players connected. Ready player N".
 *
 * Notas:
 *    1. `Broadcast()` despierta a todas las goroutines que esperan, a
 *       diferencia de `Signal()`, que despierta solo a una.
 *    2. Cada jugador comprueba `*playersRemaining > 0` en un bucle `for`, así
 *       que no importa si llega antes o después de la señal.
 *    3. `main` no espera explícitamente a los jugadores: termina 1 segundo
 *       después de lanzar al último, tiempo que normalmente basta para ver
 *       todos los mensajes.
 */

package main

import (
    "fmt"
    "sync"
    "time"
)

func main() {
    cond := sync.NewCond(&sync.Mutex{})
    playersInGame := 4
    for playerId := 0; playerId < 4; playerId++ {
        go playerHandler(cond, &playersInGame, playerId)
        time.Sleep(1 * time.Second)
    }
}

func playerHandler(cond *sync.Cond, playersRemaining *int, playerId int) {
    cond.L.Lock()
    fmt.Println(playerId, ": Connected")
    *playersRemaining--
    if *playersRemaining == 0 {
        cond.Broadcast()
    }
    for *playersRemaining > 0 {
        fmt.Println(playerId, ": Waiting for more players")
        cond.Wait()
    }
    cond.L.Unlock()
    fmt.Println("All players connected. Ready player", playerId)
    //Game started
}
