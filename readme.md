
# Framework Multi-Agents

  

## Description du projet

  

Ce projet implémente un framework simple pour la création et la gestion de systèmes multi-agents.

Il permet la communication synchrone et asynchrone entre agents à travers des conteneurs.

Chaque agent peut posséder plusieurs comportements et communiquer avec d'autres agents pour coordonner des actions

ou partager des informations. L'idée est de fournir un outil permettant à l'utilisateur de faire abstraction des

détails de communication et de gestion des agents pour se concentrer sur la modélisation de systèmes multi-agents.

  

## Installation

  

1. Assurez-vous que Go est installé sur votre machine.

2. Clonez le dépôt du projet dans votre espace de travail Go.

3. Naviguez dans le dossier du projet cloné.

4. Utilisez la commande `go build` pour compiler le projet.

  

## Utilisation

  

Cette version du framework nécéssite de lancer 2 instances de l'application:

un conteneur principal et un ou plusieurs conteneurs secondaires.

Les différents comportements sont implémentés dans le main. (Version MainContainer et version Conteneur).

Les différents behaviours sont aussi implémenté dans le fichier main.

  

### Démarrer le conteneur principal

  

Exécutez la commande suivante pour démarrer le conteneur principal sur le port 8080 (vous pouvez choisir un autre port si nécessaire) :

  

```bash

./executable  -main=true  -port=8080

```

  

### Démarrer un conteneur secondaire

  

Exécutez la commande suivante pour démarrer un conteneur secondaire qui se connectera au conteneur principal.

Remplacez `<port>` par le port sur lequel vous souhaitez exécuter le conteneur secondaire :

  

```bash

./executable  -main=false  -port=<port>

```

  

Assurez-vous que le port du conteneur principal est correctement configuré pour permettre la communication entre les conteneurs.

  

## Démonstration

Le projet comprend plusieurs fichiers de démonstration pour illustrer les fonctionnalités du framework :

  

- DemoAsyncDistant.go : Démonstration de la communication asynchrone entre agents sur des conteneurs distants.

Pour exécuter cette démo, lancez deux instances :

go run DemoAsyncDistant.go -main=true -port=8080

go run DemoAsyncDistant.go -main=false -port=8081

DemoSyncDistant.go : Démonstration de la communication synchrone entre agents sur des conteneurs distants.

Pour exécuter cette démo, lancez deux instances :

go run DemoSyncDistant.go -main=true -port=8080

go run DemoSyncDistant.go -main=false -port=8081

DemoAsyncLocal.go : Démonstration de la communication asynchrone entre agents sur un seul conteneur.

Pour exécuter cette démo, lancez une seule instance :

go run DemoAsyncLocal.go

DemoSyncLocal.go : Démonstration de la communication synchrone entre agents sur un seul conteneur.

Pour exécuter cette démo, lancez une seule instance :

go run DemoSyncLocal.go

Ces fichiers de démonstration vous permettent de voir le framework en action et de comprendre comment les agents communiquent de manière synchrone et asynchrone, à la fois localement et sur des conteneurs distants.

  
  

## Fonctionnalités clés

  

-  **Gestion des agents :** Création et gestion d'agents avec des comportements personnalisables.

-  **Communication asynchrone :** Les agents peuvent envoyer des messages asynchrones à d'autres agents sans attendre une réponse.

-  **Communication synchrone :** Permet aux agents d'engager des dialogues synchrones, nécessitant une réponse avant de continuer.

-  **Conteneurs :** Les agents sont organisés dans des conteneurs, facilitant leur gestion et leur communication.

  

Ce framework sert de base pour développer des applications complexes utilisant des systèmes multi-agents pour la simulation, l'automatisation de tâches, etc.# Framework Multi-Agents

  

## Description du projet

  

Ce projet implémente un framework simple pour la création et la gestion de systèmes multi-agents.

Il permet la communication synchrone et asynchrone entre agents à travers des conteneurs.

Chaque agent peut posséder plusieurs comportements et communiquer avec d'autres agents pour coordonner des actions

ou partager des informations. L'idée est de fournir un outil permettant à l'utilisateur de faire abstraction des

détails de communication et de gestion des agents pour se concentrer sur la modélisation de systèmes multi-agents.

  

## Installation

  

1. Assurez-vous que Go est installé sur votre machine.

2. Clonez le dépôt du projet dans votre espace de travail Go.

3. Naviguez dans le dossier du projet cloné.

4. Utilisez la commande `go build` pour compiler le projet.

  

## Utilisation

  

Cette version du framework nécéssite de lancer 2 instances de l'application:

un conteneur principal et un ou plusieurs conteneurs secondaires.

Les différents comportements sont implémentés dans le main. (Version MainContainer et version Conteneur).

Les différents behaviours sont aussi implémenté dans le fichier main.

  

### Démarrer le conteneur principal

  

Exécutez la commande suivante pour démarrer le conteneur principal sur le port 8080 (vous pouvez choisir un autre port si nécessaire) :

  

```bash

./executable  -main=true  -port=8080

```

  

### Démarrer un conteneur secondaire

  

Exécutez la commande suivante pour démarrer un conteneur secondaire qui se connectera au conteneur principal.

Remplacez `<port>` par le port sur lequel vous souhaitez exécuter le conteneur secondaire :

  

```bash

./executable  -main=false  -port=<port>

```

  

Assurez-vous que le port du conteneur principal est correctement configuré pour permettre la communication entre les conteneurs.

  

## Démonstration

Le projet comprend plusieurs fichiers de démonstration pour illustrer les fonctionnalités du framework :

  

- DemoAsyncDistant.go : Démonstration de la communication asynchrone entre agents sur des conteneurs distants.

	- Pour exécuter cette démo, lancez deux instances :

		- `go run DemoAsyncDistant.go -main=true -port=8080`

		- `go run DemoAsyncDistant.go -main=false -port=8081`

- DemoSyncDistant.go : Démonstration de la communication synchrone entre agents sur des conteneurs distants.

	- Pour exécuter cette démo, lancez deux instances :

		- `go run DemoSyncDistant.go -main=true -port=8080`
		- `go run DemoSyncDistant.go -main=false -port=8081`

- DemoAsyncLocal.go : Démonstration de la communication asynchrone entre agents sur un seul conteneur.

	- Pour exécuter cette démo, lancez une seule instance :

		- `go run DemoAsyncLocal.go`

- DemoSyncLocal.go : Démonstration de la communication synchrone entre agents sur un seul conteneur.

	- Pour exécuter cette démo, lancez une seule instance :

		- `go run DemoSyncLocal.go`

Ces fichiers de démonstration vous permettent de voir le framework en action et de comprendre comment les agents communiquent de manière synchrone et asynchrone, à la fois localement et sur des conteneurs distants.

  
  

## Fonctionnalités clés

  

-  **Gestion des agents :** Création et gestion d'agents avec des comportements personnalisables.

-  **Communication asynchrone :** Les agents peuvent envoyer des messages asynchrones à d'autres agents sans attendre une réponse.

-  **Communication synchrone :** Permet aux agents d'engager des dialogues synchrones, nécessitant une réponse avant de continuer.

-  **Conteneurs :** Les agents sont organisés dans des conteneurs, facilitant leur gestion et leur communication.

  

Ce framework sert de base pour développer des applications complexes utilisant des systèmes multi-agents pour la simulation, l'automatisation de tâches, etc.