package v1alpha1

import (
	"context"

	"knative.dev/pkg/apis"
)

// Validate implements apis.Validatable
func (as *Fire) Validate(ctx context.Context) *apis.FieldError {
	return as.Spec.Validate(ctx).ViaField("spec")
}

// Validate implements apis.Validatable
func (ass *FireSpec) Validate(ctx context.Context) *apis.FieldError {
	return nil
}
