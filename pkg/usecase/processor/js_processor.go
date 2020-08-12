package processor

import (
	"github.com/swipe-io/swipe/pkg/domain/model"
	"github.com/swipe-io/swipe/pkg/importer"
	ug "github.com/swipe-io/swipe/pkg/usecase/generator"
)

type jsClient struct {
	info      model.GenerateInfo
	option    model.ServiceOption
	importers map[string]*importer.Importer
}

func (p *jsClient) SetOption(option interface{}) bool {
	o, ok := option.(model.ServiceOption)
	p.option = o
	return ok
}

func (p *jsClient) Generators() []ug.Generator {
	generators := []ug.Generator{
		ug.NewJsonrpcMarkdownDoc(p.info, p.option),
	}
	if p.option.Transport.Protocol == "http" && p.option.Transport.Client.Enable && p.option.Transport.JsonRPC.Enable {
		generators = append(
			generators,
			ug.NewJsonRPCJSClient("client_jsonrpc_gen.js", p.info, p.option),
		)
	}
	if p.option.Transport.Openapi.Enable {
		generators = append(generators, ug.NewOpenapi(p.info, p.option))
	}
	return generators
}

func NewJsClient(info model.GenerateInfo) Processor {
	return &jsClient{info: info, importers: map[string]*importer.Importer{}}
}
