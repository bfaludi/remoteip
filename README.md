# remoteip

Detect the first valid remote IP from a request in Go. It filters out the private networks and try to find the first non-proxy address.

```go
import (
  "github.com/bfaludi/remoteip"
)


func ControllerHandler(w http.ResponseWriter, r *http.Request) {
	IP := remoteip.GetIPv4Address(r)
  ...
}

```


## License

Copyright © 2016, Bence Faludi

Distributed under the MIT License.