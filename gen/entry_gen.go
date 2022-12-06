package gen

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"os"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/frc42dispatch"

	"golang.org/x/tools/imports"

	"github.com/filecoin-project/go-state-types/cbor"

	typegen "github.com/whyrusleeping/cbor-gen"

	"html/template"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"
)

var (
	unMarshallerT       = reflect.TypeOf((*cbor.Unmarshaler)(nil)).Elem()
	errorT              = reflect.TypeOf((*error)(nil)).Elem()
	marshallerT         = reflect.TypeOf((*cbor.Marshaler)(nil)).Elem()
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

func GenEntry(stateT reflect.Type, output string) error {
	entryMeta, err := getEntryPackageMeta("main", stateT)
	if err != nil {
		return err
	}

	render, err := template.New("gen entry").Funcs(funcs).Parse(tml)
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

func formateAndWriteCode(code []byte, output string) (err error) {
	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			//output error file for debug
			fmt.Println(string(code))
		}
	}()

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

func isContext(r reflect.Type) bool {
	return r.PkgPath() == "context" && r.Name() == "Context"
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
	exports, ok := returns[0].Interface().([]interface{})
	if !ok {
		return nil, errors.New("assert Export return type fail")
	}

	var methodsArr []*methodMap
	hasParam := false
	typesToImport := []reflect.Type{stateT}
	for _, actorMethod := range exports {
		var method = &methodMap{}

		actorMethodT := reflect.TypeOf(actorMethod)
		actorFunc := actorMethod
		if actorMethodT.Kind() == reflect.Struct {
			methodValue := reflect.ValueOf(actorMethod)
			method.AliasName = methodValue.FieldByName("Name").String()
			actorFunc = methodValue.FieldByName("Func").Interface()
		}

		functionT := reflect.TypeOf(actorFunc)
		if functionT.Kind() != reflect.Func {
			return nil, fmt.Errorf("export must be function ")
		}
		method.FuncT = reflect.ValueOf(actorFunc)
		method.PkgName, method.FuncName = getFunctionName(method.FuncT)
		//	functionT := function.Type

		if functionT.NumIn() > 0 {
			isCtx := isContext(functionT.In(0))
			if isCtx == true {
				method.HasContext = true
			}
		}
		if (method.HasContext == false && functionT.NumIn() > 1) || (method.HasContext == true && functionT.NumIn() > 2) {
			return nil, fmt.Errorf("func %v can not have params more than 1", method.FuncName)
		}
		if (method.HasContext == false && functionT.NumIn() == 1) || (method.HasContext == true && functionT.NumIn() == 2) {
			if method.HasContext == true && functionT.NumIn() == 2 {
				if !functionT.In(1).AssignableTo(unMarshallerT) {
					return nil, fmt.Errorf("func %v arg type at index 2 must can be unmarshaller", method)
				}
				method.HasParam = true
				method.ParamsType = functionT.In(1)
				hasParam = true
				typesToImport = append(typesToImport, functionT.In(1))

			} else {
				if !functionT.In(0).AssignableTo(unMarshallerT) {
					return nil, fmt.Errorf("func %v arg type at index 1 must can be unmarshaller", method.FuncName)
				}

				method.HasParam = true
				method.ParamsType = functionT.In(0)
				hasParam = true
				typesToImport = append(typesToImport, functionT.In(0))
			}
		}

		if functionT.NumOut() > 2 {
			return nil, fmt.Errorf("func %v can not have return value more than 2", method.FuncName)
		}

		if functionT.NumOut() == 2 {
			if !functionT.Out(0).AssignableTo(marshallerT) {
				return nil, fmt.Errorf("func %v return value at index 0 must be marshaller", method.FuncName)
			}

			if !functionT.Out(1).AssignableTo(errorT) {
				return nil, fmt.Errorf("func %v return value at index 1 must be error", method.FuncName)
			}
			method.HasReturn = true
			method.HasError = true
			method.ReturnType = functionT.Out(0)
			typesToImport = append(typesToImport, functionT.Out(0))
		} else if functionT.NumOut() == 1 {
			if functionT.Out(0).AssignableTo(errorT) {
				method.HasReturn = false
				method.HasError = true
			} else {
				typesToImport = append(typesToImport, functionT.Out(0))
				method.ReturnType = functionT.Out(0)
				method.HasReturn = true
				method.HasError = false
			}
		} else {
			//no return
			method.HasReturn = false
			method.HasError = false
		}

		methodsArr = append(methodsArr, method)

	}

	//resolve package and name
	imports := defaultClientImport
	for _, importType := range typesToImport {
		imports = append(imports, ImportsForType(pkg, importType)...)
	}
	imports = dedupImports(imports)

	stateName := typeName(pkg, stateT)
	for _, m := range methodsArr {
		m.StateName = stateName
		if m.HasParam {
			m.ParamsTypeName = typeName(pkg, m.ParamsType)
		}
		if m.HasReturn {
			m.ReturnTypeName = typeName(pkg, m.ReturnType)
			m.DefaultReturn = defaultValue(pkg, m.ReturnType)
		}
		name := m.AliasName
		if len(name) == 0 {
			name = m.FuncName
		}
		hashNumber, err := frc42dispatch.GenMethodNumber(name)
		if err != nil {
			return nil, fmt.Errorf("function name %s not validate, may need change another", m.FuncName)
		}
		m.MethodNum = uint64(hashNumber)
		fmt.Println("Method:", m.FuncName, " MethodNumber: ", hashNumber)
	}

	return &entryMeta{
		Imports: dedupImports(imports),
		//PkgName:   ,
		HasParam:  hasParam,
		Methods:   methodsArr,
		StateName: stateName,
	}, nil
}

func typeName(pkg string, t reflect.Type) string {
	switch t.Kind() {
	case reflect.Array:
		if len(t.Name()) > 0 {
			return resolveTypeName(t, pkg)
		}
		return fmt.Sprintf("[%d]%s", t.Len(), typeName(pkg, t.Elem()))
	case reflect.Slice:
		if len(t.Name()) > 0 {
			return resolveTypeName(t, pkg)
		}
		return "[]" + typeName(pkg, t.Elem())
	case reflect.Ptr:
		return "*" + typeName(pkg, t.Elem())
	case reflect.Map:
		return "map[" + typeName(pkg, t.Key()) + "]" + typeName(pkg, t.Elem())
	default:
		return resolveTypeName(t, pkg)
	}
}

func resolveTypeName(t reflect.Type, pkg string) string {
	pkgPath := t.PkgPath()
	if pkgPath == "" {
		// It's a built-in.
		return t.String()
	} else if pkgPath == pkg {
		return t.Name()
	}
	return fmt.Sprintf("%s.%s", resolvePkgByFullName(pkgPath, t.String()), t.Name())
}

func resolvePkgByFullName(path, typeName string) string {
	parts := strings.Split(typeName, ".")
	if len(parts) != 2 {
		panic(fmt.Sprintf("expected type to have a package name: %s", typeName))
	}
	defaultName := parts[0]
	return resolvePkgNameByFullPath(path, defaultName)
}

func resolvePkgNameByFullPath(path, defaultName string) string {
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

func ImportsForType(currPkg string, t reflect.Type) []typegen.Import {
	switch t.Kind() {
	case reflect.Array, reflect.Slice, reflect.Ptr:
		return ImportsForType(currPkg, t.Elem())
	case reflect.Map:
		return dedupImports(append(ImportsForType(currPkg, t.Key()), ImportsForType(currPkg, t.Elem())...))
	default:
		path := t.PkgPath()
		if path == "" || path == currPkg {
			// built-in or in current package.
			return nil
		}

		return []typegen.Import{{PkgPath: path, Name: resolvePkgByFullName(path, t.String())}}
	}
}

type entryMeta struct {
	Imports   []typegen.Import
	HasParam  bool
	PkgName   string
	Methods   []*methodMap
	StateName string
	StateType reflect.Type
}

type methodMap struct {
	StateName  string
	MethodNum  uint64
	FuncT      reflect.Value
	PkgName    string
	FuncName   string
	AliasName  string
	HasError   bool
	HasParam   bool
	HasReturn  bool
	HasContext bool

	ParamsType     reflect.Type
	ParamsTypeName string

	ReturnType     reflect.Type
	ReturnTypeName string
	DefaultReturn  string
}

var tml = `// Code generated by github.com/ipfs-force-community/go-fvm-sdk. DO NOT EDIT.
package main

import (
	"bytes"
	"fmt"
	{{range .Imports}}
	 {{.Name}} "{{.PkgPath}}"
	{{end}}
)

// not support non-main wasm in tinygo at present
func main() {}

// Invoke The actor's WASM entrypoint. It takes the ID of the parameters block,
// and returns the ID of the return value block, or NO_DATA_BLOCK_ID if no
// return value.
//
// Should probably have macros similar to the ones on fvm.filecoin.io snippets.
// Put all methods inside an impl struct and annotate it with a derive macro
// that handles state serde and dispatch.
//
//go:export invoke
func Invoke(blockId uint32) uint32 {
	ctx:=context.Background()
	method, err := sdk.MethodNumber(ctx)
	if err != nil {
		sdk.Abort(ctx,ferrors.USR_ILLEGAL_STATE, "unable to get method number")
	}

	var callResult cbor.Marshaler
{{if .HasParam}}var raw *sdkTypes.ParamsRaw{{end}}
	switch method {
{{range .Methods}}case {{.MethodNum|hex}}:
{{if eq .MethodNum 1}}  // Constuctor
		{{if .HasParam}}raw, err = sdk.ParamsRaw(ctx,blockId)
						if err != nil {
							sdk.Abort(ctx,ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
						}
						var req {{trimPrefix .ParamsTypeName "*"}}
						err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
						if err != nil {
							sdk.Abort(ctx,ferrors.USR_ILLEGAL_STATE, "unable to unmarshal params raw")
						}
						err = {{.PkgName}}.{{.FuncName}}({{if .HasContext}} ctx, {{end}}&req)
						callResult = typegen.CborBool(true)
          {{else}}err = {{.PkgName}}.{{.FuncName}}({{if .HasContext}} ctx {{end}})
                callResult = typegen.CborBool(true)
          {{end}}
{{else}}
		  {{if .HasParam}}raw, err = sdk.ParamsRaw(ctx,blockId)
								if err != nil {
									sdk.Abort(ctx,ferrors.USR_ILLEGAL_STATE, "unable to read params raw")
								}
								var req {{trimPrefix .ParamsTypeName "*"}}
								err = req.UnmarshalCBOR(bytes.NewReader(raw.Raw))
								if err != nil {
									sdk.Abort(ctx,ferrors.USR_ILLEGAL_STATE, "unable to unmarshal params raw")
								}
       		 {{if .HasError}}
					 {{if .HasReturn}} // have params/return/error
								state := new({{.StateName}})
								sdk.LoadState(ctx,state)
								callResult, err = state.{{.FuncName}}({{if .HasContext}} ctx, {{end}}&req)
				     {{else}} 	// have params/error but no return val
								state := new({{.StateName}})
								sdk.LoadState(ctx,state)
								if err = state.{{.FuncName}}({{if .HasContext}} ctx, {{end}}&req); err == nil {
									callResult = typegen.CborBool(true)
								}
					{{end}}
			{{else}}
					{{if .HasReturn}}// have params/return but no error
							state := new({{.StateName}})
							sdk.LoadState(ctx,state)
							callResult = state.{{.FuncName}}({{if .HasContext}} ctx, {{end}}&req)
					{{else}}//have params but no return value and error
							state := new({{.StateName}})
							sdk.LoadState(ctx,state)
							state.{{.FuncName}}({{if .HasContext}} ctx, {{end}}&req)
							callResult = = typegen.CborBool(true)
					{{end}}
			{{end}}
    {{else}}
			{{if .HasError}}
					 {{if .HasReturn}} // no params but return value/error
							state := new({{.StateName}})
							sdk.LoadState(ctx,state)
							callResult, err = state.{{.FuncName}}({{if .HasContext}} ctx {{end}})
					{{else}}	// no params/return value but return error
							state := new({{.StateName}})
							sdk.LoadState(ctx,state)
							if err = state.{{.FuncName}}({{if .HasContext}} ctx {{end}}); err == nil {
									callResult = = typegen.CborBool(true)
								}
					{{end}}
			{{else}}
					{{if .HasReturn}}	// no params no error but have return value
						state := new({{.StateName}})
						sdk.LoadState(ctx,state)
						callResult = state.{{.FuncName}}({{if .HasContext}} ctx {{end}})
					{{else}}		// no params/return value/error
						state := new({{.StateName}})
						sdk.LoadState(ctx,state)
						state.{{.FuncName}}({{if .HasContext}} ctx {{end}})
						callResult = = typegen.CborBool(true)
					{{end}}
			{{end}}
    {{end}}
{{end}}
{{end}}
	default:
		sdk.Abort(ctx,ferrors.USR_ILLEGAL_STATE, "unsupport method")
	}

	if err != nil {
		exitCode := ferrors.USR_ILLEGAL_STATE
		errors.As(err, &exitCode)
		sdk.Abort(ctx, exitCode, fmt.Sprintf("call error %s", err))
	}

	if !sdk.IsNil(callResult) {
		buf := bytes.NewBufferString("")
		err = callResult.MarshalCBOR(buf)
		if err != nil {
			sdk.Abort(ctx,ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("marshal resp fail %s", err))
		}
		id, err := sdk.PutBlock(ctx,sdkTypes.DAGCbor, buf.Bytes())
		if err != nil {
			sdk.Abort(ctx,ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to store return value: %v", err))
		}
		return id
	} else {
		return sdkTypes.NoDataBlockID
	}
}

`

func defaultValue(pkg string, t reflect.Type) string {
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
			pkg := resolvePkgNameByFullPath(t.PkgPath(), t.Name())
			fmt.Println(pkg)
			return fmt.Sprintf("%s.%s{}", pkg, t.Name())
		}
	default:
		panic("unsupprt type")
	}
}

// hellocontract/contract.Constructor
// hellocontract/contract.(*State).SayHello
func getFunctionName(temp reflect.Value) (string, string) {
	fullName := runtime.FuncForPC(temp.Pointer()).Name()
	fullName = strings.TrimSuffix(fullName, "-fm")

	split := strings.Split(fullName, ".")
	name := split[len(split)-1]

	split2 := strings.Split(split[0], "/")
	pkgName := split2[len(split2)-1]
	return resolvePkgNameByFullPath(split[0], pkgName), name
}
