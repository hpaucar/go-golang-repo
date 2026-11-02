/*
 * Archivo:     latesttemp.go
 *
 * Propósito:   Solución del ejercicio 8.1 del libro *Learn Concurrent
 *              Programming with Go*, que pide mostrar siempre la temperatura
 *              más reciente. Un generador produce una lectura cada 200 ms y el
 *              visor la muestra cada 2 s; un `select` en `main` actualiza el
 *              último valor o se lo entrega al visor.
 *
 * Compilación: go build latesttemp.go
 * Ejecución:   ./latesttemp
 *              (Alternativamente: go run latesttemp.go)
 *
 * Entrada:     Ninguna.
 * Salida:      "Current temp: <valor>" cada 2 segundos, en grados Fahrenheit a
 *              partir de 50. El programa no termina solo; deténgalo con Ctrl+C.
 *
 * Notas:
 *    1. El `select` evita que el generador rápido quede bloqueado por el visor
 *       lento: las lecturas intermedias se descartan.
 *    2. Es un patrón útil cuando solo importa el dato más reciente, como en
 *       sensores o paneles de estado.
 */

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func generateTemp() chan int {
	output := make(chan int)
	go func() {
		temp := 50 //fahrenheit
		for {
			output <- temp
			temp += rand.Intn(3) - 1
			time.Sleep(200 * time.Millisecond)
		}
	}()
	return output
}

func outputTemp(input chan int) {
	go func() {
		for {
			fmt.Println("Current temp:", <-input)
			time.Sleep(2 * time.Second)
		}
	}()
}

func main() {
	temps := generateTemp()
	display := make(chan int)
	outputTemp(display)
	t := <-temps
	for {
		select {
		case t = <-temps:
		case display <- t:
		}
	}
}
