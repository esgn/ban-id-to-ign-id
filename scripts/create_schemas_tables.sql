CREATE SCHEMA ban_ign;

CREATE TABLE ban_ign.ban (
    id_ban_position TEXT,
    id_ban_adresse TEXT,
    cle_interop TEXT,
    id_ban_group TEXT,
    id_fantoir TEXT,
    -- Should be int be errors are present in ban-ign datasets
    -- numero INT,
    numero TEXT,
    suffixe TEXT,
    nom_voie VARCHAR(80),
    code_postal INT,
    nom_commune TEXT,
    code_insee VARCHAR(10),
    nom_complementaire TEXT,
    pos_name TEXT,
    x float,
    y float,
    lon float,
    lat float,
    typ_loc VARCHAR(20),
    source TEXT,
    date_der_maj_group DATE,
    date_der_maj_hn DATE,
    date_der_maj_pos DATE
);

CREATE TABLE ban_ign.housenumber_id_ign (
    id_ban_adresse TEXT,
    id_ign TEXT
);