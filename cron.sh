#!/bin/sh
git pull
go run list2md.go
git commit -m "Auto update" -a
git push origin
