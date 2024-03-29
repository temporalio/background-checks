package dataconverter

import (
	"context"

	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/workflow"
)

type (
	contextKey   struct{}
	propagator   struct{}
	CryptContext struct {
		KeyID string `json:"keyId"`
	}
)

var PropagateKey = contextKey{}

const propagationKey = "encryption"

// NewContextPropagator returns a context propagator that propagates a set of
// string key-value pairs across a workflow
func NewContextPropagator() workflow.ContextPropagator {
	return &propagator{}
}

// Inject injects values from context into headers for propagation
func (s *propagator) Inject(ctx context.Context, writer workflow.HeaderWriter) error {
	value := ctx.Value(PropagateKey)
	payload, err := converter.GetDefaultDataConverter().ToPayload(value)
	if err != nil {
		return err
	}
	writer.Set(propagationKey, payload)
	return nil
}

// InjectFromWorkflow injects values from context into headers for propagation
func (s *propagator) InjectFromWorkflow(ctx workflow.Context, writer workflow.HeaderWriter) error {
	value := ctx.Value(PropagateKey)
	payload, err := converter.GetDefaultDataConverter().ToPayload(value)
	if err != nil {
		return err
	}
	writer.Set(propagationKey, payload)
	return nil
}

// Extract extracts values from headers and puts them into context
func (s *propagator) Extract(ctx context.Context, reader workflow.HeaderReader) (context.Context, error) {
	if value, ok := reader.Get(propagationKey); ok {
		var cryptContext CryptContext
		if err := converter.GetDefaultDataConverter().FromPayload(value, &cryptContext); err != nil {
			return ctx, nil
		}
		ctx = context.WithValue(ctx, PropagateKey, cryptContext)
	}

	return ctx, nil
}

// ExtractToWorkflow extracts values from headers and puts them into context
func (s *propagator) ExtractToWorkflow(ctx workflow.Context, reader workflow.HeaderReader) (workflow.Context, error) {
	if value, ok := reader.Get(propagationKey); ok {
		var cryptContext CryptContext
		if err := converter.GetDefaultDataConverter().FromPayload(value, &cryptContext); err != nil {
			return ctx, nil
		}
		ctx = workflow.WithValue(ctx, PropagateKey, cryptContext)
	}

	return ctx, nil
}
