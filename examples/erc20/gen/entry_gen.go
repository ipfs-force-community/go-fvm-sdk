package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"

	"golang.org/x/tools/imports"

	typegen "github.com/whyrusleeping/cbor-gen"

	"html/template"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"
)

var (
	knownPackageNamesMu sync.Mutex
	pkgNameToPkgPath    = make(map[string]string)
	pkgPathToPkgName    = make(map[string]string)
)

func init() {
	for _, imp := range defaultClientImport {
		if was, conflict := pkgNameToPkgPath[imp.Name]; conflict {
			panic(fmt.Sprintf("reused pkg name %s for %s and %s", imp.Name, imp.PkgPath, was))
		}
		if _, conflict := pkgPathToPkgName[imp.Name]; conflict {
			panic(fmt.Sprintf("duplicate default import %s", imp.PkgPath))
		}
		pkgNameToPkgPath[imp.Name] = imp.PkgPath
		pkgPathToPkgName[imp.PkgPath] = imp.Name
	}
}

func gen_entry(stateT reflect.Type, output string) error {
	entryMeta, err := getEntryPackageMeta("main", stateT)
	if err != nil {
		return err
	}

	render, err := template.New("gen entry").Funcs(map[string]interface{}{
		"trimPrefix": strings.TrimPrefix,
	}).Parse(tml)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	err = render.Execute(buf, entryMeta)
	if err != nil {
		return err
	}

	return formateAndWriteCode(buf.Bytes(), output)
}

func formateAndWriteCode(code []byte, output string) error {
	f, err := os.Create(output)
	if err != nil {
		return err
	}

	fmtCode, err := format.Source(code)
	if err != nil {
		return err
	}
	fmtCode, err = imports.Process(output, fmtCode, nil)
	if err != nil {
		return err
	}
	if _, err = f.Write(fmtCode); err != nil {
		return err
	}

	return f.Close()
}

func getEntryPackageMeta(pkg string, stateT reflect.Type) (*entryMeta, error) {
	if stateT.Kind() == reflect.Ptr {
		stateT = stateT.Elem()
	}

	stateV := reflect.New(stateT)
	exportFunc, found := stateV.Type().MethodByName("Export")
	if !found {
		return nil, fmt.Errorf("state must have export function")
	}
	returns := exportFunc.Func.Call([]reflect.Value{stateV})
	exports, ok := returns[0].Interface().(map[int]interface{})
	if !ok {
		return nil, fmt.Errorf("assert Export return type fail,")
	}

	stateName := typeName(pkg, stateT)
	var methodsMap []methodMap
	var imports []typegen.Import
	imports = append(imports, typegen.ImportsForType(pkg, stateT)...)
	sortedFuncs := sortMap(exports)
	for _, sortedFunc := range sortedFuncs {
		functionT := reflect.TypeOf(sortedFunc.funcT)

		if functionT.Kind() != reflect.Func {
			return nil, fmt.Errorf("export must be function ")
		}
		var method = methodMap{}
		v := reflect.ValueOf(sortedFunc.funcT)
		method.FuncName = getFunctionName(v)
		method.MethodNum = sortedFunc.method_num
		method.StateName = stateName
		//	functionT := function.Type

		if functionT.NumIn() > 1 {
			return nil, fmt.Errorf("func %s can not have params more than 1", method.FuncName)
		}
		if functionT.NumIn() == 1 {
			if !functionT.In(0).AssignableTo(unMarshallerT) {
				return nil, fmt.Errorf("func %s return value at index 1 must be error", method.FuncName)
			}
			method.HasParams = true
			method.ParamsTypeName = typeName(pkg, functionT.In(0))
			imports = append(imports, typegen.ImportsForType(pkg, functionT.In(0))...)
		}

		if functionT.NumOut() > 2 {
			return nil, fmt.Errorf("func %s can not have return value more than 2", method.FuncName)
		}

		if functionT.NumOut() == 2 {
			if !functionT.Out(0).AssignableTo(marshallerT) {
				return nil, fmt.Errorf("func %s return value at index 0 must be marshaller", method.FuncName)
			}

			if !functionT.Out(1).AssignableTo(errorT) {
				return nil, fmt.Errorf("func %s return value at index 1 must be error", method.FuncName)
			}
			method.HasReturn = true
			method.HasError = true
			method.ReturnTypeName = typeName(pkg, functionT.Out(0))
			method.DefaultReturn = defaultValue(functionT.Out(0))
			imports = append(imports, typegen.ImportsForType(pkg, functionT.Out(0))...)
		} else if functionT.NumOut() == 1 {
			if functionT.Out(0).AssignableTo(errorT) {
				method.HasReturn = false
				method.HasError = true
			} else {
				method.ReturnTypeName = typeName(pkg, functionT.Out(0))
				method.DefaultReturn = defaultValue(functionT.Out(0))
				method.HasReturn = true
				method.HasError = false
			}
		} else {
			//no return
			method.HasReturn = false
			method.HasError = false
		}

		methodsMap = append(methodsMap, method)

	}
	return &entryMeta{
		Imports:   dedupImports(imports),
		Methods:   methodsMap,
		StateName: stateName,
	}, nil
}

