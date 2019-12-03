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


##Rest package