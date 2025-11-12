#### **Nivel 1: Fundamentos - Desplegando Goroutines y Esperándolas**

El primer paso es aprender a lanzar tareas concurrentes y, lo más importante, a esperar a que terminen.

**🎯 Objetivo:** Dominar el uso de `go` para iniciar goroutines y `sync.WaitGroup` para sincronizar su finalización.

**Ejercicio 1: "Lanzamiento Sincronizado de Tareas"**

1.  **Misión:** Escribe un programa que simule el procesamiento de 5 tareas. Cada "tarea" será una función que:
    *   Acepta un ID (un número entero) como parámetro.
    *   Imprime "Iniciando tarea [ID]".
    *   Duerme durante un tiempo aleatorio (ej. `time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)`).
    *   Imprime "Tarea [ID] completada".

2.  **Requisito Clave:** La función `main` debe lanzar estas 5 tareas como goroutines y **no debe terminar** hasta que las 5 tareas hayan impreso su mensaje de "completada".

3.  **Conceptos a Usar:**
    *   `go miFuncion(id)`
    *   `sync.WaitGroup`
    *   `wg.Add(1)`
    *   `wg.Done()`
    *   `wg.Wait()`

**Pista:** Recuerda llamar a `wg.Add(1)` en la goroutine principal *antes* de lanzar la goroutine trabajadora. Dentro de la goroutine, usa `defer wg.Done()` para asegurarte de que siempre se notifique la finalización.

----
#### **Nivel 2: Comunicación Básica - El Diálogo entre Goroutines**

**🎯 Objetivo:** Entender el funcionamiento de los canales sin búfer (unbuffered channels) como mecanismo de comunicación y sincronización.

**Ejercicio 2: "El Partido de Ping-Pong"**

1.  **Misión:** Crea un programa con dos goroutines que se pasen un "mensaje" de un lado a otro, como en un partido de ping-pong.
    *   Crea dos canales: `pingCh` y `pongCh`.
    *   La goroutine `pinger` envía la palabra "ping" por `pingCh`, luego espera a recibir "pong" por `pongCh`.
    *   La goroutine `ponger` espera a recibir "ping" por `pingCh`, lo imprime, y luego envía "pong" por `pongCh`.
    *   El `main` debe iniciar el partido y hacer que jueguen 5 rondas.

2.  **Requisito Clave:** Las dos goroutines deben estar perfectamente sincronizadas por los canales. El `pinger` no puede enviar un nuevo "ping" hasta que no haya recibido el "pong" de la ronda anterior.

3.  **Conceptos a Usar:**
    *   `make(chan string)`
    *   `canal <- "mensaje"` (enviar)
    *   `mensaje := <- canal` (recibir)

**Pista:** Piensa en el bloqueo. Cuando una goroutine envía a un canal sin búfer, se detiene hasta que otra goroutine recibe de ese canal. Esta es la "magia" que sincroniza el partido.

----
#### **Nivel 3: Productores y Consumidores - Desacoplando el Trabajo**

**🎯 Objetivo:** Utilizar canales con búfer (buffered channels) y el dúo `close`/`for range` para crear una cola de trabajo.

**Ejercicio 3: "La Línea de Ensamblaje"**

1.  **Misión:** Simula una línea de ensamblaje.
    *   Crea una goroutine **productora** que "fabrica" 10 "productos" (pueden ser `structs` simples o incluso `int`).
    *   La productora pone cada producto en un canal `trabajos`. Este canal debe tener un **búfer de 3**.
    *   Cuando la productora termina de fabricar los 10 productos, debe **cerrar** el canal `trabajos` para indicar que no hay más.
    *   Crea una goroutine **consumidora** que procesa los productos. Debe usar un bucle `for range` sobre el canal `trabajos`.
    *   El procesamiento de cada producto debe tomar un tiempo (ej. `time.Sleep(500 * time.Millisecond)`).

2.  **Requisito Clave:** El programa debe terminar cuando todos los productos hayan sido fabricados y consumidos. El bucle del consumidor debe terminar de forma natural gracias al `close` del canal.

3.  **Conceptos a Usar:**
    *   `make(chan Producto, 3)`
    *   `close(canal)`
    *   `for producto := range canal`

**Pista:** Observa el comportamiento. La productora podrá fabricar 3 productos rápidamente y ponerlos en el búfer antes de tener que esperar a que la consumidora, más lenta, tome uno.

----
#### **Nivel 4: Escalando la Operación - El Pool de Workers**

