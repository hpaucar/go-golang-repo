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
