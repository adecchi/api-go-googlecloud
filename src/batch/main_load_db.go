package main
import (
	"database/sql"
	"fmt"
	"math"
	_ "github.com/lib/pq"
)
type pos struct {
	x, y float64
}
/* Analizo Signo */
func signo(p float64) bool {
	if p > 0 {
		return (true)
	} else {
		return (false)
	}
}
/* Defino Contadores Globales. */
var alin float64 = 0
var alinsol float64 = 0
var tri float64 = 0
var trisol float64 = 0
var trisolmax float64=0
// Variable para uso del area maxima.
var maxp float64 = 0
/* Analizo si hay punto dentro del triangulo */
func analiza_tri(Pos1 pos, Pos2 pos, Pos3 pos, dia float64) {
	// Variable para uso con los puntos originales del Triangulo.
	var ori float64
	// Defino coordenada (0,0) para saber si esta dentro del Triango. Podria buscar cualquier otra modificando esta coordenada.
	var Pos4 = pos{0, 0}
	// Guardo el resultado, en relaidad solo me interesa el signo.
	ori = (Pos1.x-(Pos3.x))*(Pos2.y-(Pos3.y)) - (Pos1.y-(Pos3.y))*(Pos2.x-(Pos3.x))
	/* a1 por a4 */
	// Guardo el resultado, en relaidad solo me interesa el signo.
	var a1 float64 = (Pos4.x-(Pos3.x))*(Pos2.y-(Pos3.y)) - (Pos4.y-(Pos3.y))*(Pos2.x-(Pos3.x))
	if signo(ori) != signo(a1) {
		tri++
		insertar(dia, 3)
		return
	}
	/* a2 por a4 */
	// Guardo el resultado, en relaidad solo me interesa el signo.
	var a2 float64 = (Pos1.x-(Pos3.x))*(Pos4.y-(Pos3.y)) - (Pos1.y-(Pos3.y))*(Pos4.x-(Pos3.x))
	if signo(ori) != signo(a2) {
		tri++
		insertar(dia, 3)
		return
	}
	/* a3 por a4 */
	// Guardo el resultado, en relaidad solo me interesa el signo.
	var a3 float64 = (Pos1.x-(Pos4.x))*(Pos2.y-(Pos4.y)) - (Pos1.y-(Pos4.y))*(Pos2.x-(Pos4.x))
	if signo(ori) != signo(a3) {
		tri++
		insertar(dia, 3)
		return
	}
	trisol++
	//show("trisol", dia, Pos1, Pos2, Pos3)
	/* Hay TRIANGULO con SOL dentro.*/
	calcular_perimetro(Pos1,Pos2,Pos3,dia)
	insertar(dia, 4)
	//return trisol
}

func calcular_perimetro(Pos1 pos, Pos2 pos, Pos3 pos,dia float64) {
	math.Pow(3,2)
	var d12 float64 = math.Sqrt(math.Pow((Pos2.x-Pos1.x),2)+math.Pow((Pos2.y-Pos1.y),2))
	var d13 float64 = math.Sqrt(math.Pow((Pos3.x-Pos1.x),2)+math.Pow((Pos3.y-Pos1.y),2))
	var d23 float64 = math.Sqrt(math.Pow((Pos3.x-Pos2.x),2)+math.Pow((Pos3.y-Pos2.y),2))
	var per float64 = d12 + d13 + d23
	if per > maxp{
		maxp=per
		trisolmax=dia
	}
}

