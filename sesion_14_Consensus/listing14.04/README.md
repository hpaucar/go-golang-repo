# Listado 14.04: Algoritmo del Rey

`rey.go` implementa el **algoritmo del Rey** de Berman y Garay (Ben-Ari,
Algoritmo 12.4). Tolera **t traidores bizantinos** con **n ≥ 4t+1** generales:
uno más por traidor que los Generales Bizantinos, pero con **muchos menos
mensajes**.

## La idea

- Un traidor no puede cambiar una mayoría **abrumadora** (4 a 1 o 5 a 0).
  Solo puede influir cuando el voto está cerrado (3 a 2).
- En cada fase un general distinto es el **rey**. Si un general no tiene una
  mayoría abrumadora, obedece al rey.
- Con dos fases y un solo traidor, **al menos un rey es leal**. Ese rey deja a
  todos los leales con el mismo plan, y desde ahí la mayoría ya es abrumadora.

## El algoritmo

| Ronda | Qué pasa | Ben-Ari |
|---|---|---|
| 1 (y 3) | todos envían su plan a todos; cada uno calcula `myMajority` y `votesMajority` | p3–p8 |
| 2 (y 4) | el rey envía su `myMajority`; los demás lo adoptan **solo si** `votesMajority ≤ n/2 + t` | p9–p16 |
| final | `finalPlan = plan[myID]` | p17 |

```go
if votos > umbral || planRey == "" { // mayoría abrumadora: la mantengo
	plan[yo] = miMayoria
} else {                             // si no, obedezco al rey
	plan[yo] = planRey
}
```

Con n = 5 y t = 1, `umbral = 5/2 + 1 = 3`: hacen falta **más de 3 votos**.

Mensajes: (t+1)(n+1)(n−1) = 2 · 6 · 4 = **48** para t = 1.

| t | Generales Bizantinos | Rey |
|---|---|---|
| 1 | 4 generales, 36 mensajes | 5 generales, 48 mensajes |
| 2 | 7 generales, 392 mensajes | 9 generales, 240 mensajes |
| 3 | 10 generales, 1790 mensajes | 13 generales, 672 mensajes |
| 4 | 13 generales, 5408 mensajes | 17 generales, 1440 mensajes |

## Flags

| Flag | Por defecto | Significado |
|---|---|---|
| `-yo` | `Leo` | nombre de este general |
| `-plan` | `A` | plan inicial |
| `-reyes` | `Mike,Zoe` | rey de la fase 1 y de la fase 2 (igual en todos) |
| `-t` | `1` | traidores tolerados |
| `-generales` | `Basil,John,Leo,Mike,Zoe` | todos los generales |
| `-espera` | `2s` | tiempo límite por ronda |
| `-mentira` | (vacío) | **traidor**: plan que envía a cada general, en todos sus mensajes |

## Paso 0: abrir 5 terminales

```bash
cd sesion_14_Consensus
```

## Pasos 1 a 5: el traidor es el primer rey (Ben-Ari, sección 12.10)

Los planes están divididos: Basil y John eligen A, Leo y Zoe eligen R. El
traidor Mike envía A a dos generales y R a los otros dos, para dividirlos.

```bash
go run listing14.04/rey.go -yo Basil -plan A                                   # terminal 1
go run listing14.04/rey.go -yo John  -plan A                                   # terminal 2
go run listing14.04/rey.go -yo Leo   -plan R                                   # terminal 3
go run listing14.04/rey.go -yo Zoe   -plan R                                   # terminal 4
go run listing14.04/rey.go -yo Mike  -mentira Basil=R,John=A,Leo=A,Zoe=R       # terminal 5
```

## Resultado esperado

Terminal 2 (John):

```
  ronda  Basil John  Leo   Mike  Zoe   | myMajority votesMajority kingPlan
  1      A     A     R     A     R     | A          3
  2            A                       |                          A
[John]   (ronda 2: rey Mike; mi mayoría no es abrumadora (3 votos): adopto A)
  3      R     A     A     A     R     | A          3
  4            R                       |                          R
[John]   (ronda 4: rey Zoe; mi mayoría no es abrumadora (3 votos): adopto R)
[John] DECISIÓN FINAL: R
```

- **Ronda 1**: John ve A = 3 (con la mentira de Mike), Basil ve R = 3. Ninguna
  mayoría es abrumadora. Son las tablas de la página 276 del libro.
- **Ronda 2**: el rey traidor Mike envía A a John y Leo, R a Basil y Zoe. Todos
  obedecen y siguen divididos (p. 277).
- **Ronda 4**: el rey **leal** Zoe envía R y todos lo adoptan.

Basil, Leo y Zoe también terminan con **R**. **Consenso.**

## Paso 6: el rey leal va primero

Detenga todo y agregue `-reyes Zoe,Mike` a los **cinco**:

```bash
go run listing14.04/rey.go -yo Basil -plan A -reyes Zoe,Mike    # y así los demás
```

Terminal 2 (John):

```
  1      A     A     R     A     R     | A          3
  2            R                       |                          R
  3      R     R     R     A     R     | R          4
  4            R                       |                          A
[John]   (ronda 4: rey Mike; mi mayoría R es abrumadora (4 votos): la mantengo)
[John] DECISIÓN FINAL: R
```

Después de la ronda 2 todos los leales tienen R. En la ronda 3 la mayoría es
abrumadora (4 o 5 votos), así que el rey traidor Mike envía A y **nadie le hace
caso** (lema 12.1, caso 1).

## Paso 7: detener todo

`Ctrl+C` en las cinco terminales.

## Errores comunes

- **`Se necesitan 2 reyes y al menos 4t+1 generales`**: con `-t 1` hacen falta
  5 generales.
- **Los generales no coinciden en quién es rey**: todos deben usar el mismo
  `-reyes`.
- **`address already in use`**: quedó un general de 14.01 a 14.03.

## Para pensar y ejercicios

1. **Ejercicio 9 del libro.** Tres leales eligen A y uno R. ¿Puede el traidor
   cambiar algo? Pruébelo.
2. **Ejercicio 10 del libro.** Dos traidores con 5 generales: ponga
   `-mentira` en Mike y en Zoe, con `-reyes Mike,Zoe`. Muestre que no hay
   consenso.
3. **Ejercicio 11 del libro.** Construya los tres casos del lema 12.1
   (planes 4-0, 3-1 y 2-2 entre los leales).
4. **Contar mensajes.** Cuente los mensajes enviados y compare con 48.
