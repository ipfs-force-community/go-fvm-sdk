package main

import (
	"bytes"
	"html/template"
	"io"
	"reflect"
	"strings"

	typegen "github.com/whyrusleeping/cbor-gen"
)

var defaultClientImport = []typegen.Import{
	{
		Name:    "bytes",
		PkgPath: "bytes",
	},
	{
		Name:    "context",
		PkgPath: "context",
	},
	{
		Name:    "fmt",
		PkgPath: "fmt",
	},
	{
		Name:    "actors",
		PkgPath: "github.com/filecoin-project/venus/venus-shared/actors",
	},
	{
		Name:    "init8",
		PkgPath: "github.com/filecoin-project/specs-actors/v8/actors/builtin/init",
	},
	{
		Name:    "types2",
		PkgPath: "github.com/ipfs-force-community/go-fvm-sdk/sdk/types",
	},
	{
		Name:    "builtin",
		PkgPath: "github.com/filecoin-project/go-state-types/builtin",
	},
	{
		Name:    "address",
		PkgPath: "github.com/filecoin-project/go-address",
	},
	{
		Name:    "big",
		PkgPath: "github.com/filecoin-project/go-state-types/big",
	},
	{
		Name:    "types",
		PkgPath: "github.com/filecoin-project/venus/venus-shared/types",
	},
	{
		Name:    "cid",
		PkgPath: "github.com/ipfs/go-cid",
	},
	{
		Name:    "abi",
		PkgPath: "github.com/filecoin-project/go-state-types/abi",
	},
}

var funcs = map[string]interface{}{
	"trimPackage": func(str string) string {
		splites := strings.Split(str, ".")
		return splites[len(splites)-1]
	},
	"trimPtr": func(str string) string {
		return strings.Trim(str, "*")
	},
	"isPtr": func(str string) bool {
		return str[0] == '*'
	},
	"raw": func(str string) template.HTML {
		return template.HTML(str)
	},
}

func genContractClient(stateT reflect.Type, output string) error {
	entryMeta, err := getEntryPackageMeta("client", stateT)
	if err != nil {
		return err
	}
	imports := dedupImports(append(defaultClientImport, entryMeta.Imports...))
	buf := bytes.NewBufferString("")
	//write imports
	if err = genClientHeader(buf, imports); err != nil {
		return err
	}
	//write interface
	if err = genClientInterface(buf, *entryMeta); err != nil {
		return err
	}

	//write install/exec
	if err = genClientImplemnt(buf, *entryMeta); err != nil {
		return err
	}
	//write implement
	for _, method := range entryMeta.Methods {
		if method.HasReturn {
			if method.HasParams {
				if err = genClientParamsReturnMethod(buf, method); err != nil {
					return err
				}
			} else {
				if err = genClientNoParamsReturnMethod(buf, method); err != nil {
					return err
				}
			}
		} else {
			if method.HasParams {
				if err = genClientParamsNOReturnMethod(buf, method); err != nil {
					return err
				}
			} else {
				if err = genClientNoParamsNoReturnMethod(buf, method); err != nil {
					return err
				}
			}
		}
	}

	return formateAndWriteCode(buf.Bytes(), output)
}

func genClientHeader(w io.Writer, imports []typegen.Import) error {
	tpl := `
package client

import ({{range .}}{{.Name}} "{{.PkgPath}}"
{{end}}
)

type FullNode interface {
	MpoolPushMessage(ctx context.Context, msg *types.Message, spec *types.MessageSendSpec) (*types.SignedMessage, error)
	StateWaitMsg(ctx context.Context, cid cid.Cid, confidence uint64) (*types.MsgLookup, error)
}
`

	render, err := template.New("gen client header").Funcs(funcs).Parse(tpl)
	if err != nil {
		return err
	}

	return render.Execute(w, imports)
}

