##Mongo package

```
type Picture struct {
	Id int
	Type string
	Value  string
	Zone string
	Children string
	Parent string
	Url string
	Annotated bool
	Corrected bool
	SentToReco bool
	SentToUser bool
	Unreadable bool
}
```

```
type Modification struct {
	Id int
	Flag string
	Value bool
}
```

```
type Annotation struct {
	Id int
	Value string
}
```

TODO :

- Donne une liste d'entrée en fonction de key/value
- Donne toute la base
- Donne moi X entrées avec "SentToUser" == false + passe SentToUser à true
- Delete key/value
- Delete tout


## Commits

The title of a commit must follow this pattern : \<type>(\<scope>): \<subject>

### Type
Commits must specify their type among the following :
- build: Changes that affect the build system or external dependencies
- docs: Documentation only changes
- feat: A new feature
- fix: A bug fix
- perf: A code change that improves performance
- refactor: Modification of the code without adding features nor bugs (rename, white-space, ...)
- style: CSS or layout modifications or debug
- test: Adding missing tests or correcting existing tests
- ci: Changes to our CI configuration

### Scope
Your commits name should also precise which part of the project they concern.
You can do so by naming them using the following scopes :
- General
- RestAPI
- Database
