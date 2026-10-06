/**** P inicia el intercambio enviando el primer mensaje.
 * Q puede recibirlo y responder para que ambos terminen. ****/
chan a = [1] of { byte };
chan b = [1] of { byte };
active proctype P() {
    byte x;
    b!1;
    a?x;
}
active proctype Q() {
    byte y;
    b?y;
    a!1;
}
