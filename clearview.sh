#!/bin/sh
apppath=./appsite

go build ./gof.go

cp ./gof ${apppath}/

viewpath=${apppath}/view

${apppath}/gof -action clearview -viewpath ${viewpath}
