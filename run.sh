#!/bin/bash

go build -o bed-and-breakfast cmd/web/*
./bed-and-breakfast -dbname=bedandbreakfast -dbuser=postgres -dbpass=usman123 -cache=false production=false