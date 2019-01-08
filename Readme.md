# elastic docker url https://github.com/maxyermayank/docker-compose-elasticsearch-kibana

curl -H "Accept: application/json" -X POST http://localhost:9200/examples/_search -d '{
      "author": "admin"
    }'

curl -i -H "Accept: application/json"  -X POST http://localhost:9200/examples/_search -d '{"author": "admin"}'

curl http://localhost:9200/examples/_search?q=author:admin&pretty=true

curl http://localhost:9200/users/_search?q=name:udin&pretty=true

curl http://127.0.0.1:9200/examples/_search/?size=1000&pretty=true

https://stackoverflow.com/questions/8829468/elasticsearch-query-to-return-all-records

https://github.com/olivere/elastic/blob/release-branch.v6/delete_test.go

https://mindmajix.com/elasticsearch/curl-syntax-with-examples

https://stackoverflow.com/questions/32052507/representing-a-kibana-query-in-a-rest-curl-form

https://dzone.com/articles/23-useful-elasticsearch-example-queries

# Note : jika pake elastic memory di preference docker mesti minimal 3 Gb, klo gak gak bakalan jalan

use rabbit MQ
# https://github.com/rabbitmq/rabbitmq-tutorials/blob/master/go/receive.go

SLATE