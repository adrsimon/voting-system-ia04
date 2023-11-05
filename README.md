# [IA04] TD3 - Prise de décision collective et vote

Projet réalisé dans le cadre d'un TD d'IA04, cours sur les systèmes multi-agents, enseigné en Go à l'UTC, par Quentin Fitte-Rey et Adrien Simon, étudiants en Intelligence Artificielle et Sciences des données.

### Sujet

Le but de ce TD est d'implémenter du bureau de vote. Les méthodes de vote implémentées sont disponibles dans le dossier comsoc. Sujet disponible [ici](https://github.com/adrsimon/voting-system-ia04/blob/main/sujet.md).

### Structure
- comsoc : contient les utilitaires à la construction des méthodes de vote
    - methods : contient les méthodes de votes
    - tests : contient des tests sur les méthodes de votes
- agt : contient les agents
- cmd : contient une implémentation du système avec une mise en pratique

### Installation

Méthode 1 :
- Cloner le projet : `git clone https://github.com/adrsimon/voting-system-ia04`
- Installer le serveur Rest : `go install cmd/launchBallot`
- Lancer le serveur de vote : `launchBallot`
Si vous souhaitez lancer des votes en local via un script go, utilisez la commande `go run cmd/launchBallot/launch.go`.

Méthode 2 :
- Créer un nouveau projet go : `mkdir projet && cd projet && go mod init projet`
- Récupérer le package dans le projet : `go get github.com/adrsimon/voting-system-ia04@latest`
- Utilisez le package comme n'importe quel autre package go.

### Méthodes de vote implémentées
Vote par majorité, vote de Borda, vote par approbation, vote simple transférable, vote de Copeland.

### API Map

| **Endpoint**   | **Description**                                                                                                                  | **Méthode** | **Entrée**                                                                             | **Sortie**                       |
|----------------|----------------------------------------------------------------------------------------------------------------------------------|-------------|----------------------------------------------------------------------------------------|----------------------------------|
| **methods**    | Permet de récupérer la liste des méthodes de vote implémentées.                                                                  | GET         | /                                                                                      | {"methods" : [...]}              |
| **new_ballot** | Permet la création d'une nouvelle session de vote basée sur la règle "rule". "rule" doit être une valeur retournée par /methods. | POST        | {"rule": string,"deadline":string,"voter-ids":[]string,"#alts":int,"tie-break":[]int,} | {"ballot-id": string,}           |
| **vote**       | Permet à l'agent "agent-id" de voter avec les préférences "prefs" sur le vote "ballot-id".                                       | POST        | {"agent-id": string,"ballot-id":string,"prefs":[]int,"options":[]int,}                 | /                                |
| **result**     | Permet de récupérer le résultat du vote "ballot-id".                                                                             | POST        | {"ballot-id": string,}                                                                 | {"winner":int,"ranking":[]int, } |

### Details d'implémentation

- new_ballot génère automatiquement les alternatives disponibles à partir du nombre d'alternatives envoyées, dans le but de minimiser les erreurs. Elles sont numérotées de `0` à `#alts-1`
- Les ID des ballots sont générés à partir d'un compteur global géré au niveau du serveur, permettant d'avoir des ballots uniques.
<br><br>
- Afin d'avoir une seule fonction de factory, efficace et fonctionnelle pour toutes les méthodes de vote, nous avons ajouté un argument optionel aux méthodes, permettant pour celles le nécessitant d'envoyer un tableau d'entier afin de préciser une fonction de tiebreak pour stv, ou un seuil d'approbation pour approval.
- Les personnes qui créent les votes, et votent, sont tous des voteragent.
<br><br>
- Les principaux types utilisés dans ce projet pour définir les agents sont les suivants :
  - ```go
    type ServerRest struct {
	       sync.Mutex
        id           string
        addr         string
        ballotAgents map[string]*ballotAgent
        count        int64
    }
    ```
    Le type ServerRest représente le serveur de vote. Il contient la liste des votes crées sur le serveur, ainsi qu'un compteur global permettant de générer des ID uniques pour les votes. Il contient aussi son adresse.
  - ```go
    type ballotAgent struct {
        sync.Mutex
        ballotID     string
        rulename     string
        rule         func(comsoc.Profile, ...int64) ([]comsoc.Alternative, error)
        deadline     time.Time
        voterID      []AgentID
        profile      comsoc.Profile
        alternatives []comsoc.Alternative
        tiebreak     []comsoc.Alternative
        thresholds   []int64
    }
    ```
    Le type ballotAgent représente une session de vote. Il contient entre autres la liste des votants pouvant encore voter, leurs préférences, la liste des candidats, la règle de vote, le tiebreak, et la deadline.
  - ```go
    type Agent struct {
        agentId AgentID
        prefs   []comsoc.Alternative
        options []int64
    }
    ```
    Le type Agent représente un votant. Il contient ses préférences, et un tableau représentant les éventuelles options nécessaires pour voter, comme un seuil d'approbation.