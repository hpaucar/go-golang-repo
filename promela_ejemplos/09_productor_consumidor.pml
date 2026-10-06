/**** Un productor envia dos valores por un canal con buffer.
 * El consumidor comprueba que conserva el orden de llegada. ****/
chan q = [1] of { byte };
active proctype Productor() {
    q!10;
    q!20;
}
active proctype Consumidor() {
    byte x;
    q?x;
    assert(x == 10);
    q?x;
    assert(x == 20);
}
