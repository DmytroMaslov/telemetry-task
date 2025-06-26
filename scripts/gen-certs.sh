#!/bin/bash

set -e

rm -rf ../certs

mkdir -p ../certs

cd ../certs

echo "Generating CA certificate and key..."
openssl req -x509 \
        -newkey rsa:4096 \
        -days 365 \
        -nodes \
        -keyout ca-key.pem \
        -out ca-cert.pem \
        -subj "/C=UA/ST=Dnipro/L=Dnipro/CN=*.demo.test/emailAddress=test@demo.test"

openssl x509 -in ca-cert.pem -noout -text

echo "Generating server certificate and key..."
openssl req -newkey \
        rsa:4096 \
        -nodes \
        -keyout server-key.pem \
        -out server-req.pem \
        -subj "/C=UA/ST=Dnipro/L=Dnipro/CN=*.demo.test/emailAddress=test@demo.test"

echo "Generating server certificate signed by CA..."

cat > server.conf <<EOF
subjectAltName=DNS:localhost,IP:0.0.0.0,IP:127.0.0.1
EOF
cat server.conf


openssl x509 \
        -req \
        -in server-req.pem \
        -days 60 \
        -CA ca-cert.pem \
        -CAkey ca-key.pem \
        -CAcreateserial \
        -out server-cert.pem \
        -extfile server.conf

echo "Generating client certificate signed by CA..."
openssl x509 -in server-cert.pem -noout -text