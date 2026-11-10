# Listado 11.03: Centinela JSON para los patrones distribuidos

El **centinela JSON** (`centinela_json.go`) inicia el trabajo de los patrones
de los listados 11.04 a 11.10. Igual que el centinela del listado 11.01, por
sí solo no hace nada: envía datos de prueba al nodo de entrada (por defecto
`localhost:9001`) y termina. Los nodos del patrón tienen que estar escuchando
antes.

**Regla: primero los nodos, al final el centinela.**

## Diferencias con el listado 11.01

| | 11.01 `centinela.go` | 11.03 `centinela_json.go` |
|---|---|---|
| Formato | texto plano, un número por línea | mensajes JSON (`Mensaje`) |
| Patrones | solo el ordenamiento en cadena | Conway, papa caliente, pipeline y fan-out |
| Elección | fija | flag `-modo` |

### ¿Por qué hay un flag `-modo`?

Los cuatro patrones usan el **mismo puerto de entrada, 9001**, y en un puerto
solo puede escuchar un programa a la vez. Por eso se prueba **un patrón por
vez**: se levantan sus nodos, se ejecuta el centinela con el modo
correspondiente y luego se detiene todo con `Ctrl+C` antes de pasar al
siguiente.

| `-modo` | Patrón | Nodos que necesita | Datos que envía |
|---|---|---|---|
| `conway` | Conway (compresión RLE) | 11.04, 11.05 | `"AAAAABBBCCDAAAAA"` |
| `potato` | Papa caliente (anillo) | 11.06 × 3 | un token con `-pasos` (por defecto 5) |
| `pipeline` | Ordenamiento en pipeline | 11.07 × 5 | 45, 12, 89, 3, 27 y luego `FIN` |
| `fanout` | Fan-out / fan-in | 11.08, 11.09 × 2, 11.10 | `[1, 2, …, 10]` |

Otros flags: `-destino host:puerto` (por defecto `localhost:9001`) y
`-pasos N` (solo para el modo `potato`).

## El mensaje JSON: el "contrato" entre nodos

Todos los programas de la sesión usan la misma estructura:

```go
type Mensaje struct {
	Tipo           string `json:"tipo"`                // DATA, TOKEN, TASK, RESULT, FIN
	Origen         string `json:"origen,omitempty"`    // quién envía
	Contenido      string `json:"contenido,omitempty"` // carga en texto
	Valores        []int  `json:"valores,omitempty"`   // carga numérica
	PasosRestantes int    `json:"pasos,omitempty"`     // papa caliente
}
```

- `omitempty` omite los campos vacíos, así cada mensaje solo lleva lo que usa.
  Por ejemplo, en modo `conway` viaja
  `{"tipo":"DATA","origen":"CENTINELA","contenido":"AAAAABBBCCDAAAAA"}`.
- Un nodo que no conoce un campo simplemente lo ignora. Por eso cada nodo
  puede declarar solo los campos que necesita.

## Paso 0: preparar las terminales

Todas las órdenes de esta guía se ejecutan desde la carpeta de la sesión
(donde está `go.mod`). En **cada** terminal:

```bash
cd sesion_11_Distributed_Systems
```

Cada nodo va en su propia terminal. El centinela va siempre en la última.

---

## Modo 1: `conway` (3 terminales)

**Topología:** un pipeline de dos etapas.

```
centinela ──► compresor (9001) ──► salida (9002)
               RLE                 líneas de K caracteres
```

```bash
go run listing11.05/conway_output.go              # terminal 1 (salida, 9002)
go run listing11.04/conway_compressor.go          # terminal 2 (compresor, 9001)
go run listing11.03/centinela_json.go -modo conway   # terminal 3
```

Se levanta primero la salida porque el compresor le envía su resultado.

**Qué ocurre:** el compresor recibe `AAAAABBBCCDAAAAA`, lo comprime con RLE
(cinco A, tres B, dos C, una D, cinco A) y envía `5A3B2C1D5A` al nodo de
salida, que lo imprime en líneas de 4 caracteres.

