/**** La lectura y la escritura forman una region atomica.
 * La comprobacion ocurre despues de recibir los dos avisos de fin. ****/
byte x = 0;
chan fin = [2] of { bit };
proctype Sumar() {
    byte copia;
    atomic {
        copia = x;
        copia++;
        x = copia;
    }
    fin!1;
}
init {
    atomic { run Sumar(); run Sumar(); }
    fin?1;
    fin?1;
    assert(x == 2);
}
