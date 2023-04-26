package ddnethttp

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// SpanTags is tags of span.
type SpanTags map[string]interface{}

// DDTraceResult provides result.
type DDTraceResult func() (SpanTags, error)

// FromContext - datadog from context
func FromContext(ctx *gin.Context) ddtrace.Span {
	span, ok := tracer.SpanFromContext(ctx.Request.Context())
	if !ok {
		return nil
	}
	return span
}

// StartChildSpan child span
func StartChildSpan(ctx *gin.Context, operationName string, tags SpanTags) ddtrace.Span {
	txn := FromContext(ctx)
	if txn == nil {
		return nil
	}
	return StartDDSpan(operationName, txn, "", tags)
}

// StartDDSpan starts a datadog span.
func StartDDSpan(operationName string, parentSpan tracer.Span, spanType string, tags SpanTags) tracer.Span {
	var span tracer.Span
	if parentSpan != nil {
		span = tracer.StartSpan(operationName, tracer.ChildOf(parentSpan.Context()))
	} else {
		span = tracer.StartSpan(operationName)
	}
	if len(spanType) > 0 {
		tags[ext.SpanType] = spanType
	}
	setSpanTags(span, tags)
	return span
}

// EndSpan finishes a datadog span.
func EndSpan(span tracer.Span) {
	if !isNil(span) {
		span.Finish()
	}
}

// EndSpanError finishes a datadog span.
func EndSpanError(span tracer.Span, e error) {
	if !isNil(span) && e != nil {
		span.Finish(tracer.WithError(e))
	}
	if !isNil(span) && e == nil {
		span.Finish()
	}
}

// EndSpanTags finishes a datadog span.
func EndSpanTags(span tracer.Span, tags SpanTags) {
	setSpanTags(span, tags)
	if !isNil(span) {
		span.Finish()
	}
}

// EndSpanTagsError finishes a datadog span.
func EndSpanTagsError(span tracer.Span, tags SpanTags, e error) {
	setSpanTags(span, tags)
	if !isNil(span) && e != nil {
		span.Finish(tracer.WithError(e))
	}
	if !isNil(span) && e == nil {
		span.Finish()
	}
}

func setSpanTags(span tracer.Span, tags SpanTags) {
	if isNil(span) {
		return
	}
	if len(tags) > 0 {
		for k, v := range tags {
			span.SetTag(k, v)
		}
	}
}

func isNil(value interface{}) bool {
	if value == nil {
		return true
	}
	if reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil() {
		return true
	}
	if reflect.ValueOf(value).Kind() == reflect.Invalid {
		return true
	}
	return false
}
