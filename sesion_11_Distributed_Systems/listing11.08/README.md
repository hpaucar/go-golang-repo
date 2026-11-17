# Listado 11.08: Fan-out (master) + fan-in

`fanout_master.go` es el **master** del patrón **fan-out / fan-in**
(scatter-gather). Recibe una lista de tareas, la parte en bloques (uno por
worker) y envía los bloques **en paralelo**. Los workers (listado 11.09)
elevan cada número al cuadrado y envían los resultados al **agregador**
(listado 11.10), que los suma.

Esta guía levanta el patrón completo: agregador, dos workers, master y
centinela.

**Regla: primero los nodos, al final el centinela.** Aquí se levanta de atrás
hacia adelante: agregador, workers y por último el master.

## Topología: scatter-gather

```
                           ┌──► Worker-1 (9002) ──┐
centinela ──► master (9001)┤                      ├──► agregador (9010)
 (11.03)      fan-out      └──► Worker-2 (9003) ──┘     fan-in
              [1..10]          [1..5]  [6..10]          suma de cuadrados
```

| Proceso | Escucha | Envía a | Rol |
|---|---|---|---|
| centinela (11.03) | no | 9001 | solo cliente |
| master (11.08) | 9001 | 9002, 9003 | reparte (scatter) |
| worker (11.09) × 2 | 9002, 9003 | 9010 | calcula |
| agregador (11.10) | 9010 | no | junta y suma (gather) |

Fan-out + fan-in es la base de **MapReduce**: los workers hacen el *map*
(x → x²) y el agregador hace el *reduce* (suma).

## Cómo reparte el master

```go
// dividir reparte v en n bloques contiguos de tamaño casi igual.
func dividir(v []int, n int) [][]int {
	bloques := make([][]int, n)
	for i := 0; i < n; i++ {
		ini := i * len(v) / n
		fin := (i + 1) * len(v) / n
		bloques[i] = v[ini:fin]
	}
	return bloques
}
```

Con 10 valores y 2 workers: `[1..5]` y `[6..10]`. Cada bloque se envía en su
propia goroutine, y un `sync.WaitGroup` espera a que terminen todos los envíos
antes de imprimir `Scatter completo`.

## Flags

| Programa | Flag | Por defecto | Significado |
|---|---|---|---|
| master | `-local` | `localhost:9001` | dirección donde escucha |
| master | `-workers` | `localhost:9002,localhost:9003` | workers, separados por coma |
| worker (11.09) | `-puerto` | `9002` | puerto donde escucha |
| worker (11.09) | `-nombre` | `Worker-1` | nombre en pantalla |
| worker (11.09) | `-agregador` | `localhost:9010` | dirección del agregador |
| agregador (11.10) | `-local` | `localhost:9010` | dirección donde escucha |
| agregador (11.10) | `-esperados` | `0` | resultados que espera (0 = sin límite) |

## Paso 0: abrir 5 terminales

En **cada** terminal:

```bash
cd sesion_11_Distributed_Systems
```

## Paso 1: levantar el agregador (terminal 1)

```bash
go run listing11.10/fanin_aggregator.go -esperados 10
```

```
[FAN-IN] Escuchando resultados en localhost:9010
```

## Pasos 2 y 3: levantar los workers (terminales 2 y 3)

```bash
go run listing11.09/fanout_worker.go -puerto 9002 -nombre Worker-1   # terminal 2
go run listing11.09/fanout_worker.go -puerto 9003 -nombre Worker-2   # terminal 3
```

```
[Worker-1] Activo en el puerto 9002
```

## Paso 4: levantar el master (terminal 4)

```bash
go run listing11.08/fanout_master.go
```

```
[FAN-OUT MASTER] Escuchando en localhost:9001; workers = [localhost:9002 localhost:9003]
```

## Paso 5: ejecutar el centinela (terminal 5)

```bash
go run listing11.03/centinela_json.go -modo fanout
```

```
[CENTINELA SENT] TASK -> localhost:9001
```

## Resultado esperado

Terminal 4 (master):

```
[FAN-OUT MASTER] Carga recibida: [1 2 3 4 5 6 7 8 9 10]
[FAN-OUT MASTER -> WORKER] Bloque [6 7 8 9 10] enviado a localhost:9003
[FAN-OUT MASTER -> WORKER] Bloque [1 2 3 4 5] enviado a localhost:9002
[FAN-OUT MASTER] Scatter completo.
```

Terminal 2 (Worker-1):

```
[Worker-1] Bloque [1 2 3 4 5] recibido de MASTER
    [Worker-1] 1^2 = 1
    [Worker-1] 2^2 = 4
    ...
    [Worker-1] 5^2 = 25
```

Terminal 1 (agregador):

```
[FAN-IN] +36 de Worker-2 -> suma = 36 (1 recibidos)
[FAN-IN] +1 de Worker-1 -> suma = 37 (2 recibidos)
...
[FAN-IN] +25 de Worker-1 -> suma = 385 (10 recibidos)
[FAN-IN] ===== TOTAL FINAL = 385 =====
```

1² + 2² + … + 10² = **385**. Los resultados de los dos workers llegan
**intercalados** y en otro orden en cada ejecución: los workers trabajan en
paralelo. El total no cambia, porque la suma es asociativa y conmutativa.

## Paso 6: detener todo

Presione `Ctrl+C` en las terminales 1 a 4. Reinicie el agregador antes de
repetir la prueba: guarda la suma de la ejecución anterior.

## Errores comunes

- **`Worker localhost:9003 no disponible`** (master): un worker no estaba
  levantado. Su bloque se pierde y el total final nunca aparece.
- **`Agregador no disponible; resultado X perdido`** (worker): el agregador no
  estaba levantado.
- **Nunca aparece `TOTAL FINAL`**: el valor de `-esperados` no coincide con la
  cantidad de resultados (10 para la prueba del centinela).
- **`panic: listen tcp 127.0.0.1:9001: bind: address already in use`**: otro
  patrón sigue usando el puerto 9001.

## Para pensar y ejercicios

1. **Tres workers.** Levante un tercer worker en 9004 y lance el master con
   `-workers localhost:9002,localhost:9003,localhost:9004`. ¿Cómo quedan los
   bloques?
2. **Medir el paralelismo.** Cada cuadrado tarda 200 ms. Compare el tiempo
   total con 1, 2 y 5 workers.
3. **Round-robin.** La guía original enviaba cada número en un mensaje
   aparte. ¿Cuántos mensajes se envían con cada enfoque?
4. **Barrera distribuida.** Use este patrón como base del ejercicio 7 de la
   sesión.
