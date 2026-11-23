# Listado 14.06: Centinela de la votación

`centinela_consenso.go` envía la orden **`TRIGGER`** (`code = 3`) a todos los
nodos de la malla (listado 14.05) **a la vez**, para que empiece una ronda de
consenso. Es un cliente: no escucha, solo envía y termina.

Parte del `centinela.go` de la guía de laboratorio.

## Cómo funciona

```go
for _, n := range lista {
	wg.Add(1)
	go func(addr string) {      // una goroutine por nodo
		defer wg.Done()
		conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
		...
		json.NewEncoder(conn).Encode(msg)
	}(destino)
}
wg.Wait()                       // espera a que terminen todos los envíos
```

- Los envíos son **concurrentes**: así los nodos empiezan a votar casi al
  mismo tiempo.
- `sync.WaitGroup` evita que `main` termine antes de que salgan los mensajes.
- `ctrigger` y el formato `tmsg` deben coincidir con `consenso_mesh.go`.

## Flags

| Flag | Por defecto | Significado |
|---|---|---|
| `-nodes` | `nodo1:8000,...,nodo5:8000` | nodos de la malla (los nombres de Docker) |

En el host use las direcciones locales; en Docker, los valores por defecto.

## Paso 1: levantar los nodos

Siga los pasos 0 a 5 del [listado 14.05](../listing14.05/README.md).

## Paso 2: ejecutar el centinela

```bash
cd sesion_14_Consensus
go run listing14.06/centinela_consenso.go \
  -nodes localhost:8001,localhost:8002,localhost:8003,localhost:8004,localhost:8005
```

## Resultado esperado

```
[CENTINELA] Enviando TRIGGER (code = 3) a 5 nodos...
[CENTINELA OK] Ronda iniciada en localhost:8004
[CENTINELA OK] Ronda iniciada en localhost:8002
[CENTINELA OK] Ronda iniciada en localhost:8001
[CENTINELA OK] Ronda iniciada en localhost:8005
[CENTINELA OK] Ronda iniciada en localhost:8003
[CENTINELA] Listo.
```

El orden de las líneas OK cambia: cada envío es una goroutine. En los nodos
aparece el conteo final (ver 14.05). Cada ejecución inicia una ronda nueva.

## Errores comunes

- **`[CENTINELA ERROR] No se pudo contactar nodo1:8000: ... no such host`**:
  se ejecutó en el host con los valores por defecto (nombres de Docker). Use
  `-nodes localhost:8001,...`.
- **`connection refused`**: ese nodo no está levantado.

## Para pensar y ejercicios

1. **Sin goroutines.** Envíe los `TRIGGER` uno por uno en un bucle normal. ¿Se
   nota la diferencia en los tiempos de los nodos?
2. **Disparo parcial.** Envíe el `TRIGGER` solo a 3 nodos. ¿Qué pasa con los
   otros 2? (Pista: igual reciben votos de sus pares.)
3. **Esperar el resultado.** Haga que cada nodo le responda al centinela con
   su decisión, y que el centinela verifique que las 5 coinciden.