**🎯 Objetivo:** Construir un pool de workers funcional desde cero, combinando todos los conceptos anteriores.

**Ejercicio 4: "El Centro de Procesamiento de Datos"**

1.  **Misión:** Escribe un programa que procese 50 "lotes de datos".
    *   Crea un canal `tareas` para encolar los 50 lotes (usa `int` del 1 al 50 como IDs de lote).
    *   Crea un canal `resultados` para guardar los resultados del procesamiento.
    *   Lanza **4 goroutines `worker`**. Cada `worker` debe:
        *   Leer un ID de lote del canal `tareas`.
        *   Simular el procesamiento (ej. `time.Sleep`).
        *   Enviar un "resultado" (ej. un `string` como `fmt.Sprintf("Lote %d procesado", id)`) al canal `resultados`.
    *   La goroutine `main` debe:
        *   Llenar el canal `tareas` con los 50 lotes y luego cerrarlo.
        *   Esperar a que todos los lotes sean procesados. (¡Usa un `WaitGroup` para esto!)
        *   Leer los 50 resultados del canal `resultados` e imprimirlos.

2.  **Requisito Clave:** Utiliza un `WaitGroup` para sincronizar la finalización de todas las tareas, no para esperar a los workers directamente. La goroutine que cierra el canal `tareas` debe ser la misma que las crea.

3.  **Conceptos a Usar:**
    *   Todos los anteriores.

**Pista:** La lógica es: `wg.Add(1)` por cada tarea que pones en la cola. El `worker` llama a `defer wg.Done()` al principio. La `main` goroutine espera con `wg.Wait()` antes de cerrar el canal de resultados.

----
#### **Nivel 5: Coordinación Avanzada - Manejando Múltiples Eventos**

**🎯 Objetivo:** Usar la sentencia `select` para manejar múltiples canales y timeouts.

**Ejercicio 5: "El Worker con Límite de Tiempo"**

1.  **Misión:** Modifica el ejercicio anterior. Ahora, algunas tareas son "defectuosas" y tardan demasiado.
    *   Dentro de tu `worker`, el "procesamiento" de un lote debe tardar un tiempo aleatorio. Si el ID del lote es múltiplo de 10, haz que tarde 3 segundos; si no, que tarde 1 segundo.
    *   Usa `select` para que el worker espere en dos canales a la vez:
        1.  El resultado del procesamiento (puedes simularlo haciendo que una goroutine interna envíe el resultado a un canal local tras el `sleep`).
        2.  Un canal de timeout: `time.After(2 * time.Second)`.
    *   Si el resultado llega primero, perfecto. Si el timeout se dispara, el worker debe imprimir un mensaje de "Timeout procesando lote [ID]" y pasar a la siguiente tarea.

2.  **Requisito Clave:** El `select` es el corazón del worker. Debe manejar correctamente ambos casos: el éxito y el timeout.

3.  **Conceptos a Usar:**
    *   `select`
    *   `case`
    *   `time.After`

----
#### **Nivel 6: El Desafío Final - Sincronización de Estado y Cancelación**

**🎯 Objetivo:** Proteger un recurso compartido con `sync.Mutex` y manejar la cancelación a través del paquete `context`.

**Ejercicio 6: "El Contador Concurrente con Cancelación"**

1.  **Misión:** Crea un programa que lanza 10 `workers`. Cada `worker` corre un bucle infinito donde:
    *   Incrementa un contador global compartido.
    *   Duerme durante un corto periodo.
    *   **Protege** el acceso al contador usando un `sync.Mutex` para evitar *race conditions*.
    *   **Verifica** si se ha recibido una señal de cancelación.

2.  **Requisitos Clave:**
    *   El `main` debe crear un `context` que se cancele cuando el programa reciba una señal de interrupción del sistema (como `Ctrl+C`). Para esto, usa `signal.NotifyContext`.
    *   Pasa este `context` a cada worker.
    *   Dentro del bucle del worker, usa un `select` para comprobar `ctx.Done()` o para esperar el `time.Sleep`. Si `ctx.Done()` se activa, el worker debe imprimir un mensaje de despedida y terminar.
    *   Ejecuta tu programa con `go run -race .` para verificar que tu `Mutex` está funcionando y no hay condiciones de carrera.

3.  **Conceptos a Usar:**
    *   `context.Context`
    *   `context.WithCancel` o `signal.NotifyContext`
    *   `ctx.Done()`
    *   `sync.Mutex`
    *   `mu.Lock()` y `mu.Unlock()`