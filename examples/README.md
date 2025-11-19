# Examples

This directory contains example programs demonstrating how to use go-vcenter-auth.

## Basic Login Example

`basic_login.go` - Demonstrates username/password authentication.

### Running:

```bash
export VCENTER_HOST="vcenter.example.com"
export VCENTER_USERNAME="administrator@vsphere.local"
export VCENTER_PASSWORD="your-password"

go run basic_login.go
```

## SSPI Login Example (Windows only)

`sspi_login.go` - Demonstrates Windows integrated authentication using current user's credentials.

### Running:

```bash
export VCENTER_HOST="vcenter.example.com"

go run sspi_login.go
```

**Note:** This example only works on Windows and requires the current user to have access to the vCenter server via Kerberos/SSPI.
