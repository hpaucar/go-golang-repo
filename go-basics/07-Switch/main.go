/*
 * Archivo:     main.go
 *
 * Propósito:   Mostrar la sentencia `switch` en tres formas: sobre el valor de
 *              una variable, con varios valores en un mismo `case` (día de la
 *              semana con `time.Now().Weekday()`) y sin expresión, como
 *              alternativa a una cadena de `if`/`else`.
 *
 * Compilación: go build main.go
 * Ejecución:   ./main
 *              (Alternativamente: go run main.go)
 *
 * Entrada:     Ninguna (usa la fecha y la hora del sistema).
 * Salida:      "Tres" para el número 3; luego si es fin de semana o día de
 *              semana, la hora actual y un saludo ("Buenos días" antes de las
 *              12 y "Buenas tardes" después). La salida cambia según el día y
 *              la hora.
 *
 * Notas:
 *    1. En Go cada `case` termina solo; no hace falta `break`. Para pasar al
 *       caso siguiente se usa `fallthrough`.
 *    2. Un `case` puede listar varios valores separados por comas (`case
 *       time.Saturday, time.Sunday:`).
 *    3. Un `switch` sin expresión evalúa cada `case` como una condición
 *       booleana, de arriba hacia abajo.
 */

package main

import (
	"fmt"
	"time"
)

func main() {
	// 1. Usando switch para evaluar una variable
	fmt.Println("1. Usando switch para evaluar una variable:")

	number := 3
	fmt.Println("🔢 El número:", number, ", se escribe cómo:")

	switch number {
	case 1:
		fmt.Println("Uno")
	case 2:
		fmt.Println("Dos")
	case 3:
		fmt.Println("Tres")
	default:
		fmt.Println("❓ Número desconocido")
	}

	// 2. Usando switch con múltiples casos
	fmt.Println("\n2. Usando switch con múltiples casos:")

	switch time.Now().Weekday() {
	case time.Saturday, time.Sunday:
		fmt.Println("🏖️ Es fin de semana")
	default:
		fmt.Println("💼 Es un día de semana")
	}

	// 3. Usando switch sin una expresión
	fmt.Println("\n3. Usando switch sin una expresión:")

	t := time.Now()
	fmt.Println("La hora actual es:", t)

	switch {
	case t.Hour() < 12:
		fmt.Println("☀️ Buenos días")
	default:
		fmt.Println("🌇 Buenas tardes")
	}
}
