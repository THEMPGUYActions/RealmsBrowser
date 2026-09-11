package license

type ProcessManager struct {
	Close func()
}

func (p *ProcessManager) TerminateBrowser() {
	if p != nil && p.Close != nil {
		p.Close()
	}
}
