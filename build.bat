@echo off
@color 06

SET CGO_ENABLED=0

SET GOOS=linux

SET GOARCH=amd64

go build ./service/idc.go


pause