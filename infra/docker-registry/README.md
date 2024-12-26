regenerate certs:

```
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=registry.tymbaca/O=registry" -addext "subjectAltName = DNS:registry.tymbaca"
kubectl create secret tls registry-tls --cert=tls.crt --key=tls.key
```
