/*
 * Archivo:     httpserver.go
 *
 * Propósito:   Poner en marcha el grupo de trabajadores de `httputils.go`
 *              (10.09) con 3 goroutines y un servidor TCP en el puerto 8080.
 *              El bucle principal acepta conexiones y las envía al canal
 *              `incomingConnections`, de donde las toman los trabajadores.
 *
 * Compilación: go build httpserver.go
 * Ejecución:   ./httpserver
 *              (Alternativamente: go run httpserver.go)
 *              Luego, en otra terminal: curl http://localhost:8080/index.html
 *
 * Entrada:     Peticiones HTTP en localhost:8080.
 * Salida:      No imprime nada en la consola; responde a cada cliente con
 *              el archivo pedido de chapter10/resources. El servidor no
 *              termina hasta que se detiene con Ctrl+C.
 *
 * Notas:
 *    1. El canal no tiene búfer: si los 3 trabajadores están ocupados, el
 *       envío `incomingConnections <- conn` bloquea el bucle y no se aceptan
 *       más conexiones hasta que un trabajador quede libre.
 *    2. Hay que ejecutarlo desde chapter10/listing10.11, porque los
 *       trabajadores buscan los archivos en la ruta relativa `../resources/`.
 *    3. El puerto 8080 debe estar libre; los errores se ignoran, así que si
 *       `net.Listen` falla, el programa termina con un pánico.
 *    4. La ruta de importación `.../chapter10/listing10.9` no coincide con
 *       la carpeta real `listing10.09`; hay que corregirla para compilar.
 */

package main

import (
    "github.com/cutajarj/ConcurrentProgrammingWithGo/chapter10/listing10.9"
    "net"
)

/*
If you're having problems running this, please make sure that the port 8080 is available
Note that error handling is omitted so that listing is brief
*/
func main() {
    incomingConnections := make(chan net.Conn)
    listing10_9.StartHttpWorkers(3, incomingConnections)
    server, _ := net.Listen("tcp", "localhost:8080")
    defer server.Close()
    for {
        conn, _ := server.Accept()
        incomingConnections <- conn
    }
}
