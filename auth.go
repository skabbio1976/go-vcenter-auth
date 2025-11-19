package vcenterauth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/vmware/govmomi/session"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/soap"
)

var (
	clientMu   sync.Mutex
	cachedHost string
	cachedUser string
	cachedSess string
	cachedVim  *vim25.Client
)

// Client represents a vCenter client
type Client struct {
	vim *vim25.Client
}

// Login performs authentication to vCenter with username and password.
// Parameters:
//   - ctx: Context for timeout/cancellation
//   - host: vCenter hostname or IP (e.g., "vcenter.example.com")
//   - username: vCenter username
//   - password: vCenter password
//   - insecure: If true, skip TLS certificate verification
//
// Returns a Client instance or error if login fails.
func Login(ctx context.Context, host, username, password string, insecure bool) (*Client, error) {
	if host == "" || username == "" {
		return nil, fmt.Errorf("host or username missing")
	}

	// Cache check
	clientMu.Lock()
	if cachedVim != nil && cachedHost == host && cachedUser == username && cachedSess != "" {
		vim := cachedVim
		clientMu.Unlock()
		return &Client{vim: vim}, nil
	}
	clientMu.Unlock()

	// Build URL (always add https:// and /sdk if not present)
	raw := host
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
	u.User = url.UserPassword(username, password)

	soapClient := soap.NewClient(u, insecure)

	// Use provided context with fallback timeout
	loginCtx := ctx
	if ctx == nil {
		var cancel context.CancelFunc
		loginCtx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	}

	vimClient, err := vim25.NewClient(loginCtx, soapClient)
	if err != nil {
		return nil, fmt.Errorf("could not create vim25 client: %w", err)
	}
	sm := session.NewManager(vimClient)
	if err := sm.Login(loginCtx, u.User); err != nil {
		return nil, fmt.Errorf("login failed: %w", err)
	}
	us, err := sm.UserSession(loginCtx)
	if err != nil || us == nil {
		return nil, fmt.Errorf("could not fetch user session: %w", err)
	}

	// Cache session
	sessID := randomID()
	clientMu.Lock()
	cachedHost = host
	cachedUser = username
	cachedSess = sessID
	cachedVim = vimClient
	clientMu.Unlock()

	return &Client{vim: vimClient}, nil
}

// GetVim returns the underlying vim25.Client for advanced operations
func (c *Client) GetVim() *vim25.Client {
	return c.vim
}

// GetCachedClient returns the cached vim25.Client if available
func GetCachedClient() *vim25.Client {
	clientMu.Lock()
	defer clientMu.Unlock()
	return cachedVim
}

// ClearCache clears the cached session
func ClearCache() {
	clientMu.Lock()
	defer clientMu.Unlock()
	cachedHost = ""
	cachedUser = ""
	cachedSess = ""
	cachedVim = nil
}

func hasSDKPath(u string) bool {
	return len(u) >= 4 && (u[len(u)-4:] == "/sdk" || u[len(u)-4:] == "sdk/")
}

func randomID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
