# Listado 14.01: Consenso de una ronda

`consenso_una_ronda.go` implementa el **algoritmo de consenso de una ronda**
(Ben-Ari, Algoritmo 12.1). Es el algoritmo "obvio": cada general envía su plan
a todos y decide por mayoría. Funciona si nadie falla, pero **un solo
traidor** basta para que dos generales leales decidan distinto.

**Regla: un general por terminal.** No hay centinela: cada general espera a
que todos estén escuchando y empieza solo.

## Topología: todos con todos

```
        Basil (9201)
         ▲      ▲
         │      │
         ▼      ▼
  Leo (9203) ◄──► Zoe (9205)
```

## El algoritmo

```
p1: plan[myID] <- chooseAttackOrRetreat
p2: for all other generals G
p3:     send(G, myID, plan[myID])
p4: for all other generals G
p5:     receive(G, plan[G])
p6: finalPlan <- majority(plan)
```

- `A` = atacar, `R` = retirarse.
- La mayoría incluye el plan propio. **Empate: gana R.**
- Si un mensaje no llega en `-espera` (2 s), el emisor se considera caído y su
  plan queda como `-`.

## Flags

| Flag | Por defecto | Significado |
|---|---|---|
| `-yo` | `Leo` | nombre de este general |
| `-plan` | `A` | plan inicial: `A` o `R` |
| `-generales` | `Basil,Leo,Zoe` | todos los generales (puertos fijos) o `nombre=host:puerto,...` |
| `-host` | `localhost` | host cuando se usan nombres sin dirección |
| `-espera` | `2s` | tiempo límite por ronda |
| `-caida` | `-1` | **traidor**: deja de enviar después de k mensajes |
| `-mentira` | (vacío) | **traidor**: plan falso por destino, ej. `Leo=A,Zoe=R` |

## Paso 0: abrir 3 terminales

En **cada** terminal:

```bash
cd sesion_14_Consensus
```

## Pasos 1 a 3: escenario con falla por caída (Ben-Ari, sección 12.3)

Basil es el traidor: elige A, envía su plan a Leo y **cae** antes de
enviárselo a Zoe (`-caida 1`). Leo elige R y Zoe elige A.

```bash
go run listing14.01/consenso_una_ronda.go -yo Basil -plan A -caida 1   # terminal 1
go run listing14.01/consenso_una_ronda.go -yo Leo   -plan R            # terminal 2
go run listing14.01/consenso_una_ronda.go -yo Zoe   -plan A            # terminal 3
```

## Resultado esperado

Terminal 2 (Leo):

```
[Leo] mi plan es R
[Leo]   recibe de Basil: A
[Leo]   recibe de Zoe: A

  general  plan
  Basil    A
  Leo      R
  Zoe      A
  mayoría  A

[Leo] DECISIÓN FINAL: A
```

Terminal 3 (Zoe), 2 segundos después:

```
[Zoe] mi plan es A
[Zoe]   ronda 1: llegaron 1 de 2 mensajes (los demás, caídos)
[Zoe]   recibe de Leo: R

  general  plan
  Basil    -
  Leo      R
  Zoe      A
  mayoría  R

[Zoe] DECISIÓN FINAL: R
```

**Leo decide A y Zoe decide R: no hay consenso.** Son las mismas tablas del
libro. Basil muestra `*** CAÍDO ***` y su decisión no cuenta.

## Paso 4: repetir con un traidor bizantino

Detenga los tres con `Ctrl+C`. Ahora Basil no cae: **miente**. Envía A a Leo y
R a Zoe (Ben-Ari, sección 12.7).

```bash
go run listing14.01/consenso_una_ronda.go -yo Basil -mentira Leo=A,Zoe=R   # terminal 1
go run listing14.01/consenso_una_ronda.go -yo Leo   -plan R                # terminal 2
go run listing14.01/consenso_una_ronda.go -yo Zoe   -plan A                # terminal 3
```

```
[Leo]   recibe de Basil: A
[Leo] DECISIÓN FINAL: A
...
[Zoe]   recibe de Basil: R
[Zoe] DECISIÓN FINAL: R
```

Otra vez sin consenso, y esta vez sin esperar: el traidor responde rápido.

## Paso 5: detener todo

Los generales siguen vivos al decidir, por si otro aún les envía mensajes.
Presione `Ctrl+C` en las tres terminales.

## Errores comunes

- **`Zoe no está en -generales`**: el `-yo` no aparece en la lista.
- **`General desconocido`**: solo hay puertos fijos para Basil, John, Leo, Mike
  y Zoe. Para otros nombres use `nombre=host:puerto`.
- **El general se queda en `escuchando` y no avanza**: falta lanzar a otro
  general de la lista.
- **`Error al escuchar: ... address already in use`**: quedó un general de
  otra ejecución (14.01 a 14.04 usan los mismos puertos).

## Para pensar y ejercicios

1. **Sin traidor.** Ejecute los tres sin `-caida` ni `-mentira`. ¿Hay
   consenso? ¿Por qué?
2. **Más generales.** Con `-generales Basil,John,Leo,Mike,Zoe` y un traidor
   que cae, ¿siempre falla? Busque votos iniciales con los que no importa.
3. **El problema de fondo.** Leo no sabe que Basil es traidor. ¿Qué
   información adicional necesitaría? (Esa es la idea de 14.02.)