func typeName(pkg string, t reflect.Type) string {
	switch t.Kind() {
	case reflect.Array:
		return fmt.Sprintf("[%d]%s", t.Len(), typeName(pkg, t.Elem()))
	case reflect.Slice:
		return "[]" + typeName(pkg, t.Elem())
	case reflect.Ptr:
		return "*" + typeName(pkg, t.Elem())
	case reflect.Map:
		return "map[" + typeName(pkg, t.Key()) + "]" + typeName(pkg, t.Elem())
	default:
		pkgPath := t.PkgPath()
		if pkgPath == "" {
			// It's a built-in.
			return t.String()
		} else if pkgPath == pkg {
			return t.Name()
		}
		return fmt.Sprintf("%s.%s", resolvePkgName(pkgPath, t.String()), t.Name())
	}
}

func resolvePkgName(path, typeName string) string {
	parts := strings.Split(typeName, ".")
	if len(parts) != 2 {
		panic(fmt.Sprintf("expected type to have a package name: %s", typeName))
	}
	defaultName := parts[0]

	knownPackageNamesMu.Lock()
	defer knownPackageNamesMu.Unlock()

	// Check for a known name and use it.
	if name, ok := pkgPathToPkgName[path]; ok {
		return name
	}

	// Allocate a name.
	for i := 0; ; i++ {
		tryName := defaultName
		if i > 0 {
			tryName = fmt.Sprintf("%s%d", defaultName, i)
		}
		if _, taken := pkgNameToPkgPath[tryName]; !taken {
			pkgNameToPkgPath[tryName] = path
			pkgPathToPkgName[path] = tryName
			return tryName
		}
	}

}

type sortMethod struct {
	method_num int
	funcT      interface{}
}

func sortMap(v map[int]interface{}) []sortMethod {
	var sortMethods []sortMethod
	for index, method := range v {
		sortMethods = append(sortMethods, sortMethod{
			method_num: index,
			funcT:      method,
		})
	}
	sort.Slice(sortMethods, func(i, j int) bool {
		return sortMethods[i].method_num < sortMethods[j].method_num
	})
	return sortMethods
}

func dedupImports(imps []typegen.Import) []typegen.Import {
	impSet := make(map[string]string, len(imps))
	for _, imp := range imps {
		impSet[imp.PkgPath] = imp.Name
	}
	deduped := make([]typegen.Import, 0, len(imps))
	for pkg, name := range impSet {
		deduped = append(deduped, typegen.Import{Name: name, PkgPath: pkg})
	}
	sort.Slice(deduped, func(i, j int) bool {
		return deduped[i].PkgPath < deduped[j].PkgPath
	})
	return deduped
}

type entryMeta struct {
	Imports   []typegen.Import
	Methods   []methodMap
	StateName string
}

type methodMap struct {
	StateName      string
	MethodNum      int
	FuncName       string
	HasError       bool
	HasParams      bool
	HasReturn      bool
	ParamsTypeName string
	ReturnTypeName string
	DefaultReturn  string
}

