/**** Un contador aumenta hasta llegar a tres.
 * Las guardas controlan la repeticion y break permite salir. ****/
init {
    byte i = 0;
    do
    :: (i < 3) -> i++
    :: (i == 3) -> break
    od;
    assert(i == 3);
}
