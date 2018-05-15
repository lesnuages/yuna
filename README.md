# Yuna

## Description

Simple TLS covert channel implementation.

Works on both Linux and Windows.

Inspired by [this gist from Casey Smith](https://twitter.com/Oddvarmoe/status/996147947975962624).

## Quickstart

Start the server:

```
go run server.go -command "whoami"
```

Then launch the client:

```
go run client.go
```

## Example output

```
$ go run server.go -command "/bin/ls -laht"
2018/05/15 11:14:28 Command: /bin/ls -laht
2018/05/15 11:14:28 Listenning on 0.0.0.0:4444

total 24K
drwxrwxr-x. 3 rkva rkva 4,0K 15 mai   11:14 .
-rw-rw-r--. 1 rkva rkva  329 15 mai   11:14 README.md
-rw-rw-r--. 1 rkva rkva 1,2K 15 mai   11:11 server.go
-rw-rw-r--. 1 rkva rkva 1,6K 15 mai   11:07 client.go
drwxrwxr-x. 2 rkva rkva 4,0K  9 avril 21:28 network
drwxrwxr-x. 4 rkva rkva 4,0K  8 mars  15:40 ..
```
