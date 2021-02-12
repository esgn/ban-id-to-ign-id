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
	host          = "127.0.0.1"
	port          = 5432
	user          = "adr"
	password      = "adr"
	dbname        = "adr"
	max_ids       = 3
	error_message = `{"error":{"message":"%s"}}`
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

func HandlePath(idsBan string) ([]string, error) {

	// Suppression espaces éventuels
	idsBan = strings.TrimSpace(idsBan)

	// Suppression virgule parasite
	idsBan = strings.TrimSuffix(idsBan, ",")

	// Passage en minuscule
	idsBan = strings.ToLower(idsBan)

	// Découpage de la chaine suivant les virgules
	idsBanArray := strings.Split(idsBan, ",")

	// Vérification de la taille de la liste
	l := len(idsBanArray)
	if l > max_ids {
		return nil, errors.New("Liste d'identifiant dépassant la limite de " + strconv.Itoa(max_ids))
	}

	for i := 0; i < l; i++ {

		idBan := idsBanArray[i]
		idBan = strings.TrimSpace(idBan)

		idBanTokens := strings.Split(idBan, "_")

		if (len(idBan)) < 3 {
			return nil, errors.New(idBan + " est un identifiant invalide")
		}

		n, err := strconv.Atoi(idBanTokens[2])
		if err != nil {
			return nil, errors.New(idBan + " est un identifiant invalide")
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
		message := fmt.Sprintf(error_message, err)
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
	FROM (SELECT b.*,h.id_ign FROM ban_ign.ban b, 
	ban_ign.housenumber_id_ign h 
	WHERE LOWER(cle_interop) IN (?)
	AND h.id_ban_adresse=b.id_ban_adresse) as f;`

	var result string
	q, args, err := sqlx.In(query, idsBanArray)
	q = db.Rebind(q)
	err = db.QueryRow(q, args...).Scan(&result)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadGateway)
		message := fmt.Sprintf(error_message, err)
		fmt.Fprintf(w, message)
		defer db.Close()
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(result))

	defer db.Close()

}
