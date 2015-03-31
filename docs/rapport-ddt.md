# PComp - Rapport 2nd partie

Rapport du groupe DaJuNiBa (D. Galichet, J. Sagot, N. Cailloux, B. Fontaine)
sur le travail du groupe DDT.

## Rapport

Les instructions sont un peu légère : est-ce que le projet marche tout seul ?
Si non, de quoi ai-je besoin ? Comment on le lance ?

Bonne séparation des concepts : communication, mise en place, jeu. Les
explications sur le contenu des fichiers sont plutôt claires, sauf pour le
format du fichier de configuration. Par exemple celui-ci mentionne une ligne
qui doit « contenir la liste des paramètres demandés pour créer une partie »
sans préciser quels sont ces paramètres ni comment il faut les spécifier.

Le rapport indique qu’une des contraintes était de pouvoir tourner sur un
serveur Web alors qu’en fait c’était de pouvoir tourner sur un *client* Web,
donc un navigateur, comme le rappelle le README du projet.

## Parcours du code

Les noms de fichiers sont inconsistants (français, anglais, majuscules,
minuscules).

Concernant le code :

* indentation irrégulière mais globalement pas catastrophique.
  Certaines lignes font toutefois un peu mal aux yeux (`play.php`), pensez à
  faire des variables locales intermédiaires :

        while($st->response->status->status->status=="playing"){

        (…)

        && $resultat->response->observations[$num][1][$i]->x
            ==$resultat->response->observations[$num][0]->x + $dx

  Le second morceau de code a été coupé en deux pour préserver vos yeux, dans
  le code il tient sur une seule ligne.

* commentaires insuffisants, quelqu’un n’étant pas familier avec php ou curl a
  du mal à comprendre
* langue inconsistante dans les identificateurs, commentaires, et messages :
  `resultat`, `response`, `choix`, `inst`. Soit en français, soit en anglais.
  Vu que le langage a une base anglophone, l’anglais s’impose pour les
  identificateurs. Les commentaires et messages sont libres, mais doivent être
  consistants !
* multiples définitions de l’url dans `BIblio.php` (globalement, puis
  localement, tantôt avec la même valeur, tantôt avec une autre). Si on voulait
  jouer avec un autre serveur il faudrait du coup changer l’adresse à plusieurs
  endroits différents.
* utilisation abusive du `<br>` (au lieu d’un alias, si jamais on doit le
  changer)
* IA limitée (mais c’est le cas de tout le monde, et le rapport le mentionne)
* Le code PHP et le HTML sont complètement mélangés, il n’y a pas de séparation
  entre le programme et son interface.
* Certains morceaux de code sont morts : ni utilisés, ni utiles, comme la
  fonction `doHackcode` dans `play.php`.
* Certaines parties n’ont visiblement pas été testées, comme l’appel à `substr`
  dans `play.php`, qui ne respecte même pas le type des arguments (si on peut
  parler de type en PHP) et provoquera une erreur à l’exécution. Un autre
  exemple est cette ligne dans `BIblio.php` :

        curl_setopt($curl, CURLOPT_URL, $url.http_build_query($url));

  Ici, `$url` est le tableau associatif qui correspond aux paramètres de
  l’appel. Outre le nommage qui ne correspond pas à ce qu’est la variable, on
  voit ici qu’on concatène `$url` (donc à priori une chaîne) avec la version
  encodée des paramètres fournis par `$url` (donc un tableau associatif). Comme
  pour l’exemple précédent, celui-ci provoquera également une erreur à
  l’exécution.
* Les exceptions ne sont pas du tout utilisées, à la place les fonctions
  affichent une erreur et retournent `null`.
* Beaucoup de code est répété, notamment dans `BIblio.php`, il aurait pu être
  factorisé.
* La bibliothèque expose complètement l’API, donc si l’API proposée par le
  serveur distant change, tout le code devra être modifié, pas seulement le
  fichier `BIblio.php`.
* Certaines valeurs devraient être définies dans des constantes au nom
  évoquateur, comme le mystérieux `-15` de la ligne 71 de `start.php`, qui ne
  parle pas au lecteur qui n’a pas déjà lu le rapport.

## Exécution

Un avantage non-négligeable aurait été de pouvoir fournir la configuration en
ligne de commande, en plus de la possibilité d’utiliser un fichier.

### Correction

Comme on l’a vu précédemment, il n’est pas nécessaire de lancer le programme
pour voir que certaines parties ne marchent pas.
