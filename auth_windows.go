//go:build windows

package vcenterauth

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/alexbrainman/sspi/negotiate"
	"github.com/vmware/govmomi/vim25"
	vimmethods "github.com/vmware/govmomi/vim25/methods"
	vimsoap "github.com/vmware/govmomi/vim25/soap"
	vimtypes "github.com/vmware/govmomi/vim25/types"
)

// LoginSSPI establishes a vCenter session using Windows integrated authentication (Kerberos/SSPI).
// This function is only available on Windows platforms.
//
// Parameters:
//   - ctx: Context for timeout/cancellation
//   - host: vCenter hostname or IP (e.g., "vcenter.example.com")
//   - insecure: If true, skip TLS certificate verification
//
// The function uses the current Windows user's credentials for authentication.
// Returns a Client instance or error if login fails.
func LoginSSPI(ctx context.Context, host string, insecure bool) (*Client, error) {
	h := normalizeServerHost(host)
	if h == "" {
		return nil, fmt.Errorf("empty host - no vCenter specified")
	}

	// Cache check
	clientMu.Lock()
	if cachedVim != nil && cachedHost == h && cachedSess != "" {
		vim := cachedVim
		clientMu.Unlock()
		return &Client{vim: vim}, nil
	}
	clientMu.Unlock()

	// Build URL with https:// and /sdk
	raw := h
	if !(len(raw) >= 8 && (raw[:8] == "https://" || raw[:7] == "http://")) {
		raw = "https://" + raw
	}
	// Remove trailing slash if present
	if raw[len(raw)-1] == '/' {
		raw = raw[:len(raw)-1]
	}
	if !hasSDKPath(raw) {
		raw = raw + "/sdk"
	}
	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	soapClient := vimsoap.NewClient(u, insecure)

	// Use provided context with fallback timeout
	loginCtx := ctx
	if ctx == nil {
		var cancel context.CancelFunc
		loginCtx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	c, err := vim25.NewClient(loginCtx, soapClient)
	if err != nil {
		return nil, fmt.Errorf("vim25.NewClient failed: %w", err)
	}

	// Prepare SSPI (Kerberos) context for SPN host/<vcenter>
	cred, err := negotiate.AcquireCurrentUserCredentials()
	if err != nil {
		return nil, fmt.Errorf("AcquireCurrentUserCredentials: %w", err)
	}
	defer cred.Release()

	target := "host/" + h
	secctx, outToken, err := negotiate.NewClientContext(cred, target)
	if err != nil {
		return nil, fmt.Errorf("NewClientContext: %w", err)
	}
	defer secctx.Release()

	var sess vimtypes.UserSession
	for {
		req := vimtypes.LoginBySSPI{
			This:        *c.ServiceContent.SessionManager,
			Locale:      "en_US",
			Base64Token: base64.StdEncoding.EncodeToString(outToken),
		}

		resp, err := vimmethods.LoginBySSPI(loginCtx, c, &req)
		if err == nil {
			sess = resp.Returnval
			break
		}

		// Handle SSPIChallenge: update client security context and continue
		if vimsoap.IsSoapFault(err) {
			if vf := vimsoap.ToVimFault(err); vf != nil {
				if ch, ok := vf.(*vimtypes.SSPIChallenge); ok {
					in, _ := base64.StdEncoding.DecodeString(ch.Base64Token)
					done, next, uerr := secctx.Update(in)
					if uerr != nil {
						return nil, fmt.Errorf("SSPI Update: %w", uerr)
					}
					outToken = next
					if done && len(outToken) == 0 {
						outToken = []byte{}
					}
					continue
				}
			}
		}
		return nil, fmt.Errorf("LoginBySSPI: %w", err)
	}

	if sess.Key == "" {
		return nil, fmt.Errorf("SSPI login missing session")
	}

	// Cache session
	clientMu.Lock()
	cachedVim = c
	cachedHost = h
	cachedUser = sess.UserName
	cachedSess = sess.Key
	clientMu.Unlock()

	return &Client{vim: c}, nil
}

// normalizeServerHost removes schema (http/https), port and path so only FQDN/IP remains
func normalizeServerHost(h string) string {
	hs := strings.TrimSpace(h)
	if hs == "" {
		return hs
	}
	if strings.HasPrefix(strings.ToLower(hs), "http://") || strings.HasPrefix(strings.ToLower(hs), "https://") {
		if u, err := url.Parse(hs); err == nil {
			hostPart := u.Host
			// Remove port
			if idx := strings.Index(hostPart, ":"); idx != -1 {
				hostPart = hostPart[:idx]
			}
			return hostPart
		}
	}
	// Remove path if user wrote e.g. vcenter.local/sdk
	if slash := strings.IndexRune(hs, '/'); slash != -1 {
		hs = hs[:slash]
	}
	// Remove port if no schema
	if idx := strings.Index(hs, ":"); idx != -1 {
		hs = hs[:idx]
	}
	return hs
}
