/**** Cada proceso separa la lectura del contador y su escritura.
 * La asercion detecta una actualizacion perdida despues de ambos avisos. ****/
byte x = 0;
chan fin = [2] of { bit };
proctype Sumar() {
    byte copia;
    copia = x;
    copia++;
    x = copia;
    fin!1;
}
init {
    atomic { run Sumar(); run Sumar(); }
    fin?1;
    fin?1;
    assert(x == 2);
}
