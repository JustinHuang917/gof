#!/bin/sh
apppath=./appsite

go build ./gof.go

cp ./gof ${apppath}/

viewpath=${apppath}/view

echo $viewpath

${apppath}/gof -action clearview -viewpath ${viewpath}

${apppath}/gof -action compileview -viewpath ${viewpath} -outviewpath ${viewpath} -other ${viewpath}/helper.go

go build ${apppath}/m.go

cp ./m ${apppath}/