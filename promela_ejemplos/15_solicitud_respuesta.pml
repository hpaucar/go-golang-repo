/**** El cliente pide incrementar 41 y espera una respuesta.
 * El servidor devuelve 42 por un canal diferente. ****/
mtype = { REQ, ACK };
chan pedido = [1] of { mtype, byte };
chan respuesta = [1] of { mtype, byte };
active proctype Cliente() {
    byte r;
    pedido!REQ,41;
    respuesta?ACK,r;
    assert(r == 42);
}
active proctype Servidor() {
    byte v;
    pedido?REQ,v;
    respuesta!ACK,(v + 1);
}
