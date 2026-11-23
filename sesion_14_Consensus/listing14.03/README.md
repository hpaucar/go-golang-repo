# Listado 14.03: Algoritmo de inundación (flooding)

`inundacion.go` implementa el **algoritmo de inundación** (Ben-Ari, Algoritmo
12.3) para consenso con fallas **por caída**. Es muy simple: durante **t+1
rondas**, cada general envía a todos el conjunto de planes que conoce y lo une
con lo que recibe. Al final decide por mayoría.

La idea: con t+1 rondas y a lo sumo t caídas, hay al menos **una ronda sin
caídas nuevas**. En esa ronda todo lo que sabe un leal llega a todos los
demás.

## El algoritmo

```
    plan <- {chooseAttackOrRetreat}
p1: do t+1 times
p2:     for all other generals G
p3:         send(G, plan)
p4:     for all other generals G
p5:         receive(receivedPlan)
p6:         plan <- plan ∪ receivedPlan
p7: finalPlan <- majority(plan)
```

El "conjunto de planes" es un mapa `general -> plan` (para votar hay que saber
de quién es cada plan). Viaja completo en el JSON:

```json
{"ronda":2,"origen":"Basil","planes":{"Basil":"A","John":"A","Leo":"R","Zoe":"R"}}
```

## Flags

| Flag | Por defecto | Significado |
|---|---|---|
| `-yo` | `Leo` | nombre de este general |
| `-plan` | `A` | plan inicial |
| `-t` | `2` | número máximo de caídas; se hacen **t+1** rondas |
| `-generales` | `Basil,John,Leo,Zoe` | todos los generales |
| `-espera` | `2s` | tiempo límite por ronda |
| `-caida` | `-1` | **traidor**: deja de enviar después de k mensajes |

## El escenario (Ben-Ari, sección 12.9)

Cuatro generales, **dos** traidores que caen (t = 2):

| General | Plan | Comportamiento |
|---|---|---|
| Basil | A | ronda 1 completa (3 mensajes) y en la ronda 2 solo alcanza a John (`-caida 4`) |
| John | A | leal |
| Leo | R | envía su plan **solo a Basil** y cae (`-caida 1`) |
| Zoe | R | leal |

El plan de Leo viaja así: Leo → Basil (ronda 1) → John (ronda 2) → Zoe
(ronda 3). Es el último diagrama de la sección 12.9.

## Paso 0: abrir 4 terminales

```bash
cd sesion_14_Consensus
```

## Pasos 1 a 4: con t+1 = 3 rondas

```bash
go run listing14.03/inundacion.go -yo Basil -plan A -caida 4   # terminal 1
go run listing14.03/inundacion.go -yo John  -plan A            # terminal 2
go run listing14.03/inundacion.go -yo Leo   -plan R -caida 1   # terminal 3
go run listing14.03/inundacion.go -yo Zoe   -plan R            # terminal 4
```

## Resultado esperado

Terminal 4 (Zoe). Cada ronda tarda unos 2 s porque Zoe espera a los caídos:

```
[Zoe] mi plan es R; se harán 3 rondas (t = 2)
[Zoe]   ronda 1: por John se entera de que John eligió A
[Zoe]   ronda 1: por Basil se entera de que Basil eligió A
[Zoe] fin de la ronda 1: conoce {Basil=A, John=A, Zoe=R}
[Zoe] fin de la ronda 2: conoce {Basil=A, John=A, Zoe=R}
[Zoe]   ronda 3: por John se entera de que Leo eligió R
[Zoe] fin de la ronda 3: conoce {Basil=A, John=A, Leo=R, Zoe=R}

  mayoría  R

[Zoe] DECISIÓN FINAL: R
```

Terminal 2 (John): se entera de Leo en la ronda 2, por Basil, y también decide
**R**. **Consenso.**

## Paso 5: repetir con solo t = 1 (2 rondas)

Detenga todo y agregue `-t 1` a los **cuatro** generales:

```bash
go run listing14.03/inundacion.go -yo Basil -plan A -caida 4 -t 1   # terminal 1
go run listing14.03/inundacion.go -yo John  -plan A -t 1            # terminal 2
go run listing14.03/inundacion.go -yo Leo   -plan R -caida 1 -t 1   # terminal 3
go run listing14.03/inundacion.go -yo Zoe   -plan R -t 1            # terminal 4
```

```
[John] fin de la ronda 2: conoce {Basil=A, John=A, Leo=R, Zoe=R}
[John] DECISIÓN FINAL: R          <- 2 contra 2, empate -> R
...
[Zoe] fin de la ronda 2: conoce {Basil=A, John=A, Zoe=R}
[Zoe] DECISIÓN FINAL: A           <- nunca supo el plan de Leo
```

**Sin consenso**: faltó la ronda en la que John le cuenta a Zoe el plan de
Leo. Es el ejercicio 8 del libro.

## Paso 6: detener todo

`Ctrl+C` en las cuatro terminales.

## Errores comunes

- **Todo tarda mucho**: cada ronda espera `-espera` por cada general caído.
  Con `-espera 1s` va más rápido.
- **Resultados distintos**: revise que los cuatro usen el mismo `-t`. Si no,
  hacen distinto número de rondas.
- **`address already in use`**: quedó un general de 14.01, 14.02 o 14.04.

## Para pensar y ejercicios

1. **¿Por qué no sirve contra bizantinos?** Un traidor podría enviar
   `{"Leo":"A"}` aunque Leo eligió R. ¿Qué haría falta para detectarlo?
2. **Ejercicio 6 del libro.** ¿Cuántos mensajes se envían como máximo con n
   nodos y t+1 rondas? Compruébelo contando los mensajes.
3. **Ejercicio 7 del libro.** Haga que Leo caiga con `-caida 0` (no envía
   nada). ¿Hay consenso con `-t 1`?
4. **Optimización.** Un general solo necesita enviar los planes **nuevos** de
   cada ronda. Modifique el código y cuente cuántos bytes ahorra.
