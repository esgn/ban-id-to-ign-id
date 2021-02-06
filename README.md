# BAN ID TO IGN ID

Début d'API expédiée et non finalisée permettant d'obtenir les informations présentes dans les [fichiers CSV de l'export de la BAN IGN](https://adresse.data.gouv.fr/data/ban/export-api-gestion/latest/) à partir d'une cle_interop retournée par la BAN Etalab.

Technologies utilisées :

* Docker / Docker Compose
* Go
* PostgreSQL/PostGIS

## Récupération des données

Un script permettant de télécharger en parallèle les CSV nécessaires doit d'abord être lancé.

```
	pip install -r requirements.txt
	python get_csv.py
```

## Import des données CSV en base de données

Commencer par démarrer les containers

```
	docker-compose up -d
```

La commande suivante permet ensuite de lancer l'import des données. Cet import peut prendre un certain temps. La configuration de la base de données n'a pas été optimisée pour minimiser le temps d'import pour le moment.

```
	docker exec -ti adr-postgis /bin/bash /tmp/scripts/import_ban_data.sh
```

## Utilisation des services

La seule opération de l'API est à ce jour la suivante.

```
	Methode : GET
	Port : 80 (par défaut, peut être modifié au niveau du docker-compose)
	URL : /position/cle_interop (ex : /position/73008_1700_00008)
	Réponse : Ensemble des positions présentes dans les fichiers d'export CSV de la BAN IGN correspondant à la cle_interop utilisée
	Format de la réponse : GeoJSON
	Format des erreurs : A terminer. Message texte et code HTTP pour le moment.
```

Une instance d'[Adminer](https://github.com/vrana/adminer) est également disponible par défaut sur le port 8080 pour gérer la base de données si nécessaire.

## Mise à jour de la base

Pour mettre à jour la base, la seule solution possible, pour le moment, est de la reconstruire.

Commencer par éteindre et supprimer le contenu de la base de donnée via la commande suivante :

```
	docker-compose down -v
```

Répéter les opérations ci-dessus pour reconstruire la base de données.
