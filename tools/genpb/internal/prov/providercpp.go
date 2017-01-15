package prov

type cppRrovider struct {
}

func (p *cppRrovider) GenCmdFile(protoFile string) (string, error) {
	return "Hello", nil
}

func init() {
	Register("cpp", &goRrovider{})
}
