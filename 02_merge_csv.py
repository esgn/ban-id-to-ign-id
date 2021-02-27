
#!/usr/bin/env python
# coding: utf-8

import pandas as pd 
import os
import sys
import re
import shutil

ban_p = re.compile("^ban-([AB0-9]+).csv.gz")
hn_t = "housenumber-id-ign-{}.csv.gz"
out_dir = "ban-ign-id"

if os.path.exists(out_dir):
    shutil.rmtree(out_dir)
os.mkdir(out_dir)

for ban_f in os.listdir("ban"):
    if ban_p.match(ban_f):
        print("Traitement de " + str(ban_f))
        n = re.search(ban_p,ban_f).group(1)
        hn_f = hn_t.format(n)
        df_ban = pd.read_csv("ban/"+ban_f,sep=';',encoding='utf-8',dtype=object)
        df_hn = pd.read_csv("housenumber-id-ign/"+hn_f,sep=';',encoding='utf-8')
        result = df_ban.merge(df_hn,how='left',on='id_ban_adresse')
        # Vérifications rapides
        if result['ign'].isnull().values.any():
            print("Valeur manquante dans la colonne ign")
            print(str(len(result[result['ign'].isna()].index)) + " valeurs nulles")
        if len(df_ban.index) != len(result.index):
            print("Le fichier source et résultat ne sont pas de la même taille")
            print("source : " + str(len(df_ban.index)) + " lignes")
            print("source : " + str(len(result.index)) + " lignes")
            sys.exit(0)
        out_file = 'ban-ign-id-'+n+'.csv.gz'
        result.to_csv(os.path.join(out_dir,out_file),index=False,encoding='utf-8',sep=';',compression='gzip')
