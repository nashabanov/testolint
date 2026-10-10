package run_discovery

import tt "github.com/ozontech/testo"

func aliased(t *tt.T) {
	tt.RunSuite(t, S{})    // want "TESTO_RUN_DISCOVERY: RunSuite run_discovery.S"
	tt.RunSubSuite(t, S{}) // want "TESTO_RUN_DISCOVERY: RunSubSuite run_discovery.S"
}
