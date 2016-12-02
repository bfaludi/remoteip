# remoteip

Detect the first valid remote IPv4 address from a request in Go. During the process, it filters out the private networks and try to find the first non-proxy address.

```go
import (
	"github.com/bfaludi/remoteip"
	"net/http"
)

func ControllerHandler(w http.ResponseWriter, r *http.Request) {
	IP := remoteip.GetIPv4Address(r)
}
```

## License

Copyright © 2016, Bence Faludi

Distributed under the MIT License.