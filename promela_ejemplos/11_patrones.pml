/**** La recepcion distingue una etiqueta de una variable destino.
 * eval exige coincidencia con el valor actual de una variable. ****/
mtype = { DATO, ACK };
chan q = [2] of { mtype, byte };
init {
    byte valor;
    byte esperado = 7;
    q!DATO,7;
    q?DATO,valor;
    assert(valor == 7);
    q!ACK,7;
    q?ACK,eval(esperado);
}
