/*
 * Archivo:     pipeline_sort_node.go
 *
 * Propósito:   Ordenamiento distribuido en pipeline ("extrusión de mínimos").
 *              Cada nodo guarda el menor valor que ha visto y expulsa el mayor
 *              al siguiente nodo. Al recibir FIN, imprime su valor y propaga
 *              el FIN. Al terminar, el nodo k guarda el k-ésimo menor.
 *
 * Compilación: go build pipeline_sort_node.go
 * Ejecución:   5 valores => 5 nodos, una terminal por nodo:
 *                go run pipeline_sort_node.go -local localhost:9001 -siguiente localhost:9002
 *                go run pipeline_sort_node.go -local localhost:9002 -siguiente localhost:9003
 *                go run pipeline_sort_node.go -local localhost:9003 -siguiente localhost:9004
 *                go run pipeline_sort_node.go -local localhost:9004 -siguiente localhost:9005
 *                go run pipeline_sort_node.go -local localhost:9005        (último)
 *              Luego: go run ../listing11.03/centinela_json.go -modo pipeline
 *
 * Entrada:     Mensajes DATA {"valores":[x]} seguidos de un mensaje FIN.
 * Salida:      Las decisiones de cada nodo y, tras el FIN, una línea
 *              "[PIPELINE ...] RESULTADO: X". Con 45 12 89 3 27 los nodos
 *              9001..9005 terminan con 3, 12, 27, 45, 89.
 *
 * Notas:
 *    1. Orden de los mensajes: la versión de la guía abría una conexión TCP
 *       nueva por cada número. TCP conserva el orden DENTRO de una conexión,
 *       pero no ENTRE conexiones distintas, así que un FIN podía adelantarse
 *       a un dato. Aquí cada nodo mantiene UNA conexión abierta con el
 *       siguiente y lee su entrada como un flujo de mensajes.
 *    2. El sync.Mutex protege `minVal` si llegan varias conexiones a la vez
 *       (sincronización DENTRO del nodo). La coordinación ENTRE nodos es solo
 *       por mensajes. Es la misma separación que hace Ben-Ari en la 10.1.
 *    3. Ejercicio: ¿qué pasa si hay más valores que nodos? Haga que el último
 *       nodo guarde una lista ordenada en lugar de un solo valor.
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"sync"
)

type Mensaje struct {
	Tipo    string `json:"tipo"` // DATA o FIN
	Origen  string `json:"origen,omitempty"`
	Valores []int  `json:"valores,omitempty"`
}

var (
	mu        sync.Mutex
	minVal    int
	tieneVal  bool
	yo        string
	siguiente string
	salida    *json.Encoder // conexión persistente con el siguiente nodo
)

func main() {
	flag.StringVar(&yo, "local", "localhost:9001", "dirección de escucha")
	flag.StringVar(&siguiente, "siguiente", "", "siguiente nodo (vacío = último)")
	flag.Parse()

	ln, err := net.Listen("tcp", yo)
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	fmt.Printf("[PIPELINE %s] Listo. Siguiente: %q\n", yo, siguiente)

	for {
		conn, err := ln.Accept()
		if err == nil {
			go atender(conn)
		}
	}
}

// atender lee mensajes de la conexión hasta que el emisor la cierra.
func atender(conn net.Conn) {
	defer conn.Close()
	dec := json.NewDecoder(conn)
	for {
		var msg Mensaje
		if err := dec.Decode(&msg); err != nil {
			if err != io.EOF {
				fmt.Printf("[PIPELINE %s] error de lectura: %v\n", yo, err)
			}
			return
		}
		procesar(msg)
	}
}

func procesar(msg Mensaje) {
	mu.Lock()
	defer mu.Unlock()

	switch msg.Tipo {
	case "FIN":
		if tieneVal {
			fmt.Printf("[PIPELINE %s] RESULTADO: %d\n", yo, minVal)
		} else {
			fmt.Printf("[PIPELINE %s] RESULTADO: (sin valor)\n", yo)
		}
		reenviar(msg)
	case "DATA":
		if len(msg.Valores) == 0 {
			return
		}
		nuevo := msg.Valores[0]
		switch {
		case !tieneVal:
			minVal, tieneVal = nuevo, true
			fmt.Printf("[PIPELINE %s] Primer valor guardado: %d\n", yo, minVal)
		case nuevo < minVal:
			fmt.Printf("[PIPELINE %s] Nuevo MIN %d; expulsa %d\n", yo, nuevo, minVal)
			reenviar(Mensaje{Tipo: "DATA", Valores: []int{minVal}})
			minVal = nuevo
		default:
			fmt.Printf("[PIPELINE %s] Mantiene MIN %d; expulsa %d\n", yo, minVal, nuevo)
			reenviar(Mensaje{Tipo: "DATA", Valores: []int{nuevo}})
		}
	}
}

// reenviar envía por la conexión persistente; se llama con `mu` tomado, así
// los mensajes salen en el mismo orden en que se decidieron.
func reenviar(msg Mensaje) {
	if siguiente == "" {
		if msg.Tipo == "DATA" {
			fmt.Printf("[PIPELINE %s] Último nodo: se descarta %v (faltan nodos)\n", yo, msg.Valores)
		}
		return
	}
	if salida == nil {
		conn, err := net.Dial("tcp", siguiente)
		if err != nil {
			fmt.Printf("[PIPELINE %s] No se pudo conectar a %s: %v\n", yo, siguiente, err)
			return
		}
		salida = json.NewEncoder(conn)
	}
	msg.Origen = yo
	salida.Encode(msg)
}
