/*
 * Archivo:     fanout_worker.go
 *
 * Propósito:   Worker de la granja Fan-Out. Recibe un bloque de tareas en JSON,
 *              calcula el cuadrado de cada valor y envía cada resultado al
 *              agregador Fan-In (puerto 9010).
 *
 * Compilación: go build fanout_worker.go
 * Ejecución:   go run fanout_worker.go -puerto 9002 -nombre Worker-1
 *              go run fanout_worker.go -puerto 9003 -nombre Worker-2
 *              (orden completo de arranque: ver listing11.08)
 *
 * Entrada:     {"tipo":"TASK","origen":"MASTER","valores":[...]}.
 * Salida:      Una línea por valor procesado. Además envía
 *              {"tipo":"RESULT","origen":...,"resultado":x*x} al agregador.
 *
 * Notas:
 *    1. El mismo programa sirve para N workers: solo cambian los flags. Esto
 *       se llama SPMD (Single Program, Multiple Data).
 *    2. Si el agregador no está levantado, el resultado se pierde y se avisa
 *       en pantalla (la guía lo descartaba en silencio).
 *    3. El sleep simula un cálculo costoso; con dos workers el tiempo total es
 *       casi la mitad que con uno.
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

type Mensaje struct {
	Tipo      string `json:"tipo"`
	Origen    string `json:"origen,omitempty"`
	Valores   []int  `json:"valores,omitempty"`
	Resultado int    `json:"resultado"`
}

func main() {
	puerto := flag.String("puerto", "9002", "puerto local")
	nombre := flag.String("nombre", "Worker-1", "nombre del worker")
	agregador := flag.String("agregador", "localhost:9010", "dirección del agregador Fan-In")
	flag.Parse()

	ln, err := net.Listen("tcp", "localhost:"+*puerto)
	if err != nil {
		fmt.Printf("[%s ERROR] %v\n", *nombre, err)
		os.Exit(1)
	}
	defer ln.Close()
	fmt.Printf("[%s] Activo en el puerto %s\n", *nombre, *puerto)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go procesarTarea(conn, *nombre, *agregador)
	}
}

func procesarTarea(conn net.Conn, nombre, agregador string) {
	defer conn.Close()
	var msg Mensaje
	if err := json.NewDecoder(conn).Decode(&msg); err != nil {
		fmt.Printf("[%s ERROR] JSON inválido: %v\n", nombre, err)
		return
	}
	fmt.Printf("[%s] Bloque %v recibido de %s\n", nombre, msg.Valores, msg.Origen)

	for _, val := range msg.Valores {
		res := val * val
		time.Sleep(200 * time.Millisecond) // cálculo costoso simulado
		fmt.Printf("    [%s] %d^2 = %d\n", nombre, val, res)
		enviarResultado(agregador, Mensaje{Tipo: "RESULT", Origen: nombre, Resultado: res})
	}
}

func enviarResultado(dir string, msg Mensaje) {
	conn, err := net.Dial("tcp", dir)
	if err != nil {
		fmt.Printf("    [%s] Agregador no disponible; resultado %d perdido\n", msg.Origen, msg.Resultado)
		return
	}
	defer conn.Close()
	json.NewEncoder(conn).Encode(msg)
}
