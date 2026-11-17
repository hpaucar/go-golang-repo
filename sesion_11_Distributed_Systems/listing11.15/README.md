# Listado 11.15: Anillo de papa caliente en Docker Compose

Este listado ejecuta el **anillo de papa caliente** del listado 11.06, pero
con **cada nodo en su propio contenedor**. No hay código Go nuevo: solo un
`Dockerfile` que empaqueta `hot_potato_node.go` y un `docker-compose.yml` que
levanta tres contenedores y los conecta en anillo.

Es lo más cercano en esta sesión a un sistema distribuido real: cada nodo
tiene su propia red, su propia IP y su propio sistema de archivos.

**Regla: primero los nodos, al final el centinela.** Los nodos los levanta
Compose; el centinela se ejecuta en el host.

## Archivos

| Archivo | Para qué sirve |
|---|---|
| `Dockerfile` | multi-stage: compila `../listing11.06/hot_potato_node.go` y lo copia a una imagen Alpine |
| `docker-compose.yml` | tres servicios (`nodo1`, `nodo2`, `nodo3`) con la misma imagen y una red `anillo` |

## Topología: un anillo dentro de una red de Docker

```
  host                   red de Docker "anillo"
 ┌──────────────┐      ┌───────────────────────────────────────────┐
 │ centinela ───┼─9001─┼─► nodo1:9001 ──► nodo2:9001 ──► nodo3:9001 │
 │ (go run)     │      │      ▲                              │      │
 └──────────────┘      │      └──────────────────────────────┘      │
                       └───────────────────────────────────────────┘
```

Tres ideas clave:

1. **Los nodos se encuentran por nombre.** Compose crea un DNS interno: dentro
   de la red, `nodo2` es la dirección del contenedor `nodo2`. Por eso el
   listado 11.06 recibe direcciones completas (`-siguiente nodo2:9001`) y no
   solo puertos.
2. **Todos escuchan en el puerto 9001.** No hay conflicto: cada contenedor
   tiene su propia IP. En el host, en cambio, los tres nodos necesitaban
   puertos distintos (9001, 9002, 9003).
3. **Solo `nodo1` publica su puerto** (`ports: "9001:9001"`), para que el
   centinela pueda entrar al anillo desde el host. `nodo2` y `nodo3` solo son
   accesibles dentro de la red.

## Cómo está escrito el `docker-compose.yml`

```yaml
x-nodo: &nodo                  # bloque común (ancla YAML)
  build:
    context: ..                # la carpeta de la sesión
    dockerfile: listing11.15/Dockerfile
  image: sesion11-hotpotato
  networks: [anillo]

services:
  nodo1:
    <<: *nodo                  # copia el bloque común
    command: ["-nombre", "nodo1", "-local", ":9001", "-siguiente", "nodo2:9001"]
    ports:
      - "9001:9001"
  nodo2:
    <<: *nodo
    command: ["-nombre", "nodo2", "-local", ":9001", "-siguiente", "nodo3:9001"]
  ...
```

- `context: ..`: el contexto de construcción es la carpeta de la sesión, para
  que el Dockerfile pueda copiar `listing11.06/hot_potato_node.go`.
- `command:` pasa los flags al `ENTRYPOINT` del Dockerfile
  (`hot_potato_node`).
- `-local :9001` escucha en **todas** las interfaces del contenedor. Con
  `localhost:9001` el nodo solo aceptaría conexiones desde dentro de su propio
  contenedor.

## Requisitos

- Docker con el plugin Compose.
- Go, para ejecutar el centinela en el host.
- El puerto 9001 del host libre (detenga los nodos de los patrones anteriores).

## Paso 0: abrir 2 terminales

```bash
cd sesion_11_Distributed_Systems/listing11.15
```

## Paso 1: construir y levantar el anillo (terminal 1)

```bash
docker compose up --build
```

Sin `-d`, los logs de los tres contenedores quedan a la vista en esta
terminal:

