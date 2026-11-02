/*
 * Archivo:     httpservernonblocking.go
 *
 * Propósito:   Solución del ejercicio 10.2 del libro *Learn Concurrent
 *              Programming with Go*, que pide que el servidor HTTP con pool de
 *              trabajadores (worker pool) no se bloquee cuando está saturado.
 *              Si la cola de 10 conexiones está llena, el `select` con
 *              `default` responde 503 de inmediato.
 *
 * Compilación: go build httpservernonblocking.go
 * Ejecución:   ./httpservernonblocking
 *              (Alternativamente: go run httpservernonblocking.go)
 *
 * Entrada:     Peticiones HTTP en localhost:8080 (por ejemplo, `curl
 *              localhost:8080/index.html`).
 * Salida:      Las respuestas HTTP de los 3 trabajadores del listado 10.9 y, si
 *              el servidor está ocupado, "Server is busy" en consola y un 503
 *              al cliente. Se ejecuta hasta Ctrl+C.
 *
 * Notas:
 *    1. Actualmente no compila: importa `.../chapter10/listing10.9`, pero la
 *       carpeta se llama `listing10.09`.
 *    2. Los trabajadores leen archivos de `../resources`, relativo al
 *       directorio actual; desde esta carpeta no existe y responden 404.
 *    3. Rechazar con 503 en lugar de encolar sin límite es una forma simple de
 *       control de carga (load shedding).
 */

package main

import (
    "fmt"
    "github.com/cutajarj/ConcurrentProgrammingWithGo/chapter10/listing10.9"
    "net"
)


func main() {
    incomingConnections := make(chan net.Conn, 10)
    listing10_9.StartHttpWorkers(3, incomingConnections)
    server, _ := net.Listen("tcp", "localhost:8080")
    defer server.Close()
    for {
        conn, _ := server.Accept()
        select {
        case incomingConnections <- conn:
        default:
            fmt.Println("Server is busy")
            conn.Write([]byte("HTTP/1.1 503 Service Unavailable\r\n\r\n" +
                "<html>Busy</html>\n"))
            conn.Close()
        }
    }
}
