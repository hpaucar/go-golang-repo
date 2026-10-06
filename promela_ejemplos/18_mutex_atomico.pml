/**** La consulta y la toma del candado forman un paso protegido.
 * La asercion comprueba la exclusion mutua en dos usuarios. ****/
bool ocupado = false;
byte dentro = 0;
active [2] proctype Usuario() {
    atomic {
        (!ocupado) -> ocupado = true;
    }
    dentro++;
    assert(dentro == 1);
    dentro--;
    ocupado = false;
}
