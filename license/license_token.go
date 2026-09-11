package license

// Token is the local representation of the existing RealmsLauncher license
// response. It is not a new wire format or cryptographic token.
type Token struct {
	LicenseKey string
	Owner      string
	Status     string
}
