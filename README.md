# Signature Service - Coding Challenge

The coding challege was developed following hexagonal architecture, with the domain logic in the domain package.

## How to run

Execute `go run main.go` to run the server.

## Example usage

```
curl -X PUT localhost:8080/api/v0/devices/1 --data '{"label": "asd", "signature_alg": "RSA"}'
curl -X GET localhost:8080/api/v0/devices/1
curl -X POST localhost:8080/api/v0/devices/1/sign -H 'Content-Type: application/json' --data '{"data_to_be_signed": "c47757abe4020b9168d0776f6c91617f9290e790ac2f6ce2bd6787c74ad88199"}'
```

## Consideration

I would have written all the tests but I cut quite few corners in order to present a working solution and give an idea about how I would develop the solution.