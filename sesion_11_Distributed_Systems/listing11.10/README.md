# Listado 11.10: Agregador fan-in

`fanin_aggregator.go` es el **agregador** del patrón fan-out / fan-in. Recibe
los resultados parciales de los workers (listado 11.09) y los **suma**. Con
`-esperados N`, cuando llega el resultado número N imprime el total final.

Este programa es la última etapa del patrón. **La guía completa, con las 5
terminales, está en [`../listing11.08/README.md`](../listing11.08/README.md).**
Aquí se explica solo el papel del agregador.

## Lugar en la topología

```
master ──► Worker-1 ──┐
       └─► Worker-2 ──┴──► agregador (9010)
                            suma = 385
```

Es el **primer** programa que se levanta: los workers le envían resultados en
cuanto empiezan a trabajar.

## Una condición de carrera evitada

Cada resultado llega en su propia conexión y se atiende en su propia
goroutine. Varias goroutines cambian `sumaTotal` y `conteo` a la vez, así que
se protegen con un `sync.Mutex`:

```go
mu.Lock()
defer mu.Unlock()
sumaTotal += msg.Resultado
conteo++
if esperados > 0 && conteo == esperados {
	fmt.Printf("[FAN-IN] ===== TOTAL FINAL = %d =====\n", sumaTotal)
}
```

Los resultados llegan en cualquier orden. La suma no depende del orden
(es asociativa y conmutativa), y esa propiedad permite hacer la **reducción en
paralelo**.

## Flags

| Flag | Por defecto | Significado |
|---|---|---|
| `-local` | `localhost:9010` | dirección donde escucha |
| `-esperados` | `0` | cuántos resultados espera. `0` = sin límite (nunca imprime el total) |

## Probar el agregador solo

```bash
# terminal 1 (desde sesion_11_Distributed_Systems)
go run listing11.10/fanin_aggregator.go -esperados 2
# terminal 2
echo '{"tipo":"RESULT","origen":"A","resultado":10}' | nc -q1 localhost 9010
echo '{"tipo":"RESULT","origen":"B","resultado":5}'  | nc -q1 localhost 9010
```

Terminal 1:

```
[FAN-IN] Escuchando resultados en localhost:9010
[FAN-IN] +10 de A -> suma = 10 (1 recibidos)
[FAN-IN] +5 de B -> suma = 15 (2 recibidos)
[FAN-IN] ===== TOTAL FINAL = 15 =====
```

(En macOS, use `nc localhost 9010` sin `-q1`.)

## Ver la condición de carrera

1. Comente las líneas `mu.Lock()` y `defer mu.Unlock()` en `recolectar`.
2. Levante el agregador con el detector de carreras:
   ```bash
   go run -race listing11.10/fanin_aggregator.go -esperados 10
   ```
3. Repita la prueba completa de [11.08](../listing11.08/README.md). Go muestra
   `WARNING: DATA RACE`.
4. Vuelva a poner el mutex.

## Errores comunes

- **Nunca aparece `TOTAL FINAL`**: `-esperados` no coincide con la cantidad
  de resultados, o se dejó en 0.
- **El total es mayor que 385**: el agregador no se reinició y sumó encima de
  una ejecución anterior.
- **`panic: listen tcp 127.0.0.1:9010: bind: address already in use`**: hay
  otro agregador activo.

## Para pensar y ejercicios

1. **Otras reducciones.** Cambie la suma por el máximo. ¿Sigue sin importar
   el orden? ¿Y con la resta?
2. **Reiniciar.** Haga que, después de imprimir el total, el agregador
   reinicie `sumaTotal` y `conteo` para aceptar un trabajo nuevo.
