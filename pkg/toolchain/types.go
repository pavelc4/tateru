package toolchain

type Env struct {
	ClangBin string
	Version  string
}

func (e *Env) Env() []string {
	return []string{
		"CC=" + e.ClangBin,
		"CLANG=" + e.ClangBin,
	}
}
