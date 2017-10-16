GO-API Rest & Batch Job 

---

_...Pronóstico climático de una Galaxia muy muy lejana 🚀..._

Leer [aqui](https://github.com/adecchi/api-go-googlecloud/blob/master/docs/ejercio.md) el enunciado de TP.

### Preconcideraciones:

Debido a que la presente implementación requiere cálculo trigonométrico para su resolución, es factible que los resultados varien segun la precisión decimal numérica que se seleccione (float32,float64, etc). 

Para la resolución del ejercicio propuesto se usaron las formulas:


- Se propone la siguiente formula para calcular el perímetro:

``` go
var d12 = math.Sqrt(math.Pow((Pos2.x-Pos1.x),2)+math.Pow((Pos2.y-Pos1.y),2))
var d13 = math.Sqrt(math.Pow((Pos3.x-Pos1.x),2)+math.Pow((Pos3.y-Pos1.y),2))
var d23 = math.Sqrt(math.Pow((Pos3.x-Pos2.x),2)+math.Pow((Pos3.y-Pos2.y),2))
var perimetro = d12 + d13 + d23
```

- Para saber si la posición (0,0) `sol` esta dentro del triángulo:

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

Para el correcto funcionamiento de la aplicación, se deberán generar previamente dos tablas. Una de ellas las predicciones desagregadas del clima, y la otra para los conteos de los diferentes estados climaticos, correspondientes a los enunciados del ejercicio propuesto.

Dichas tablas responden al siguiente schema:

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
### Requisitos de Deploy:

- Registrado en Google Cloud.
- Contar con acceso a las API de Google Cloud.
- Descargar, instalar y configurar el `Google Cloud SDK`
- Una instancia PostgreSQL con capacidad de conexión desde `AppEngine` 
- Google SQL Proxy.


### Implementación y tecnologias usadas

El proyecto contiene un servidor desarrollado [go](https://golang.org/) y alojado en APP ENGINE de [Google Cloud](https://console.cloud.google.com). El mismo ejecuta la API Rest.
Para la carga de registro de forma batch en [Google Cloud](https://console.cloud.google.com) se utiliza [go](https://golang.org/) de manera local en la máquina. Una vez finalizada la carga, se pueden consultar el tiempo del clima, como asi también las estadísticas vía las API Rest comentadas en el párrafo anterior.


### Comentarios relevantes

_Al momento de empezar el trabajo, no me encontraba familiarizado con [go](https://golang.org/), ní con [Google Cloud](https://console.cloud.google.com), por tanto, tomé esta situación, como la oportunidad perfecta para un nuevo desafío,  e incluso también para poder aprender los conceptos básicos de este lenguaje. Utilicé como principal referencia la [documentación oficial de go](https://golang.org/doc/) junto a las guías presentadas en su sitio oficial._

En mi ambiente de desarrollo, el servidor de GO tardaba un tiempo considerablemente rápido para la realizacion de los insert en modo batch, no asi en[Google Cloud](https://console.cloud.google.com). Para mitigar un poco este delay, decidí utilizar la herramineta de Google SQL Proxy, con lo cual el insert de 3650 registro me llevo un tiempo promedio de 40 minutos.

### Pendientes

Luego del desarrollo realizado, y en base a mi personalidad proactiva, inquieta y creativa, siento que me quedaron algunos pendientes como: 

- Mantener la conexión abierta a la BD, y cerrar la misma al terminar todos los inserts.
- Unificar métodos de conexion, desconexión a la BD.
- Omitir inicilización de variables, ya que GO las hace por DEFAULT.
- Remover el guión bajo `_`de separación de dos palabras en el nombre de funciones, ya que es mejor práctica usar primera letra en mayúscula de la inicial de la segunda palabra, por `

EJ:`analiza_tri() `pasaria` analizaTri()
- Revisar usabilidad.
- Realizar una mejor documentación del código. Agregar comentarios a todos los metodos para dejar en claro su funcionamiento esperado, parametros que reciben y contexto de ejcución.
- Implementar cache de respuestas, sobretodo para consultas a la api.
- Integrar servicios de monitoreo, para llevar registro de uso, performance y posibles errores no atrapados de la aplicación.

### Setup

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
Previamente deberemos abrir una terminal en nuestra máquina, y estar en la carpeta dónde hemos descargado `cloud_sql_proxy`, `main_load_db.go` y `Google-Cloud-Key.json`.
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
Previamente deberemos abrir una terminal en nuestra máquina, y estar en la carpeta dónde hemos descargado `api_query.go`, `app.yaml`.
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
