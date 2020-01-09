package api

import (
	"context"
	"encoding/json"
	"fmt"

	"go.opentelemetry.io/otel/api/trace"
	"google.golang.org/grpc/codes"
)

type CalcRequest struct {
	Method   string `json:"method"`
	Operands []int  `json:"operands"`
}

func ParseCalcRequest(ctx context.Context, body []byte) (CalcRequest, error) {
	var parsedRequest CalcRequest

	trace.CurrentSpan(ctx).AddEvent(ctx, "attempting to parse body")
	trace.CurrentSpan(ctx).AddEvent(ctx, fmt.Sprintf("%s", body))
	err := json.Unmarshal(body, &parsedRequest)
	if err != nil {
		trace.CurrentSpan(ctx).SetStatus(codes.InvalidArgument)
		trace.CurrentSpan(ctx).AddEvent(ctx, err.Error())
		trace.CurrentSpan(ctx).End()
		return parsedRequest, err
	}
	trace.CurrentSpan(ctx).End()
	return parsedRequest, nil
}
