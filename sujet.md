---
tags: TD, TP, Go, comsoc
---

# [IA04] TD3 - Prise de décision collective et vote

|  Information |   Valeur              |
| :------------ | :------------- |
| **Auteurs** | Sylvain Lagrue ([sylvain.lagrue@utc.fr](mailto:sylvain.lagrue@utc.fr))|
| | Khaled Belahcene (khaled.belahcene@utc.fr) |
| **Licence** | Creative Common [CC BY-SA 3.0](https://creativecommons.org/licenses/by-sa/3.0) |
| **Version document** | 1.4.0 |

L'objectif de ce TD est de créer un bureau de vote à distance pour des agents. Le travail demandé sera :

1. Implémenter des procédures de vote usuelles (majorité, majorité à 2 tours, Borda, Copeland) ainsi que recherche de gagnant de Condorcet
2. Adapter les agents des méthodes d'appariement en agents votants (`Voters`)

Lors des prochaines séances, ce travail sera complété afin de réaliser un système multi-agent de vote accessible en ligne. 

**L'ensemble donnera lieu à un rendu en binômes.**

## I. Préliminaires

On utilisera l'arborescence suivante :

- `ia04`
    - `cmd` (qui contiendra un répertoire par exécutable dont vous aurez besoin)
    - `comsoc` (qui contiendra les implémentations des méthodes de vote)
    - `agt` (qui contiendra l'ensemble des agents)
        - `ballotagent` (l'agent qui gère la procédure de vote)
        - `voteragent` (les *voters*)

Ne pas oublier d'initialiser correctement le module au niveau de `ia04` à l'aide de la commande `go mod`.

Ne pas hésiter à décomposer la suite en fonctions.

## II. Implementation des méthodes de vote  : créaction du package `comsoc`

On placera l'ensemble des fonctions demandées dans cette section dans un package nommé `comsoc`. 

### Types de base

On utilisera les types de bases suivants.

```golang
type Alternative int
type Profile [][]Alternative
type Count map[Alternative]int
```

- Les alternatives seront représentées par des entiers.
- Les profils de préférences sont telles que si `profile` est un profil, `profile[12]` représentera les préférences du votant `12`. Les alternatives sont classée de la préférée à la moins préférée :  `profile[12][0]` represente l'alternative préférée du votant `12`.
- Enfin, les méthodes de vote renvoient un décompte sous forme d'une *map* qui associe à chaque alternative un entier : plus cet entier est élevé, plus l'alternative *a de points* et plus elle est préférée pour le groupe compte tenu de la méthode considérée.

On pourra créer quelques fonctions utilitaires :

```golang
// renvoie l'indice ou se trouve alt dans prefs
func rank(alt Alternative, prefs []Alternative) int 

// renvoie vrai ssi alt1 est préférée à alt2 
func isPref(alt1, alt2 Alternative, prefs []Alternative) bool 

// renvoie les meilleures alternatives pour un décomtpe donné
func maxCount(count Count) (bestAlts []Alternative)

// vérifie les préférences d'un agent, par ex. qu'ils sont tous complets et que chaque alternative n'apparaît qu'une seule fois
func checkProfile(prefs []Alternative, alts []Alternative) error 

// vérifie le profil donné, par ex. qu'ils sont tous complets et que chaque alternative de alts apparaît exactement une fois par préférences
func checkProfileAlternative(prefs Profile, alts []Alternative) error 
```
 
### Procédures de vote

On distingue les fonctions de bien-être social (*social welfare function*, *SWF*), qui retournent un décompte à partir d'un profil, des fonctions de choix social (choix social, *social choice function*, SCF) qui renvoient quant à elles uniquement les alternatives préférées.

```golang
func SWF(p Profile) (count Count, err error)
func SCF(p Profile) (bestAlts []Alternative, err error)
```

Chaque fonction peuvent renvoyer une erreur, par exemple dans le cas de profils malformés.

On utilisera **un fichier source différent** par méthode vote.

1. Commencer par la majorité simple.

```golang
func MajoritySWF(p Profile) (count Count, err error)
func MajoritySCF(p Profile) (bestAlts []Alternative, err error)
```

2. Puis la méthode de Borda.

```golang 
func BordaSWF(p Profile) (count Count, err error)
func BordaSCF(p Profile) (bestAlts []Alternative, err error) {
```

3. Et l'on continuera avec le vote par approbation (attention, dans ce cas il faut ajouter un nombre représentant le seuil à partir duquel les alternatives ne sont plus approuvées).

```golang
func ApprovalSWF(p Profile, thresholds []int) (count Count, err error)
func ApprovalSCF(p Profile, thresholds []int) (bestAlts []Alternative, err error)
```

4. Fonctions de départage

En cas d'égalité, les SCF doivent renvoyer un seul élement et les SWF doivent renvoyer un ordre total sans égalité. On utilisera pour cela des fonctions de *tie-break* qui étant donné un ensemble d'alternatives, renvoie la meilleure. Elles respectent la signature suivante (une erreur pouvant se produire si le slice d'alternatives est vide) :

```golang
func TieBreak([]Alternative) (Alternative, error)
```

Il peut être intéressant d'utiliser des *factories* de fonctions de *tie-break* permettant de créer une telle fonction à partir d'un ordre strict (représenté par un slice d'alternatives) telles que :

```golang
func TieBreakFactory(orderedAlts []Alternative) (func ([]Alternative) (Alternative, error))
```

Enfin, à partir d'une SWF (resp. SCF) et d'une fonction de *tie-break* il est possible de construire une autre *factory* permettant d'obtenir des SWF (resp. SCF) n'admettant aucun ex aequo.

```golang
func SWFFactory(func swf(p Profile) (Count, error), func ([]Alternative) (Alternative, error)) (func(Profile) ([]Alternative, error))

func SCFFactory(func scf(p Profile) ([]Alternative, error), func ([]Alternative) (Alternative, error)) (func(Profile) (Alternative, error))
```

5. Définir une fonction permettant de trouver le gagnant de Condorcet 

Qui renvoie un *slice* éventuellement vide ou ne contenant qu'un seul élément.

```golang
func CondorcetWinner(p Profile) (bestAlts []Alternative, err error)
```

#### Si vous êtes en avance...

6. Implémenter la procédure de Copeland.

```golang 
func CopelandSWF(p Profile) (Count, error)
func CopelandSCF(p Profile) (bestAlts []Alternative, err error) {
```

##### Rapel : Règle de Copeland

- Le meilleur candidat est celui qui bat le plus d’autres candidats
- On associe à chaque candidat $a$ le score suivant : pour chaque autre candidat $b\ne a$ $+1$ si une majorité préfère $a$ à $b$, $-1$ si une majorité préfère $b$ à $a$ et $0$ sinon
- Le candidat élu est celui qui a le plus haut score de Copeland

7. Implémenter la procédure du Vote Simple Transférable (STV).

```golang 
func STV_SWF(p Profile) (Count, error)
func STV_SCF(p Profile) (bestAlts []Alternative, err error) {
```

##### Rappel : Vote Simple Transférable (Single Transferable Vote (STV) 

- Chaque individu indique donne son ordre de préférence $>_i$
- Pour $n$ candidats, on fait $n − 1$ tours (à moins d’avoir avant une majorité stricte pour un candidat)
- On suppose qu’à chaque tour chaque individu “vote” pour son candidat préféré (parmi ceux encore en course)
- À chaque tour on élimine le plus mauvais candidat (celui qui a le moins de voix)

## II. Modélisation multi-agent : package `agt`

### Types, interfaces et fonctions de base

Un agent sera représenté par la structure suivante :

```golang
type Agent struct {
	ID    AgentID
	Name  string
	Prefs []Alternative
}
```

Avec :

```golang
type Alternative int
```

Un agent doit implémenter les méthodes de l'interface suivante :

```golang
type AgentI interface {
	Equal(ag AgentI) bool
	DeepEqual(ag AgentI) bool
	Clone() AgentI
	String() string
	Prefers(a Alternative, b Alternative) bool
	Start()
}
```

### Question

Proposer une modélisation multi-agent simple permettant d'implémenter une prise de décision pour un ensemble d'agents. En faire une première implémentation naïve.

