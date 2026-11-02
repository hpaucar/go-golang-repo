/*
 * Archivo:     passwordguesser.go
 *
 * Propósito:   Ilustrar la cancelación de goroutines con un canal que se
 *              cierra. 39 goroutines prueban por fuerza bruta rangos de 10
 *              millones de contraseñas (números en base 27); la que encuentra
 *              "go far" envía el resultado y cierra el canal `stop`.
 *
 * Compilación: go build passwordguesser.go
 * Ejecución:   ./passwordguesser
 *              (Alternativamente: go run passwordguesser.go)
 *
 * Entrada:     Ninguna (la contraseña está fija en `passwordToGuess`).
 * Salida:      "password found: go far" y, durante los 5 segundos finales,
 *              una línea "Stopped at N [desde,hasta)" por cada goroutine que
 *              se detiene (algunas pueden imprimir "Not found between ..."
 *              si terminan su rango antes). Usa mucha CPU.
 *
 * Notas:
 *    1. `toBase27` convierte un entero en una cadena con el alfabeto
 *       " abcdefghijklmnopqrstuvwxyz"; 27^6 = 387 420 489 cubre todas las
 *       cadenas de hasta 6 caracteres.
 *    2. `close(stop)` sirve de aviso a todas las goroutines a la vez: un
 *       canal cerrado siempre está listo para leer, así que `case <-stop`
 *       se activa en todas.
 *    3. `select` con `default` permite revisar `stop` en cada iteración sin
 *       bloquear la búsqueda.
 *    4. El `time.Sleep(5 * time.Second)` final solo da tiempo a ver los
 *       mensajes "Stopped at".
 */

package main

import (
    "fmt"
    "time"
)

const (
    passwordToGuess = "go far"
    alphabet        = " abcdefghijklmnopqrstuvwxyz"
)

func toBase27(n int) string {
    result := ""
    for n > 0 {
        result = string(alphabet[n%27]) + result
        n /= 27
    }
    return result
}

func guessPassword(from int, upto int, stop chan int, result chan string) {
    for guessN := from; guessN < upto; guessN += 1 {

        select {

        case <-stop:
            fmt.Printf("Stopped at %d [%d,%d)\n", guessN, from, upto)
            return

        default:
            if toBase27(guessN) == passwordToGuess {
                result <- toBase27(guessN)
                close(stop)
                return
            }
        }
    }
    fmt.Printf("Not found between [%d,%d)\n", from, upto)
}

func main() {
    finished := make(chan int)

    passwordFound := make(chan string)

    for i := 1; i <= 387_420_488; i += 10_000_000 {
        go guessPassword(i, i+10_000_000, finished, passwordFound)
    }

    fmt.Println("password found:", <-passwordFound)
    close(passwordFound)
    time.Sleep(5 * time.Second)
}
