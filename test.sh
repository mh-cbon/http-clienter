#!/bin/sh

set -ex

rm `which http-clienter`
go install

http-clienter - demo/Controller:HTTPClientController | grep -F 'func(t HTTPClientController) GetByID(urlID int)' || exit 1;
http-clienter - demo/Controller:HTTPClientController | grep -F 'surl := "/{id}"' || exit 1;
# http-clienter -mode std - demo/Controller:HTTPClientControllerRPC | grep -F 'router.HandleFunc("GetByID", t.embed.GetByID)' || exit 1;
http-clienter - demo/Controller:HTTPClientController | grep "package main" || exit 1;
http-clienter -p nop - demo/Controller:HTTPClientControllerRPC | grep "package nop" || exit 1;

# rm -fr gen_test
# goriller demo/Controller:gen_test/ControllerGoriller || exit 1;
# ls -al gen_test | grep "controllergoriller.go" || exit 1;
# cat gen_test/controllergoriller.go | grep -F 'router.HandleFunc("/{id}", t.embed.GetByID).Methods("GET")' || exit 1;
# cat gen_test/controllergoriller.go | grep "package gen_test" || exit 1;
# rm -fr gen_test
#
# rm -fr demo/*gen.go
# go generate demo/main.go
# ls -al demo | grep "controllergoriller.go" || exit 1;
# cat demo/controllergoriller.go | grep "package main" || exit 1;
# cat demo/controllergoriller.go | grep "NewControllerGoriller(" || exit 1;
# go run demo/*.go | grep "Red" || exit 1;
#
# rm -fr demo/*gen.go
# go generate github.com/mh-cbon/goriller/demo
# ls -al demo | grep "controllergoriller.go" || exit 1;
# cat demo/controllergoriller.go | grep "package main" || exit 1;
# cat demo/controllergoriller.go | grep "NewControllerGoriller(" || exit 1;
# go run demo/*.go | grep "Red" || exit 1;
# # rm -fr demo/gen # keep it for demo
#
# # go test
#
#
# echo ""
# echo "ALL GOOD!"
