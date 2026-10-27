/*
 * Archivo:     concurrent_factorization.go
 *
 * Propósito:   Ilustrar el uso de goroutines y canales para ejecutar 
 *              tareas independientes y pesadas de forma concurrente (en este caso, 
 *              el cálculo de los factores enteros de grandes números). 
 *              Demuestra cómo una goroutine en segundo plano calcula 
 *              factores mientras el hilo principal realiza otra tarea 
 *              simultánea antes de recoger el resultado.
 *
 * Compilación: go build concurrent_factorization.go
 * Ejecución:   ./concurrent_factorization
 *              (Alternativamente: go run concurrent_factorization.go)
 *
 * Entrada:     Ninguna. (Los números de prueba están codificados internamente).
 * Salida:      Imprime en la consola las listas de factores enteros 
 *              correspondientes a ambos números calculados.
 *
 * Notas:       
 *    1. La función `findFactors` implementa un algoritmo secuencial básico 
 *       que busca divisores desde 1 hasta el número dado mediante fuerza bruta.
 *    2. Se utiliza un canal sin búfer (`resultCh`) para comunicar el resultado 
 *       desde la goroutine anónima hacia el hilo principal, sincronizando 
 *       la lectura (`<-resultCh`) hasta que el cálculo finaliza.
 *    3. Mientras la goroutine calcula los factores de `3419110721` en segundo plano, 
 *       el hilo principal procede a calcular y mostrar de manera independiente 
 *       los factores de `4033836233`, ilustrando la ejecución concurrente 
 *       de tareas CPU-bound.
 */

package main

import (
    "fmt"
)

func findFactors(number int) []int {
    result := make([]int, 0)
    for i := 1; i <= number; i++ {
        if number%i == 0 {
            result = append(result, i)
        }
    }
    return result
}

func main() {
    resultCh := make(chan []int)
    
    go func() {
        resultCh <- findFactors(3419110721)
    }()
    
    fmt.Println(findFactors(4033836233))
    fmt.Println(<-resultCh)
}