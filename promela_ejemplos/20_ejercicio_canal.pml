/**** Dos envios consecutivos intentan llenar un canal de capacidad uno.
 * El proceso queda bloqueado antes de alcanzar su recepcion. ****/
chan q = [1] of { byte };
init {
    byte x;
    q!1;
    q!2;
    q?x;
}
