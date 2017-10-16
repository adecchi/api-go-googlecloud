package main
import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"google.golang.org/appengine"
	"os"
	_ "github.com/lib/pq"
)
/* Modelo de Datos */
type climastruct struct {
	id           int `json:"id"`
	dia          int `json:"dia"`
	codigo_clima int `json:"codigo_clima"`
}
/* Modelo de Datos 
Defino en Mayusculas porque tienen que ser accedidas desde librerias externas, para poder 
pasar a JSON.
*/
type clima_response struct {
	Dia   int    `json:"Dia"`
	Clima string `json:"Clima"`
}
type clima_status struct {
	Id   int    `json:"Id"`
	Codigo_status   string    `json:"Referencia"` //estaba int
	Valor string `json:"Valor"`
}
type Climas []clima_response

type StatusClima struct{
	Estado [] clima_status
}
/* Obtengo los resultados totales del clima */
func Status(w http.ResponseWriter, r *http.Request) {
	estado := StatusClima{}
	err := getStatus(&estado)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, estado)
}

func Index(w http.ResponseWriter, r *http.Request) {
	var err error
	clima := climastruct{}
	clima.dia, err = strconv.Atoi(r.FormValue("dia"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else {
		if clima.dia < 1 || clima.dia > 3650 {
			clim := Climas{
				clima_response{Dia: clima.dia, Clima: "Fuera de Rango"},
			}
			respondWithJSON(w, http.StatusOK, clim)
		}
	}
	clima_id := clima.getClima()
	switch clima_id {
	case 1:
		clim := Climas{
			clima_response{Dia: clima.dia, Clima: "Sequia"},
		}
		respondWithJSON(w, http.StatusOK, clim)
	case 2:
		clim := Climas{
			clima_response{Dia: clima.dia, Clima: "Conciciones Optimas de Pre. y Temp."},
		}
		respondWithJSON(w, http.StatusOK, clim)
	case 3:
		clim := Climas{
			clima_response{Dia: clima.dia, Clima: "Triangulo Incorrecto"},
		}
		respondWithJSON(w, http.StatusOK, clim)
	case 4:
		clim := Climas{
			clima_response{Dia: clima.dia, Clima: "Lluvia"},
		}
		respondWithJSON(w, http.StatusOK, clim)
	default:
	}
}

/* Valido conexion.*/
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
/* Buscar Clima.*/
func getStatus(clista *StatusClima) error{
	var err error
	var db *sql.DB
	datastoreName := os.Getenv("POSTGRES_CONNECTION")
	db, err = sql.Open("postgres", datastoreName)
	checkErr(err)
	defer db.Close()
	rows,err := db.Query("SELECT * FROM clima_status")
	switch {
	case err == sql.ErrNoRows:
	//log.Printf("No Informacion de Día.")
	case err != nil:
		log.Fatal(err)
	default:
	}
	defer rows.Close()
	for rows.Next() {
		climast := clima_status{}
		err = rows.Scan(
			&climast.Id,
			&climast.Codigo_status,
			&climast.Valor,
		)
		if err != nil {
			return err
		}
		switch {
		case climast.Codigo_status == strconv.Itoa(1):
			climast.Codigo_status = "Cantidad dias de SEQUIA"
		case climast.Codigo_status == strconv.Itoa(2):
			climast.Codigo_status = "Cantidad dias Condiciones optimas de presion y temperatura"
		case climast.Codigo_status == strconv.Itoa(3):
			climast.Codigo_status = "Cantidad dias Triangulos sin SOL"
		case climast.Codigo_status == strconv.Itoa(4):
			climast.Codigo_status = "Cantidad dias de Lluvias"
		case climast.Codigo_status == strconv.Itoa(5):
			climast.Codigo_status = "Dia Max. LLuvia"
		default:
		}
		clista.Estado = append(clista.Estado, climast)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

func (cli *climastruct) getClima() int {
	var err error
	var db *sql.DB
	/* Defino los datos de conexion a APPGINE GOOGLE. */
	datastoreName := os.Getenv("POSTGRES_CONNECTION")
	db, err = sql.Open("postgres", datastoreName)
	/* Verifico error */
	checkErr(err)
	/* Cierro la Conexion*/
	defer db.Close()
	err = db.QueryRow("SELECT codigo_clima FROM clima WHERE dia = $1", cli.dia).Scan(&cli.codigo_clima)
	switch {
	case err == sql.ErrNoRows:
	//log.Printf("No Informacion de Día.")
	case err != nil:
		log.Fatal(err)
	default:
	}
	return (cli.codigo_clima)
}
/* Para poder devolver JSON si o si los nombres de los campos del Struct Tienen que comenzar en mayuscula.*/
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
/* Para Manejo de errores */
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func main() {
	/* Para saber el clima del dia*/
	http.HandleFunc("/clima", Index)
	/* Para saber los Totales*/
	http.HandleFunc("/clima/status", Status)
	log.Fatal(http.ListenAndServe(":8080", nil))
	appengine.Main()

}
