# TLS upon TCP/HTTP/RPC in go

This repo gives an simple example of TLS configure upon TCP/HTTP/RPC in go.

To run these code:

```bash
# run server and client with cmd below in two terminal
go run fileName.go
```

If want to regenerate the key and cert

```bash
openssl req -x509 -nodes -days 365 -newkey rsa:4096 -keyout key_name.key -out cert_name.crt -config openssl.cnf -extensions req_ext
```

