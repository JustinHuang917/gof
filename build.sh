#!/bin/sh
apppath=./appsite
#export GOPATH=${apppath}/goftool:${apppath}/gofweb

go build ./gof.go

cp ./gof ${apppath}/

#clcp ${apppath}/goftool/src/gof ${apppath}/gofweb/src

viewpath=${apppath}/view

 echo $viewpath

 ${apppath}/gof -action clearview -viewpath ${viewpath}

 ${apppath}/gof -action compileview -viewpath ${viewpath} -outviewpath ${viewpath} -other ${viewpath}/helper.go

 go build ${apppath}/m.go
 cp ./m ${apppath}/