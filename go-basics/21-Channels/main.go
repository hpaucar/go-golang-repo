package main

import (
	"fmt"
	"time"
)

// !Channles
// Los canales (channels) son una forma de comunicación entre goroutines en Go.
// Permiten enviar y recibir datos entre goroutines de manera segura y sincronizada.
// Los canales son útiles para coordinar la ejecución de goroutines y compartir datos entre ellas.

// En este ejemplo, creamos dos canales y dos goroutines que envían mensajes a estos canales.
// Utilizamos un select para recibir mensajes de ambos canales, permitiendo que el programa

func main() {
	// Crear dos canales
	canal1 := make(chan string)
	canal2 := make(chan string)

	// Goroutine para canal1
	go func() {
		for i := 1; i <= 3; i++ {
			canal1 <- fmt.Sprintf("Mensaje %d desde canal 1", i)
			time.Sleep(300 * time.Millisecond) // Simula un retraso en el envío de mensajes
		}
		close(canal1) // Cerrar el canal1 al finalizar
	}()

	// Goroutine para canal2
	go func() {
		for i := 1; i <= 3; i++ {
			canal2 <- fmt.Sprintf("Mensaje %d desde canal 2", i)
			time.Sleep(500 * time.Millisecond)
		}
		close(canal2)
	}()

	// Recibir mensajes de ambos canales
	for canal1 != nil || canal2 != nil {
		// Utilizamos select para esperar mensajes de ambos canales
		// El select permite manejar múltiples canales de manera concurrente
		select {

		case msg, ok := <-canal1: // Recibe mensajes del canal1
			// ok, indica si el canal está abierto o cerrado
			// Si el canal está cerrado, ok será false
			if ok {
				fmt.Println("📨 " + msg)
			} else {
				// Si el canal está cerrado, lo establecemos a nil
				// para que no se vuelva a recibir mensajes de él
				canal1 = nil
				fmt.Println("❌ Canal 1 cerrado")
			}

		case msg, ok := <-canal2:
			if ok {
				fmt.Println("📬 " + msg)
			} else {
				canal2 = nil
				fmt.Println("❌ Canal 2 cerrado")
			}
		}
	}

	fmt.Println("✅ Todos los mensajes han sido recibidos")
}
