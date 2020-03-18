TODO :

- Donne une liste d'entrée en fonction de key/value
- Donne toute la base
- Donne moi X entrées avec "SentToUser" == false + passe SentToUser à true
- Delete key/value

# Micro-database
## MongoGo API
The different structs define how the JSON files must be parsed :
```go
/* Parse PiFF files */
type Meta struct {
	Type string
	URL  string
}

type Location struct {
	Type    string
	Polygon [][2]int
	Id      string
}

type Data struct {
	Type       string
	LocationId string
	Value      string
	Id         string
}

type PiFFStruct struct {
	Meta     Meta
	Location []Location
	Data     []Data
	Children []int
	Parent   int
}
```
```go
/* Our working structure */
type Picture struct {
	// Id in db
	Id primitive.ObjectID `bson:"_id" json:"Id"`
	// Piff
	PiFF PiFFStruct `json:"PiFF"`
	// Url fileserver
	Url string `json:"Url"`
	// Flags
	Annotated  bool `json:"Annotated"`
	Corrected  bool `json:"Corrected"`
	SentToReco bool `json:"SentToReco"`
	Unreadable bool `json:"Unreadable"`
}

```
Be careful, other people not working with the mongo module should use []byte instead of ObjectID
```go
/* Update flags */
type Modification struct {
	Id int
	Flag string
	Value bool
}
```
```go
/* Update PiFF value and some flags */
type Annotation struct {
	Id int
	Value string
}
```
MongoDB is case sensitive so the fields MUST begin with a lower case

In order to connect to the database, the API uses `const URI = "mongodb://mongo:27017/"` (see also the [Standard Connection String Format](#standard-connection-string-format))

## Rest API
The rest API transform rest request into mongoGo API method call. 

See the API documentation in API Blueprint format [here](api.md)

## Conventions
### Standard Connection String Format
```
mongodb://[username:password@]host1[:port1][,...hostN[:portN]][/[database][?options]]
```

# Deployment

## Dockerfile
Micro-database build from a golang image (the latest, even if we should decide on which version to use) where it download dependancies for protocole file, go.sum/go.mod (which manage almost all dependancies),...
And then build the project. This image weights more than 1Go so we build another one from alpine, using our golang image as builder. We just need to copy the exec file and to launch it (with CMD).

The exposed port for micro-database is 8080.
## Creating and running containers
To build a new docker image for this microservice, you need to execute this line, in the same directory than the Dockerfile:
```shell script
docker build -t micro-database .
```

To access mongodb container from outside you need to expose its 27017 port. When you run the docker image you should do :
```shell script
docker run --name mongo -d -p 27017:27017 mongo:3.6.3
```
mongo:3.6.3 being the latest version of mongo when I did this. Note : you might need to pull mongodb image before running this.

Now you need to link your micro-database container with mongo. 
```shell script
docker run --name micro-database -d --link mongo:mongo micro-database
```
You can test the request at `http://172.17.0.3:8080/` (I got the IPAddress by doing `docker container inspect micro-database`).

## deployment.yml
Currently two deployments are supported, the isoproduction environment (nicknamed "dev") and the production environment, they set up environment variables that the service uses to determine which database to join.  
When you run locally on your machine, as no environment variable is set, it will join you local mongodb daemon.  
If the need arises for a local test on the existing databases, see the documentation of your IDE to know how to set environment variables.  
In Golang, these variables are setin the **Edit configurations** menu.

In the future the number of replicas will be managed by Kubernetes, for now it is set to 3

# Commits

The title of a commit must follow this pattern : \<type>(\<scope>): \<subject>

## Type
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

## Scope
Your commits name should also precise which part of the project they concern.
You can do so by naming them using the following scopes :
- General
- RestAPI
- Database
