# Listado 11.05: Nodo de salida de Conway

`conway_output.go` es el **nodo 2** del problema de Conway distribuido
(Ben-Ari, Alg. 8.2). Recibe la cadena comprimida del compresor (listado
11.04) y la imprime en **líneas de ancho fijo K**, como pide el enunciado
original de Conway.

**Regla: primero los nodos, al final el centinela.** Este nodo es el primero
que se levanta, porque el compresor le envía su resultado.

## Topología: un pipeline de dos etapas

```
centinela ──► compresor (9001) ──► salida (9002)
 (11.03)      (11.04, RLE)          (11.05, líneas de K caracteres)
```

| Proceso | Escucha | Envía | Rol |
|---|---|---|---|
| centinela (11.03) | no | a 9001 | solo cliente |
| compresor (11.04) | 9001 | a 9002 | servidor y cliente |
| **salida (11.05)** | 9002 | no | solo servidor |

## Partir en líneas de K caracteres

```go
func partir(s string, k int) []string {
	runas := []rune(s)                 // por runas, no por bytes
	var lineas []string
	for i := 0; i < len(runas); i += k {
		fin := i + k
		if fin > len(runas) {
			fin = len(runas)           // la última línea puede ser más corta
		}
		lineas = append(lineas, string(runas[i:fin]))
	}
	return lineas
}
```

Se trabaja con `[]rune` por la misma razón que en el compresor: así una `Ñ`
(2 bytes en UTF-8) cuenta como un solo carácter y no se parte en dos líneas.

## El mensaje que recibe

```json
{"tipo":"RLE_OUTPUT","origen":"COMPRESSOR","contenido":"5A3B2C1D5A"}
```

## Flags

| Flag | Por defecto | Significado |
|---|---|---|
| `-local` | `localhost:9002` | dirección donde escucha |
| `-k` | `4` | caracteres por línea |

## Paso 0: abrir 3 terminales

En **cada** terminal:

```bash
cd sesion_11_Distributed_Systems
```

## Paso 1: levantar el nodo de salida (terminal 1)

```bash
go run listing11.05/conway_output.go
```

```
[NODO 2 - OUTPUT] Escuchando en localhost:9002 (K = 4)
```

## Paso 2: levantar el compresor (terminal 2)

```bash
go run listing11.04/conway_compressor.go
```

```
[CONWAY COMPRESSOR] Escuchando en localhost:9001; salida -> localhost:9002
```

## Paso 3: ejecutar el centinela (terminal 3)

```bash
go run listing11.03/centinela_json.go -modo conway
```

```
[CENTINELA SENT] DATA -> localhost:9001
```

## Resultado esperado

Terminal 1 (salida):

```
----------------------------------------
[18:33:29.120] PAQUETE RECIBIDO
  Tipo:   RLE_OUTPUT
  Origen: COMPRESSOR
  Salida en líneas de 4:
    5A3B
    2C1D
    5A
----------------------------------------
```

Repita con `go run listing11.05/conway_output.go -k 6` (detenga antes el nodo
con `Ctrl+C`) y la misma cadena sale como `5A3B2C` / `1D5A`.

## Paso 4 (opcional): probar el nodo solo

El nodo de salida no sabe quién le envía. Con solo la terminal 1 levantada,
envíe un mensaje a mano con `nc` (netcat):

```bash
echo '{"tipo":"RLE_OUTPUT","origen":"YO","contenido":"ABCDEFGHIJ"}' | nc -q1 localhost 9002
```

Sale `ABCD`, `EFGH`, `IJ`. (En macOS, use `nc localhost 9002` sin `-q1`.)

## Paso 5: detener todo

Presione `Ctrl+C` en las terminales 1 y 2.

## Errores comunes

- **`[NODO OUTPUT ERROR] No se pudo escuchar en localhost:9002: ... address already in use`**:
  otro nodo usa el puerto 9002 (por ejemplo, un nodo de papa caliente o un
  worker de fan-out). Deténgalo con `Ctrl+C`.
- **El compresor muestra `No se pudo enviar a localhost:9002`**: este nodo no
  estaba levantado. Levántelo primero.
- **`[NODO OUTPUT ERROR] JSON inválido`**: llegó algo que no es JSON.

## Para pensar y ejercicios

1. **Bloques mezclados.** Cada conexión se atiende en su propia goroutine. Si
   llegan dos mensajes casi a la vez, sus bloques pueden mezclarse en
   pantalla. Pruébelo enviando varios `nc` seguidos y evítelo con un
   `sync.Mutex`.
2. **`-k 0`.** ¿Qué pasa si se lanza con `-k 0`? Revise el bucle de `partir`
   y valide el flag.
3. **Revisar el tipo.** Haga que el nodo ignore los mensajes cuyo `tipo` no
   sea `RLE_OUTPUT`.
