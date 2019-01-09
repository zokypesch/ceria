# Lazy developer use this !!!

# Beta Version 1.0.0

full example <a href="https://github.com/zokypesch/example-ceria">here</a> with docker-compose for running environment

![Screenshot](ceria_diagram.png)

How long did it take to finish this from scratch ???

![Screenshot](ERD.png)

Don’t sad baby …..
If you are using ceria you cant do it for 5 minutes

```
SAY no to killing my time …
Speed development is a key success for your bussiness
```

Why choose us
because ceria is `LIGHT & warm`

# Package
```
- Ceria Repository (Handler like CRUD API GATEWAY)
- Ceria Core (DB Conn, Redis, RabbitMQ, Elastic)
- Ceria Util (Converting, GetValue)
- Ceria Helper (Test Helper, Wrapper)
- Docker (Environment)
- Makefile (Short hands for running your apps)
- Example (Example how to use)
```

# HOW ABOUT PERFORMANCE, QUALITY & Customize ??
```
- Ceria using trusted library & high performance
- Ceria have UNIT & INTEGRATION TEST Coverage more than 90 %
- Ceria easly customize because it’s transparant
- Ceria have documentation in code, if your using Go-Lint in visual studio, u can see description method, struct, etc.
- Ceria core is friendly you can modified it easly
- Ceria using TDD (Test design driven)
```

# Trusted Library use in Ceria Workspace
```
Gorm (ORM)
Gin (HTTP Framewrok)
GIN-JWT (JWT Auth)
Redisstore (Redis)
Assert (Unit Testing)
Go-playground (Validatior)
Gorilla-Sessions (Manage Session)
GO-Mocket(Mocking SQL for Gorm)
Ampq (Rabbit MQ)

VIPER (Read Config File)
Dep (Depedency tool vendor)
Elastic (Elastic Library)
Ceria Core
Ceria Repository(Handler Management)
Ceria Util (Utility)
Ceria Helper (Http Helper)
Ceria Wrapper
```

# How to it works ?
```
clone or deownload example in https://github.com/zokypesch/example-ceria
type "make help" in current folder example-ceria
make init mode=development
make install_docker
make install
make rundb

see example in folder example and "go run main.go"
for full example type "cd examples/wrapper && go run main.go" 
```

# Open your Postman or using Curl
```
curl -H "Accept: application/json" -X GET http://localhost:9090/articles?page=1&limit=30

curl -H "Accept: application/json" -X DELETE http://localhost:9090/articles/1

curl -i -H "Accept: application/json"  -X POST http://localhost:9090/articles -d \
'{"title": "hello welcome to ceria", "tag": "#Ceriaworkspace", "body": "lorem ipsum lorem ipsum", "Author": {"fullname": "Ceria Lover"}}'

curl -i -H "Accept: application/json"  -X POST http://localhost:9090/articles -d \
'{"title": "hello welcome to ceria", "tag": "#Ceriaworkspace", "body": "lorem ipsum lorem ipsum", "author_id": 1}'

curl -i -H "Accept: application/json"  -X PUT http://localhost:9090/articles/1 -d \
'{"data": {"title": "iam in ceria"}}'

curl -i -H "Accept: application/json"  -X POST http://localhost:9090/articles/bulkcreate -d \
'[{"title": "hello welcome to ceria", "tag": "#Ceriaworkspace", "body": "lorem ipsum lorem ipsum", "author_id": 1}]'

curl -i -H "Accept: application/json"  -X POST http://localhost:9090/articles/find -d \
'{"condition": {"author": "admin"}}'

Login to get JWT
curl -i -H "Accept: application/json"  -X POST http://localhost:9090/login -d \
'{"username": "admin", "password": "admin"}'

get token and accept to barrier
curl -i -H "Authorization: Bearer <your token>"  -X GET http://localhost:9090/auth/comments 

check your elastic search
curl -i -H "Accept: application/json"  -X POST http://localhost:9200/examples/_search -d '{"author": "admin"}'

for query using expression like
curl -H "Accept: application/json" -X GET http://localhost:9090/articles?page=1&limit=30&where=title:welcome:EQUAL|author_id:1:LIKE
```
# Happy coding guys :)