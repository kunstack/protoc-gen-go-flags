package module

import (
	pgs "github.com/lyft/protoc-gen-star/v2"
)

var _ pgs.Module = (*WrapModule)(nil)

type WrapModule struct {
	buildContext pgs.BuildContext
}

func (s *WrapModule) Name() string {
	return "flags"
}

func (s *WrapModule) InitContext(c pgs.BuildContext) {
	s.buildContext = c
}

func (s *WrapModule) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	as := make([]pgs.Artifact, 0, len(targets))
	for _, f := range targets {
		mod := Flags()
		mod.InitContext(s.buildContext)
		as = append(as, mod.Execute(f, packages)...)
	}
	return as
}
