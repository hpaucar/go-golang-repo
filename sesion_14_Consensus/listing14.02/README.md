# Listado 14.02: Generales Bizantinos (dos rondas)

`generales_bizantinos.go` implementa el **algoritmo de los Generales
Bizantinos** de Lamport, Shostak y Pease (Ben-Ari, Algoritmo 12.2). Agrega una
segunda ronda en la que cada general **informa** a los demás lo que recibió.
Así un traidor ya no puede engañar a un leal con un solo mensaje.

Con este listado se ven tres resultados del capítulo:

| Escenario | Resultado |
|---|---|
| 3 generales, traidor que **cae** | consenso |
| 3 generales, traidor **bizantino** | **sin** consenso (hacen falta 3t+1 = 4) |
| 4 generales, traidor bizantino | consenso |

## El algoritmo

| Paso | Qué hace cada general | Ben-Ari |
|---|---|---|
| ronda 1 | envía su plan a los otros n−1 y recibe los de ellos | p2–p5 |
| ronda 2 | a cada general G' le informa el plan que recibió de cada otro G | p6–p11 |
| 1.ª votación | para cada G: mayoría de `plan[G]` y los informes sobre G | p12–p13 |
| 2.ª votación | mayoría de esas mayorías (con su propio plan) | p14–p15 |

En la ronda 2 cada mensaje dice "**G me dijo que su plan es P**" (campos
`Origen`, `Sobre`, `Plan` del JSON):

```go
enviar(destino, Mensaje{Ronda: 2, Origen: yo, Sobre: g, Plan: informe})
```

Mensajes por general: (n−1) + (n−1)(n−2). Con 4 generales: 3 + 6 = 9 por
general, **36 en total**.

## La tabla de cada general

Cada general imprime la tabla del libro. Ejemplo de Leo:

```
  general  plan  <Basil <Zoe   mayoría
  Basil    A            A      A        <- Zoe informa que Basil le dijo A
  Leo      R                   R        <- mi propio plan
  Zoe      A     R             R        <- Basil informa que Zoe le dijo R
  mayoría                      R
```

`<X` es la columna "reported by X" del libro.

## Flags

Los mismos de 14.01 y uno más:

| Flag | Por defecto | Significado |
|---|---|---|
| `-generales` | `Basil,Leo,Zoe` | todos los generales |
| `-caida` | `-1` | **traidor**: deja de enviar después de k mensajes |
| `-mentira` | (vacío) | **traidor**: en la ronda 1 envía el plan P a G (`G=P,...`) |
| `-invertir` | (vacío) | **traidor**: en la ronda 2 informa lo contrario a estos generales |

Orden de envío (importa para `-caida`): ronda 1 a los otros en orden
alfabético; ronda 2, por destino en orden alfabético.

## Paso 0: abrir 3 terminales

```bash
cd sesion_14_Consensus
```

## Pasos 1 a 3: traidor que cae (Ben-Ari, sección 12.5)

Basil cae después de enviar su plan solo a Leo (`-caida 1`):

```bash
go run listing14.02/generales_bizantinos.go -yo Basil -plan A -caida 1   # terminal 1
go run listing14.02/generales_bizantinos.go -yo Leo   -plan R            # terminal 2
go run listing14.02/generales_bizantinos.go -yo Zoe   -plan A            # terminal 3
```

Terminal 3 (Zoe):

```
  general  plan  <Basil <Leo   mayoría
  Basil    -            A      A        <- Zoe no recibió nada de Basil, pero Leo sí
  Leo      R     -             R
  Zoe      A                   A
  mayoría                      A

[Zoe] DECISIÓN FINAL: A
```

Leo también decide **A**: **hay consenso**. Leo le reenvió a Zoe el plan de
Basil. Pruebe también `-caida 3` (Basil envía la ronda 1 completa y un informe
de la ronda 2): otra vez consenso en A.

## Pasos 4 a 6: traidor bizantino con 3 generales (Ben-Ari, sección 12.7)

Basil envía A a los dos, pero en la ronda 2 le **miente a Leo**: le dice que
Zoe eligió R (`-invertir Leo`).

```bash
go run listing14.02/generales_bizantinos.go -yo Basil -plan A -invertir Leo   # terminal 1
go run listing14.02/generales_bizantinos.go -yo Leo   -plan R                 # terminal 2
go run listing14.02/generales_bizantinos.go -yo Zoe   -plan A                 # terminal 3
```

```
[Leo]   ronda 2: Basil informa que el plan de Zoe es R
...
[Leo] DECISIÓN FINAL: R
...
[Zoe] DECISIÓN FINAL: A
```

**Sin consenso.** Leo tiene dos versiones del plan de Zoe (A de Zoe, R de
Basil), el empate se rompe a favor de R. Con 3 generales es imposible: Leo no
sabe si miente Basil o Zoe (teorema 12.3).

## Pasos 7 a 10: 4 generales (Ben-Ari, sección 12.8)

Abra una cuarta terminal. Ahora John se une y **Zoe** es la traidora: envía R a
Basil y Leo, A a John, y miente en todos sus informes.

```bash
G=Basil,John,Leo,Zoe
go run listing14.02/generales_bizantinos.go -generales $G -yo Basil -plan A   # terminal 1
go run listing14.02/generales_bizantinos.go -generales $G -yo John  -plan A   # terminal 2
go run listing14.02/generales_bizantinos.go -generales $G -yo Leo   -plan R   # terminal 3
go run listing14.02/generales_bizantinos.go -generales $G -yo Zoe \
    -mentira Basil=R,Leo=R,John=A -invertir Basil,John,Leo                   # terminal 4
```

Terminal 1 (Basil):

```
  general  plan  <John  <Leo   <Zoe   mayoría
  Basil    A                          A
  John     A            A      R      A
  Leo      R     R             A      R
  Zoe      R     A      R             R
  mayoría                             R

[Basil] DECISIÓN FINAL: R
```

John y Leo también deciden **R**. Los informes falsos de Zoe (columna `<Zoe`)
pierden 2 a 1 contra los informes de dos leales.

## Paso 11: detener todo

`Ctrl+C` en todas las terminales.

## Errores comunes

- **Un general decide antes que los demás y otro tarda 2 s más**: normal; ese
  esperó un mensaje de un traidor caído.
- **Resultados distintos a los de esta guía**: revise que todos usen la misma
  lista `-generales`. Si no, esperan mensajes de generales que no existen.
- **`address already in use`**: quedó un general de 14.01, 14.03 o 14.04.

## Para pensar y ejercicios

1. **Ejercicio 1 del libro.** En el escenario de 4 generales, cambie
   `Basil=R` por `Basil=A` en la mentira de Zoe. ¿Qué deciden ahora?
2. **Contar mensajes.** Con 4 generales, ¿cuántos mensajes envía cada uno?
   Compare con la fórmula n·[(n−1) + (n−1)(n−2)] = 36.
3. **El traidor que cae con 4 generales.** Use `-caida` en Zoe en vez de
   `-mentira`. ¿Cambia la decisión?
4. **Ejercicio 5 del libro.** ¿Cuántas rondas hacen falta para 2 traidores?
   ¿Cuántos generales? Esboce el cambio en el código.
