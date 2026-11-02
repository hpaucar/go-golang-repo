/*
 * Archivo:     main.go
 *
 * Propósito:   Mostrar el patrón de grupo de trabajadores (worker pool). Tres
 *              goroutines `worker` leen tareas de un canal con búfer, simulan 3
 *              s de trabajo por tarea y envían el doble del número a un canal
 *              de resultados que `main` recoge.
 *
 * Compilación: go build main.go
 * Ejecución:   ./main
 *              (Alternativamente: go run main.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Diez mensajes "📥 Tarea N enviada al canal", los mensajes de
 *              inicio y fin de cada worker y diez resultados (el doble de cada
 *              tarea, en orden variable). Con 3 workers y 10 tareas tarda unos
 *              12 segundos (4 rondas de 3 s).
 *
 * Notas:
 *    1. `for tarea := range tareasChannel` termina cuando el canal se cierra y
 *       queda vacío; por eso `main` llama a `close(tareasChannel)`.
 *    2. Los canales con búfer de tamaño `numTareas` permiten enviar todas las
 *       tareas sin esperar a los workers.
 *    3. `main` sabe cuándo terminar porque recibe exactamente `numTareas`
 *       resultados; no espera a que los workers salgan de su bucle. Para eso se
 *       usa `sync.WaitGroup`, como en el ejercicio 23.1.
 */

package main

import (
	"fmt"
	"time"
)

// ! workers
// Este ejemplo muestra cómo crear un worker que procesa tareas de un canal.
// Los workers son goroutines que realizan tareas específicas y pueden ser utilizados
// para procesar datos en paralelo, mejorando la eficiencia y el rendimiento de las
// aplicaciones en Go.

// worker procesa tareas recibidas desde un canal de entrada y envía los resultados a un canal de salida.
//
// Parámetros:
//   - id: identificador único del worker para su seguimiento en los logs
//   - tareasChannel: canal de solo lectura que recibe las tareas a procesar
//   - resultadossChannel: canal de solo escritura donde se envían los resultados
//
// El worker multiplica cada tarea por 2 y simula un procesamiento que toma 3 segundos.
// Imprime mensajes de inicio y finalización para cada tarea procesada.
func worker(id int, tareasChannel <-chan int, resultadossChannel chan<- int) {
	for tarea := range tareasChannel {
		fmt.Printf("👷 Worker %d procesando tarea %d...\n", id, tarea)
		time.Sleep(3 * time.Second) // Simula un trabajo que toma tiempo
		fmt.Printf("✅ Worker %d completó tarea %d\n", id, tarea)
		// Envía el resultado al canal de resultados
		resultadossChannel <- tarea * 2 // Envía resultado
	}
}

func main() {
	fmt.Println("🚀 Iniciando worker...")
	const numWorkers = 3
	const numTareas = 10

	tareasChannel := make(chan int, numTareas)     // Canal para tareas
	resultadosChannel := make(chan int, numTareas) // Canal para resultados

	// Crear y lanzar workers
	for i := 1; i <= numWorkers; i++ {
		go worker(i, tareasChannel, resultadosChannel)
	}

	// Enviar tareas al canal
	for i := 1; i <= numTareas; i++ {
		tareasChannel <- i
		fmt.Printf("📥 Tarea %d enviada al canal\n", i)
	}
	// Cerrar el canal de tareas para indicar que no hay más tareas
	close(tareasChannel)

	// Recoger resultados
	fmt.Println("📊 Recogiendo resultados...")
	for i := 1; i <= numTareas; i++ {
		result := <-resultadosChannel
		fmt.Printf("📈 Resultado recibido: %d\n", result)
	}
	fmt.Println("🏁 Todos los workers han terminado su trabajo")
	close(resultadosChannel) // Cierra el canal de resultados

	fmt.Println("🚀 Worker finalizado")
}
