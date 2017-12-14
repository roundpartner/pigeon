[![Build Status](https://travis-ci.org/roundpartner/pigeon.svg?branch=master)](https://travis-ci.org/roundpartner/pigeon)
# Pigeon
A Comms Micro Service
## Abstract
Provides end points for sending communication such as Email
# Usage
```bash
export DOMAIN="mailgun domain"
export API_KEY="mailgun api key"
export PUBLIC_API_KEY="mailgun api public key"
export TEMPLATES="templates"
```
Enable test mode
```
export TEST_MODE=1
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
