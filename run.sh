#!/bin/bash

go build -o bookings cmd/web/*.go
./bookings -dbname=golangbookings -dbuser=root -dbpass=root -dbport=3306 -cache=false -production=false