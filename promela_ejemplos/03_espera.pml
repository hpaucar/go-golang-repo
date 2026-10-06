/**** La condicion detiene al lector hasta que exista un dato.
 * El escritor publica primero el dato y despues la senal. ****/
byte dato = 0;
bool listo = false;
active proctype Escritor() {
    dato = 7;
    listo = true;
}
active proctype Lector() {
    (listo);
    assert(dato == 7);
}
