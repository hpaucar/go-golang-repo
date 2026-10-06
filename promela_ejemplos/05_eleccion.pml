/**** Dos guardas verdaderas permiten dos resultados posibles.
 * La seleccion no da prioridad a la primera alternativa. ****/
init {
    byte x = 2;
    byte y = 0;
    if
    :: (x > 0) -> y = 1
    :: (x < 3) -> y = 2
    fi;
    assert(y == 1 || y == 2);
}
