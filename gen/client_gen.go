package gen

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
		Name:    "types",
		PkgPath: "github.com/filecoin-project/venus/venus-shared/types",
	},
	{
		Name:    "sdkTypes",
		PkgPath: "github.com/ipfs-force-community/go-fvm-sdk/sdk/types",
	},
	{
		Name:    "cid",
		PkgPath: "github.com/ipfs/go-cid",
	},
	{
		Name:    "abi",
		PkgPath: "github.com/filecoin-project/go-state-types/abi",
	},
	{
		Name:    "typegen",
		PkgPath: "github.com/whyrusleeping/cbor-gen",
	},
	{
		Name:    "cbor",
		PkgPath: "github.com/filecoin-project/go-state-types/cbor",
	},
	{
		Name:    "sdk",
		PkgPath: "github.com/ipfs-force-community/go-fvm-sdk/sdk",
	},
	{
		Name:    "ferrors",
		PkgPath: "github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors",
	},
	{
		Name:    "context",
		PkgPath: "context",
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

func GenContractClient(stateT reflect.Type, output string) error {
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

	if err = genClientConstructor(buf, *entryMeta); err != nil {
		return err
	}

	//write implement
	for _, method := range entryMeta.Methods {
		if method.FuncName == "Constructor" {
			//skip constructor function, because this only called by init actor
			continue
		}
		if method.HasReturn {
			if method.HasParam {
				if err = genClientParamsReturnMethod(buf, method); err != nil {
					return err
				}
			} else {
				if err = genClientNoParamsReturnMethod(buf, method); err != nil {
					return err
				}
			}
		} else {
			if method.HasParam {
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
	tpl := `// Code generated by github.com/ipfs-force-community/go-fvm-sdk. DO NOT EDIT.
package client

import ({{range .}}{{.Name}} "{{.PkgPath}}"
{{end}}
v0 "github.com/filecoin-project/venus/venus-shared/api/chain/v0"
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
	node v0.FullNode
	cfg ClientOption
}


// Option option func
type Option func(opt *ClientOption)
type SendOption func(opt *ClientOption)

// ClientOption option for set client config
type ClientOption struct {
	fromAddress address.Address
	actor       address.Address
	codeCid     cid.Cid
}

// SetFromAddressOpt used to set from address who send actor messages
func SetFromAddressOpt(fromAddress address.Address) Option {
	return func(opt *ClientOption) {
		opt.fromAddress = fromAddress
	}
}

// SetFromAddrSendOpt used to set from address who send actor messages
func SetFromAddrSendOpt(fromAddress address.Address) SendOption {
	return func(opt *ClientOption) {
		opt.fromAddress = fromAddress
	}
}

// SetActorOpt used to set actor address
func SetActorOpt(actor address.Address) Option {
	return func(opt *ClientOption) {
		opt.actor = actor
	}
}

// SetCodeCid used to set actor code cid
func SetCodeCid(codeCid cid.Cid) Option {
	return func(opt *ClientOption) {
		opt.codeCid = codeCid
	}
}

func New{{trimPackage .StateName}}Client(fullNode v0.FullNode, opts ...Option) *{{trimPackage .StateName}}Client {
	cfg := ClientOption{}
	for _, opt := range opts {
		opt(&cfg)
	}
	return &{{trimPackage .StateName}}Client{
		node: fullNode,
		cfg: cfg,
	}
}

func (c *{{trimPackage .StateName}}Client) ChangeFromAddress(addr address.Address) {
	c.cfg.fromAddress = addr
}

func (c *{{trimPackage .StateName}}Client) Install(ctx context.Context, code []byte, opts ...SendOption) (*sdkTypes.InstallReturn, error) {
	cfg_copy := c.cfg
	for _, opt := range opts {
		opt(&cfg_copy)
	}
	params, aerr := actors.SerializeParams(&sdkTypes.InstallParams{
		Code: code,
	})
	if aerr != nil {
		return nil, fmt.Errorf("failed to serialize params: %w", aerr)
	}

	msg := &types.Message{
		To:     builtin.InitActorAddr,
		From:   cfg_copy.fromAddress,
		Value:  big.Zero(),
		Method: 4,
		Params: params,
	}

	smsg, err := c.node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.node.StateWaitMsg(ctx, smsg.Cid(), 0)
	if err != nil {
		return nil, fmt.Errorf("error waiting for message: %w", err)
	}

	// check it executed successfully
	if wait.Receipt.ExitCode != 0 {
		return nil, fmt.Errorf("actor installation failed")
	}

	var result sdkTypes.InstallReturn
	r := bytes.NewReader(wait.Receipt.Return)
	if err := result.UnmarshalCBOR(r); err != nil {
		return nil, fmt.Errorf("error unmarshaling return value: %w", err)
	}
	c.cfg.codeCid = result.CodeCid
	return &result, nil
}
`

	render, err := template.New("gen client interface").Funcs(funcs).Parse(tpl)
	if err != nil {
		return err
	}

	return render.Execute(w, meta)
}

func genClientConstructor(w io.Writer, meta entryMeta) error {
	tpl := `func (c *{{trimPackage .StateName}}Client) CreateActor(ctx context.Context, codeCid cid.Cid{{if .HasParam}}, ctorReq {{.ParamsTypeName}}{{else}}{{end}}, opts ...SendOption) (*init8.ExecReturn, error) {
	cfg_copy := c.cfg
	for _, opt := range opts {
		opt(&cfg_copy)
	}
    {{if .HasParam}}buf := bytes.NewBufferString("")
	if err := ctorReq.MarshalCBOR(buf); err != nil {
		return nil, err
	}
	params, aErr := actors.SerializeParams(&init8.ExecParams{
		CodeCID:           codeCid,
		ConstructorParams: buf.Bytes(),
	})
	{{else}}
	params, aErr := actors.SerializeParams(&init8.ExecParams{
		CodeCID:           codeCid,
	})
	{{end}}

	if aErr != nil {
		return nil, fmt.Errorf("failed to serialize params: %w", aErr)
	}

	msg := &types.Message{
		To:     builtin.InitActorAddr,
		From:   cfg_copy.fromAddress,
		Value:  big.Zero(),
		Method: 2,
		Params: params,
	}

	smsg, err := c.node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.node.StateWaitMsg(ctx, smsg.Cid(), 0)
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
	c.cfg.actor = result.IDAddress
	return &result, nil
}
`

	render, err := template.New("gen client interface").Funcs(funcs).Parse(tpl)
	if err != nil {
		return err
	}

	var ctrFunc *methodMap
	for _, method := range meta.Methods {
		if method.FuncName == "Constructor" {
			ctrFunc = method
		}
	}
	return render.Execute(w, ctrFunc)
}

func genClientInterface(w io.Writer, entry entryMeta) error {
	tpl := `
type I{{trimPackage .StateName}}Client interface {
	Install( context.Context,  []byte, ...SendOption) (*sdkTypes.InstallReturn, error)
	{{range .Methods}}
    {{if ne .FuncName "Constructor"}}
		{{if .HasParam}}
			{{if .HasReturn}}
						{{.FuncName}}(context.Context, {{.ParamsTypeName}}, ...SendOption) ({{.ReturnTypeName}}, error)
			{{else}}
						{{.FuncName}}(context.Context, {{.ParamsTypeName}}, ...SendOption) error
			{{end}}
		{{else}}
			{{if .HasReturn}}
				{{.FuncName}}(context.Context, ...SendOption) ({{.ReturnTypeName}}, error)
			{{else}}
				{{.FuncName}}(context.Context, ...SendOption) error
			{{end}}
		{{end}}
	{{end}}
    {{if eq .FuncName "Constructor"}}
		CreateActor( context.Context,  cid.Cid{{if .HasParam}}, {{.ParamsTypeName}}{{else}}{{end}}, ...SendOption) (*init8.ExecReturn, error)
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

func genClientParamsReturnMethod(w io.Writer, entry *methodMap) error {
	tpl := `
func (c *{{trimPackage .StateName}}Client) {{.FuncName}}(ctx context.Context, p0 {{.ParamsTypeName}}, opts ...SendOption) ({{.ReturnTypeName}}, error) {
	cfg_copy := c.cfg
	for _, opt := range opts {
		opt(&cfg_copy)
	}

	if c.cfg.actor == address.Undef {
		return {{.DefaultReturn|raw}}, fmt.Errorf("unset actor address for call")
	}

	buf := bytes.NewBufferString("")
	if err := p0.MarshalCBOR(buf); err != nil {
		return {{.DefaultReturn|raw}}, err
	}
	msg := &types.Message{
		To:     cfg_copy.actor,
		From:   cfg_copy.fromAddress,
		Value:  big.Zero(),
		Method: abi.MethodNum({{.MethodNum}}),
		Params: buf.Bytes(),
	}

	smsg, err := c.node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return {{.DefaultReturn|raw}}, fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.node.StateWaitMsg(ctx, smsg.Cid(), 0)
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

func genClientParamsNOReturnMethod(w io.Writer, entry *methodMap) error {
	tpl := `
func (c *{{trimPackage .StateName}}Client) {{.FuncName}}(ctx context.Context, p0 {{.ParamsTypeName}}, opts ...SendOption) error {
	cfg_copy := c.cfg
	for _, opt := range opts {
		opt(&cfg_copy)
	}

	if c.cfg.actor == address.Undef {
		return  fmt.Errorf("unset actor address for call")
	}

	buf := bytes.NewBufferString("")
	if err := p0.MarshalCBOR(buf); err != nil {
		return  err
	}
	msg := &types.Message{
		To:     cfg_copy.actor,
		From:   cfg_copy.fromAddress,
		Value:  big.Zero(),
		Method: abi.MethodNum({{.MethodNum}}),
		Params: buf.Bytes(),
	}

	smsg, err := c.node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.node.StateWaitMsg(ctx, smsg.Cid(), 0)
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

func genClientNoParamsReturnMethod(w io.Writer, entry *methodMap) error {
	tpl := `
func (c *{{trimPackage .StateName}}Client) {{.FuncName}}(ctx context.Context, opts ...SendOption) ({{.ReturnTypeName}}, error) {
	cfg_copy := c.cfg
	for _, opt := range opts {
		opt(&cfg_copy)
	}
	if c.cfg.actor == address.Undef {
		return {{.DefaultReturn|raw}}, fmt.Errorf("unset actor address for call")
	}

	msg := &types.Message{
		To:     cfg_copy.actor,
		From:   cfg_copy.fromAddress,
		Value:  big.Zero(),
		Method: abi.MethodNum({{.MethodNum}}),
		Params: nil,
	}

	smsg, err := c.node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return {{.DefaultReturn|raw}}, fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.node.StateWaitMsg(ctx, smsg.Cid(), 0)
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

func genClientNoParamsNoReturnMethod(w io.Writer, entry *methodMap) error {
	tpl := `
func (c *{{trimPackage .StateName}}Client) {{.FuncName}}(ctx context.Context, opts ...SendOption) error {
	cfg_copy := c.cfg
	for _, opt := range opts {
		opt(&cfg_copy)
	}
	if c.cfg.actor == address.Undef {
		return fmt.Errorf("unset actor address for call")
	}

	msg := &types.Message{
		To:     cfg_copy.actor,
		From:   cfg_copy.fromAddress,
		Value:  big.Zero(),
		Method: abi.MethodNum({{.MethodNum}}),
		Params: nil,
	}

	smsg, err := c.node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.node.StateWaitMsg(ctx, smsg.Cid(), 0)
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
