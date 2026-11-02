/*
 * Archivo:     channelcondvar.go
 *
 * Propósito:   Solución del ejercicio 7.4 del libro *Learn Concurrent
 *              Programming with Go*, que pide implementar un canal propio con
 *              variables de condición en lugar de semáforos. `Channel[M]` usa
 *              una `sync.Cond` y una lista; `Receive` aumenta la capacidad
 *              temporalmente, por lo que también funciona con capacidad 0.
 *
 * Compilación: go build channelcondvar.go
 * Ejecución:   ./channelcondvar
 *              (Alternativamente: go run channelcondvar.go)
 *
 * Entrada:     Ninguna.
 * Salida:      "Sending: 1" a "Sending: 6" y "Received: 1" a "Received: -1". El
 *              emisor se adelanta hasta llenar el búfer de 2 y luego avanza al
 *              ritmo del receptor (uno por segundo).
 *
 * Notas:
 *    1. `Send` espera mientras el búfer esté lleno; `Receive` espera mientras
 *       esté vacío. Ambos usan `Broadcast` para despertar al otro lado.
 *    2. Al incrementar `maxCapacity` antes de esperar, `Receive` deja pasar un
 *       envío aunque la capacidad sea 0, como en un canal sin búfer.
 *    3. El valor -1 es una señal de fin; `main` espera al receptor con un
 *       `sync.WaitGroup`.
 */

package main

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type Channel[M any] struct {
	cond        *sync.Cond
	maxCapacity int
	buffer      *list.List
}

func NewChannel[M any](capacity int) *Channel[M] {
	return &Channel[M]{
		cond:        sync.NewCond(&sync.Mutex{}),
		maxCapacity: capacity,
		buffer:      list.New(),
	}
}

func (c *Channel[M]) Send(message M) {
	c.cond.L.Lock()
	for c.buffer.Len() == c.maxCapacity {
		c.cond.Wait()
	}
	c.buffer.PushBack(message)
	c.cond.Broadcast()
	c.cond.L.Unlock()
}

func (c *Channel[M]) Receive() M {
	c.cond.L.Lock()
	c.maxCapacity++
	c.cond.Broadcast()
	for c.buffer.Len() == 0 {
		c.cond.Wait()
	}
	c.maxCapacity--
	v := c.buffer.Remove(c.buffer.Front()).(M)
	c.cond.L.Unlock()
	return v
}

func receiver(messages *Channel[int], wGroup *sync.WaitGroup) {
	msg := 0
	for msg != -1 {
		time.Sleep(1 * time.Second)
		msg = messages.Receive()
		fmt.Println("Received:", msg)
	}
	wGroup.Done()
}

func main() {
	channel := NewChannel[int](2)
	wGroup := sync.WaitGroup{}
	wGroup.Add(1)
	go receiver(channel, &wGroup)
	for i := 1; i <= 6; i++ {
		fmt.Println("Sending: ", i)
		channel.Send(i)
	}
	channel.Send(-1)
	wGroup.Wait()
}
