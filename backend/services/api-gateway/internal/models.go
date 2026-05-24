package internal

// ProxyTarget maps a route prefix to a downstream service address.
type ProxyTarget struct {
	Prefix  string
	Address string
}
