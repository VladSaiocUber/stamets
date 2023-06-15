package stamets

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/tools/go/packages"
)

// PackagesLoad loads packages according to the specified configuration and further
// filters them with `query`. It performs additional filtering when the configuration includes
// test packages.
func PackagesLoad(config *packages.Config, query string) (BaseMetrics[[]*packages.Package], error) {
	start := time.Now()

	pkgs, err := packages.Load(config, query)
	if err != nil {
		return BaseMetrics[[]*packages.Package]{}, err
	} else if packages.PrintErrors(pkgs) > 0 {
		return BaseMetrics[[]*packages.Package]{}, errors.New("errors encountered while loading packages")
	}
	if config.Tests {
		// Deduplicate packages that have test functions (such packages are
		// returned twice, once with no tests and once with tests. We discard
		// the package without tests.) This prevents duplicate versions of the
		// same types, functions, ssa values, etc., which can be very confusing
		// when debugging.
		packageIDs := map[string]bool{}
		for _, pkg := range pkgs {
			packageIDs[pkg.ID] = true
		}

		filteredPkgs := []*packages.Package{}
		for _, pkg := range pkgs {
			if !packageIDs[fmt.Sprintf("%s [%s.test]", pkg.ID, pkg.ID)] {
				filteredPkgs = append(filteredPkgs, pkg)
			}
		}
		pkgs = filteredPkgs
	}
	return BaseMetrics[[]*packages.Package]{
		Payload:  pkgs,
		Duration: time.Since(start),
	}, nil
}
