/*
 * Archivo:     scoreupdate.go
 *
 * Propósito:   Solución del ejercicio 11.1 del libro *Learn Concurrent
 *              Programming with Go*: elimina el interbloqueo (deadlock)
 *              ordenando los jugadores por nombre con `sort.Slice` antes de
 *              tomar sus mutex. Así todas las goroutines bloquean los recursos
 *              en el mismo orden global.
 *
 * Compilación: go build scoreupdate.go
 * Ejecución:   ./scoreupdate
 *              (Alternativamente: go run scoreupdate.go)
 *
 * Entrada:     Ninguna.
 * Salida:      El puntaje final de cada uno de los 6 jugadores. La suma de los
 *              puntajes cambia en cada ejecución, pero el programa siempre
 *              termina.
 *
 * Notas:
 *    1. Ordenar los recursos rompe la condición de espera circular, una de las
 *       cuatro condiciones de Coffman.
 *    2. Cada goroutine recibe una copia del sub-slice, por lo que ordenarlo no
 *       afecta al slice `players` de `main`.
 */

package main

import (
    "fmt"
    "math/rand"
    "sort"
    "sync"
)

type Player struct {
    name  string
    score int
    mutex sync.Mutex
}

func incrementScores(players []*Player, increment int) {
    sort.Slice(players, func(a, b int) bool {
        return players[a].name < players[b].name
    })
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