func genClientImplemnt(w io.Writer, meta entryMeta) error {
	tpl := `
var _ I{{trimPackage .StateName}}Client = (*{{trimPackage .StateName}}Client)(nil)

type {{trimPackage .StateName}}Client struct {
	Node        FullNode
	FromAddress address.Address
	Actor       address.Address
}

func (c *{{trimPackage .StateName}}Client) CreateActor(ctx context.Context, codeCid cid.Cid, execParams []byte) (*init8.ExecReturn, error) {
	params, aErr := actors.SerializeParams(&init8.ExecParams{
		CodeCID:           codeCid,
		ConstructorParams: execParams,
	})
	if aErr != nil {
		return nil, fmt.Errorf("failed to serialize params: %w", aErr)
	}

	msg := &types.Message{
		To:     builtin.InitActorAddr,
		From:   c.FromAddress,
		Value:  big.Zero(),
		Method: 2,
		Params: params,
	}

	smsg, err := c.Node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.Node.StateWaitMsg(ctx, smsg.Cid(), 0)
	if err != nil {
		return nil, fmt.Errorf("error waiting for message: %w", err)
	}

	// check it executed successfully
	if wait.Receipt.ExitCode != 0 {
		return nil, fmt.Errorf("actor execution failed")
	}

	var result init8.ExecReturn
	r := bytes.NewReader(wait.Receipt.Return)
	if err := result.UnmarshalCBOR(r); err != nil {
		return nil, fmt.Errorf("error unmarshaling return value: %w", err)
	}
	return &result, nil
}

func (c *{{trimPackage .StateName}}Client) Install(ctx context.Context, code []byte) (*init8.InstallReturn, error) {
	params, aerr := actors.SerializeParams(&init8.InstallParams{
		Code: code,
	})
	if aerr != nil {
		return nil, fmt.Errorf("failed to serialize params: %w", aerr)
	}

	msg := &types.Message{
		To:     builtin.InitActorAddr,
		From:   c.FromAddress,
		Value:  big.Zero(),
		Method: 3,
		Params: params,
	}

	smsg, err := c.Node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.Node.StateWaitMsg(ctx, smsg.Cid(), 0)
	if err != nil {
		return nil, fmt.Errorf("error waiting for message: %w", err)
	}

	// check it executed successfully
	if wait.Receipt.ExitCode != 0 {
		return nil, fmt.Errorf("actor installation failed")
	}

	var result init8.InstallReturn
	r := bytes.NewReader(wait.Receipt.Return)
	if err := result.UnmarshalCBOR(r); err != nil {
		return nil, fmt.Errorf("error unmarshaling return value: %w", err)
	}
	return &result, nil
}
`

	render, err := template.New("gen client interface").Funcs(funcs).Parse(tpl)
	if err != nil {
		return err
	}

	return render.Execute(w, meta)
}

func genClientInterface(w io.Writer, entry entryMeta) error {
	tpl := `
type I{{trimPackage .StateName}}Client interface {
	Install( context.Context,  []byte) (*init8.InstallReturn, error)
	CreateActor( context.Context,  cid.Cid,  []byte) (*init8.ExecReturn, error)
	{{range .Methods}}
		{{if .HasParams}}
			{{if .HasReturn}}
						{{.FuncName}}(context.Context, {{.ParamsTypeName}}) ({{.ReturnTypeName}}, error)
			{{else}}
						{{.FuncName}}(context.Context, {{.ParamsTypeName}}) error
			{{end}}
		{{else}}
			{{if .HasReturn}}
				{{.FuncName}}(context.Context) ({{.ReturnTypeName}}, error)
			{{else}}
				{{.FuncName}}(context.Context, {{.ParamsTypeName}}) error
			{{end}}
		{{end}}
	{{end}}
}
`

	render, err := template.New("gen client interface").Funcs(funcs).Parse(tpl)
	if err != nil {
		return err
	}

	return render.Execute(w, entry)
}

func genClientParamsReturnMethod(w io.Writer, entry methodMap) error {
	tpl := `
func (c *{{trimPackage .StateName}}Client) {{.FuncName}}(ctx context.Context, p0 {{.ParamsTypeName}}) ({{.ReturnTypeName}}, error) {
	if c.Actor == address.Undef {
		return {{.DefaultReturn|raw}}, fmt.Errorf("unset actor address for call")
	}

	buf := bytes.NewBufferString("")
	if err := p0.MarshalCBOR(buf); err != nil {
		return {{.DefaultReturn|raw}}, err
	}
	msg := &types.Message{
		To:     c.Actor,
		From:   c.FromAddress,
		Value:  big.Zero(),
		Method: abi.MethodNum({{.MethodNum}}),
		Params: buf.Bytes(),
	}

	smsg, err := c.Node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return {{.DefaultReturn|raw}}, fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.Node.StateWaitMsg(ctx, smsg.Cid(), 0)
	if err != nil {
		return {{.DefaultReturn|raw}}, fmt.Errorf("error waiting for message: %w", err)
	}

	// check it executed successfully
	if wait.Receipt.ExitCode != 0 {
		return {{.DefaultReturn|raw}}, fmt.Errorf("actor execution failed")
	}
	if len(wait.Receipt.Return) == 0 {
		return {{.DefaultReturn|raw}}, fmt.Errorf("expect get result for call")
	}
	
	result := new({{.ReturnTypeName|trimPtr}})
	result.UnmarshalCBOR(bytes.NewReader(wait.Receipt.Return))

	{{if isPtr .ReturnTypeName}} 
			return result, nil
	{{else}}
			return *result, nil
	{{end}}
}
`
	render, err := template.New("gen client interface").Funcs(funcs).Parse(tpl)
	if err != nil {
		return err
	}

	return render.Execute(w, entry)
}

