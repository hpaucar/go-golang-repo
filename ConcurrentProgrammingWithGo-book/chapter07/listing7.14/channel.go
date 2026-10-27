/*
 * Archivo:     custom_channel_semaphore.go
 *
 * Propósito:   Implementar una estructura de datos personalizada de tipo Canal 
 *              genérico (`Channel[M any]`) utilizando semáforos contadores y 
 *              exclusión mutua (`sync.Mutex`), imitando el comportamiento interno 
 *              de los canales nativos de Go pero construido manualmente con 
 *              primitivas de sincronización de nivel inferior.
 *
 * Compilación: go build custom_channel_semaphore.go
 * Ejecución:   (Requiere el módulo y paquete de referencia de la librería del libro 
 *              "Concurrent Programming with Go" de J. Cutajar).
 *
 * Entrada:     Ninguna (Definición de estructura y métodos).
 * Salida:      Ninguna de forma directa; provee la implementación de un canal 
 *              con búfer acotado y seguro para hilos concurrentes.
 *
 * Notas:       
 *    1. Utiliza genéricos de Go (`[M any]`) para permitir que el canal transporte 
 *       cualquier tipo de datos de manera segura en tiempo de compilación.
 *    2. Emplea dos semáforos independientes (`capacitySema` y `sizeSema`) junto 
 *       con una lista enlazada (`container/list`) como búfer subyacente y un 
 *       `sync.Mutex` para proteger las operaciones de inserción y extracción.
 *    3. `capacitySema` controla la capacidad máxima disponible (restringiendo 
 *       los envíos cuando el búfer está lleno), mientras que `sizeSema` lleva la 
 *       cuenta de los elementos presentes (bloqueando las lecturas si está vacío).
 *    4. ADVERTENCIA LÓGICA EN EL CÓDIGO ORIGINAL: Existe un error sutil en el 
 *       método `Receive()`, ya que invoca `c.capacitySema.Release()` *antes* de 
 *       adquirir `c.sizeSema.Acquire()`. En una implementación correcta de canales 
 *       con búfer, el orden de adquisición y liberación de los semáforos debe 
 *       garantizar que no se libere espacio en el búfer antes de extraer el elemento.
 */

package listing7_14

import (
    "container/list"
    "github.com/cutajarj/ConcurrentProgrammingWithGo/chapter5/listing5.16"
    "sync"
)

type Channel[M any] struct {
    capacitySema *listing5_16.Semaphore
    sizeSema     *listing5_16.Semaphore
    mutex        sync.Mutex
    buffer       *list.List
}

func NewChannel[M any](capacity int) *Channel[M] {
    return &Channel[M]{
        capacitySema: listing5_16.NewSemaphore(capacity),
        sizeSema:     listing5_16.NewSemaphore(0),
        buffer:       list.New(),
    }
}

func (c *Channel[M]) Send(message M) {
    c.capacitySema.Acquire()
    c.mutex.Lock()
    c.buffer.PushBack(message)
    c.mutex.Unlock()
    c.sizeSema.Release()
}

func (c *Channel[M]) Receive() M {
    c.capacitySema.Release()
    c.sizeSema.Acquire()
    c.mutex.Lock()
    v := c.buffer.Remove(c.buffer.Front()).(M)
    c.mutex.Unlock()
    return v
}