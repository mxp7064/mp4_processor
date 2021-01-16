## Documentation

### Initialization segment

- https://www.w3.org/2013/12/byte-stream-format-registry/isobmff-byte-stream-format.html

- https://www.file-recovery.com/mp4-signature-format.htm

- An ISO BMFF initialization segment is defined as a single File Type Box (ftyp) followed by a single Movie Header Box (moov)

- Valid top-level boxes such as pdin, free, and sidx are allowed to appear before the moov box. These boxes must be accepted and ignored by the user agent and are not considered part of the initialization segment

- video.mp4 byte breakdown
  ```
  0-3 = 36 (0 0 0 36)
  4-7 = ftyp (102 116 121 112)
  8-35 = data (36-4-4=28 of data, data end is at 7+28=35)
  36-39 = 787 (0 0 3 19)
  40-43 = moov (109 111 111 118)
  44-822 = data (onda 787-4-4=779, data end is at 43+779=822)
  823-826 = 36 (0 0 0 36) (ftypSize+moovSize=36+787=823)
  827-830 = styp
  830... = media data
  ```

### Nats

- nats docker tutorial
https://docs.nats.io/nats-server/nats_docker/nats-docker-tutorial

- nats docker image documentation:
https://hub.docker.com/_/nats

- nats docs:
https://docs.nats.io/

- remove previous nats-server container and run
```sh
docker rm -f $(docker container ls -a | grep nats-server | awk '{print $1}');
docker run -d -p 4222:4222 -p 8222:8222 -p 6222:6222 --name nats-server nats:latest --user nats-user --pass pass123;
```

- check nats-server container logs
```sh
docker logs $(docker container ls -a | grep nats-server | awk '{print $1;}') -f
```

### Go Modules

- initialize a project using Modules, it’ll generate a go.mod module config file in your project’s root directory
```sh
go mod init extractor-service
```

- fetch the missing dependencies automatically and include them in go.mod file
```sh
go build
```

- remove any unused dependencies in your project and update the go.mod file, run the go mod tidy command after making any changes to your code. It’ll ensure your module file is accurate and clean
```sh
go mod tidy
```

- list all dependencies
```sh
go list -m all
```

+ references:
  + https://www.honeybadger.io/blog/golang-go-package-management/
  + https://www.whitesourcesoftware.com/free-developer-tools/blog/golang-dependency-management/

## Build & Run

+ system dependencies
  + tested with go 1.15.6 and node 14.15.4
  + docker
  
### Nats

- run:
```sh
docker run -d -p 4222:4222 -p 8222:8222 -p 6222:6222 --name nats-server nats:latest --user nats-user --pass pass123;
```

### Node.js app

- run:
```sh
cd app
npm install
./app_start.sh
```

### Go extractor_service

- in new terminal run:
```sh
cd extractor_service
go build
./es_start.sh
```

## Test use cases

+ success case 
  + provide valid mp4 file path (which has an initialization segment) to the nodejs app command line 
  + if your mp4 file path is at .../videos/video.mp4, your initialization segment file will be created in the same folder and it will have a random UUID name
+ exception cases
  + input file path which doesn't exist
    + nodejs app prints error/validation message
    + continues working and waits for new input
  + input file which is not mp4
    + nodejs app prints error/validation message
    + continues working and waits for new input
  + input file which doesn't have an initialization segment
    + go extractor service returns error message via nats to nodejs app
    + nodejs app continues working and waits for new input
    + go extractor service continues waiting for new file path via nats
  + nodejs app and nats working but go extractor service is down
    + when providing the file path to the nodejs app, it will timeout after 5 sec, print the error message and continue waiting for new input
    + if the go extractor service is started, everything will continue to work properly