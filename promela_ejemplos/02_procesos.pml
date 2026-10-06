/**** init crea dos instancias de un mismo tipo de proceso.
 * El orden de sus mensajes puede cambiar entre ejecuciones. ****/
proctype Alumno(byte id) {
    printf("Soy el alumno %d\n", id);
}
init {
    run Alumno(1);
    run Alumno(2);
}
