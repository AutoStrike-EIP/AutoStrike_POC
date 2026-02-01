#!/bin/bash
set -e

# AutoStrike Certificate Generation Script
# Generates CA, server, and agent certificates for mTLS

CERT_DIR="${1:-./certs}"
DAYS_VALID=365
COUNTRY="FR"
STATE="IDF"
LOCALITY="Paris"
ORG="AutoStrike"
ORG_UNIT="Security"

mkdir -p "$CERT_DIR"
cd "$CERT_DIR"

echo "=== Generating CA Certificate ==="
openssl genrsa -out ca.key 4096
openssl req -new -x509 -days $DAYS_VALID -key ca.key -out ca.crt \
    -subj "/C=$COUNTRY/ST=$STATE/L=$LOCALITY/O=$ORG/OU=$ORG_UNIT/CN=AutoStrike CA"

echo "=== Generating Server Certificate ==="
openssl genrsa -out server.key 2048

cat > server.ext <<EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = server
DNS.3 = autostrike-server
IP.1 = 127.0.0.1
IP.2 = ::1
EOF

openssl req -new -key server.key -out server.csr \
    -subj "/C=$COUNTRY/ST=$STATE/L=$LOCALITY/O=$ORG/OU=$ORG_UNIT/CN=AutoStrike Server"

openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
    -out server.crt -days $DAYS_VALID -extfile server.ext

echo "=== Generating Agent Certificate Template ==="
openssl genrsa -out agent-template.key 2048

cat > agent.ext <<EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
extendedKeyUsage = clientAuth
EOF

openssl req -new -key agent-template.key -out agent-template.csr \
    -subj "/C=$COUNTRY/ST=$STATE/L=$LOCALITY/O=$ORG/OU=Agents/CN=AutoStrike Agent"

openssl x509 -req -in agent-template.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
    -out agent-template.crt -days $DAYS_VALID -extfile agent.ext

# Cleanup
rm -f *.csr *.ext *.srl

echo ""
echo "=== Certificates Generated Successfully ==="
echo "Directory: $CERT_DIR"
echo ""
echo "Files created:"
echo "  - ca.crt / ca.key         : Certificate Authority"
echo "  - server.crt / server.key : Server certificate"
echo "  - agent-template.crt/key  : Agent certificate template"
echo ""
echo "To verify: openssl verify -CAfile ca.crt server.crt"