func genClientParamsNOReturnMethod(w io.Writer, entry methodMap) error {
	tpl := `
func (c *{{trimPackage .StateName}}Client) {{.FuncName}}(ctx context.Context, p0 {{.ParamsTypeName}}) error {
	if c.Actor == address.Undef {
		return  fmt.Errorf("unset actor address for call")
	}

	buf := bytes.NewBufferString("")
	if err := p0.MarshalCBOR(buf); err != nil {
		return  err
	}
	msg := &types.Message{
		To:     c.Actor,
		From:   c.FromAddress,
		Value:  big.Zero(),
		Method: abi.MethodNum({{.MethodNum}}),
		Params: buf.Bytes(),
	}

	smsg, err := c.Node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.Node.StateWaitMsg(ctx, smsg.Cid(), 0)
	if err != nil {
		return  fmt.Errorf("error waiting for message: %w", err)
	}

	// check it executed successfully
	if wait.Receipt.ExitCode != 0 {
		return fmt.Errorf("actor execution failed")
	}
	return nil
}
`

	render, err := template.New("gen client interface").Funcs(funcs).Parse(tpl)
	if err != nil {
		return err
	}

	return render.Execute(w, entry)
}

func genClientNoParamsReturnMethod(w io.Writer, entry methodMap) error {
	tpl := `
func (c *{{trimPackage .StateName}}Client) {{.FuncName}}(ctx context.Context) ({{.ReturnTypeName}}, error) {
	if c.Actor == address.Undef {
		return {{.DefaultReturn|raw}}, fmt.Errorf("unset actor address for call")
	}

	msg := &types.Message{
		To:     c.Actor,
		From:   c.FromAddress,
		Value:  big.Zero(),
		Method: abi.MethodNum({{.MethodNum}}),
		Params: nil,
	}

	smsg, err := c.Node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return {{.DefaultReturn|raw}}, fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.Node.StateWaitMsg(ctx, smsg.Cid(), 0)
	if err != nil {
		return {{.DefaultReturn|raw}}, fmt.Errorf("error waiting for message: %w", err)
	}

	// check it executed successfully
	if wait.Receipt.ExitCode != 0 {
		return {{.DefaultReturn|raw}}, fmt.Errorf("actor execution failed")
	}
	if len(wait.Receipt.Return) == 0 {
		return {{.DefaultReturn|raw}}, fmt.Errorf("expect get result for call")
	}

	result := new({{.ReturnTypeName|trimPtr}})
	result.UnmarshalCBOR(bytes.NewReader(wait.Receipt.Return))
	{{if isPtr .ReturnTypeName}} 
			return result, nil
	{{else}}
			return *result, nil
	{{end}}
}
`

	render, err := template.New("gen client interface").Funcs(funcs).Parse(tpl)
	if err != nil {
		return err
	}

	return render.Execute(w, entry)
}

func genClientNoParamsNoReturnMethod(w io.Writer, entry methodMap) error {
	tpl := `
func (c *{{trimPackage .StateName}}Client) {{.FuncName}}(ctx context.Context) error {
	if c.Actor == address.Undef {
		return fmt.Errorf("unset actor address for call")
	}

	msg := &types.Message{
		To:     c.Actor,
		From:   c.FromAddress,
		Value:  big.Zero(),
		Method: abi.MethodNum({{.MethodNum}}),
		Params: nil,
	}

	smsg, err := c.Node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.Node.StateWaitMsg(ctx, smsg.Cid(), 0)
	if err != nil {
		return fmt.Errorf("error waiting for message: %w", err)
	}

	// check it executed successfully
	if wait.Receipt.ExitCode != 0 {
		return fmt.Errorf("actor execution failed")
	}
	return nil
}
`

	render, err := template.New("gen client interface").Funcs(funcs).Parse(tpl)
	if err != nil {
		return err
	}

	return render.Execute(w, entry)
}
