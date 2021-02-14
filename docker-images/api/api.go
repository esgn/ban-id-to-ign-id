package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	http.HandleFunc("/position/", BanToIgn)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

const (
	host         = "adr-postgis"
	port         = 5432
	user         = "adr"
	password     = "$PASSWORD"
	dbname       = "adr"
	maxIds       = 500
	errorMessage = `{"error":{"message":"%s"}}`
)

func OpenConnection() *sqlx.DB {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if strings.TrimSpace(str) != "" {
			r = append(r, str)
		}
	}
	return r
}

func HandlePath(idsBan string) ([]string, error) {

	// Passage en minuscule
	idsBan = strings.ToLower(idsBan)

	// Découpage de la chaine suivant les virgules
	idsBanArray := strings.Split(idsBan, ",")

	// Suppression des chaines de caractères vides
	idsBanArray = deleteEmpty(idsBanArray)

	// Vérification de la taille de la liste
	l := len(idsBanArray)
	if l > maxIds {
		return nil, errors.New("Liste de cle_interop dépassant la limite de " + strconv.Itoa(maxIds))
	}
	if l == 0 {
		return nil, errors.New("Veuillez indiquer au moins une cle_interop")
	}

	for i := 0; i < l; i++ {

		idBan := idsBanArray[i]
		idBan = strings.TrimSpace(idBan)
		idBanTokens := strings.Split(idBan, "_")

		if (len(idBanTokens)) < 3 {
			return nil, errors.New(idBan + " est une cle_interop invalide")
		}

		n, err := strconv.Atoi(idBanTokens[2])
		if err != nil {
			return nil, errors.New(idBan + " est une cle_interop invalide")
		}

		idBanTokens[2] = strconv.Itoa(n)

		idBan = strings.Join(idBanTokens, "_")

		idsBanArray[i] = idBan
	}

	return idsBanArray, nil

}

func BanToIgn(w http.ResponseWriter, r *http.Request) {

	idsBan := strings.TrimPrefix(r.URL.Path, "/position/")

	idsBanArray, err := HandlePath(idsBan)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		message := fmt.Sprintf(errorMessage, err)
		fmt.Fprintf(w, message)
		return
	}

	db := OpenConnection()

	query := `SELECT json_build_object(
		'type', 'FeatureCollection',
		'crs',  json_build_object(
			'type',      'name', 
			'properties', json_build_object(
				'name', 'EPSG:4326'  
			)
		), 
		'features', json_agg(
			json_build_object(
				'type',       'Feature',
				'id',        id_ban_position,
				'geometry',   ST_AsGeoJSON(geom)::json,
				'properties', json_build_object(
					'id_ban_position', id_ban_position,
					'id_ban_adresse', id_ban_adresse,
					'id_ign', id_ign,
					'cle_interop', cle_interop,
					'id_ban_group', id_ban_group,
					'id_fantoir', id_fantoir,
					'numero', numero,
					'suffixe', suffixe,
					'nom_voie', nom_voie,
					'code_postal', code_postal,
					'nom_commune', nom_commune,
					'code_insee', code_insee,
					'nom_complementaire', nom_complementaire,
					'pos_name', pos_name,
					'typ_loc', typ_loc,
					'source', source,
					'date_der_maj_group', date_der_maj_group,
					'date_der_maj_hn', date_der_maj_hn,
					'date_der_maj_pos', date_der_maj_pos
				)
			)
		)
	)
	FROM (SELECT * FROM ban_ign.ban
	WHERE lower(cle_interop) IN (?)) as f`

	var result string
	q, args, err := sqlx.In(query, idsBanArray)
	q = db.Rebind(q)
	err = db.QueryRow(q, args...).Scan(&result)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadGateway)
		message := fmt.Sprintf(errorMessage, err)
		fmt.Fprintf(w, message)
		defer db.Close()
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(result))

	defer db.Close()

}
