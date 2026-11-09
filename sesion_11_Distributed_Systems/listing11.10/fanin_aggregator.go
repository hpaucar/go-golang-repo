/*
 * Archivo:     fanin_aggregator.go
 *
 * Propósito:   Patrón Fan-In (gather). Recibe resultados parciales de varios
 *              workers y los acumula. Con -esperados N, al llegar el N-ésimo
 *              resultado imprime el total final (reducción).
 *
 * Compilación: go build fanin_aggregator.go
 * Ejecución:   go run fanin_aggregator.go -esperados 10
 *              (orden completo de arranque: ver listing11.08)
 *
 * Entrada:     {"tipo":"RESULT","origen":"Worker-1","resultado":x}.
 * Salida:      Una línea por resultado parcial con la suma acumulada. Con los
 *              valores 1..10 elevados al cuadrado, el total es 385.
 *
 * Notas:
 *    1. Varias goroutines (una por conexión) modifican `sumaTotal` y
 *       `conteo`: sin el sync.Mutex habría una condición de carrera. Pruebe
 *       a quitarlo y ejecutar con `go run -race`.
 *    2. Los resultados llegan en cualquier orden; la suma no depende del
 *       orden porque es asociativa y conmutativa. Esa propiedad es la que
 *       permite "reducir" en paralelo.
 *    3. Cambio respecto a la guía: el flag -esperados permite saber cuándo
 *       terminó el trabajo (antes la suma nunca se daba por cerrada).
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"sync"
)

type Mensaje struct {
	Tipo      string `json:"tipo"`
	Origen    string `json:"origen,omitempty"`
	Resultado int    `json:"resultado"`
}

var (
	mu        sync.Mutex
	sumaTotal int
	conteo    int
	esperados int
)

func main() {
	local := flag.String("local", "localhost:9010", "dirección de escucha")
	flag.IntVar(&esperados, "esperados", 0, "resultados esperados (0 = sin límite)")
	flag.Parse()

	ln, err := net.Listen("tcp", *local)
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	fmt.Printf("[FAN-IN] Escuchando resultados en %s\n", *local)

	for {
		conn, err := ln.Accept()
		if err == nil {
			go recolectar(conn)
		}
	}
}

func recolectar(conn net.Conn) {
	defer conn.Close()
	var msg Mensaje
	if err := json.NewDecoder(conn).Decode(&msg); err != nil {
		fmt.Println("[FAN-IN] JSON inválido:", err)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	sumaTotal += msg.Resultado
	conteo++
	fmt.Printf("[FAN-IN] +%d de %s -> suma = %d (%d recibidos)\n",
		msg.Resultado, msg.Origen, sumaTotal, conteo)
	if esperados > 0 && conteo == esperados {
		fmt.Printf("[FAN-IN] ===== TOTAL FINAL = %d =====\n", sumaTotal)
	}
}
