# go-vcenter-auth

A lightweight, standalone Go package for vCenter authentication supporting both username/password and Windows integrated authentication (SSPI/Kerberos).

## Features

- **Simple API** - Clean interface with loose parameters for easy integration
- **Multiple auth methods**:
  - Username/Password authentication
  - Windows integrated authentication (SSPI/Kerberos) - Windows only
- **Session caching** - Automatic caching to avoid repeated logins
- **Cross-platform** - Works on all platforms (SSPI on Windows only)
- **Context support** - Proper context handling for timeouts and cancellation
- **Built on govmomi** - Uses the official VMware vSphere API Go bindings

## Installation

```bash
go get github.com/skabbio1976/go-vcenter-auth
```

## Usage

### Basic Authentication (Username/Password)

```go
package main

import (
    "context"
    "fmt"
    "time"

    vcauth "github.com/skabbio1976/go-vcenter-auth"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    client, err := vcauth.Login(
        ctx,
        "vcenter.example.com",  // host
        "administrator@vsphere.local",  // username
        "password",  // password
        true,  // insecure (skip TLS verification)
    )
    if err != nil {
        panic(err)
    }

    fmt.Println("Successfully logged in to vCenter!")

    // Get underlying vim25.Client for advanced operations
    vim := client.GetVim()
    // ... use vim client for vSphere operations
}
```

### Windows Integrated Authentication (SSPI/Kerberos)

```go
package main

import (
    "context"
    "fmt"
    "time"

    vcauth "github.com/skabbio1976/go-vcenter-auth"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Uses current Windows user's credentials
    client, err := vcauth.LoginSSPI(
        ctx,
        "vcenter.example.com",  // host
        false,  // insecure
    )
    if err != nil {
        panic(err)
    }

    fmt.Println("Successfully logged in via SSPI!")

    vim := client.GetVim()
    // ... use vim client
}
```

### Working with the Cached Client

```go
// Get the cached client (returns nil if no cached session)
cachedVim := vcauth.GetCachedClient()
if cachedVim != nil {
    // Use cached session
}

// Clear the cache when needed
vcauth.ClearCache()
```

## API Reference

### `Login(ctx, host, username, password, insecure) (*Client, error)`

Authenticates to vCenter using username and password.

**Parameters:**
- `ctx` - Context for timeout/cancellation (can be nil for default 30s timeout)
- `host` - vCenter hostname or IP (e.g., "vcenter.example.com")
- `username` - vCenter username (e.g., "administrator@vsphere.local")
- `password` - vCenter password
- `insecure` - If true, skip TLS certificate verification

**Returns:** `*Client` or error

### `LoginSSPI(ctx, host, insecure) (*Client, error)` (Windows only)

Authenticates to vCenter using Windows integrated authentication (Kerberos/SSPI).
Uses the current Windows user's credentials.

**Parameters:**
- `ctx` - Context for timeout/cancellation (can be nil for default 30s timeout)
- `host` - vCenter hostname or IP
- `insecure` - If true, skip TLS certificate verification

**Returns:** `*Client` or error

**Note:** Returns `ErrSSPINotSupported` on non-Windows platforms.

### `Client.GetVim() *vim25.Client`

Returns the underlying govmomi vim25.Client for advanced vSphere operations.

### `GetCachedClient() *vim25.Client`

Returns the cached vim25.Client if available, nil otherwise.

### `ClearCache()`

Clears the cached session.

## Dependencies

- [govmomi](https://github.com/vmware/govmomi) - VMware vSphere API Go bindings
- [sspi](https://github.com/alexbrainman/sspi) - Windows SSPI bindings (Windows only)

## Platform Support

- **All platforms**: Username/Password authentication
- **Windows only**: SSPI/Kerberos authentication

## Session Caching

The package automatically caches successful sessions to avoid repeated authentication.
The cache is checked based on:
- Host
- Username (for Login)
- Session key

To clear the cache, call `ClearCache()`.

## Examples

See the [examples](examples/) directory for complete working examples.

## License

MIT License

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.
