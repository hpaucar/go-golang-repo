package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// Estructura que representa un email a procesar
type Email struct {
	ID           int
	Destinatario string
	Asunto       string
	Contenido    string
}

// Estructura que representa el resultado del procesamiento
type ResultadoProceso struct {
	EmailID int
	Exito   bool
	Mensaje string
	Tiempo  time.Duration
}

// Worker que simula el envío de emails
func workerEnvioEmail(workerID int, emails <-chan Email, resultados chan<- ResultadoProceso, wg *sync.WaitGroup) {
	defer wg.Done() // Marca este worker como terminado cuando la función termine

	for email := range emails {
		fmt.Printf("📧 Worker %d: Procesando email %d para %s...\n", workerID, email.ID, email.Destinatario)

		inicio := time.Now()

		// Simula el tiempo de procesamiento (envío real del email)
		// En la vida real aquí iría la lógica de conexión SMTP, autenticación, etc.
		tiempoProceso := time.Duration(rand.Intn(2000)+500) * time.Millisecond
		time.Sleep(tiempoProceso)

		// Simula que algunos emails fallan (10% de probabilidad)
		exito := rand.Float32() > 0.1

		resultado := ResultadoProceso{
			EmailID: email.ID,
			Exito:   exito,
			Tiempo:  time.Since(inicio),
		}

		if exito {
			resultado.Mensaje = fmt.Sprintf("Email enviado exitosamente a %s", email.Destinatario)
			fmt.Printf("✅ Worker %d: Email %d enviado correctamente\n", workerID, email.ID)
		} else {
			resultado.Mensaje = fmt.Sprintf("Error al enviar email a %s", email.Destinatario)
			fmt.Printf("❌ Worker %d: Error enviando email %d\n", workerID, email.ID)
		}

		// Envía el resultado al canal de resultados
		resultados <- resultado
	}

	fmt.Printf("🏁 Worker %d: Terminando trabajo\n", workerID)
}

func main() {
	fmt.Println("🚀 Sistema de Envío Masivo de Emails - Iniciando...")

	// Configuración del sistema
	const numWorkers = 3 // Número de workers (hilos de procesamiento)
	const numEmails = 10 // Número total de emails a enviar

	// Crear canales
	emailsChannel := make(chan Email, numEmails)                // Canal buffer para emails
	resultadosChannel := make(chan ResultadoProceso, numEmails) // Canal buffer para resultados

	// WaitGroup para esperar que todos los workers terminen
	var wg sync.WaitGroup

	// ========================================
	// PASO 1: Crear e iniciar los workers
	// ========================================
	fmt.Printf("👷 Iniciando %d workers...\n", numWorkers)

	// Lanzar los workers
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1) // Incrementa el contador del WaitGroup
		go workerEnvioEmail(i, emailsChannel, resultadosChannel, &wg)
	}

	// ========================================
	// PASO 2: Generar emails y enviarlos al canal
	// ========================================
	fmt.Printf("📨 Generando %d emails...\n", numEmails)

	go func() {
		for i := 1; i <= numEmails; i++ {
			email := Email{
				ID:           i,
				Destinatario: fmt.Sprintf("usuario%d@ejemplo.com", i),
				Asunto:       fmt.Sprintf("Oferta especial #%d", i),
				Contenido:    "¡No te pierdas nuestra increíble oferta!",
			}
			emailsChannel <- email
		}
		close(emailsChannel) // Cierra el canal cuando no hay más emails
		fmt.Println("📪 Todos los emails han sido enviados al canal")
	}()

	// ========================================
	// PASO 3: Recopilar resultados en una goroutine separada
	// ========================================
	go func() {
		wg.Wait()                // Espera que todos los workers terminen
		close(resultadosChannel) // Cierra el canal de resultados
	}()

	// ========================================
	// PASO 4: Procesar y mostrar resultados
	// ========================================
	fmt.Println("\n📊 Procesando resultados...")

	var exitosos, fallidos int
	var tiempoTotal time.Duration

	for resultado := range resultadosChannel {
		fmt.Printf("📋 Email %d: %s (Tiempo: %v)\n",
			resultado.EmailID, resultado.Mensaje, resultado.Tiempo)

		if resultado.Exito {
			exitosos++
		} else {
			fallidos++
		}
		tiempoTotal += resultado.Tiempo
	}

	// ========================================
	// PASO 5: Mostrar estadísticas finales
	// ========================================
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("📈 ESTADÍSTICAS FINALES")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("✅ Emails enviados exitosamente: %d\n", exitosos)
	fmt.Printf("❌ Emails fallidos: %d\n", fallidos)
	fmt.Printf("📧 Total de emails procesados: %d\n", exitosos+fallidos)
	fmt.Printf("⏱️  Tiempo promedio por email: %v\n", tiempoTotal/time.Duration(numEmails))
	fmt.Printf("📊 Tasa de éxito: %.1f%%\n", float64(exitosos)/float64(numEmails)*100)

	fmt.Println("\n🎉 Sistema de envío completado!")
}
