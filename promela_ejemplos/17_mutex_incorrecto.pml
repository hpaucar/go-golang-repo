/**** Consultar el candado y ocuparlo ocurre en pasos separados.
 * Dos procesos pueden superar la consulta y entrar a la vez. ****/
bool ocupado = false;
byte dentro = 0;
active [2] proctype Usuario() {
    (!ocupado);
    ocupado = true;
    dentro++;
    assert(dentro == 1);
    dentro--;
    ocupado = false;
}
