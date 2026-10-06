/**** Ambos procesos esperan recibir antes de enviar.
 * Como los canales empiezan vacios, ninguno puede avanzar. ****/
chan a = [1] of { byte };
chan b = [1] of { byte };
active proctype P() {
    byte x;
    a?x;
    b!1;
}
active proctype Q() {
    byte y;
    b?y;
    a!1;
}
