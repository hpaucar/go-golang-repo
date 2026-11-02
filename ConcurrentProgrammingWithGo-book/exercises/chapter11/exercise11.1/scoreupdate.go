/*
 * Archivo:     scoreupdate.go
 *
 * Propósito:   Código de partida del ejercicio 11.1 del libro *Learn Concurrent
 *              Programming with Go*, que pide encontrar y corregir el
 *              interbloqueo (deadlock). Mil goroutines suman puntos a
 *              subconjuntos aleatorios de jugadores y bloquean el mutex de cada
 *              uno en un orden aleatorio.
 *
 * Compilación: go build scoreupdate.go
 * Ejecución:   ./scoreupdate
 *              (Alternativamente: go run scoreupdate.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Si termina, el puntaje de cada jugador. Lo más probable es que
 *              se detenga con "fatal error: all goroutines are asleep -
 *              deadlock!".
 *
 * Notas:
 *    1. El deadlock ocurre cuando dos goroutines toman los mismos mutex en
 *       orden distinto y cada una espera el que tiene la otra.
 *    2. La solución, en `../exercise11.1solution/scoreupdate.go`, ordena los
 *       jugadores antes de bloquearlos.
 */

package main

import (
    "fmt"
    "math/rand"
    "sync"
)

type Player struct {
    name  string
    score int
    mutex sync.Mutex
}

func incrementScores(players []*Player, increment int) {
    for _, player := range players {
        player.mutex.Lock()
    }
    for _, player := range players {
        player.score += increment
    }
    for _, player := range players {
        player.mutex.Unlock()
    }
}

func main() {
    players := []*Player{
        {"James", 0, sync.Mutex{}},
        {"Ann", 0, sync.Mutex{}},
        {"Paul", 0, sync.Mutex{}},
        {"Isabel", 0, sync.Mutex{}},
        {"Peter", 0, sync.Mutex{}},
        {"Jane", 0, sync.Mutex{}},
    }
    wg := sync.WaitGroup{}
    for i := 0; i < 1000; i++ {
        n := rand.Intn(len(players)) + 1
        rand.Shuffle(len(players), func(i, j int) { players[i], players[j] = players[j], players[i] })
        wg.Add(1)
        sublist := make([]*Player, n)
        copy(sublist, players[:n])
        go func(players []*Player) {
            incrementScores(players, 10)
            wg.Done()
        }(sublist)
    }
    wg.Wait()
    for _, player := range players {
        fmt.Printf("Score for %s is %d\n", player.name, player.score)
    }
}
