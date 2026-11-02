/*
 * Archivo:     countdown.go
 *
 * Propósito:   Ilustrar la comunicación mediante memoria compartida: la
 *              goroutine `countdown` decrementa cada segundo la variable
 *              `count` a través de un puntero, mientras `main` la lee e imprime
 *              cada 500 ms hasta que llega a 0.
 *
 * Compilación: go build countdown.go
 * Ejecución:   ./countdown
 *              (Alternativamente: go run countdown.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Los valores 5, 5, 4, 4, 3, ... (cada valor suele aparecer dos
 *              veces) hasta que `count` llega a 0. Dura unos 5 segundos.
 *
 * Notas:
 *    1. El programa tiene una carrera de datos (data race) a propósito: una
 *       goroutine escribe `count` y otra la lee sin sincronización. Ejecútelo
 *       con `go run -race countdown.go` para verla.
 *    2. Aunque la salida parezca correcta, el modelo de memoria de Go no
 *       garantiza que `main` vea los cambios; el error se corrige en un
 *       ejercicio del capítulo 4 con un mutex.
 */

package main

import (
    "fmt"
    "time"
)
/*
Note: this program has a race condition for demonstration purposes
This is then fixed later as an exercise in chapter 4
 */
func countdown(seconds *int) {
    for *seconds > 0 {
        time.Sleep(1 * time.Second)
        *seconds -= 1
    }
}

func main() {
    count := 5
    go countdown(&count)
    for count > 0 {
        time.Sleep(500 * time.Millisecond)
        fmt.Println(count)
    }
}
