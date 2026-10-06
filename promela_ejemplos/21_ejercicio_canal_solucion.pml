/**** Cada envio va seguido de una recepcion que libera espacio.
 * Las aserciones comprueban ambos valores y el canal termina vacio. ****/
chan q = [1] of { byte };
init {
    byte x;
    q!1;
    q?x;
    assert(x == 1);
    q!2;
    q?x;
    assert(x == 2);
}
