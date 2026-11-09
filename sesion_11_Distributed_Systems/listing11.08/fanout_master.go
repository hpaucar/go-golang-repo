/*
 * Archivo:     fanout_master.go
 *
 * Propósito:   Patrón Fan-Out (scatter). El master recibe una lista de tareas
 *              en un solo mensaje JSON, la divide en bloques (uno por worker)
 *              y envía los bloques en paralelo, una goroutine por worker.
 *
 * Compilación: go build fanout_master.go
 * Ejecución:   Cinco terminales, en este orden:
 *                go run ../listing11.10/fanin_aggregator.go -esperados 10
 *                go run ../listing11.09/fanout_worker.go -puerto 9002 -nombre Worker-1
 *                go run ../listing11.09/fanout_worker.go -puerto 9003 -nombre Worker-2
 *                go run fanout_master.go
 *                go run ../listing11.03/centinela_json.go -modo fanout
 *
 * Entrada:     {"tipo":"TASK","valores":[1,...,10]} en el puerto 9001.
 * Salida:      Los bloques asignados a cada worker y la confirmación de envío.
 *
 * Notas:
 *    1. Cambio respecto a la guía: la guía decía "dividir en subconjuntos",
 *       pero enviaba cada número en un mensaje aparte (round-robin). Aquí se
 *       envía un bloque contiguo por worker: menos mensajes por la red.
 *    2. sync.WaitGroup espera a que terminen todos los envíos antes de
 *       cerrar la conexión del centinela.
 *    3. Fan-Out + Fan-In (listing11.10) forman el esquema scatter-gather, la
 *       base de MapReduce.
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Mensaje struct {
	Tipo    string `json:"tipo"`
	Origen  string `json:"origen,omitempty"`
	Valores []int  `json:"valores,omitempty"`
}

var workers []string

func main() {
	local := flag.String("local", "localhost:9001", "dirección de escucha")
	lista := flag.String("workers", "localhost:9002,localhost:9003", "workers separados por coma")
	flag.Parse()
	workers = strings.Split(*lista, ",")

	ln, err := net.Listen("tcp", *local)
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	fmt.Printf("[FAN-OUT MASTER] Escuchando en %s; workers = %v\n", *local, workers)

	for {
		conn, err := ln.Accept()
		if err == nil {
			go distribuirTareas(conn)
		}
	}
}

func distribuirTareas(conn net.Conn) {
	defer conn.Close()
	var msg Mensaje
	if err := json.NewDecoder(conn).Decode(&msg); err != nil {
		fmt.Println("[FAN-OUT MASTER] JSON inválido:", err)
		return
	}
	fmt.Printf("[FAN-OUT MASTER] Carga recibida: %v\n", msg.Valores)

	bloques := dividir(msg.Valores, len(workers))
	var wg sync.WaitGroup
	for i, bloque := range bloques {
		if len(bloque) == 0 {
			continue
		}
		wg.Add(1)
		go func(dir string, b []int) {
			defer wg.Done()
			enviarAWorker(dir, Mensaje{Tipo: "TASK", Origen: "MASTER", Valores: b})
		}(workers[i], bloque)
	}
	wg.Wait()
	fmt.Println("[FAN-OUT MASTER] Scatter completo.")
}

// dividir reparte v en n bloques contiguos de tamaño casi igual.
func dividir(v []int, n int) [][]int {
	bloques := make([][]int, n)
	for i := 0; i < n; i++ {
		ini := i * len(v) / n
		fin := (i + 1) * len(v) / n
		bloques[i] = v[ini:fin]
	}
	return bloques
}

func enviarAWorker(dir string, msg Mensaje) {
	conn, err := net.Dial("tcp", dir)
	if err != nil {
		fmt.Printf("[FAN-OUT MASTER] Worker %s no disponible: %v\n", dir, err)
		return
	}
	defer conn.Close()
	json.NewEncoder(conn).Encode(msg)
	fmt.Printf("[FAN-OUT MASTER -> WORKER] Bloque %v enviado a %s\n", msg.Valores, dir)
}
