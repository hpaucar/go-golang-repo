# Listado 14.07: Clúster de consenso en Docker Compose

Este listado ejecuta la **votación en malla completa** de los listados 14.05 y
14.06 con **cada nodo en su propio contenedor**. No hay código Go nuevo: dos
Dockerfiles y un `docker-compose.yml` con 5 nodos y un centinela.

**Regla: primero los nodos, al final el centinela.** El centinela solo se
levanta cuando se pide el perfil `test`.

## Archivos

| Archivo | Para qué sirve |
|---|---|
| `Dockerfile` | multi-stage: compila `../listing14.05/consenso_mesh.go` (Go 1.26) y lo copia a Alpine |
| `Dockerfile.centinela` | lo mismo con `../listing14.06/centinela_consenso.go` |
| `docker-compose.yml` | 5 servicios `nodo1`…`nodo5`, el servicio `centinela` y la red `red_consenso_mesh` |

## Topología: malla dentro de una red de Docker

```
                   red de Docker "red_consenso_mesh"
 ┌──────────────────────────────────────────────────────────────┐
 │   nodo1:8000 ◄──► nodo2:8000 ◄──► nodo3:8000                  │
 │       ▲  ╲            ▲            ╱  ▲                        │
 │       ▼   ╲           ▼           ╱   ▼     (todos con todos)  │
 │   nodo5:8000 ◄────────────────► nodo4:8000                    │
 │                                                              │
 │   centinela ──TRIGGER──► nodo1..nodo5   (solo con --profile)  │
 └──────────────────────────────────────────────────────────────┘
```

1. **Los nodos se encuentran por nombre** (DNS interno de Compose). Por eso no
   se usa `-host`: la dirección propia es `nodo1:8000`, la misma que usan los
   demás.
2. **Todos escuchan en el puerto 8000.** Cada contenedor tiene su propia IP.
3. **Ningún puerto se publica al host.** El centinela corre **dentro** de la
   red, como otro contenedor.

## Cómo está escrito el `docker-compose.yml`

```yaml
x-nodo: &nodo                       # bloque común (ancla YAML)
  build:
    context: ..                     # la carpeta de la sesión
    dockerfile: listing14.07/Dockerfile
  image: sesion14-consenso
  networks: [red_consenso_mesh]
  environment:                      # leídas por consenso_mesh.go
    PEERS: nodo1:8000,nodo2:8000,nodo3:8000,nodo4:8000,nodo5:8000
    NODE_PORT: "8000"
    DELAY_MS: "300"

services:
  nodo1:
    <<: *nodo
    command: ["-name=nodo1", "-vote=1"]
  ...
  centinela:
    build: { context: .., dockerfile: listing14.07/Dockerfile.centinela }
    profiles: [test]                # no arranca con "docker compose up"
    depends_on: [nodo1, nodo2, nodo3, nodo4, nodo5]
```

Diferencias con la guía:

- **Ancla YAML** `x-nodo`: los 5 nodos comparten imagen, red y variables. La
  guía repetía el bloque 5 veces.
- **La misma lista `PEERS` para todos**, incluido el propio nodo:
  `consenso_mesh.go` se descarta a sí mismo. Así solo cambian `-name` y
  `-vote`.
- **Una sola imagen** (`image: sesion14-consenso`) para los 5 nodos: se compila
  una vez.

## Requisitos

- Docker con el plugin Compose.
- Internet la primera vez (descarga `golang:1.26-alpine` y `alpine:3.21`).

## Paso 0: abrir 2 terminales

```bash
cd sesion_14_Consensus/listing14.07
```

## Paso 1: construir y levantar los nodos (terminal 1)

```bash
docker compose up --build -d
docker compose logs -f
```

```
nodo1-1  | [nodo1] NODO LISTO | Addr: nodo1:8000 | Voto: OPCIÓN A (1)
nodo2-1  | [nodo2] NODO LISTO | Addr: nodo2:8000 | Voto: OPCIÓN B (2)
...
```

## Paso 2: lanzar el centinela (terminal 2)

```bash
docker compose --profile test run --rm centinela
```

```
[CENTINELA] Enviando TRIGGER (code = 3) a 5 nodos...
[CENTINELA OK] Ronda iniciada en nodo4:8000
...
[CENTINELA] Listo.
```

`--rm` borra el contenedor del centinela al terminar.

## Resultado esperado

Terminal 1 (logs de Compose), fragmento:

```
nodo1-1  | [nodo1] RONDA 1, CONTEO FINAL (5 nodos): Opción A = 3 | Opción B = 2
nodo1-1  | >>> [nodo1] RESULTADO DEL CONSENSO: ¡GO A! <<<
nodo5-1  | [nodo5] RONDA 1, CONTEO FINAL (5 nodos): Opción A = 3 | Opción B = 2
nodo5-1  | >>> [nodo5] RESULTADO DEL CONSENSO: ¡GO A! <<<
...
```

Los **5 nodos** deciden **¡GO A!** Para verlo resumido:

```bash
docker compose logs | grep RESULTADO
```

## Paso 3: detener todo

```bash
docker compose --profile test down
```

## Variantes útiles

```bash
docker compose ps                 # estado de los 5 nodos
docker compose stop nodo5         # "apagar" un nodo (ejercicio de quórum)
docker compose start nodo5
```

Con `nodo5` detenido, el centinela muestra `No se pudo contactar nodo5:8000` y
los otros cuatro se quedan esperando el quinto voto.

## Errores comunes

- **`COPY failed: ... listing14.05/consenso_mesh.go not found`**: se cambió
  `context: ..` o se ejecutó `docker build` a mano desde otra carpeta.
- **El centinela no aparece en `docker compose up`**: es lo esperado; tiene
  `profiles: [test]`.
- **`no such host`**: ese contenedor no está corriendo.
- **No aparece `CONTEO FINAL`**: algún nodo no recibió el `TRIGGER`. Revise
  `docker compose ps`.

## Para pensar y ejercicios

1. **Votos aleatorios.** Cambie `-vote=1` por `-vote=0` en los 5 servicios,
   `docker compose up -d` y lance el centinela 3 veces.
2. **Siete nodos.** Agregue `nodo6` y `nodo7`. ¿Qué más hay que cambiar?
3. **Los generales en Docker.** Levante los 5 generales de 14.04 con Compose.
   Pista: `-generales Basil=basil:9201,John=john:9202,Leo=leo:9203,Mike=mike:9204,Zoe=zoe:9205`.
