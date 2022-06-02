# URL Shortener
Highly Sacalable, distributed URL shortener fully writen in Go (Golang). It uses MongoDB, etcd as NoSQL database and Redis as cache database.

# About app
It uses Base62 encoding to create unique strings for each new url. Etcd is used to provide consistency among every running server and solve concurrency issues. It holds lastCounterEnd which is used to create instances of counter for multiple instances of app. Documents containing URL Models that consist of Long URL, Short URL, Created_At and Expired_AT are held in MongoDB and in Redis for cache. 
When POST request is made, app validates if long url from body is valid with http standars, current counter value in app is incremented by 1 and passed to Base62 encoder,which generates new unique Short URL string. Then new URL model is created containing Long URL from request body, Short URL that was made, Created_At and Expired_At. Next that model is saved to MongoDB and Redis cache. Client gets that URL model in response. 
WHEN GET request is made on the short URL, app first searches Redis cache for matching result if it cant find it, then it checks MongoDB. If that record is present then client is being redirected to correspodning Long URL, if not the gets response that there is no matching result.
![](https://github.com/pgorczyca/url-shortener/blob/main/architecture.jpg)
# Installation
```
$ git clone https://github.com/pgorczyca/url-shortener
$ cd url-shortener
$ docker-compose up -d
```
# Quick start

## curl requests
Create URL
```
$ curl -X POST -H "Content-Type: application/json" \
    -d '{"long": "https://github.com/pgorczyca/url-shortener"}' \
    http://localhost:8081/url
```
Get URL

```
$ curl -X GET --location 'localhost:8081/1'

```

TODO
refactor range provider initialize
context
tests