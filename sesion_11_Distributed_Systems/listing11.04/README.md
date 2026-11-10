# Listado 11.04: Compresor de Conway (RLE)

`conway_compressor.go` es el **nodo 1** del problema de Conway distribuido
(Ben-Ari, Alg. 8.2). Recibe una cadena en un mensaje JSON, la comprime con
**Run-Length Encoding** (RLE) y envía el resultado al **nodo 2**, el nodo de
salida del listado 11.05, que la imprime en líneas de ancho fijo.

**Regla: primero los nodos, al final el centinela.** Aquí, además, primero el
nodo de salida y después el compresor, porque el compresor le envía su
resultado.

## Topología: un pipeline de dos etapas

```
centinela ──► compresor (9001) ──► salida (9002)
 (cliente)    (servidor y cliente)  (solo servidor)
              RLE                   líneas de K caracteres
```

| Proceso | Escucha (`net.Listen`) | Envía (`net.Dial`) | Rol |
|---|---|---|---|
| centinela (11.03) | no | a 9001 | solo cliente |
| compresor (11.04) | 9001 | a 9002 (`-salida`) | servidor y cliente |
| salida (11.05) | 9002 | no | solo servidor |

La compresión y la presentación quedan **desacopladas** en dos procesos. Cada
etapa hace una sola cosa y se comunica con la otra solo por mensajes.

## ¿Qué es RLE?

Run-Length Encoding reemplaza cada **corrida** (secuencia de caracteres
iguales y seguidos) por su longitud y el carácter:

```
AAAAA BBB CC D AAAAA   ->   5A 3B 2C 1D 5A   ->   5A3B2C1D5A
```

La función `rle` recorre la cadena una sola vez: **O(n)**.

```go
func rle(s string) string {
	runas := []rune(s)          // por runas, no por bytes
	if len(runas) == 0 {
		return ""
	}
	var b strings.Builder
	actual, cuenta := runas[0], 1
	for _, r := range runas[1:] {
		if r == actual {        // la corrida continúa
			cuenta++
			continue
		}
		b.WriteString(strconv.Itoa(cuenta)) // la corrida terminó: se escribe
		b.WriteRune(actual)
		actual, cuenta = r, 1
	}
	b.WriteString(strconv.Itoa(cuenta))     // la última corrida
	b.WriteRune(actual)
	return b.String()
}
```

Dos detalles de Go:

- **`[]rune(s)`**: un `string` de Go guarda bytes en UTF-8. Letras como `Ñ`
  ocupan 2 bytes. Si se recorriera byte por byte, `Ñ` se partiría en dos
  "caracteres" inválidos. Con runas, `ÑÑÑAB` da `3Ñ1A1B`.
- **`strings.Builder`**: concatenar con `+` crea un string nuevo en cada paso;
  `Builder` va llenando un solo búfer.

## Los mensajes

Entrada (del centinela, puerto 9001):

```json
{"tipo":"DATA","origen":"CENTINELA","contenido":"AAAAABBBCCDAAAAA"}
```

Salida (al nodo 2, puerto 9002):

```json
{"tipo":"RLE_OUTPUT","origen":"COMPRESSOR","contenido":"5A3B2C1D5A"}
```

El compresor solo declara los campos que usa (`tipo`, `origen`, `contenido`);
el resto del JSON se ignora.

## Flags

| Programa | Flag | Por defecto | Significado |
|---|---|---|---|
| compresor | `-local` | `localhost:9001` | dirección donde escucha |
| compresor | `-salida` | `localhost:9002` | dirección del nodo de salida |
| salida (11.05) | `-local` | `localhost:9002` | dirección donde escucha |
| salida (11.05) | `-k` | `4` | caracteres por línea |

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

Terminal 2 (compresor):

```
[NODO COMPRESS] Entrada: AAAAABBBCCDAAAAA -> Salida RLE: 5A3B2C1D5A
```

Terminal 1 (salida):

```
----------------------------------------
[21:00:37.959] PAQUETE RECIBIDO
  Tipo:   RLE_OUTPUT
  Origen: COMPRESSOR
  Salida en líneas de 4:
    5A3B
    2C1D
    5A
----------------------------------------
```

Con `go run listing11.05/conway_output.go -k 6` la misma cadena sale como:

```
    5A3B2C
    1D5A
```

## Paso 4 (opcional): enviar sus propias cadenas

El centinela siempre envía `AAAAABBBCCDAAAAA`. Para probar otras cadenas sin
cambiar el código, envíe el JSON a mano con `nc` (netcat), con los nodos
levantados:

```bash
echo '{"tipo":"DATA","contenido":"ÑÑÑAB"}' | nc -q1 localhost 9001
```

```
[NODO COMPRESS] Entrada: ÑÑÑAB -> Salida RLE: 3Ñ1A1B
```

Esto muestra que el centinela no tiene nada especial: cualquier programa que
abra una conexión TCP y envíe el JSON correcto puede usar el nodo. (En macOS,
use `nc localhost 9001` sin `-q1`.)

## Paso 5: detener todo

Los nodos siguen escuchando. Presione `Ctrl+C` en las terminales 1 y 2.

## Errores comunes

- **`[CENTINELA ERROR] No se pudo conectar a localhost:9001`**: el compresor no
  estaba levantado.
- **`[NODO COMPRESS] No se pudo enviar a localhost:9002: ... connection refused`**:
  el compresor sí comprimió (se ve la línea `Salida RLE`), pero el nodo de
  salida no estaba escuchando y el resultado se perdió. Levante primero el
  nodo de salida.
- **`[NODO COMPRESS] JSON inválido: invalid character ...`**: llegó algo que no
  es JSON (por ejemplo, `echo 'hola' | nc ...`).
- **No pasa nada**: el JSON no trae `contenido` (o viene vacío). El compresor
  no tiene nada que comprimir y no envía ningún mensaje.
- **`panic: listen tcp 127.0.0.1:9001: bind: address already in use`**: otro
  nodo (por ejemplo, de papa caliente o de pipeline) sigue usando el puerto
  9001. Deténgalo con `Ctrl+C`.

## Para pensar y ejercicios

1. **Ben-Ari vs. esta versión.** En el Alg. 8.2 de Ben-Ari, una letra aislada
   se deja tal cual (`D`, no `1D`) y las corridas se cortan en 9 (`12A` se
   escribe `9A3A`). Aquí se sigue la guía del curso (`1D`, `12A`).
   Implemente la regla de Ben-Ari y compare las salidas.
2. **¿Por qué cortar en 9?** Comprima `"111"` y luego `"1"` repetido 12
   veces. ¿Se puede descomprimir siempre sin ambigüedad cuando la entrada
   contiene dígitos? ¿Ayuda que cada número tenga un solo dígito?
3. **El campo `tipo` no se revisa.** El compresor comprime cualquier
   `contenido`, aunque `tipo` no sea `DATA`. Haga que ignore los mensajes de
   otro tipo y lo informe en pantalla.
4. **Tercera etapa.** Agregue un nodo descompresor entre el compresor y la
   salida (o en paralelo) que reconstruya la cadena original y verifique que
   coincide. ¿Qué cambia en la topología y en el orden de arranque?
