# Listado 11.11: Recurso compartido (observador de exclusión mutua)

`recurso_compartido.go` simula un recurso que **solo un nodo puede usar a la
vez**, por ejemplo una impresora o un registro de una base de datos. Cuenta
cuántos nodos están "dentro" y avisa con **`VIOLACIÓN DE EXCLUSIÓN MUTUA`**
si hay más de uno.

No forma parte del algoritmo: **solo observa**. Sirve para comprobar si los
algoritmos de exclusión mutua distribuida de los listados 11.12
(Ricart-Agrawala) y 11.13 (RA con token) funcionan. Es como un `assert` de
Promela, pero en un sistema real.

Guías completas: [`../listing11.12/README.md`](../listing11.12/README.md) y
[`../listing11.13/README.md`](../listing11.13/README.md).

## Lugar en la topología

```
   nodo 1 (9101) ◄──────► nodo 2 (9102)
        ▲   \            /   ▲
        │    \          /    │       mensajes del algoritmo
        │     ▼        ▼     │
        └──── nodo 3 (9103) ─┘
              │   │   │
              ▼   ▼   ▼   ENTRA / SALE
          recurso (9100)
```

## El protocolo con el recurso

Cada vez que un nodo usa el recurso abre **una** conexión y envía:

```json
{"tipo":"ENTRA","origen":2}
{"tipo":"SALE","origen":2}
```

El recurso responde `{"tipo":"ACK"}` a cada mensaje. El nodo no sale de su
sección crítica hasta recibir el ACK de `SALE`. Sin ese ACK, la entrada de
otro nodo podría llegar antes que esta salida, y el recurso marcaría una
violación falsa.

## Flags

| Flag | Por defecto | Significado |
|---|---|---|
| `-local` | `localhost:9100` | dirección donde escucha |

## Paso 1: levantarlo (antes que los nodos)

```bash
cd sesion_11_Distributed_Systems
go run listing11.11/recurso_compartido.go
```

```
[RECURSO] Escuchando en localhost:9100
```

## Paso 2: probarlo solo con dos "nodos" manuales

En otras dos terminales, simule dos nodos que entran a la vez y no salen
(`nc` mantiene la conexión abierta hasta `Ctrl+C`):

```bash
nc localhost 9100      # terminal 2: escriba {"tipo":"ENTRA","origen":1} y Enter
nc localhost 9100      # terminal 3: escriba {"tipo":"ENTRA","origen":2} y Enter
```

Terminal 1:

```
18:40:01.120 ENTRA nodo 1 (dentro = 1)
18:40:05.734 ENTRA nodo 2 (dentro = 2)
18:40:05.734 *** VIOLACIÓN DE EXCLUSIÓN MUTUA #1 ***
```

Las terminales 2 y 3 reciben `{"tipo":"ACK","origen":0}` por cada mensaje.
Escriba `{"tipo":"SALE","origen":1}` en la terminal 2 para ver la salida, y
cierre `nc` con `Ctrl+C`.

## Salida con un algoritmo correcto

Con los tres nodos de 11.12 (3 rondas cada uno) se ven 9 pares ENTRA/SALE, sin
solapamiento y siempre con `dentro = 1`:

```
18:33:39.114 ENTRA nodo 1 (dentro = 1)
18:33:39.615 SALE  nodo 1
18:33:39.735 ENTRA nodo 2 (dentro = 1)
18:33:40.236 SALE  nodo 2
...
```

El programa no imprime un resumen al pulsar `Ctrl+C`: cuente las violaciones
en pantalla. El número `#N` de la última violación es el total.

## Errores comunes

- **Los nodos muestran `recurso no disponible`**: el recurso no estaba
  levantado. Los nodos siguen funcionando, pero nadie verifica la exclusión
  mutua.
- **`panic: listen tcp 127.0.0.1:9100: bind: address already in use`**: hay
  otro recurso activo.
- **`dentro` queda en 1 o más tras detener los nodos**: un nodo se detuvo con
  `Ctrl+C` dentro de su sección crítica y no envió `SALE`. Reinicie el recurso.

## Para pensar y ejercicios

1. **Resumen final.** Capture `Ctrl+C` (`os/signal`) e imprima el total de
   entradas y de violaciones antes de salir.
2. **Sin ACK.** Quite el ACK en el recurso y en los nodos. ¿Aparecen
   violaciones falsas con 11.12?
