/*
 * Archivo:     main.go
 *
 * Propósito:   Mostrar cómo simular enumeraciones (enums) con un tipo propio
 *              (`ServerState`) y constantes generadas con `iota`. Implementa el
 *              método `String()` para imprimir cada estado y una máquina de
 *              estados simple con `updateServerState`.
 *
 * Compilación: go build main.go
 * Ejecución:   ./main
 *              (Alternativamente: go run main.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Los estados "🚀 Starting", "✅ Running", "🛑 Stopping" y "⏹️
 *              Stopped". Después el programa termina con un pánico (panic)
 *              intencional: "No hay un siguiente estado después de ⏹️ Stopped",
 *              con la traza de la pila y código de salida 2.
 *
 * Notas:
 *    1. `iota` vale 0 en la primera constante del bloque `const` y aumenta en 1
 *       en cada línea.
 *    2. Al implementar `String() string`, el tipo cumple la interfaz
 *       `fmt.Stringer` y `fmt.Println` usa ese método automáticamente.
 *    3. El pánico final es parte del ejemplo. En código real, un estado
 *       inválido suele devolverse como `error` en lugar de usar `panic`.
 */

package main

import "fmt"

type ServerState int // Definición de un tipo de enumeración para los estados del servidor

// Definición de constantes para los diferentes estados del servidor
// Cada constante representa un estado específico del servidor y se asigna un valor entero único
const (
	Starting ServerState = iota // 0
	Running                     // 1
	Stopping                    // 2
	Stopped                     // 3
)

// map que asocia cada estado del servidor con su representación en cadena
var serverStateName = map[ServerState]string{
	Starting: "🚀 Starting",
	Running:  "✅ Running",
	Stopping: "🛑 Stopping",
	Stopped:  "⏹️ Stopped",
}

// String convierte el estado del servidor a su representación en cadena
func (serverState ServerState) String() string {
	return serverStateName[serverState]
}

// updateServerState recibe un estado del servidor y devuelve el siguiente estado en la secuencia
func updateServerState(state ServerState) ServerState {
	switch state {
	case Starting:
		return Running
	case Running:
		return Stopping
	case Stopping:
		return Stopped
	default:
		panic(fmt.Errorf("No hay un siguiente estado después de %s", state))
	}
}

func main() {
	redServer := Starting
	fmt.Println(redServer)

	redServer = updateServerState(redServer)
	fmt.Println(redServer)

	redServer = updateServerState(redServer)
	fmt.Println(redServer)

	redServer = updateServerState(redServer)
	fmt.Println(redServer)

	// Esto provocará un pánico porque no hay un estado siguiente después de Stopped
	redServer = updateServerState(redServer)
}
