package api

import (
	"context"
	"encoding/json"
	"fmt"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
)

type CalcRequest struct {
	Method   string `json:"method"`
	Operands []int  `json:"operands"`
}

func ParseCalcRequest(ctx context.Context, body []byte) (CalcRequest, error) {
	tracer := global.TraceProvider().Tracer("calcRequest")
	var parsedRequest CalcRequest
	tracer.Start(ctx, "parse")
	trace.CurrentSpan(ctx).AddEvent(ctx, "attempting to parse body")
	trace.CurrentSpan(ctx).AddEvent(ctx, fmt.Sprintf("%s", body))
	err := json.Unmarshal(body, &parsedRequest)
	if err != nil {
		trace.CurrentSpan(ctx).AddEvent(ctx, err.Error())
		trace.CurrentSpan(ctx).End()
		return parsedRequest, err
	}
	trace.CurrentSpan(ctx).End()
	return parsedRequest, nil
}
