/**** Dos procesos modifican una variable compartida.
 * La asercion puede fallar si Q cambia x antes de comprobarla. ****/
byte x = 1;
active proctype P() {
    x++;
    assert(x == 2);
}
active proctype Q() {
    x--;
}
