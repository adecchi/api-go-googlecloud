# GO API Rest and Batch Job 
Ejercicio práctico de API Rest y Batch Job.
Se puede ver una version estable del proyecto cuando me lo soliciten.

- [Ejercicio](#ejercicio)
  - [Especificaciones](#especificaciones)
  - [Implementación y tecnologias usadas](#implementaci%C3%B3n-y-tecnologias-usadas)
  - [Comentarios relevantes](#comentarios-relevantes)
  - [Pendientes](#pendientes)
- [Setup](#setup)
  - [Ejecutar el batch local y poblar BD PostgreSql en Google Cloud](#ejecutar-el-batch-local-y-poblar-bd-postgresql-en-google-cloud)
  - [Levantar el servidor API Rest en APP ENGINE](#levantar-el-servidor-api-rest-en-app-engine)


## Ejercicio
En una galaxia lejana, existen tres civilizaciones. Vulcanos, Ferengis y Betasoides. Cada civilización vive en paz en su respectivo planeta.
Dominan la predicción del clima mediante un complejo sistema informático.
A continuación el diagrama del sistema solar.
`Premisas`:
- El planeta Ferengi se desplaza con una velocidad angular de 1 grados/día en sentido horario. Su distancia con respecto al sol es de 500Km.
- El planeta Betasoide se desplaza con una velocidad angular de 3 grados/día en sentido horario. Su distancia con respecto al sol es de 2000Km.
- El planeta Vulcano se desplaza con una velocidad angular de 5 grados/día en sentido anti­horario, su distancia con respecto al sol es de 1000Km.
- Todas las órbitas son circulares.

Cuando los tres planetas están alineados entre sí y a su vez alineados con respecto al sol, el sistema solar experimenta un período de sequía.
Cuando los tres planetas no están alineados, forman entre sí un triángulo. Es sabido que en el momento en el que el sol se encuentra dentro del triángulo, el sistema solar experimenta un período de lluvia, teniendo éste, un pico de intensidad cuando el perímetro del triángulo está en su máximo.
Las condiciones óptimas de presión y temperatura se dan cuando los tres planetas están alineados entre sí pero no están alineados con el sol.
Realizar un programa informático para poder predecir en los próximos 10 años:
1. ¿Cuántos períodos de sequía habrá?
2. ¿Cuántos períodos de lluvia habrá y qué día será el pico máximo de lluvia?
3. ¿Cuántos períodos de condiciones óptimas de presión y temperatura habrá?

`Bonus:`
Para poder utilizar el sistema como un servicio a las otras civilizaciones, los Vulcanos requieren tener una base de datos con las condiciones meteorológicas de todos los días y brindar una API REST de consulta sobre las condiciones de un día en particular.
1) Generar un modelo de datos con las condiciones de todos los días hasta 10 años en adelante utilizando un job para calcularlas.
2) Generar una API REST la cual devuelve en formato JSON la condición climática del día consultado.
3) Hostear el modelo de datos y la API REST en un cloud computing  de Google APP Engine.

`Ej:`   GET → http://....../clima?dia=566   → Respuesta: {“dia”:566, “clima”:”lluvia”}
### Especificaciones
Para la realización del ejercicio se partio de que los 3 planetas y el sol arrancan alineados desde la posición 90° alineados. Donde se tienen 3 puntos P1(x,y),P2(x,y),P3(x,y)

- Para el cálculo del perimetro se utilizo la siguiente fórmula:
``` go
var d12 = math.Sqrt(math.Pow((Pos2.x-Pos1.x),2)+math.Pow((Pos2.y-Pos1.y),2))
var d13 = math.Sqrt(math.Pow((Pos3.x-Pos1.x),2)+math.Pow((Pos3.y-Pos1.y),2))
var d23 = math.Sqrt(math.Pow((Pos3.x-Pos2.x),2)+math.Pow((Pos3.y-Pos2.y),2))
var perimetro = d12 + d13 + d23
```
- Para el cálculo de ver si la posición (0,0) de ahora en más el sol esta dentro del triángulo se utilizo la siguiente fórmula:
``` go
var ori float64
var Pos4 = pos{0, 0}
ori = (Pos1.x-(Pos3.x))*(Pos2.y-(Pos3.y)) - (Pos1.y-(Pos3.y))*(Pos2.x-(Pos3.x))
var a1 float64 = (Pos4.x-(Pos3.x))*(Pos2.y-(Pos3.y)) - (Pos4.y-(Pos3.y))*(Pos2.x-(Pos3.x))
	if signo(ori) != signo(a1) {
		tri++
		return
	}
var a2 float64 = (Pos1.x-(Pos3.x))*(Pos4.y-(Pos3.y)) - (Pos1.y-(Pos3.y))*(Pos4.x-(Pos3.x))
	if signo(ori) != signo(a2) {
		tri++
		return
	}
var a3 float64 = (Pos1.x-(Pos4.x))*(Pos2.y-(Pos4.y)) - (Pos1.y-(Pos4.y))*(Pos2.x-(Pos4.x))
	if signo(ori) != signo(a3) {
		tri++
		return
	}
trisol++
}
```
- Se deberán generar previamente 2 tablas, dónde se guardaran los resultados de las operaciones, como se muestra a continuación:
``` sql
CREATE TABLE clima (  
  id SERIAL PRIMARY KEY,
  dia INT,
  codigo_clima INT
);

CREATE TABLE clima_status (  
  id SERIAL PRIMARY KEY,
  codigo_status INT,
  valor INT
);
```
- Se deberá haber registrado en Google Cloud, contar con acceso a las API de Google Cloud, haber descargado,instalado y configurado Google Cloud SDK, haber generado una instancia SQL Postgresql dentro del mismo lugar donde se deployará la aplicación API Rest GO, como asi tambien haber descargado Google SQL Proxy y haber creado las dos tablas arriba mencionada.
### Implementación y tecnologias usadas

El proyecto contiene un servidor montado en [go](https://golang.org/) ejecutando en APP ENGINE de [Google Cloud](https://console.cloud.google.com) dónde se ejecutan las API Rest.
Para la carga de registro en forma batch en [Google Cloud](https://console.cloud.google.com) se utiliza [go](https://golang.org/) de manera local en la máquina. Una vez finalizada la carga, se pueden consultar el tiempo del clima, como asi también las estadísticas vía las API Rest comentadas en el párrafo anterior.
### Comentarios relevantes

Al momento de empezar el trabajo, no me encontraba familiarizado con [go](https://golang.org/), ní con [Google Cloud](https://console.cloud.google.com), con lo cual tomé la oportunidad como desafío y también para poder aprender los conceptos básicos de este lenguaje. Utilicé como principal referencia la [documentación oficial de go](https://golang.org/doc/) junto a las guías presentadas en su sitio oficial.
En mi ambiente de desarrollo, el servidor de go tardaba un tiempo considerable rápido para la realizacion de los insert en modo batch, no asi en[Google Cloud](https://console.cloud.google.com). Para mitigar un poco este delay, decidí utilizar la herramineta de Google SQL Proxy, con lo cual el insert de 3650 registro me llevo un tiempo promedio de 40 minutos.
### Pendientes
Me quedaron pendientes al momento de cerrar este trabajo, las siguientes mejoras:
- Mantener la conexión abierta a la BD, y cerrar la misma al terminar todos los inserts.
- Unificar métodos de conexion, desconexión a la BD.
- Omitir inicilización de variables, ya que GO las hace por DEFAULT.
- Remover el guión bajo `_`de separación de dos palabras en el nombre de funciones, ya que es mejor práctica usar primera letra en mayúscula de la inicial de la segunda palabra, por `EJ:`analiza_tri() `pasaria` analizaTri()
- Revisar usabilidad.
- Realizar una mejor documentación del código. Agregar comentarios a todos los metodos para dejar en claro su funcionamiento esperado, parametros que reciben y contexto de ejcución.
- Implementar cache de respuestas, sobretodo para consultas a la api.
- Integrar servicios de monitoreo, para llevar registro de uso, performance y posibles errores no atrapados de la aplicación.

## Setup

Como dependencia del proyecto se encuentra [go](https://golang.org/), como asi también contar con los siguientes imports:
-	database/sql
-	fmt
-	math
-	github.com/lib/pq
-	encoding/json
-	log
-	net/http
-	strconv
-	google.golang.org/appengine
-	os
	
También requiere la instalación para poder realizar los insert remotos en Google Cloud en la máquina local de los siguientes paquetes:

``` bash
$ go get -u github.com/lib/pq
$ go get -u google.golang.org/appengine
```
### Ejecutar el batch local y poblar BD PostgreSql en Google Cloud
Previamente deveremos abrir una terminal en nuestra máquina, y estar en la carpeta dónde hemos descargado `cloud_sql_proxy`, `main_load_db.go` y `Google-Cloud-Key.json`.
Luego ejecutamos los siguientes comandos, dónde se nos solicitara la credencial de acceso a Google Cloud y los datos de conexión a la BD.
``` bash
$ gcloud auth login
$ gcloud config set project "Nombre-Proyecto-Google-Cloud"
$ ./cloud_sql_proxy -instances="Instance-Name"=tcp:5432 -credential_file=/"Google-Cloud-Key.json" &
$ psql -h localhost -U "postgres"
```
**Nota:**
- Deberemos tener bajo los servicios de `postgresql` en nuestra máquina local si lo estamos ejecutando. 

Una vez realizado el procedimiento anterior, abrimos una nueva terminal para ejecutar el codigo sin compilar y comenzar a poblar la BD de Google Cloud.
Podremos observar en la terminal anterior como se van insertando los registros, viendo conexiones y desconexiones desde nuestra máquina.
``` bash
$ go build main_load_db.go
```
### Levantar el servidor API Rest en APP ENGINE
Previamente deveremos abrir una terminal en nuestra máquina, y estar en la carpeta dónde hemos descargado `api_query.go`, `app.yaml`.
Luego ejecutamos los siguientes comandos, dónde se nos solicitara la credencial de acceso a Google Cloud.
**Nota:**
No usar la misma carpeta, para la ejecucion de ambos codigos GO.
``` bash
$ gcloud auth login
$ gcloud config set project "Nombre-Proyecto-Google-Cloud"
$ gcloud app deploy
```
Una vez que la terminal nos devuelve el prompt, podremos ejecutar los siguientes comandos, uno para abrir el browser en la direccion de nuestra API REST y otro para ver el log de acceso a nuestra API REST.
``` bash
$ gcloud app browse
$ gcloud app logs tail -s default
```
