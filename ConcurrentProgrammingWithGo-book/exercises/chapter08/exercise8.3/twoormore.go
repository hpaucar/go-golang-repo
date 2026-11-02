/*
 * Archivo:     twoormore.go
 *
 * Propósito:   Solución del ejercicio 8.3 del libro *Learn Concurrent
 *              Programming with Go*, que pide un juego donde cuatro jugadores
 *              envían movimientos por sus canales y la partida termina cuando
 *              queda uno solo. Cuando un jugador cierra su canal, se reemplaza
 *              por `nil` para que el `select` lo ignore.
 *
 * Compilación: go build twoormore.go
 * Ejecución:   ./twoormore
 *              (Alternativamente: go run twoormore.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Movimientos "Player N: UP/DOWN/LEFT/RIGHT", mensajes "Player N
 *              left the game" con los jugadores restantes y, al final, "Game
 *              finished".
 *
 * Notas:
 *    1. Un caso de `select` sobre un canal `nil` nunca se elige; es la forma
 *       estándar de desactivar un caso.
 *    2. `move, moreData := <-ch` permite distinguir un valor real del cierre
 *       del canal.
 *    3. Cada jugador envía un número aleatorio de movimientos (0 a 99) con
 *       pausas aleatorias, así que la salida cambia en cada ejecución.
 */

package main

import (
    "fmt"
    "math/rand"
    "time"
)

func player() chan string {
    output := make(chan string)
    count := rand.Intn(100)
    move := []string{"UP", "DOWN", "LEFT", "RIGHT"}
    go func() {
        defer close(output)
        for i := 0; i < count; i++ {
            output <- move[rand.Intn(4)]
            d := time.Duration(rand.Intn(200))
            time.Sleep(d * time.Millisecond)
        }
    }()
    return output
}

func handlePlayer(id int, moreData bool, move string,
    players []chan string, totalPlayers *int) {
    if moreData {
        fmt.Printf("Player %d: %s\n", id, move)
    } else {
        players[id] = nil
        *totalPlayers--
        fmt.Printf("Player %d left the game. Remaining players: %d\n", id, *totalPlayers)
    }
}

func main() {
    players := []chan string {player(), player(), player(), player()}
    totalPlayers := 4
    for totalPlayers > 1 {
        select {
        case move, moreData := <-players[0]:
            handlePlayer(0, moreData, move, players, &totalPlayers)
        case move, moreData := <-players[1]:
            handlePlayer(1, moreData, move, players, &totalPlayers)
        case move, moreData := <-players[2]:
            handlePlayer(2, moreData, move, players, &totalPlayers)
        case move, moreData := <-players[3]:
            handlePlayer(3, moreData, move, players, &totalPlayers)
        }
    }
    fmt.Println("Game finished")
}

