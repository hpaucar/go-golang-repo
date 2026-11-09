/*
 * Archivo:     recurso_compartido.go
 *
 * Propósito:   Simula un recurso que solo un nodo puede usar a la vez (por
 *              ejemplo, una impresora o un registro de base de datos). Cuenta
 *              cuántos nodos están "dentro" y avisa si hay más de uno: así se
 *              comprueba si un algoritmo de exclusión mutua distribuida
 *              funciona.
 *
 * Compilación: go build recurso_compartido.go
 * Ejecución:   go run recurso_compartido.go            (puerto 9100)
 *              Se usa con listing11.12 y listing11.13.
 *
 * Entrada:     Por cada uso: una conexión con {"tipo":"ENTRA","origen":id},
 *              luego {"tipo":"SALE","origen":id}. El recurso responde
 *              {"tipo":"ACK"} a cada mensaje.
 * Salida:      "ENTRA nodo X (dentro = 1)", "SALE nodo X", y
 *              "*** VIOLACIÓN DE EXCLUSIÓN MUTUA ***" si dentro > 1.
 *              Al pulsar Ctrl+C no imprime resumen; cuente las violaciones
 *              en pantalla.
 *
 * Notas:
 *    1. El ACK es importante: el nodo no sale de su sección crítica hasta que
 *       el recurso registró la salida. Sin él, la siguiente entrada podría
 *       llegar antes que la salida anterior y daría una falsa violación.
 *    2. Este programa NO participa en el algoritmo: solo observa. Es como un
 *       assert de Promela, pero en un sistema real.
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"sync"
	"time"
)

type Mensaje struct {
	Tipo   string `json:"tipo"`
	Origen int    `json:"origen"`
}

var (
	mu          sync.Mutex
	dentro      int
	violaciones int
)

func main() {
	local := flag.String("local", "localhost:9100", "dirección de escucha")
	flag.Parse()

	ln, err := net.Listen("tcp", *local)
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	fmt.Printf("[RECURSO] Escuchando en %s\n", *local)

	for {
		conn, err := ln.Accept()
		if err == nil {
			go atender(conn)
		}
	}
}

func atender(conn net.Conn) {
	defer conn.Close()
	dec := json.NewDecoder(conn)
	enc := json.NewEncoder(conn)
	for {
		var msg Mensaje
		if err := dec.Decode(&msg); err != nil {
			return
		}
		hora := time.Now().Format("15:04:05.000")
		mu.Lock()
		switch msg.Tipo {
		case "ENTRA":
			dentro++
			fmt.Printf("%s ENTRA nodo %d (dentro = %d)\n", hora, msg.Origen, dentro)
			if dentro > 1 {
				violaciones++
				fmt.Printf("%s *** VIOLACIÓN DE EXCLUSIÓN MUTUA #%d ***\n", hora, violaciones)
			}
		case "SALE":
			dentro--
			fmt.Printf("%s SALE  nodo %d\n", hora, msg.Origen)
		}
		mu.Unlock()
		enc.Encode(Mensaje{Tipo: "ACK"})
	}
}
