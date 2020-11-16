[![Go Report Card](https://goreportcard.com/badge/github.com/roundpartner/pigeon)](https://goreportcard.com/report/github.com/roundpartner/pigeon)
# Pigeon
A Comms Micro Service
## Abstract
Provides end points for sending communication such as Email
## Prerequisite
```bash
dep ensure
```
# Usage
```bash
export DOMAIN="mailgun domain"
export API_KEY="mailgun api key"
export PUBLIC_API_KEY="mailgun api public key"
export TEMPLATES="templates"
export TO_EMAIL="test@mailinator.com"
```
Enable test mode
```
export TEST_MODE=1
```
Set api url to eu region
```bash
export MG_URL=https://api.eu.mailgun.net/v3
```

## Send an email
```bash
curl -X POST http://localhost:3411/email \
-d "{
        \"to\":\"receipient@mailinator.com\",
        \"from\":\"sender@mailinator.com\",
        \"subject\":\"Cool Subject\",
        \"text\":\"Interesting Message\"
    }"
```
Use a template
```bash
curl -X POST http://localhost:3411/email \
-d "{
        \"to\": \"receipient@mailinator.com\",
        \"from\": \"sender@mailinator.com\",
        \"template\": \"test\",
        \"params\": {
            \"name\": \"Cuthbert\",
            \"colour\": \"Purple\"
        }
    }"
```
### Get a template
```bash
curl -X POST http://localhost:3411/template \
-d "{
        \"template\": \"test\",
        \"params\": {
            \"name\": \"Cuthbert\",
            \"colour\": \"Purple\"
        }
    }"
```
### Verify Email
```bash
curl -X POST http://localhost:3411/verify \
-d "{
        \"email\": \"tester@mailinator.com\",
        \"ip\": \"127.0.0.1\"
    }"
```
## Upgrading Packages
Upgrade packages by specifying the package to upgrade.
```bash
go get -u github.com/aws/aws-sdk-go
```
