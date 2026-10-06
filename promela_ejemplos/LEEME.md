# Laboratorio: Promela y SPIN desde cero

Material para Programación Concurrente y Distribuida. Los 21 modelos son independientes y acompañan las diapositivas. Cada archivo contiene una cabecera que explica su propósito.

## Preparación

Necesitas SPIN y un compilador C, por ejemplo GCC, disponibles en la terminal. Consulta la distribución oficial en https://spinroot.com/spin/whatispin.html. Los comandos siguientes usan una terminal Linux o WSL. En Windows nativo, el ejecutable del verificador puede llamarse pan.exe.

Comprueba las herramientas:

```sh
spin -V
gcc --version
```

Trabaja en una carpeta independiente por modelo. SPIN genera archivos llamados pan.c, pan.h y otros que se reemplazan al generar otro verificador. No cambies el modelo entre la generación de una traza y su reproducción.

## Simular

```sh
spin 01_hola.pml
spin -p -g -l 07_carrera.pml
```

`-p` muestra los pasos, `-g` las variables globales y `-l` las locales. Una simulación recorre un camino y puede no encontrar un error existente.

## Verificar y reproducir

```sh
spin -a 07_carrera.pml
gcc -O2 -o pan pan.c
./pan
spin -t -p -g -l 07_carrera.pml
```

El último comando reproduce el archivo `.trail` producido cuando se detecta un error. Para exigir que no queden mensajes en canales al terminar, ejecuta `./pan -q`.

Para la propiedad LTL del modelo 16:

```sh
spin -a 16_ltl_entrega.pml
gcc -O2 -o pan pan.c
./pan -a
```

`spin -a` genera el verificador. `pan -a` activa la búsqueda de ciclos de aceptación. Son opciones de programas distintos.

## Resultados esperados

| Archivo | Propósito | Resultado esperado |
|---|---|---|
| 01_hola | Primer proceso | Saludo y terminación |
| 02_procesos | init, run y parámetros | Dos procesos, distintos órdenes posibles |
| 03_espera | Condición bloqueante | Aserción satisfecha |
| 04_assert_falla | Intercalado antes de assert | Violación de aserción |
| 05_eleccion | Guardas simultáneamente habilitadas | y vale 1 o 2 |
| 06_bucle | do y break | i termina en 3 |
| 07_carrera | Lectura y escritura separadas | Violación de aserción |
| 08_carrera_atomic | Actualización protegida | x termina en 2 |
| 09_productor_consumidor | Canal con buffer | Recepción de 10 y luego 20 |
| 10_rendezvous | Canal sin buffer | Recepción de 7 |
| 11_patrones | mtype y eval | Coincidencia de etiquetas y valores |
| 12_deadlock | Ambos reciben primero | invalid end state |
| 13_deadlock_corregido | Un proceso envía primero | Ambos terminan |
| 14_timeout | Bloqueo global | Mensaje “Sin avance posible” |
| 15_solicitud_respuesta | Cliente y servidor | Respuesta 42 |
| 16_ltl_entrega | Propiedad temporal | Se cumple <> recibido |
| 17_mutex_incorrecto | Consulta y toma separadas | Violación de exclusión mutua |
| 18_mutex_atomico | Adquisición indivisible | Exclusión mutua en dos visitas |
| 19_ejercicio_if | Dos alternativas habilitadas | Violación de aserción |
| 20_ejercicio_canal | Buffer lleno antes de recibir | invalid end state |
| 21_ejercicio_canal_solucion | Liberación de espacio | Ambos valores recibidos |

Los modelos marcados con errores son intencionales. La salida `errors: 0` se interpreta junto con la propiedad, opciones y finalización de la búsqueda. Una verificación del modelo no demuestra por sí sola la corrección de un programa real.

## Práctica

1. Predice qué ocurrirá antes de ejecutar.
2. Describe los pasos que permiten un error.
3. Compara la versión incorrecta con la corregida.
4. Explica qué supuesto representa cada `atomic`.
5. Cambia la capacidad de un canal y revisa su efecto.
6. En el cliente/servidor, permite omitir una respuesta con una elección no determinista y explica por qué el cliente puede quedar bloqueado.

## Fuentes

- Promela/SPIN, PDF de 50 páginas proporcionado por el docente, con reconocimiento a material de Flavio Lerda y al curso de Ed Clarke.
- Wishnu Prasetya, Modeling and Verification with Promela, PDF de 47 páginas proporcionado por el docente.
- Manual oficial: https://spinroot.com/spin/Man/promela.html

Los ejemplos son adaptaciones introductorias. Los canales básicos suponen comunicación confiable. `timeout` representa falta global de avance y no tiempo físico. `x++` es una instrucción indivisible en Promela; el ejemplo de carrera separa explícitamente lectura y escritura.
