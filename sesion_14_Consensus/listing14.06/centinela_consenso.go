/*
 * Archivo:     centinela_consenso.go
 *
 * Propósito:   Centinela (disparador) de la votación de listing14.05. Envía
 *              el mensaje TRIGGER (code = 3) a todos los nodos de la malla a
 *              la vez, para que empiece una ronda de consenso.
 *
 * Compilación: go build centinela_consenso.go
 * Ejecución:   Después de levantar los 5 nodos (listing14.05):
 *                go run listing14.06/centinela_consenso.go \
 *                  -nodes localhost:8001,localhost:8002,localhost:8003,localhost:8004,localhost:8005
 *              En Docker Compose lo lanza el servicio "centinela" (listing14.07).
 *
 * Entrada:     Ninguna (flag -nodes).
 * Salida:      Una línea OK o ERROR por nodo.
 *
 * Notas:
 *    1. Parte del centinela.go de la guía de laboratorio. La constante
 *       ctrigger y el formato tmsg deben coincidir con consenso_mesh.go.
 *    2. Los envíos se hacen en goroutines y un sync.WaitGroup espera a que
 *       terminen todas antes de salir.
 *    3. Cada ejecución inicia una ronda nueva: los nodos reinician su mapa
 *       de votos al decidir.
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

const ctrigger = 3 // igual que en consenso_mesh.go

type tmsg struct {
	Code int    `json:"code"`
	Addr string `json:"addr"`
	Op   int    `json:"op"`
}

func main() {
	nodes := flag.String("nodes", "nodo1:8000,nodo2:8000,nodo3:8000,nodo4:8000,nodo5:8000",
		"nodos de la malla, separados por coma")
	flag.Parse()

	lista := strings.Split(*nodes, ",")
	fmt.Printf("[CENTINELA] Enviando TRIGGER (code = %d) a %d nodos...\n", ctrigger, len(lista))

	var wg sync.WaitGroup
	msg := tmsg{Code: ctrigger, Addr: "centinela"}
	for _, n := range lista {
		destino := strings.TrimSpace(n)
		if destino == "" {
			continue
		}
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
			if err != nil {
				fmt.Printf("[CENTINELA ERROR] No se pudo contactar %s: %v\n", addr, err)
				return
			}
			defer conn.Close()
			json.NewEncoder(conn).Encode(msg)
			fmt.Printf("[CENTINELA OK] Ronda iniciada en %s\n", addr)
		}(destino)
	}
	wg.Wait()
	fmt.Println("[CENTINELA] Listo.")
}
