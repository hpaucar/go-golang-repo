/*
 * Archivo:     centinela_json.go
 *
 * Propósito:   Centinela con mensajes JSON. Inyecta datos de prueba en el
 *              primer nodo de cada patrón distribuido. El flag -modo elige
 *              qué prueba enviar, porque cada patrón es un programa distinto
 *              que escucha en el mismo puerto 9001.
 *
 * Compilación: go build centinela_json.go
 * Ejecución:   El centinela es el ÚLTIMO paso. Cada orden va en su propia
 *              terminal, desde sesion_11_Distributed_Systems. Los nodos se
 *              dejan corriendo: tras su primera línea esperan mensajes.
 *              Un modo por vez (todos usan el puerto 9001); Ctrl+C en todos
 *              los nodos antes de cambiar de modo.
 *
 *              Modo conway (3 terminales):
 *                # Terminal 1: nodo de salida (9002). Déjelo corriendo.
 *                go run listing11.05/conway_output.go
 *                # Terminal 2: compresor (9001). Déjelo corriendo.
 *                go run listing11.04/conway_compressor.go
 *                # Terminal 3: espere a que la terminal 2 muestre
 *                # "Escuchando en localhost:9001", luego:
 *                go run listing11.03/centinela_json.go -modo conway
 *                (resultado en la terminal 1)
 *
 *              Modo potato (4 terminales):
 *                # Terminales 1 a 3: nodos del anillo. Déjelos corriendo.
 *                go run listing11.06/hot_potato_node.go -local localhost:9001 -siguiente localhost:9002
 *                go run listing11.06/hot_potato_node.go -local localhost:9002 -siguiente localhost:9003
 *                go run listing11.06/hot_potato_node.go -local localhost:9003 -siguiente localhost:9001
 *                # Terminal 4: espere a que las terminales 1 a 3 muestren
 *                # "Escuchando en ...", luego:
 *                go run listing11.03/centinela_json.go -modo potato -pasos 5
 *
 *              Modo pipeline (6 terminales, del último nodo al primero):
 *                # Terminales 1 a 5: nodos 9005..9001. Déjelos corriendo.
 *                go run listing11.07/pipeline_sort_node.go -local localhost:9005
 *                go run listing11.07/pipeline_sort_node.go -local localhost:9004 -siguiente localhost:9005
 *                go run listing11.07/pipeline_sort_node.go -local localhost:9003 -siguiente localhost:9004
 *                go run listing11.07/pipeline_sort_node.go -local localhost:9002 -siguiente localhost:9003
 *                go run listing11.07/pipeline_sort_node.go -local localhost:9001 -siguiente localhost:9002
 *                # Terminal 6: espere a que las terminales 1 a 5 muestren
 *                # "Listo. Siguiente: ...", luego:
 *                go run listing11.03/centinela_json.go -modo pipeline
 *
 *              Modo fanout (5 terminales, de atrás hacia adelante):
 *                # Terminal 1: agregador fan-in (9010). Déjelo corriendo.
 *                go run listing11.10/fanin_aggregator.go -esperados 10
 *                # Terminales 2 y 3: workers (9002, 9003). Déjelos corriendo.
 *                go run listing11.09/fanout_worker.go -puerto 9002 -nombre Worker-1
 *                go run listing11.09/fanout_worker.go -puerto 9003 -nombre Worker-2
 *                # Terminal 4: master fan-out (9001). Déjelo corriendo.
 *                go run listing11.08/fanout_master.go
 *                # Terminal 5: espere a que las terminales 1 a 4 muestren
 *                # su primera línea, luego:
 *                go run listing11.03/centinela_json.go -modo fanout
 *                (total final en la terminal 1)
 *
 *              Opcional: -destino host:puerto (por defecto localhost:9001).
 *              Guía paso a paso: listing11.03/README.md
 *
 * Entrada:     Ninguna.
 * Salida:      Una línea "[CENTINELA SENT] ..." por mensaje enviado.
 *
 * Notas:
 *    1. Revisión del centinela de la guía: allí se enviaban las pruebas de
 *       Conway, Hot Potato y Pipeline al mismo puerto 9001 en una sola
 *       ejecución, pero solo un programa puede escuchar en ese puerto. Por
 *       eso ahora se elige el modo.
 *    2. `Mensaje` es el "contrato" entre nodos: todos los programas de la
 *       sesión usan los mismos nombres JSON (tipo, origen, contenido, ...).
 *       Un nodo que no conoce un campo simplemente lo ignora.
 *    3. Modo pipeline: todos los números van por UNA sola conexión TCP, y al
 *       final se envía un mensaje FIN. Dentro de una conexión TCP el orden se
 *       conserva; entre conexiones distintas, no (ver listing11.07).
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
)

// Mensaje es el contrato de datos común entre los nodos distribuidos.
type Mensaje struct {
	Tipo           string `json:"tipo"`                // DATA, TOKEN, TASK, RESULT, FIN
	Origen         string `json:"origen,omitempty"`    // quién envía
	Contenido      string `json:"contenido,omitempty"` // carga en texto
	Valores        []int  `json:"valores,omitempty"`   // carga numérica
	PasosRestantes int    `json:"pasos,omitempty"`     // Hot Potato
}

func main() {
	modo := flag.String("modo", "conway", "conway | potato | pipeline | fanout")
	destino := flag.String("destino", "localhost:9001", "nodo de entrada")
	pasos := flag.Int("pasos", 5, "pases del token (modo potato)")
	flag.Parse()

	switch *modo {
	case "conway":
		enviarJSON(*destino, Mensaje{Tipo: "DATA", Origen: "CENTINELA",
			Contenido: "AAAAABBBCCDAAAAA"})
	case "potato":
		enviarJSON(*destino, Mensaje{Tipo: "TOKEN", Origen: "CENTINELA",
			PasosRestantes: *pasos})
	case "pipeline":
		enviarFlujo(*destino, []int{45, 12, 89, 3, 27})
	case "fanout":
		enviarJSON(*destino, Mensaje{Tipo: "TASK", Origen: "CENTINELA",
			Valores: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}})
	default:
		fmt.Println("Modo desconocido:", *modo)
		os.Exit(1)
	}
}

// enviarJSON abre una conexión, envía un mensaje y la cierra.
func enviarJSON(dir string, msg Mensaje) {
	conn, err := net.Dial("tcp", dir)
	if err != nil {
		fmt.Printf("[CENTINELA ERROR] No se pudo conectar a %s: %v\n", dir, err)
		return
	}
	defer conn.Close()
	if err := json.NewEncoder(conn).Encode(msg); err != nil {
		fmt.Println("[CENTINELA ERROR]", err)
		return
	}
	fmt.Printf("[CENTINELA SENT] %s -> %s\n", msg.Tipo, dir)
}

// enviarFlujo envía varios mensajes DATA y un FIN por la MISMA conexión.
func enviarFlujo(dir string, numeros []int) {
	conn, err := net.Dial("tcp", dir)
	if err != nil {
		fmt.Printf("[CENTINELA ERROR] No se pudo conectar a %s: %v\n", dir, err)
		return
	}
	defer conn.Close()
	enc := json.NewEncoder(conn)
	for _, n := range numeros {
		enc.Encode(Mensaje{Tipo: "DATA", Origen: "CENTINELA", Valores: []int{n}})
		fmt.Printf("[CENTINELA SENT] DATA %d -> %s\n", n, dir)
	}
	enc.Encode(Mensaje{Tipo: "FIN", Origen: "CENTINELA"})
	fmt.Printf("[CENTINELA SENT] FIN -> %s\n", dir)
}
