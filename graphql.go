package graphql

import (
	"context"

	"github.com/quan-to/graphql/gqlerrors"
	"github.com/quan-to/graphql/language/parser"
	"github.com/quan-to/graphql/language/source"
)

type Params struct {
	// The GraphQL type system to use when validating and executing a query.
	Schema Schema

	// A GraphQL language formatted string representing the requested operation.
	RequestString string

	// The value provided as the first argument to resolver functions on the top
	// level type (e.g. the query object type).
	RootObject map[string]interface{}

	// A mapping of variable name to runtime value to use for all variables
	// defined in the requestString.
	VariableValues map[string]interface{}

	// The name of the operation to use if requestString contains multiple
	// possible operations. Can be omitted if requestString contains only
	// one operation.
	OperationName string

	// Context may be provided to pass application-specific per-request
	// information to resolve functions.
	Context context.Context

	// CustomErrorFormatter may be a function that processes errors individually
	// to resolve fields or to return different content
	CustomErrorFomatter func(err error) gqlerrors.FormattedError
}

func Do(p Params) *Result {
	src := source.NewSource(&source.Source{
		Body: []byte(p.RequestString),
		Name: "GraphQL request",
	})
	AST, err := parser.Parse(parser.ParseParams{Source: src})
	if err != nil {
		return &Result{
			Errors: gqlerrors.FormatErrorsFunc(p.CustomErrorFomatter, err),
		}
	}
	validationResult := ValidateDocument(&p.Schema, AST, nil)

	if !validationResult.IsValid {
		return &Result{
			Errors: validationResult.Errors,
		}
	}

	return Execute(ExecuteParams{
		Schema:               p.Schema,
		Root:                 p.RootObject,
		AST:                  AST,
		OperationName:        p.OperationName,
		Args:                 p.VariableValues,
		Context:              p.Context,
		CustomErrorFormatter: p.CustomErrorFomatter,
	})
}
