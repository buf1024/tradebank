package prov

type goRrovider struct {
}

func (p *goRrovider) GenCmdFile(protoFile string) (string, error) {
	return "Hello", nil
}

func init() {
	Register("go", &goRrovider{})
}
