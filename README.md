# http-echo
A simple http-server which can be used as a drop-in to inspect http requests.

## Usage
- `--listenaddr`: listen address; port only: `:8080` or with interface to bind on `127.0.0.1:8080` (default `:8080`)
- `--keepalive`: enable keepalive; true or false (default true)

## Example Response
The response is sent as a plaintext in the response body. It mimics the received request
```text
GET / HTTP/1.1
Host: localhost:8080
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7
Accept-Encoding: gzip, deflate, br, zstd
Accept-Language: de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7
Connection: keep-alive
Sec-Ch-Ua: "Not:A-Brand";v="99", "Google Chrome";v="145", "Chromium";v="145"
Sec-Ch-Ua-Mobile: ?0
Sec-Ch-Ua-Platform: "macOS"
Sec-Fetch-Dest: document
Sec-Fetch-Mode: navigate
Sec-Fetch-Site: none
Sec-Fetch-User: ?1
Upgrade-Insecure-Requests: 1
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36

``` 

## Docker Images
- Dockerhub: `gw2auth/http-echo`
- GitHub Container Registry: `ghcr.io/gw2auth/http-echo`