#!/bin/bash

set -e

NAME="example"
REQUEST="csr"
DURATION="365"
TYPE="RSA"

openssl genpkey -algorithm $TYPE -out $NAME.key

openssl req -new -key $NAME.key -out $REQUEST.pem # Distinguished Name (DN)

openssl req -x509 -days $DURATION -key $NAME.key -in $REQUEST.pem -out $NAME.crt
