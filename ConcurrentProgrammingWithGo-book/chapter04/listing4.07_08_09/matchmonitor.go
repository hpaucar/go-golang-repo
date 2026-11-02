/*
 * Archivo:     matchmonitor.go
 *
 * Propósito:   Simular un servidor de resultados deportivos: `matchRecorder`
 *              agrega un evento cada 200 ms y 5000 goroutines `clientHandler`
 *              copian 100 veces la lista completa de eventos. Todo se protege
 *              con un `sync.Mutex`, de modo que incluso las lecturas se
 *              ejecutan una a la vez.
 *
 * Compilación: go build matchmonitor.go
 * Ejecución:   ./matchmonitor
 *              (Alternativamente: go run matchmonitor.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Muchas líneas "N events copied in <tiempo>" y, de vez en cuando,
 *              "Appended match event". El programa dura 100 segundos.
 *
 * Notas:
 *    1. Con un mutex normal, los lectores se excluyen entre sí aunque no
 *       modifiquen los datos; el tiempo acumulado muestra el costo.
 *    2. Compare los tiempos con matchmonitorrw.go, que usa `sync.RWMutex` y
 *       permite lectores simultáneos.
 *    3. `matchRecorder` es un bucle infinito; termina solo cuando acaba `main`.
 */

package main

import (
    "fmt"
    "strconv"
    "sync"
    "time"
)

func matchRecorder(matchEvents *[]string, mutex *sync.Mutex) {
    for i := 0; ; i++ {
        mutex.Lock()
        *matchEvents = append(*matchEvents,
            "Match event " + strconv.Itoa(i))
        mutex.Unlock()
        time.Sleep(200 * time.Millisecond)
        fmt.Println("Appended match event")
    }
}

func clientHandler(mEvents *[]string, mutex *sync.Mutex, st time.Time) {
    for i := 0; i < 100; i ++ {
        mutex.Lock()
        allEvents := copyAllEvents(mEvents)
        mutex.Unlock()

        timeTaken := time.Since(st)
        fmt.Println(len(allEvents), "events copied in", timeTaken)
    }
}

func copyAllEvents(matchEvents *[]string) []string {
    allEvents := make([]string, 0, len(*matchEvents))
    for _, e := range *matchEvents {
        allEvents = append(allEvents, e)
    }
    return allEvents
}

func main() {
    mutex := sync.Mutex{}
    var matchEvents = make([]string, 0, 10000)
    for j := 0; j < 10000; j++ {
        matchEvents = append(matchEvents, "Match event")
    }
    go matchRecorder(&matchEvents, &mutex)
    start := time.Now()
    for j := 0; j < 5000; j++ {
        go clientHandler(&matchEvents, &mutex, start)
    }
    time.Sleep(100 * time.Second)
}
