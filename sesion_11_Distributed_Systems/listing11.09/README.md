# Listado 11.09: Worker de fan-out

`fanout_worker.go` es un **worker** del patrón fan-out / fan-in. Recibe un
bloque de números del master (listado 11.08), calcula el cuadrado de cada uno
y envía **cada resultado** al agregador (listado 11.10).

Este programa no se usa solo. **La guía completa del patrón, con las 5
terminales, está en [`../listing11.08/README.md`](../listing11.08/README.md).**
Aquí se explica solo el papel del worker.

## Lugar en la topología

```
                           ┌──► Worker-1 (9002) ──┐
centinela ──► master (9001)┤                      ├──► agregador (9010)
                           └──► Worker-2 (9003) ──┘
```

## Un programa, muchos workers (SPMD)

No hay un archivo por worker. Se lanza **el mismo programa** dos veces con
flags distintos:

```bash
go run listing11.09/fanout_worker.go -puerto 9002 -nombre Worker-1
go run listing11.09/fanout_worker.go -puerto 9003 -nombre Worker-2
```

A esto se le llama **SPMD** (*Single Program, Multiple Data*): mismo código,
distintos datos.

## Qué hace con cada bloque

```go
for _, val := range msg.Valores {
	res := val * val
	time.Sleep(200 * time.Millisecond) // cálculo costoso simulado
	enviarResultado(agregador, Mensaje{Tipo: "RESULT", Origen: nombre, Resultado: res})
}
```

Mensaje que recibe del master y mensaje que envía al agregador:

```json
{"tipo":"TASK","origen":"MASTER","valores":[1,2,3,4,5]}
{"tipo":"RESULT","origen":"Worker-1","resultado":25}
```

## Flags

| Flag | Por defecto | Significado |
|---|---|---|
| `-puerto` | `9002` | puerto donde escucha |
| `-nombre` | `Worker-1` | nombre en pantalla y en el campo `origen` |
| `-agregador` | `localhost:9010` | dirección del agregador fan-in |

## Probar el worker solo

Con el agregador y un worker levantados (sin master ni centinela), envíe un
bloque a mano con `nc` (netcat):

```bash
# terminal 1
go run listing11.10/fanin_aggregator.go
# terminal 2
go run listing11.09/fanout_worker.go -puerto 9002 -nombre Worker-1
# terminal 3
echo '{"tipo":"TASK","origen":"YO","valores":[2,3]}' | nc -q1 localhost 9002
```

Terminal 2:

```
[Worker-1] Bloque [2 3] recibido de YO
    [Worker-1] 2^2 = 4
    [Worker-1] 3^2 = 9
```

Terminal 1: `+4 ... suma = 4` y `+9 ... suma = 13`. (En macOS, use
`nc localhost 9002` sin `-q1`.)

## Errores comunes

- **`[Worker-1 ERROR] listen tcp ... address already in use`**: otro programa
  usa el puerto (por ejemplo, el nodo de salida de Conway en 9002).
- **`Agregador no disponible; resultado X perdido`**: el agregador no estaba
  levantado; el resultado no se reintenta.
- **Dos workers con el mismo `-puerto`**: el segundo falla al escuchar.

## Para pensar y ejercicios

1. **Un mensaje por bloque.** El worker envía un mensaje por cada resultado.
   Cámbielo para enviar la suma parcial del bloque en un solo mensaje. ¿Qué
   tiene que cambiar en el agregador (`-esperados`)?
2. **Calcular en paralelo dentro del worker.** Procese los valores del bloque
   en goroutines. ¿El orden de los resultados importa?
