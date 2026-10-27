/*
 * Archivo:     producer_consumer_unbuffered.go
 *
 * Propósito:   Ilustrar un patrón clásico de productor-consumidor continuo 
 *              utilizando canales unidireccionales y sin búfer en Go. El 
 *              emisor (`sender`) genera y envía números enteros de manera 
 *              indefinida con un intervalo de 1 segundo, mientras el receptor 
 *              (`receiver`) los procesa de forma estrictamente sincrónica.
 *
 * Compilación: go build producer_consumer_unbuffered.go
 * Ejecución:   ./producer_consumer_unbuffered
 *              (Alternativamente: go run producer_consumer_unbuffered.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Mensajes impresos con marcas de tiempo que muestran la secuencia 
 *              alternada de envío y recepción de valores enteros durante el 
 *              tiempo que el hilo principal se mantiene activo.
 *
 * Notes:       
 *    1. El uso de parámetros con restricciones direccionales (`<-chan int` y 
 *       `chan<- int`) asegura la separación de responsabilidades en tiempo de 
 *       compilación: el emisor solo puede escribir y el receptor solo puede leer.
 *    2. Al tratarse de un canal sin búfer, la operación de envío en el emisor 
 *       se bloquea hasta que el receptor toma el elemento, logrando una 
 *       sincronización perfecta paso a paso.
 *    3. La llamada a `time.Sleep(5 * time.Second)` en el hilo principal finaliza 
 *       abruptamente el programa después de 5 segundos, terminando con todas 
 *       las goroutines en ejecución. En aplicaciones reales, se recomienda 
 *       gestionar la parada mediante canales de cancelación o `context`.
 */

package main

import (
    "fmt"
    "time"
)

func receiver(messages <-chan int) {
    for {
        msg := <-messages
        fmt.Println(time.Now().Format("15:04:05"), "Received:", msg)
    }
}

func sender(messages chan<- int) {
    for i := 1; ; i++ {
        fmt.Println(time.Now().Format("15:04:05"), "Sending:", i)
        messages <- i
        time.Sleep(1 * time.Second)
    }
}

func main() {
    msgChannel := make(chan int)
    go receiver(msgChannel)
    go sender(msgChannel)
    time.Sleep(5 * time.Second)
}