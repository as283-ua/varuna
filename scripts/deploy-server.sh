#!/bin/bash
go build -o ../server/server ../server
sudo setcap 'cap_net_bind_service=+ep' ../server/server
