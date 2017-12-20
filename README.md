# Celtra Data Engineer Challenge solution

This is a solution for Celtra Data Engineer Challenge, writen in Go.

## Building procedure

1. Install Go - https://golang.org/doc/install.
2. Set environment variable `GOOGLE_APPLICATION_CREDENTIALS` to a JSON
previously stored to enable system-wide secure access to Google
Cloud Platform:
```
$ export GOOGLE_APPLICATION_CREDENTIALS=<path_to_google_app_credentials>.json
```
3. Execute:
```
$ go get github.com/matejl/challenge
$ go build github.com/matejl/challenge
$ ./challenge [-port=<port>]
```

App should be up and running now.

## Demo queries

In following queries, `<root_url>` is usually `localhost:8080` (depending on configuration).

- HTTP request to get number of `impressions`, `interactions`, and `swipes` for each ad in a specific campaign:
```
<root_url>/campaign?id=4&dimensions[]=adId&dimensions[]=adName&metrics[]=impressions&metrics[]=swipes&metrics[]=pinches&metrics[]=touches
```

- HTTP request to get number of `uniqueUsers` and `impressions` for each ad in the last week
```
<root_url>/campaign?dateRange=lastWeek&dimensions[]=adId&dimensions[]=adName&metrics[]=uniqueUsers&metrics[]=impressions
```

## Test data

Test data is already stored in preconfigured database but you can still
fill it with more dummy values using `testdata` command:

```
$ cd $GOPATH/src/github.com/matejl/challenge/testdata
... configure values (constants) on top of main.go
$ go build
$ ./testdata
```

## Scaling

Scaling can be performed using deployment of multiple instanes of API
and using a load balancer (for example, nginx).

Sample nginx configuration inside one of `/etc/nginx/sites-available/` files:
```
http {
    upstream stat_api {
        server server1.domain.com:<port>;
        server server2.domain.com:<port>;
        ...
    }

    server {
        listen 80;

        location / {
            proxy_pass http://stat_api;
        }
    }
}
```