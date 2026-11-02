/*
 * Archivo:     main.go
 *
 * Propósito:   Implementar una "carrera de consultas": dos goroutines simulan
 *              consultas a dos servidores con un retraso aleatorio de 500 a
 *              1499 ms, y `select` procesa solo la primera respuesta que llega.
 *              Un `time.After` de 2 segundos actúa como tiempo límite
 *              (timeout).
 *
 * Compilación: go build main.go
 * Ejecución:   ./main
 *              (Alternativamente: go run main.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Los mensajes "🔍 Consultando Servidor A..." y "🔍 Consultando
 *              Servidor B..." y luego la primera respuesta, con el nombre del
 *              servidor y su tiempo. El ganador cambia entre ejecuciones; el
 *              timeout nunca se activa, porque el retraso máximo es menor que 2
 *              s.
 *
 * Notas:
 *    1. `select` espera varias operaciones de canal y ejecuta la primera que
 *       esté lista; si hay varias listas, elige una al azar.
 *    2. `time.After(d)` devuelve un canal que recibe un valor después de `d`;
 *       es la forma habitual de poner un timeout.
 *    3. La goroutine perdedora queda bloqueada para siempre al enviar en un
 *       canal sin búfer (goroutine leak). Aquí no importa porque el programa
 *       termina, pero en un servidor se usaría un canal con búfer o `context`
 *       para cancelarla.
 */

package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Este ejercicio implementa:
// !channels:
// Para manejar la comunicación entre goroutines
// !goroutines:
// Para ejecutar tareas simultáneamente
// !select:
// Para esperar múltiples operaciones
// !timeout:
// Para manejar tiempos de espera

// Simula una consulta a un servidor externo
func consultarServidor(nombre string, canal chan<- string) {
	fmt.Printf("🔍 Consultando %s...\n", nombre)
	// Simula un tiempo de respuesta aleatorio
	tiempo := time.Duration(rand.Intn(1000)+500) * time.Millisecond
	time.Sleep(tiempo)
	canal <- fmt.Sprintf("Respuesta de %s en %v", nombre, tiempo)
}

// Implementa un patrón de "carrera de consultas" donde se realizan peticiones
// simultáneas a dos servidores virtuales y se procesa solo la primera respuesta recibida.
//
// Este patrón es útil para mejorar la latencia en aplicaciones distribuidas,
// ya que permite obtener la respuesta del servidor más rápido sin esperar
// a que todos los servidores respondan.
func main() {
	// 1. Crea dos canales para recibir respuestas de cada servidor
	servidor1 := make(chan string)
	servidor2 := make(chan string)

	// 2. Inicia dos goroutines que simulan consultas a servidores distintos
	go consultarServidor("Servidor A", servidor1)
	go consultarServidor("Servidor B", servidor2)

	// 3. Utiliza una declaración select para manejar el primer resultado que llegue
	select {
	case respuesta := <-servidor1:
		fmt.Println("✅ Primer respuesta recibida:", respuesta)
	case respuesta := <-servidor2:
		fmt.Println("✅ Primer respuesta recibida:", respuesta)
	// 4. Incluye un timeout de 2 segundos para evitar bloqueos indefinidos
	case <-time.After(2 * time.Second):
		fmt.Println("⏰ Timeout: ningún servidor respondió a tiempo")
	}
}
