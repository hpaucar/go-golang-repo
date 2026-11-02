/*
 * Archivo:     cupcakefactory.go
 *
 * Propósito:   Ilustrar el patrón de tubería (pipeline) con canales. La
 *              función genérica `AddOnPipe` crea una goroutine que lee de un
 *              canal de entrada, aplica una etapa `f` y envía el resultado a
 *              un canal de salida. Encadenando cinco etapas, los 10 cupcakes
 *              avanzan por la línea de producción al mismo tiempo.
 *
 * Compilación: go build cupcakefactory.go
 * Ejecución:   ./cupcakefactory
 *              (Alternativamente: go run cupcakefactory.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Los mensajes de las etapas mezclados (varias etapas trabajan
 *              a la vez) y una línea "... received" por cupcake. Tarda unos
 *              58 segundos frente a unos 130 de la versión secuencial
 *              (10.14), porque el horno, de 5 segundos, marca el ritmo.
 *
 * Notas:
 *    1. Cada etapa es una goroutine conectada a la siguiente por un canal
 *       sin búfer; mientras una etapa procesa un cupcake, la anterior ya
 *       puede trabajar en el siguiente.
 *    2. El canal `quit` permite detener todas las etapas (al cerrarlo, cada
 *       `select` elige `case <-q` y la goroutine cierra su salida). En este
 *       `main` nunca se cierra: el programa termina al leer 10 resultados.
 *    3. Usa genéricos (`[X, Y any]`), disponibles desde Go 1.18.
 */

package main

import (
    "fmt"
    "github.com/cutajarj/ConcurrentProgrammingWithGo/chapter10/listing10.13"
)

func AddOnPipe[X, Y any](q <-chan int, f func(X) Y, in <-chan X) chan Y {
    output := make(chan Y)
    go func() {
        defer close(output)
        for {
            select {
            case <-q:
                return
            case input := <-in:
                output <- f(input)
            }
        }
    }()
    return output
}

func main() {
    input := make(chan int)
    quit := make(chan int)
    output := AddOnPipe(quit, listing10_13.Box,
        AddOnPipe(quit, listing10_13.AddToppings,
            AddOnPipe(quit, listing10_13.Bake,
                AddOnPipe(quit, listing10_13.Mixture,
                    AddOnPipe(quit, listing10_13.PrepareTray, input)))))
    go func() {
        for i := 0; i < 10; i++ {
            input <- i
        }
    }()
    for i := 0; i < 10; i++ {
        fmt.Println(<-output, "received")
    }
}