/* Calculo pendiente de la recta */
func pend(pos1 pos, pos2 pos) float64 {
	return (redondeo((pos2.y-pos1.y)/(pos2.x-pos1.x), 8))
}
/* Redondeo a la cantidad de digitos que me indiquen. */
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
func redondeo(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
/* Funcion Imprime por Pantalla String y Puntos */
func show(s string, t float64, p1 pos, p2 pos, p3 pos) {
	fmt.Println(s, t, p1, p2, p3)
}
/* Valido conexion.*/
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
/* Funcion Insertar, inserta los dias con su respectivo codigo de clima.
codigo_clima 1 Sequia, 2 condiciones optimas, 3 triangulo incorrecto ,4 lluvia
*/
func insertar(d float64, codclima float64) {
	/* Defino datos de conexion */
	const (
		host     = "localhost"
		port     = 5432
		user     = "user_postgres"
		password = "Passwrod_user_postgres"
		dbname   = "Db_Name"
	)
	/* Conexion a la BD */
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	/* Abro la conexion */
	db, err := sql.Open("postgres", psqlInfo)
	/* Verifico error */
	checkErr(err)
	/* Cierro la Conexion*/
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO clima(dia,codigo_clima) VALUES($1,$2);")
	res,err := stmt.Exec(d, codclima)
	fmt.Println(res.LastInsertId())
}

/* Funcion Insertar los estados finales*/
func insertar_status(cod_status float64, valor float64) {
	/* Defino datos de conexion */
	const (
		host     = "localhost"
		port     = 5432
		user     = "user_postgres"
		password = "Passwrod_user_postgres"
		dbname   = "Db_Name"
	)
	/* Conexion a la BD */
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	/* Abro la conexion */
	db, err := sql.Open("postgres", psqlInfo)
	/* Verifico error */
	checkErr(err)
	/* Cierro la Conexion*/
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO clima_status(codigo_status,valor) VALUES($1,$2);")
	res, err := stmt.Exec(cod_status, valor)
	fmt.Println(res.LastInsertId())
}

func main() {
	/* Defino Cant. dias a jecutar el FOR */
	var MAXT float64 = 10 * 365
	/* Especifico las distancias.
	d1: Ferengi. Sentido Horario. (w) 1°/dia.
	d2: Vulcano. Sentido AntiHorario. (w) 5°/dia.
	d3: Betasoide. Sentido Horario. (w) 3°/dia.
	*/
	var d1 float64 = 500
	var d2 float64 = 1000
	var d3 float64 = 2000
	/* Como no especifica, supongo que arrancan alineados a los 90°*/
	var Pos1 = pos{0, d1}
	var Pos2 = pos{0, d2}
	var Pos3 = pos{0, d3}
	/* Paso a radianes los 90° que es el punto inicial y lo que vamos a sumar a lo que se mueba el planeta */
	var phi1 float64 = math.Pi / 2
	var phi2 float64 = math.Pi / 2
	var phi3 float64 = math.Pi / 2
	/* Velocidades Angulares (w). Tomo como + para el lado de los - ya que los angulos aumentan*/
	var w1 float64 = -math.Pi / 180
	var w2 float64 = -3 * math.Pi / 180
	var w3 float64 = 5 * math.Pi / 180
	/* Imprimo posiciones en consola */
	/* show("Show", 1, Pos1, Pos2, Pos3) */
	var t float64
	/* Realizamos un for hasta que se cumpla la cantidad especificada en MAXT */
	for t = 1; t <= MAXT; t++ {
		Pos1.x = d1 * math.Cos(w1*t+phi1)
		Pos1.y = d1 * math.Sin(w1*t+phi1)

		Pos2.x = d2 * math.Cos(w2*t+phi2)
		Pos2.y = d2 * math.Sin(w2*t+phi2)

		Pos3.x = d3 * math.Cos(w3*t+phi3)
		Pos3.y = d3 * math.Sin(w3*t+phi3)

		if redondeo(Pos1.x, 8) == redondeo(Pos2.x, 8) && redondeo(Pos2.x, 8) == redondeo(Pos3.x, 8) {
			if redondeo(Pos1.x, 8) == 0 {
				/* Aliñado con el SOL Cuando X iguales y X1 es 0*/
				alinsol++
				insertar(t, 1)
			} else {
				/* Aliñador solo Planetas, sin alinear con SOL*/
				alin++
				insertar(t, 2)
				continue
			}
		} else {
			if redondeo(Pos1.x, 8) == redondeo(Pos2.x, 8) || redondeo(Pos2.x, 8) == redondeo(Pos3.x, 8) || redondeo(Pos1.x, 8) == redondeo(Pos3.x, 8) {
				/* No hay Alineacion */
				/* Analizo Triangulo paso los 3 puntos*/
				analiza_tri(Pos1, Pos2, Pos3, t)
			} else {
				if pend(Pos2, Pos1) == pend(Pos3, Pos2) {
					/* Hay alineacion */
					/* Varifico Alineacion con SOL*/
					if redondeo(Pos3.y, 8) == redondeo((pend(Pos3, Pos2)*Pos3.x), 8) {
						alinsol++
						insertar(t, 1)
					} else {
						/* Alineacion planetas solamente */
						alin++
						insertar(t, 2)
					}
				} else {
					/* Analizo Triangulo paso los 3 puntos*/
					analiza_tri(Pos1, Pos2, Pos3, t)
				}
			}
		}
	}
	/*
	1 alinsol = Sequia.
	2 alin = Condiciones óptimas de presión y temperatura.
	3 tri = Triangulos Incorrectos.
	4 Trisol= LLuvia.
	5 trisolmax = Dia LLuvia Max.
	*/
	insertar_status(1,alinsol)
	//fmt.Println("Días de Sequía: ", alinsol)
	insertar_status(2,alin)
	//fmt.Println("Condiciones óptimas de presión y temperatura: ", alin)
	insertar_status(3,tri)
	//fmt.Println("Triángulados sin SOL: ", tri)
	insertar_status(4,trisol)
	//fmt.Println("Días de LLuvias", trisol)
	insertar_status(5,trisolmax)
	//fmt.Println("El día de lluvia Max. es: ",trisolmax)
}