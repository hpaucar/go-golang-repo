/*
 * Archivo:     cpucheck.go
 *
 * Propósito:   Mostrar cuántos núcleos lógicos de CPU detecta Go con
 *              `runtime.NumCPU()` y cuántos hilos del sistema operativo pueden
 *              ejecutar código Go en paralelo, consultando
 *              `runtime.GOMAXPROCS(0)`.
 *
 * Compilación: go build cpucheck.go
 * Ejecución:   ./cpucheck
 *              (Alternativamente: go run cpucheck.go)
 *
 * Entrada:     Ninguna. (Se puede fijar la variable de entorno `GOMAXPROCS`,
 *              por ejemplo: GOMAXPROCS=2 ./cpucheck).
 * Salida:      Dos líneas: "Number of CPUs: N" y "GOMAXPROCS: M".
 *
 * Notas:
 *    1. Llamar a `runtime.GOMAXPROCS(0)` solo consulta el valor actual; un
 *       argumento mayor que 0 lo cambiaría.
 *    2. Por defecto `GOMAXPROCS` es igual al número de CPU lógicas (en
 *       contenedores, versiones recientes de Go también consideran el límite de
 *       CPU del cgroup).
 */

package main

import (
    "fmt"
    "runtime"
)

func main() {
    fmt.Println("Number of CPUs:", runtime.NumCPU())

    fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0))
}
