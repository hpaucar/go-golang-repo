/**** El ejemplo permite que se ejecute cualquiera de dos ramas.
 * La propiedad propuesta es falsa si la segunda rama asigna dos. ****/
init {
    byte x = 0;
    if
    :: (x == 0) -> x = 1
    :: true -> x = 2
    fi;
    assert(x == 1);
}
