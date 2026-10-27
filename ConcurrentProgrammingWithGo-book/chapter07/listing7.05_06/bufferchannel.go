/*
 * Archivo:     buffered_channel_sentinel.go
 *
 * Propósito:   Ilustrar el uso de canales con búfer (buffered channels) en Go, 
 *              combinados con una variable de control `sync.WaitGroup` y un 
 *              valor centinela (`-1`) para coordinar la terminación segura de una 
 *              goroutine receptora. El programa muestra cómo los canales con 
 *              búfer permiten enviar múltiples elementos sin bloqueo inmediato 
 *              hasta que su capacidad máxima se agota.
 *
 * Compilación: go build buffered_channel_sentinel.go
 * Ejecución:   ./buffered_channel_sentinel
 *              (Alternativamente: go run buffered_channel_sentinel.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Mensajes impresos con marcas de tiempo (timestamps) que muestran 
 *              el estado dinámico del búfer durante el envío, la recepción 
 *              secuencial de los elementos con un retraso de 1 segundo por cada 
 *              lectura, y la finalización coordinada del programa.
 *
 * Notas:       
 *    1. `make(chan int, 3)` crea un canal con un búfer de capacidad 3. Esto 
 *       permite que la goroutine emisora (`main`) coloque hasta 3 elementos 
 *       sin bloquearse, antes de que se requiera que el receptor libere espacio.
 *    2. Se utiliza el valor `-1` como un "valor centinela" (sentinel value) 
 *       para indicarle a la función `receiver` que el flujo de datos ha concluido 
 *       y que debe salir de su bucle de procesamiento.
 *    3. Se emplea correctamente un `sync.WaitGroup` para asegurar que el hilo 
 *       principal espere (`wGroup.Wait()`) a que el receptor llame a `Done()`, 
 *       evitando que el programa principal finalice prematuramente.
 */

package main

import (
    "fmt"
    "sync"
    "time"
)

func receiver(messages chan int, wGroup *sync.WaitGroup) {
    msg := 0
    for msg != -1 {
        time.Sleep(1 * time.Second)
        msg = <-messages
        fmt.Println(time.Now().Format("15:04:05"), "Received:", msg)
    }
    wGroup.Done()
}

func main() {
    msgChannel := make(chan int, 3)
    wGroup := sync.WaitGroup{}
    wGroup.Add(1)
    
    go receiver(msgChannel, &wGroup)
    
    for i := 1; i <= 6; i++ {
        size := len(msgChannel)
        fmt.Printf("%s Sending: %d. Buffer Size: %d\n",
            time.Now().Format("15:04:05"), i, size)
        msgChannel <- i
    }
    
    msgChannel <- -1
    wGroup.Wait()
}