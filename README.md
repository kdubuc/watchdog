# Watchdog

```bash
./watchdog -f "php index.php" -p 8080
```

Tiny Golang webserver as a generic stdio interface : marshal a HTTP request and pass it into a function's fork stdin, wait for function stdout and return HTTP response.

Huge inspiration from : https://github.com/openfaas/faas/tree/master/watchdog
