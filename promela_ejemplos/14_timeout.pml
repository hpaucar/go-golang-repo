/**** timeout se habilita cuando el sistema no puede avanzar.
 * Representa un bloqueo global y no una duracion en segundos. ****/
chan q = [1] of { byte };
init {
    byte x;
    if
    :: q?x -> printf("Recibido %d\n", x)
    :: timeout -> printf("Sin avance posible\n")
    fi;
}
