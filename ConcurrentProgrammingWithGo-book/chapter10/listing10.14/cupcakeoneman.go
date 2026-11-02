/*
 * Archivo:     cupcakeoneman.go
 *
 * Propósito:   Fabricar 10 cupcakes de forma secuencial, como si una sola
 *              persona hiciera todas las etapas: cada cupcake pasa por las
 *              cinco funciones de `cupcakepipeline.go` (10.13) antes de
 *              empezar el siguiente. Sirve de base para comparar con 10.15.
 *
 * Compilación: go build cupcakeoneman.go
 * Ejecución:   ./cupcakeoneman
 *              (Alternativamente: go run cupcakeoneman.go)
 *
 * Entrada:     Ninguna.
 * Salida:      Los mensajes de cada etapa en orden estricto y una línea
 *              "Accepting ..." por cupcake terminado. Tarda unos 130
 *              segundos (10 cupcakes x 13 segundos).
 *
 * Notas:
 *    1. Las llamadas anidadas expresan la tubería, pero aquí no hay
 *       concurrencia: en cada momento solo se trabaja en una etapa.
 *    2. El tiempo total es la suma de todas las etapas de todos los
 *       cupcakes: mientras el horno trabaja, nadie prepara la bandeja
 *       siguiente.
 */

package main

import (
    "fmt"
    "github.com/cutajarj/ConcurrentProgrammingWithGo/chapter10/listing10.13"
)

func main() {
    for i := 0; i < 10; i++ {
        result := listing10_13.Box(
            listing10_13.AddToppings(
                listing10_13.Bake(
                    listing10_13.Mixture(
                        listing10_13.PrepareTray(i)))))
        fmt.Println("Accepting", result)
    }
}
