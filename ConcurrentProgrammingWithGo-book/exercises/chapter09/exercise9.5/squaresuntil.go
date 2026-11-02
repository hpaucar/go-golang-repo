/*
 * Archivo:     squaresuntil.go
 *
 * Propósito:   Solución del ejercicio 9.5 del libro *Learn Concurrent
 *              Programming with Go*, que pide unir las etapas de los ejercicios
 *              9.1 a 9.4 en un pipeline: generar cuadrados, tomarlos mientras
 *              sean menores o iguales a 1 000 000, imprimirlos y descartarlos.
 *
 * Compilación: go build squaresuntil.go
 * Ejecución:   ./squaresuntil
 *              (Alternativamente: go run squaresuntil.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Los cuadrados 1, 4, 9, ..., 1000000, uno por línea. El programa
 *              termina cuando `TakeUntil` cierra `quit`.
 *
 * Notas:
 *    1. Actualmente no compila: importa `.../exercises/chapter9/...`, pero la
 *       carpeta del repositorio se llama `chapter09`.
 *    2. `main` espera con `<-quit`; leer de un canal cerrado lo desbloquea y el
 *       programa termina.
 *    3. Las funciones se anidan de adentro hacia afuera: la salida de cada
 *       etapa es la entrada de la siguiente.
 */

package main

import (
	"github.com/cutajarj/ConcurrentProgrammingWithGo/exercises/chapter9/exercise9.1"
	"github.com/cutajarj/ConcurrentProgrammingWithGo/exercises/chapter9/exercise9.2"
	"github.com/cutajarj/ConcurrentProgrammingWithGo/exercises/chapter9/exercise9.3"
	"github.com/cutajarj/ConcurrentProgrammingWithGo/exercises/chapter9/exercise9.4"
)

func main() {
	quit := make(chan int)
	exercise9_4.Drain(quit,
		exercise9_3.Print(quit,
			exercise9_2.TakeUntil(quit, func(s int) bool { return s <= 1000000 } ,
				exercise9_1.GenerateSquares(quit))))
	<-quit
}
