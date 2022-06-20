package gen

import (
	"path/filepath"

	gen "github.com/whyrusleeping/cbor-gen"
)

func GenCborType(outputdir string, pkg string, types ...interface{}) error {
	if pkg == "" {
		pkg = filepath.Base(outputdir)
	}
	return gen.WriteTupleEncodersToFile(filepath.Join(outputdir, "cbor_gen.go"), pkg, types...)
}
