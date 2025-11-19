# Examples

This directory contains example programs demonstrating how to use go-vcenter-auth.

## Basic Login Example

`basic-login/` - Demonstrates username/password authentication.

### Running:

```bash
export VCENTER_HOST="vcenter.example.com"
export VCENTER_USERNAME="administrator@vsphere.local"
export VCENTER_PASSWORD="your-password"

cd basic-login
go run main.go
```

Or from the examples directory:
```bash
go run ./basic-login
```

## SSPI Login Example (Windows only)

`sspi-login/` - Demonstrates Windows integrated authentication using current user's credentials.

### Running:

```bash
export VCENTER_HOST="vcenter.example.com"

cd sspi-login
go run main.go
```

Or from the examples directory:
```bash
go run ./sspi-login
```

**Note:** This example only works on Windows and requires the current user to have access to the vCenter server via Kerberos/SSPI.
