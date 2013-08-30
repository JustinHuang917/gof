#!/bin/sh
apppath=./appsite

go build ./gof.go

cp ./gof ${apppath}/

viewpath=${apppath}/view

${apppath}/gof -action clearview -viewpath ${viewpath}

${apppath}/gof -action compileview -viewpath ${viewpath} -outviewpath ${viewpath} 

go build ${apppath}/m.go

cp ./m ${apppath}/