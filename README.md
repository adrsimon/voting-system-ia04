# [IA04] TD3 - Prise de décision collective et vote

Projet réalisé dans le cadre d'un TD d'IA04, cours sur les systèmes multi-agents, enseigné en Go à l'UTC, par Quentin Fitte-Rey et Adrien Simon, étudiants en Intelligence Artificielle et Sciences des données.

### Sujet

Le but de ce TD est d'implémenter du bureau de vote. Les méthodes de vote implémentées sont disponibles dans le dossier comsoc. Sujet disponible [ici](https://github.com/adrsimon/voting-system-ia04/blob/main/sujet.md).

### Structure
- comsoc : contient les méthodes de votes
- agt : contient les agents
- cmd : contient une implémentation du système avec une mise en pratique

### Installation
- Cloner le projet : `git clone https://github.com/adrsimon/voting-system-ia04`
- Installer le projet : `go install cmd/launchBallot`
- Lancer le serveur de vote : `launchBallot`

Si vous souhaitez lancer des votes en local via un script go, utilisez la commande `go run cmd/launchBallot/launch.go`.


### Méthodes de vote implémentées
Vote par majorité, vote de Borda, vote par approbation, vote simple transférable, vote de Copeland.

### Details d'implémentation

.....

### API Map

| **Endpoint**   | **Description**                                                                                                                  | **Méthode** | **Entrée**                                                                             | **Sortie**                       |
|----------------|----------------------------------------------------------------------------------------------------------------------------------|-------------|----------------------------------------------------------------------------------------|----------------------------------|
| **methods**    | Permet de récupérer la liste des méthodes de vote implémentées.                                                                  | GET         | /                                                                                      | {"methods" : [...]}              |
| **new_ballot** | Permet la création d'une nouvelle session de vote basée sur la règle "rule". "rule" doit être une valeur retournée par /methods. | POST        | {"rule": string,"deadline":string,"voter-ids":[string],"#alts":int,"tie-break":[int],} | {"ballot-id": string,}           |
| **vote**       | Permet à l'agent "agent-id" de voter avec les préférences "prefs" sur le vote "ballot-id".                                       | POST        | {"agent-id": string,"ballot-id":string,"prefs":[int],"options":[int],}                 | /                                |
| **result**     | Permet de récupérer le résultat du vote "ballot-id".                                                                             | POST        | {"ballot-id": string,}                                                                 | {"winner":int,"ranking":[int], } |

