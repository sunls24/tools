### hmac-sha1

```shell
go build -o $GOPATH/bin/hmac-sha1 encryption/hmac.go
```

### Scaleway

```
curl https://raw.githubusercontent.com/sunls24/tools/master/script/sw-init.sh > sw-init.sh
chmod +x sw-init.sh
./sw-init.sh

systemctl enable --now wg-quick@wgcf
```
