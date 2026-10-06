/**** Una propiedad temporal exige recibir la respuesta.
 * El modelo finito supone canales confiables y un servidor disponible. ****/
mtype = { REQ, ACK };
chan pedido = [1] of { mtype, byte };
chan respuesta = [1] of { mtype, byte };
bool recibido = false;
active proctype Cliente() {
    byte r;
    pedido!REQ,41;
    respuesta?ACK,r;
    assert(r == 42);
    recibido = true;
}
active proctype Servidor() {
    byte v;
    pedido?REQ,v;
    respuesta!ACK,(v + 1);
}
ltl entrega { <> recibido }