var tml = `// Code generated by github.com/ipfs-force-community/go-fvm-sdk. DO NOT EDIT.
package main

import (
	"bytes"
	"fmt"
	"github.com/filecoin-project/go-state-types/cbor"

	typegen "github.com/whyrusleeping/cbor-gen"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"

	{{range .Imports}}
	 {{.Name}} "{{.PkgPath}}"
	{{end}}
)

//not support non-main wasm in tinygo at present
func main() {}

/// The actor's WASM entrypoint. It takes the ID of the parameters block,
/// and returns the ID of the return value block, or NO_DATA_BLOCK_ID if no
/// return value.
///
/// Should probably have macros similar to the ones on fvm.filecoin.io snippets.
/// Put all methods inside an impl struct and annotate it with a derive macro
/// that handles state serde and dispatch.
//go:export invoke
func Invoke(blockId uint32) uint32 {
	method, err := sdk.MethodNumber()
	if err != nil {
		sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to get method number")
	}

	var callResult cbor.Marshaler
	var raw *types.ParamsRaw
	switch method {
{{range .Methods}}case {{.MethodNum}}:
{{if eq .FuncName  "Constructor"}}  //Constuctor
		{{if .HasParams}}raw, err = sdk.ParamsRaw(blockId)
						if err != nil {
							sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
						}
						var req {{trimPrefix .ParamsTypeName "*"}}
						err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
						if err != nil {
							sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to unmarshal params raw")
						}
						err = new({{.StateName}}).Constructor(&req)
						callResult = typegen.CborBool(true)
          {{else}}err = new({{.StateName}}).Constructor()
                callResult = typegen.CborBool(true)
          {{end}}
{{else}}
		  {{if .HasParams}}raw, err = sdk.ParamsRaw(blockId)
								if err != nil {
									sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
								}
								var req {{trimPrefix .ParamsTypeName "*"}}
								err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
								if err != nil {
									sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to unmarshal params raw")
								}
       		 {{if .HasError}}
					 {{if .HasReturn}} //have params/return/error
								state := new({{.StateName}})
								sdk.LoadState(state)
								callResult, err = state.{{.FuncName}}(&req)
				     {{else}} 	//have params/error but no return val
								state := new({{.StateName}})
								sdk.LoadState(state)
								if err = state.{{.FuncName}}(&req); err == nil {
									callResult = typegen.CborBool(true)
								}
					{{end}}
			{{else}}
					{{if .HasReturn}}//have params/return but no error
							state := new({{.StateName}})
							sdk.LoadState(state)
							callResult = state.{{.FuncName}}(&req)
					{{else}}//have params but no return value and error
							state := new({{.StateName}})
							sdk.LoadState(state)
							state.{{.FuncName}}(&req)
							callResult = = typegen.CborBool(true)
					{{end}}
			{{end}}
    {{else}}
			{{if .HasError}}
					 {{if .HasReturn}} //no params but return value/error
							state := new({{.StateName}})
							sdk.LoadState(state)
							callResult, err = state.{{.FuncName}}()
					{{else}}	//no params/return value but return error
							state := new({{.StateName}})
							sdk.LoadState(state)
							if err = state.{{.FuncName}}(); err == nil {
									callResult = = typegen.CborBool(true)
								}
					{{end}}
			{{else}}
					{{if .HasReturn}}	//no params no error but have return value
						state := new({{.StateName}})
						sdk.LoadState(state)
						callResult = state.{{.FuncName}}()
					{{else}}		//no params/return value/error
						state := new({{.StateName}})
						sdk.LoadState(state)
						state.{{.FuncName}}()
						callResult = = typegen.CborBool(true)
					{{end}}
			{{end}}
    {{end}}
{{end}}
{{end}}
	default:
		sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unsupport method")
	}

	if err != nil {
		sdk.Abort(ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("call error %s", err))
	}

	if !sdk.IsNil(callResult) {
		buf := bytes.NewBufferString("")
		err = callResult.MarshalCBOR(buf)
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("marshal resp fail %s", err))
		}
		id, err := sdk.PutBlock(types.DAG_CBOR, buf.Bytes())
		if err != nil {
			sdk.Abort(ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to store return value: %v", err))
		}
		return id
	} else {
		return types.NO_DATA_BLOCK_ID
	}
}

`

func defaultValue(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Bool:
		return "false"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fallthrough
	case reflect.Float32, reflect.Float64:
		fallthrough
	case reflect.Complex64, reflect.Complex128:
		fallthrough
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "0"
	case reflect.String:
		return `""`
	case reflect.Map:
		fallthrough
	case reflect.Ptr:
		fallthrough
	case reflect.Slice:
		return "nil"
	case reflect.Struct:
		switch t.Name() {
		case "Address":
			return "address.Undef"
		case "cid":
			return "cid.Undef"
		default:
			return fmt.Sprintf("%s{}", t.Name())
		}
	default:
		panic("unsupprt type")
	}
}
func getFunctionName(temp reflect.Value) string {
	strs := strings.Split((runtime.FuncForPC(temp.Pointer()).Name()), ".")
	strs = strings.Split(strs[len(strs)-1], "-")
	return strs[0]
}