```
nodo1-1  | [NODO nodo1] Escuchando en :9001. Siguiente en el anillo: nodo2:9001
nodo2-1  | [NODO nodo2] Escuchando en :9001. Siguiente en el anillo: nodo3:9001
nodo3-1  | [NODO nodo3] Escuchando en :9001. Siguiente en el anillo: nodo1:9001
```

La primera vez tarda más: descarga las imágenes base y compila.

## Paso 2: ejecutar el centinela (terminal 2)

```bash
go run ../listing11.03/centinela_json.go -modo potato -pasos 5
```

```
[CENTINELA SENT] TOKEN -> localhost:9001
```

## Resultado esperado

Terminal 1 (los logs de Compose):

```
nodo1-1  | [NODO nodo1] Token recibido de CENTINELA con pasos = 5
nodo1-1  | [NODO nodo1] Reenviando token con pasos = 4 a nodo2:9001
nodo2-1  | [NODO nodo2] Token recibido de nodo1 con pasos = 4
nodo2-1  | [NODO nodo2] Reenviando token con pasos = 3 a nodo3:9001
nodo3-1  | [NODO nodo3] Token recibido de nodo2 con pasos = 3
nodo3-1  | [NODO nodo3] Reenviando token con pasos = 2 a nodo1:9001
nodo1-1  | [NODO nodo1] Token recibido de nodo3 con pasos = 2
nodo1-1  | [NODO nodo1] Reenviando token con pasos = 1 a nodo2:9001
nodo2-1  | [NODO nodo2] Token recibido de nodo1 con pasos = 1
nodo2-1  | [NODO nodo2] Reenviando token con pasos = 0 a nodo3:9001
nodo3-1  | [NODO nodo3] Token recibido de nodo2 con pasos = 0
nodo3-1  | [NODO nodo3] ¡PAPA CALIENTE! El contador llegó a 0: este nodo pierde.
```

Es el mismo resultado que en el listado 11.06 (pierde el tercer nodo), pero
ahora los mensajes viajan entre contenedores.

## Paso 3: detener todo

En la terminal 1, `Ctrl+C` detiene los contenedores. Para borrarlos junto con
la red:

```bash
docker compose down
```

## Variantes útiles

```bash
docker compose up -d --build     # en segundo plano
docker compose ps                # estado: solo nodo1 muestra 0.0.0.0:9001->9001
docker compose logs -f           # seguir los logs
docker compose stop nodo2        # "romper" el anillo
```

Con `nodo2` detenido, ejecute el centinela otra vez. `nodo1` muestra
`¡Token perdido! No se pudo conectar a nodo2:9001`. Vuelva a levantarlo con
`docker compose start nodo2`.

## Errores comunes

- **`failed to bind host port 0.0.0.0:9001/tcp: address already in use`**:
  un nodo de otro patrón (Conway, papa caliente local, pipeline, fan-out)
  sigue usando el puerto 9001 del host. Deténgalo.
- **`[CENTINELA ERROR] No se pudo conectar a localhost:9001`**: los
  contenedores aún no arrancaron, o se ejecutó `docker compose down`.
- **`¡Token perdido! No se pudo conectar a nodo2:9001: ... no such host`**:
  ese contenedor no está corriendo, así que su nombre no existe en el DNS de
  la red.
- **`COPY failed: ... listing11.06/hot_potato_node.go not found`**: se cambió
  `context: ..` o se ejecutó `docker build` a mano desde otra carpeta.

## Para pensar y ejercicios

1. **Cuatro nodos.** Agregue `nodo4` al anillo. ¿Qué `command:` cambia?
2. **Sin publicar puertos.** Quite `ports:` de `nodo1`. Ahora el centinela
   del host no puede entrar. Agregue un servicio `centinela` en la misma red,
   con su propio Dockerfile, que ejecute
   `centinela_json -modo potato -destino nodo1:9001`.
3. **Ricart-Agrawala en Docker.** Levante los nodos de 11.12 y el recurso de
   11.11 con Compose. Pista: `-nodos 1=nodo1:9101,2=nodo2:9102,3=nodo3:9103`
   y `-recurso recurso:9100`.
