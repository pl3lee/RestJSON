#!/bin/bash

if [ -f .env ]; then
	source .env
fi

cd sql/schema
goose postgres $SHARED_DB_URL up
