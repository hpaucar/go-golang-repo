/*
 * Archivo:     hot_potato_node.go
 *
 * Propósito:   Papa caliente en un anillo de nodos (token ring). Cada nodo
 *              recibe un token JSON con un contador, lo decrementa y lo pasa
 *              al siguiente. Pierde el nodo que recibe el token con 0.
 *
 * Compilación: go build hot_potato_node.go
 * Ejecución:   Anillo 9001 -> 9002 -> 9003 -> 9001, tres terminales:
 *                go run hot_potato_node.go -local localhost:9001 -siguiente localhost:9002
 *                go run hot_potato_node.go -local localhost:9002 -siguiente localhost:9003
 *                go run hot_potato_node.go -local localhost:9003 -siguiente localhost:9001
 *              Cuarta terminal:
 *                go run ../listing11.03/centinela_json.go -modo potato -pasos 5
 *
 * Entrada:     {"tipo":"TOKEN","pasos":N}.
 * Salida:      Cada nodo muestra el token recibido y a quién lo reenvía. Con
 *              5 pasos y entrada en 9001, el token vuelve a 9001 y pierde
 *              9003 (recibe 0).
 *
 * Notas:
 *    1. Cambio respecto a la guía: los flags reciben direcciones completas
 *       (host:puerto) en vez de solo puertos. Así el mismo programa funciona
 *       dentro de Docker con nombres de servicio (listing11.15).
 *    2. Hay un único token circulando: solo su dueño "actúa". Esta es la idea
 *       base de la exclusión mutua por token (Ben-Ari, sección 10.6).
 *    3. El sleep simula la latencia de red y permite seguir el token a ojo.
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
	Tipo           string `json:"tipo"`
	Origen         string `json:"origen,omitempty"`
	PasosRestantes int    `json:"pasos"`
}

func main() {
	local := flag.String("local", "localhost:9001", "dirección de escucha")
	siguiente := flag.String("siguiente", "localhost:9002", "siguiente nodo del anillo")
	nombre := flag.String("nombre", "", "nombre del nodo (por defecto, -local)")
	flag.Parse()
	if *nombre == "" {
		*nombre = *local
	}

	ln, err := net.Listen("tcp", *local)
	if err != nil {
		fmt.Println("Error al escuchar:", err)
		os.Exit(1)
	}
	defer ln.Close()
	fmt.Printf("[NODO %s] Escuchando en %s. Siguiente en el anillo: %s\n",
		*nombre, *local, *siguiente)

	for {
		conn, err := ln.Accept()
		if err == nil {
			go procesarToken(conn, *nombre, *siguiente)
		}
	}
}

func procesarToken(conn net.Conn, yo, siguiente string) {
	defer conn.Close()
	var msg Mensaje
	if err := json.NewDecoder(conn).Decode(&msg); err != nil {
		fmt.Printf("[NODO %s] JSON inválido: %v\n", yo, err)
		return
	}

	time.Sleep(500 * time.Millisecond) // latencia simulada
	fmt.Printf("[NODO %s] Token recibido de %s con pasos = %d\n", yo, msg.Origen, msg.PasosRestantes)

	if msg.PasosRestantes <= 0 {
		fmt.Printf("[NODO %s] ¡PAPA CALIENTE! El contador llegó a 0: este nodo pierde.\n", yo)
		return
	}

	msg.PasosRestantes--
	msg.Origen = yo
	fmt.Printf("[NODO %s] Reenviando token con pasos = %d a %s\n", yo, msg.PasosRestantes, siguiente)
	enviar(siguiente, msg)
}

func enviar(dir string, msg Mensaje) {
	conn, err := net.Dial("tcp", dir)
	if err != nil {
		fmt.Printf("¡Token perdido! No se pudo conectar a %s: %v\n", dir, err)
		return
	}
	defer conn.Close()
	json.NewEncoder(conn).Encode(msg)
}