Salida esperada (terminal 1):

```
[NODO 2 - OUTPUT] Escuchando en localhost:9002 (K = 4)
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

Pruebe `go run listing11.05/conway_output.go -k 6` para cambiar el ancho.

---

## Modo 2: `potato` (4 terminales)

**Topología:** un **anillo**. A diferencia de la cadena del listado 11.01, el
último nodo envía al primero.

```
centinela ──► 9001 ──► 9002 ──► 9003
               ▲                  │
               └──────────────────┘
```

```bash
go run listing11.06/hot_potato_node.go -local localhost:9001 -siguiente localhost:9002   # terminal 1
go run listing11.06/hot_potato_node.go -local localhost:9002 -siguiente localhost:9003   # terminal 2
go run listing11.06/hot_potato_node.go -local localhost:9003 -siguiente localhost:9001   # terminal 3
go run listing11.03/centinela_json.go -modo potato -pasos 5                              # terminal 4
```

**Qué ocurre:** cada nodo recibe el token, resta 1 a `pasos` y lo pasa al
siguiente. Pierde el nodo que recibe el token con `pasos = 0`. Cada nodo
espera 0,5 s (latencia simulada), así se puede seguir el token a ojo.

| Recibe | `pasos` al llegar | Reenvía con |
|---|---|---|
| 9001 (desde el centinela) | 5 | 4 → 9002 |
| 9002 | 4 | 3 → 9003 |
| 9003 | 3 | 2 → 9001 |
| 9001 | 2 | 1 → 9002 |
| 9002 | 1 | 0 → 9003 |
| 9003 | 0 | **pierde** |

Salida esperada (terminal 3):

```
[NODO localhost:9003] Token recibido de localhost:9002 con pasos = 3
[NODO localhost:9003] Reenviando token con pasos = 2 a localhost:9001
[NODO localhost:9003] Token recibido de localhost:9002 con pasos = 0
[NODO localhost:9003] ¡PAPA CALIENTE! El contador llegó a 0: este nodo pierde.
```

Pregunta para la clase: con 3 nodos y entrada en 9001, ¿qué valor de `-pasos`
hace perder a 9001? ¿Y a 9002?

---

## Modo 3: `pipeline` (6 terminales)

**Topología:** una cadena, como en el listado 11.01, pero con 5 nodos.

```
centinela ──► 9001 ──► 9002 ──► 9003 ──► 9004 ──► 9005
```

```bash
go run listing11.07/pipeline_sort_node.go -local localhost:9005                             # terminal 1 (último)
go run listing11.07/pipeline_sort_node.go -local localhost:9004 -siguiente localhost:9005   # terminal 2
go run listing11.07/pipeline_sort_node.go -local localhost:9003 -siguiente localhost:9004   # terminal 3
go run listing11.07/pipeline_sort_node.go -local localhost:9002 -siguiente localhost:9003   # terminal 4
go run listing11.07/pipeline_sort_node.go -local localhost:9001 -siguiente localhost:9002   # terminal 5
go run listing11.03/centinela_json.go -modo pipeline                                        # terminal 6
```

**Qué ocurre:** la misma regla que en el listado 11.01: cada nodo guarda el
menor valor visto y expulsa el mayor al siguiente. Hay dos diferencias
importantes:

1. **No hace falta `-n`.** Al final llega un mensaje `FIN`: cada nodo imprime
   su resultado y propaga el `FIN` al siguiente.
2. **Una sola conexión TCP.** El centinela envía los 5 números y el `FIN` por
   la **misma** conexión, y cada nodo mantiene una conexión abierta con el
   siguiente. TCP conserva el orden **dentro** de una conexión, pero no
   **entre** conexiones distintas. Con una conexión por mensaje (como en
   11.01), el `FIN` podría adelantarse a un dato.

Traza del nodo 9001:

```
[PIPELINE localhost:9001] Primer valor guardado: 45
[PIPELINE localhost:9001] Nuevo MIN 12; expulsa 45
[PIPELINE localhost:9001] Mantiene MIN 12; expulsa 89
[PIPELINE localhost:9001] Nuevo MIN 3; expulsa 12
[PIPELINE localhost:9001] Mantiene MIN 3; expulsa 27
[PIPELINE localhost:9001] RESULTADO: 3
```

Resultado esperado:

| Nodo | 9001 | 9002 | 9003 | 9004 | 9005 |
|---|---|---|---|---|---|
| `RESULTADO` | 3 | 12 | 27 | 45 | 89 |

---

## Modo 4: `fanout` (5 terminales)

**Topología:** reparto y recolección (*scatter-gather*).

```
                         ┌──► Worker-1 (9002) ──┐
