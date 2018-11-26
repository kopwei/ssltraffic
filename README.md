# ssltraffic
Establish traffic for testing purpose

## Build
```bash
$ cd /path_to_repo
$ docker run -it --rm \
    -v $PWD:/go/src/github.com/kopwei/ssltraffic golang:1.11.2 \
    go build -o /go/src/github.com/kopwei/ssltraffic/ssltraffic \
    github.com/kopwei/ssltraffic
```

## Run
```bash
$ ./ssltraffic https -t 192.168.100.1
$ ./ssltraffic ssh -t localhost
$ ./ssltraffic sftp -t 192.168.20.3 -u user -p pass -f /tmp/file.pdf
```
