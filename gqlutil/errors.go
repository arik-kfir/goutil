package gqlutil

import (
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"reflect"
)

type UserError struct {
	err *gqlerror.Error
}

func (e *UserError) Error() string {
	return e.err.Error()
}

func (e *UserError) Is(err error) bool {
	return reflect.TypeOf(err) == reflect.TypeOf(e)
}

func (e *UserError) Unwrap() error {
	return e.err
}

func ErrorPresenter(ctx context.Context, e error) *gqlerror.Error {
	if errors.Is(e, &UserError{}) {
		return e.(*gqlerror.Error)
	}

	var ginCtx *gin.Context
	if gcv := ctx.Value(gin.ContextKey); gcv != nil {
		if gc, ok := gcv.(*gin.Context); ok {
			ginCtx = gc
		}
	}

	path := graphql.GetPath(ctx)
	if ginCtx != nil {
		_ = ginCtx.Error(e).
			SetType(gin.ErrorTypePrivate).
			SetMeta(map[string]interface{}{
				"gqlPath": path,
			})
	} else {
		log.Ctx(ctx).
			Error().
			Stack().
			Err(e).
			Interface("gqlPath", path).
			Msg("Internal error occurred")
	}

	return &gqlerror.Error{
		Message:    "An internal error has occurred.",
		Path:       path,
		Extensions: map[string]interface{}{"code": "INTERNAL_SERVER_ERROR"},
	}
}

func PanicRecoverer(ctx context.Context, p interface{}) error {
	if e, ok := p.(error); ok {
		// We assume stack trace will be retained for the wrapped error
		return fmt.Errorf("recovered panic (error type): %w", e)
	} else {
		// No stack trace will be available, so we format the panic object itself with as much information as possible
		return fmt.Errorf("recovered panic of type '%T': %+v", p, p)
	}
}
