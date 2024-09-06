#!/bin/bash

DBSTRING="host=localhost user=postgres password=postgres port=5432 dbname=calendar sslmode=disable"

goose postgres "$DBSTRING" up