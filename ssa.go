package stamets

import (
	"time"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

// AllPackages builds a list of packages as an SSA program. It
// also invokes .Build() on the produced SSA program.
func AllPackages(pkgs []*packages.Package, mode ssa.BuilderMode) BaseMetrics[*ssa.Program] {
	now := time.Now()
	ssaprog, _ := ssautil.AllPackages(pkgs, mode)

	ssaprog.Build()

	return BaseMetrics[*ssa.Program]{
		Duration: time.Since(now),
		Payload:  ssaprog,
	}
}

// AllPackagesWithTimeout builds a list of packages as an SSA program in the alloted time limit.
// It also invokes .Build() on the produced SSA program.
func AllPackagesWithTimeout(t time.Duration, pkgs []*packages.Package, mode ssa.BuilderMode) (BaseMetrics[*ssa.Program], bool) {
	return TaskWIthTimeout(t, func() BaseMetrics[*ssa.Program] {
		return AllPackages(pkgs, mode)
	})
}
