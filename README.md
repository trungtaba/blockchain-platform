# HetaChain

## Development Environment
- OS: Ubuntu 18.04, Mac OS 10.13.6
- Language: Go 1.9+
- IDE: Visual Studio Code + Go plugin / IntelliJ GoLand
- Unit test tool: built-in testing command (`go test`)

## Enviroment Settings
```
LOCAL_CLIENT_ID=1
HTTP_SERVICE_PORT=9000
ENODE_PORT=9100
BLOCK_TIME=10 # 10s
NODE_PRIVATE_KEY= # go run cli/main.go create address
```

## Visual Studio Code Go Plugin Settings
* Install Go Extension `ms-vscode.go`
* File > Preferences > Settings and set `"go.autocompleteUnimportedPackages": true,`

## Test 
* Clone project to $GOPATH/src/github.com/sotatek-dev/heta

```
# get dependencies
go get -v

# run each with seperate terminal
make reset
make test-bootnode

make test-producer-1
make test-producer-2
make test-producer-3
make test-producer-4
make test-client
```
