/**** Un canal sin buffer sincroniza envio y recepcion.
 * Ningun mensaje puede quedar almacenado dentro del canal. ****/
chan q = [0] of { byte };
active proctype Emisor() {
    q!7;
}
active proctype Receptor() {
    byte x;
    q?x;
    assert(x == 7);
}
