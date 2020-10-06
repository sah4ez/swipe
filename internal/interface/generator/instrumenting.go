package generator

import (
	"context"
	stdtypes "go/types"
	"strconv"

	"github.com/swipe-io/swipe/v2/internal/usecase/generator"

	"github.com/swipe-io/swipe/v2/internal/domain/model"
	"github.com/swipe-io/swipe/v2/internal/importer"
	"github.com/swipe-io/swipe/v2/internal/types"
	"github.com/swipe-io/swipe/v2/internal/writer"
)

type instrumentingGenerator struct {
	writer.GoLangWriter
	serviceID      string
	serviceType    stdtypes.Type
	serviceMethods []model.ServiceMethod
	methodOptions  map[string]model.MethodHTTPTransportOption
	i              *importer.Importer
}

func (g *instrumentingGenerator) Prepare(ctx context.Context) error {
	return nil
}

func (g *instrumentingGenerator) Process(ctx context.Context) error {
	var (
		timePkg string
	)
	if len(g.serviceMethods) > 0 {
		timePkg = g.i.Import("time", "time")
	}
	metricsPkg := g.i.Import("metrics", "github.com/go-kit/kit/metrics")
	typeStr := stdtypes.TypeString(g.serviceType, g.i.QualifyPkg)
	stdPrometheusPkg := g.i.Import("prometheus", "github.com/prometheus/client_golang/prometheus")
	kitPrometheusPkg := g.i.Import("prometheus", "github.com/go-kit/kit/metrics/prometheus")

	name := "instrumentingMiddleware" + g.serviceID
	constructName := "NewInstrumentingMiddleware" + g.serviceID

	g.W("type InstrumentingOption func(*%s)\n\n", name)

	g.W("func Namespace(v string) InstrumentingOption {\nreturn func(i *instrumentingMiddlewareSwipe) {\ni.namespace = v\n}\n}\n\n")
	g.W("func Subsystem(v string) InstrumentingOption {\nreturn func(i *instrumentingMiddlewareSwipe) {\ni.subsystem = v\n}\n}\n\n")

	g.W("func RequestLatency(requestLatency %s.Histogram) InstrumentingOption {\nreturn func(i *instrumentingMiddlewareSwipe) {\ni.requestLatency = requestLatency\n}\n}\n\n", metricsPkg)
	g.W("func RequestCount(requestCount %s.Counter) InstrumentingOption {\nreturn func(i *instrumentingMiddlewareSwipe) {\ni.requestCount = requestCount\n}\n}\n\n", metricsPkg)

	g.W("type %s struct {\n", name)
	g.W("next %s\n", typeStr)
	g.W("requestCount %s.Counter\n", metricsPkg)
	g.W("requestLatency %s.Histogram\n", metricsPkg)
	g.W("namespace string\n")
	g.W("subsystem string\n")

	g.W("}\n\n")

	for _, m := range g.serviceMethods {
		mopt := g.methodOptions[m.Name]

		var params []string

		if m.ParamCtx != nil {
			params = append(params, m.ParamCtx.Name(), stdtypes.TypeString(m.ParamCtx.Type(), g.i.QualifyPkg))
		}

		params = append(params, types.NameTypeParams(m.Params, g.i.QualifyPkg, nil)...)
		results := types.NameTypeParams(m.Results, g.i.QualifyPkg, nil)

		if m.ReturnErr != nil {
			results = append(results, "", "error")
		}

		g.WriteFunc(m.Name, "s *"+name, params, results, func() {
			if mopt.InstrumentingEnable {
				g.WriteDefer(
					[]string{"begin " + timePkg + ".Time"},
					[]string{timePkg + ".Now()"},
					func() {
						g.W("s.requestCount.With(\"method\", \"%s\").Add(1)\n", m.Name)
						g.W("s.requestLatency.With(\"method\", \"%s\").Observe(%s.Since(begin).Seconds())\n", m.Name, timePkg)
					},
				)
			}

			if len(m.Results) > 0 || m.ReturnErr != nil {
				g.W("return ")
			}
			g.W("s.next.%s(", m.Name)
			if m.ParamCtx != nil {
				g.W("%s,", m.ParamCtx.Name())
			}

			for i, p := range m.Params {
				if i > 0 {
					g.W(",")
				}
				g.W(p.Name())
			}

			g.W(")\n")
		})
	}

	g.W("func %[1]s(s %[2]s, opts ...InstrumentingOption) %[2]s {\n", constructName, typeStr, metricsPkg)

	g.W("i := &%s{next: s}\n", name)

	g.W("for _, o := range opts {\no(i)\n}\n")

	g.W("if i.requestCount == nil {\n")
	g.W("i.requestCount = %s.NewCounterFrom(%s.CounterOpts{\n", kitPrometheusPkg, stdPrometheusPkg)
	g.W("Namespace: i.namespace,\n")
	g.W("Subsystem: i.subsystem,\n")
	g.W("Name: %s,\n", strconv.Quote("request_count"))
	g.W("Help: %s,\n", strconv.Quote("Number of requests received."))
	g.W("}, []string{\"method\"})\n")
	g.W("\n}\n")

	g.W("if i.requestLatency == nil {\n")
	g.W("i.requestLatency = %s.NewSummaryFrom(%s.SummaryOpts{\n", kitPrometheusPkg, stdPrometheusPkg)
	g.W("Namespace: i.namespace,\n")
	g.W("Subsystem: i.subsystem,\n")
	g.W("Name: %s,\n", strconv.Quote("request_latency_microseconds"))
	g.W("Help: %s,\n", strconv.Quote("Total duration of requests in microseconds."))
	g.W("}, []string{\"method\"})\n")
	g.W("\n}\n")

	g.W("return i\n}\n")
	return nil
}

func (g *instrumentingGenerator) PkgName() string {
	return ""
}

func (g *instrumentingGenerator) OutputDir() string {
	return ""
}

func (g *instrumentingGenerator) Filename() string {
	return "instrumenting_gen.go"
}

func (g *instrumentingGenerator) SetImporter(i *importer.Importer) {
	g.i = i
}

func NewInstrumenting(
	serviceID string,
	serviceType stdtypes.Type,
	serviceMethods []model.ServiceMethod,
	methodOptions map[string]model.MethodHTTPTransportOption,
) generator.Generator {
	return &instrumentingGenerator{
		serviceID:      serviceID,
		serviceType:    serviceType,
		serviceMethods: serviceMethods,
		methodOptions:  methodOptions,
	}
}
