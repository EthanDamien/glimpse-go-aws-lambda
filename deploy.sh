#!/bin/bash
# Run these two one after the other
GOOS=linux go build main.go createUser.go createAdmin.go

zip function.zip main