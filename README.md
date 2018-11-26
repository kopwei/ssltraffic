# ssltraffic
Establish traffic for testing purpose

## Build
```bash
$ docker run -it --rm -v $PWD:/go/ssltraffic go build -o /go/ssltraffic/ssltraffic ssltraffic
```

## Run
```bash
$ ./ssltraffic https -t 192.168.100.1
$ ./ssltraffic ssh -t localhost
$ ./ssltraffic sftp -t 192.168.20.3 -u user -p pass -f /tmp/file.pdf
```