centinela ──► master (9001)                      ├──► agregador fan-in (9010)
                         └──► Worker-2 (9003) ──┘
```

Se levantan **de atrás hacia adelante**: primero quien recibe los resultados.

```bash
go run listing11.10/fanin_aggregator.go -esperados 10                    # terminal 1 (9010)
go run listing11.09/fanout_worker.go -puerto 9002 -nombre Worker-1       # terminal 2
go run listing11.09/fanout_worker.go -puerto 9003 -nombre Worker-2       # terminal 3
go run listing11.08/fanout_master.go                                     # terminal 4 (9001)
go run listing11.03/centinela_json.go -modo fanout                       # terminal 5
```

**Qué ocurre:**

1. **Fan-out:** el master divide `[1..10]` en dos bloques, `[1 2 3 4 5]` para
   Worker-1 y `[6 7 8 9 10]` para Worker-2, y los envía en paralelo.
2. Cada worker eleva al cuadrado cada valor (0,2 s por valor, cálculo
   simulado) y envía cada resultado al agregador.
3. **Fan-in:** el agregador suma los resultados. Al llegar el resultado
   número 10 (`-esperados 10`) imprime el total.

Salida esperada (terminal 1):

```
[FAN-IN] +1 de Worker-1 -> suma = 1 (1 recibidos)
[FAN-IN] +36 de Worker-2 -> suma = 37 (2 recibidos)
...
[FAN-IN] +100 de Worker-2 -> suma = 385 (10 recibidos)
[FAN-IN] ===== TOTAL FINAL = 385 =====
```

Los resultados de los dos workers llegan intercalados y en cualquier orden.
La suma no cambia porque es asociativa y conmutativa: esa propiedad permite
"reducir" en paralelo (la idea base de MapReduce).

---

## Detener todo y cambiar de modo

Los nodos siguen escuchando después de terminar. Antes de probar otro modo,
presione `Ctrl+C` en **todas** las terminales de nodos: si no, el puerto 9001
sigue ocupado.

## Errores comunes

- **`[CENTINELA ERROR] No se pudo conectar a localhost:9001: ... connection refused`**:
  se ejecutó el centinela antes que el nodo de entrada.
- **`bind: address already in use`** (en algunos nodos aparece como `panic`):
  un nodo del modo anterior sigue escuchando en ese puerto. Deténgalo con
  `Ctrl+C`.
- **`Modo desconocido: ...`**: el valor de `-modo` está mal escrito. Valores
  válidos: `conway`, `potato`, `pipeline`, `fanout`.
- **El centinela dice `SENT` pero el resultado no tiene sentido**: el modo no
  corresponde a los nodos levantados. Por ejemplo, con `-modo conway` y los
  nodos de papa caliente en 9001, el nodo ignora `contenido`, lee `pasos = 0`
  (el campo no viene en el JSON) y declara que pierde. El JSON no da error:
  los campos desconocidos se ignoran y los que faltan quedan en su valor cero.
- **Fan-in nunca imprime `TOTAL FINAL`**: algún worker no estaba levantado, o
  `-esperados` no coincide con la cantidad de valores. Revise los mensajes
  `no disponible` o `perdido` en las terminales del master y de los workers.
