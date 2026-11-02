/*
 * Archivo:     httpservernonblocking.go
 *
 * Propósito:   Mostrar un envío no bloqueante con `select` y `default`. El
 *              servidor usa los mismos 3 trabajadores que 10.11, pero si
 *              ninguno está libre para recibir la conexión, no espera:
 *              responde de inmediato con "429 Too Many Requests".
 *
 * Compilación: go build httpservernonblocking.go
 * Ejecución:   ./httpservernonblocking
 *              (Alternativamente: go run httpservernonblocking.go)
 *              Para provocar el mensaje de servidor ocupado, en otra terminal:
 *              seq 1 2000 | xargs -Iname -P100 curl -s \
 *                  "http://localhost:8080/index.html" | grep Busy
 *
 * Entrada:     Peticiones HTTP en localhost:8080.
 * Salida:      Imprime "Server is busy" cada vez que rechaza una conexión;
 *              los clientes reciben el archivo pedido o una respuesta 429.
 *              El servidor no termina hasta que se detiene con Ctrl+C.
 *
 * Notas:
 *    1. En un `select`, la rama `default` se ejecuta cuando ninguna otra
 *       operación puede realizarse en ese momento; así el envío al canal
 *       deja de ser bloqueante.
 *    2. Rechazar trabajo cuando el sistema está saturado (load shedding)
 *       mantiene el servidor rápido para las peticiones que sí acepta.
 *    3. Hay que ejecutarlo desde chapter10/listing10.12 (ruta relativa
 *       `../resources/`) y con el puerto 8080 libre.
 *    4. La ruta de importación `.../chapter10/listing10.9` no coincide con
 *       la carpeta real `listing10.09`; hay que corregirla para compilar.
 */

package main

import (
    "fmt"
    "github.com/cutajarj/ConcurrentProgrammingWithGo/chapter10/listing10.9"
    "net"
)

/*
To trigger the "Too Many Requests" message try running this:
seq 1 2000 | xargs -Iname  -P100  curl -s "http://localhost:8080/index.html" | grep Busy
*/
func main() {
    incomingConnections := make(chan net.Conn)
    listing10_9.StartHttpWorkers(3, incomingConnections)
    server, _ := net.Listen("tcp", "localhost:8080")
    defer server.Close()
    for {
        conn, _ := server.Accept()
        select {
        case incomingConnections <- conn:
        default:
            fmt.Println("Server is busy")
            conn.Write([]byte("HTTP/1.1 429 Too Many Requests\r\n\r\n" +
                "<html>Busy</html>\n"))
            conn.Close()
        }
    }
}
