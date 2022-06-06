# URL Shortener
Highly Sacalable, distributed URL shortener fully written in Go (Golang). It uses MongoDB, etcd as NoSQL database and Redis as cache database.

# Architecture
![](https://github.com/pgorczyca/url-shortener/blob/main/architecture.jpg)
###
ShortGenerator holds currently used range, backup, and current counter value. RangeProvider provides new instances of counterRange based on lastCounterEnd from etcd. Current counter value in ShortGenerator is passed to Base62 encoder, which generates new unique Short url string. Long and short urls are saved to MongoDB and Redis cache. In case of pulling up existing urls from databases, app prioritizes Redis cache over MongoDB to be the most efficient.
# Installation
```
$ git clone https://github.com/pgorczyca/url-shortener
$ cd url-shortener
$ docker-compose up -d
```
# Quick start

Create URL
```
$ curl -X POST -H "Content-Type: application/json" \
    -d '{"long": "https://github.com/pgorczyca/url-shortener"}' \
    http://localhost:8081/url
```

Get URL
```
$ curl --silent -v localhost:8081/1
```

# References
https://youtu.be/JQDHz72OA3c - inspired by
