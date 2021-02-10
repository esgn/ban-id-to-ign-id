package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	http.HandleFunc("/position/", BanToIgn)
	http.ListenAndServe(":8080", nil)
}

const (
	host     = "adr-postgis"
	port     = 5432
	user     = "adr"
	password = "adr"
	dbname   = "adr"
)

func BanToIgn(w http.ResponseWriter, r *http.Request) {

	idban := strings.TrimPrefix(r.URL.Path, "/position/")

	// Conversion cle_interop BAL1.2 vers cle_interop IGN
	idbanTokens := strings.Split(idban, "_")
	if len(idbanTokens) < 3 {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Clé interop invalide")
		return
	}

	n, err := strconv.Atoi(idbanTokens[2])
	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Le troisième élément de la clé interop ne doit contenir que des chiffres")
		return
	}

	db := OpenConnection()

	idbanTokens[2] = strconv.Itoa(n)
	idban = strings.Join(idbanTokens, "_")

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
	WHERE cle_interop=$1 
	AND h.id_ban_adresse=b.id_ban_adresse) as f;`

	var result string
	err = db.QueryRow(query, idban).Scan(&result)

	if err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Erreur lors de l'interrogation de la base de données")
		defer db.Close()
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(result))

	defer db.Close()

}

func OpenConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
