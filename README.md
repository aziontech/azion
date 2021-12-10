# azion-cli

## Azion-CLI - How to Use

```
#build project
$ make build
#then run the binary with the flag --version
$ ./azion --version

#OR

#build project
$ make build
#then run the binary with the command version
$ ./azion version

#build project defining the local test environment for token authentication
$ AUTH_ENDPOINT="http://localhost:8080/" make build

#open a new terminal and start the testing http server

$ cd tests/
$ ./httpserver.sh

#...then test the binary with valid and invalid tokens
$ ./azion configure -t tokenValid

#AND

$ ./azion configure -t tokenInvalid


```