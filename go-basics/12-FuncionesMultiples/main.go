/*
 * Archivo:     main.go
 *
 * Propósito:   Mostrar funciones que devuelven varios valores (`calcular`
 *              devuelve suma, resta, multiplicación y división) y funciones con
 *              retorno nombrado (`exponenciar`). También muestra el
 *              identificador en blanco `_` para ignorar valores de retorno.
 *
 * Compilación: go build main.go
 * Ejecución:   ./main
 *              (Alternativamente: go run main.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Para 10 y 5: suma 15.0, resta 5.0, multiplicación 50.0 y
 *              división 2.0; luego algunos de esos valores por separado y "2
 *              elevado a 3 es: 8.0".
 *
 * Notas:
 *    1. Go exige recibir todos los valores de retorno; los que no interesan se
 *       descartan con `_`.
 *    2. Con un retorno nombrado, `resultado` ya está declarada y un `return`
 *       sin valores la devuelve. Conviene usarlo solo en funciones cortas.
 *    3. `calcular` devuelve 0 como división cuando `b` es 0, y `exponenciar`
 *       solo funciona con exponentes enteros no negativos (convierte `b` a
 *       `int`). En código real, la división por cero se informaría con un
 *       `error`.
 */

package main

import "fmt"

// función que devuelve múltiples valores: fsuma, fresta, fmultiplicación y fdivisión
func calcular(a, b float64) (float64, float64, float64, float64) {
	fsuma := a + b
	fresta := a - b
	fmultiplicacion := a * b
	var fdivision float64
	if b != 0 {
		fdivision = a / b
	}

	return fsuma, fresta, fmultiplicacion, fdivision
}

// función que usa un return nombrado para calcular la exponenciación
// return a^b
func exponenciar(a, b float64) (resultado float64) {
	resultado = 1
	for range int(b) {
		resultado *= a // resultado = resultado*a
	}
	return
}

func main() {
	// Valores de ejemplo
	x, y := 10.0, 5.0

	// Llamando a la función con múltiples valores de retorno
	suma, resta, mult, div := calcular(x, y)

	fmt.Printf("Números: %.1f y %.1f\n", x, y)
	fmt.Printf("Suma: %.1f\n", suma)   // Corregí este error que tenía en el original
	fmt.Printf("Resta: %.1f\n", resta) // Corregí el nombre de la variable
	fmt.Printf("Multiplicación: %.1f\n", mult)
	fmt.Printf("División: %.1f\n", div)

	// Usando identificador en blanco para ignorar algunos valores de retorno
	soloSuma, _, _, _ := calcular(x, y)
	fmt.Printf("Solo suma: %.1f\n", soloSuma)

	// Ignorando la suma y la división
	_, _, soloMultiplicacion, _ := calcular(x, y)
	fmt.Printf("Solo multiplicación: %.1f\n", soloMultiplicacion)

	// Ignorando la multiplicación y la suma
	_, soloResta, _, soloDivision := calcular(x, y)
	fmt.Printf("Solo resta: %.1f\n", soloResta)
	fmt.Printf("Solo división: %.1f\n", soloDivision)

	// Función con return nombrado
	resultado := exponenciar(2, 3)
	fmt.Printf("2 elevado a 3 es: %.1f\n", resultado)
}
