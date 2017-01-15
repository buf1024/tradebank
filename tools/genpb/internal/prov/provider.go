package prov

// Provider interface
type Provider interface {
	GenCmdFile(protoFile string) (string, error)
}

var provider map[string]Provider

// GetProvider return the provider
func GetProvider(ext string) Provider {
	if prov, ok := provider[ext]; ok {
		return prov
	}
	return nil
}

// Register register a provider
func Register(ext string, prov Provider) {
	if provider == nil {
		provider = make(map[string]Provider)
	}
	if _, ok := provider[ext]; !ok {
		provider[ext] = prov
	}
}
