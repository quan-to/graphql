package graphql

import (
	"fmt"
	"github.com/volatiletech/null"
	"log"
)

type PageInfo struct {
	StartCursor     string
	EndCursor       string
	HasNextPage     bool
	HasPreviousPage bool
}

type EdgeData struct {
	Node   interface{}
	Cursor string
}

type ConnectionData struct {
	TotalCount int64
	PageInfo   PageInfo
	Edges      []EdgeData
}

func MakeConnectionArgs(arguments FieldConfigArgument) FieldConfigArgument {
	var out = FieldConfigArgument{
		"First": {
			Type:        Int,
			Description: "Show first N results starting from cursor",
		},
		"Last": {
			Type:        Int,
			Description: "Show last N results starting from cursor",
		},
		"After": {
			Type:        String,
			Description: "Show results after specified cursor",
		},
		"Before": {
			Type:        String,
			Description: "Show results before specified cursor",
		},
	}

	for k, v := range arguments {
		if _, ok := out[k]; ok {
			log.Printf("WARN: There is already a field named %s in argument. Skiping adding it.", k)
		} else {
			out[k] = v
		}
	}

	return out
}

var pageInfoModel = NewObject(ObjectConfig{
	Name: "EdgePageInfo",
	Fields: Fields{
		"StartCursor": {
			Type:        String,
			Description: "Cursor of the first page item",
		},
		"EndCursor": {
			Type:        String,
			Description: "Cursor of the last page item",
		},
		"HasNextPage": {
			Type:        Boolean,
			Description: "If there is a next page",
		},
		"HasPreviousPage": {
			Type:        Boolean,
			Description: "If there is a previous page",
		},
	},
})

func makeGraphQLEdge(nodeType *Object) *Object {
	return NewObject(ObjectConfig{
		Name:        fmt.Sprintf("%sConnectionEdges", nodeType.Name()),
		Description: "Edges for GraphQL Pagination",
		Fields: Fields{
			"Node": {
				Type:        nodeType,
				Description: fmt.Sprintf("Object %s", nodeType.Name()),
			},
			"Cursor": {
				Type:        String,
				Description: "Cursor of the row",
			},
		},
	})
}

func MakeGraphQLConnection(nodeType *Object) *Object {
	return NewObject(ObjectConfig{
		Name:        fmt.Sprintf("%sConnection", nodeType.Name()),
		Description: fmt.Sprintf("Connections for model %s", nodeType.Name()),
		Fields: Fields{
			"TotalCount": {
				Type:        Int,
				Description: "Total Count of Objects in Query",
			},
			"PageInfo": {
				Type:        pageInfoModel,
				Description: "Pagination Information",
			},
			"Edges": {
				Type:        NewList(makeGraphQLEdge(nodeType)),
				Description: "Rows of the query",
			},
		},
	})
}

func MakeEdges(nodes []interface{}, cursorFunc func(m interface{}) string) []EdgeData {
	var ret = make([]EdgeData, len(nodes))

	for i := 0; i < len(nodes); i++ {
		var value interface{}
		switch nodes[i].(type) {
		case null.Bool:
			value, _ = nodes[i].(null.Bool).Value()
		case null.String:
			value, _ = nodes[i].(null.String).Value()
		case null.Time:
			value, _ = nodes[i].(null.Time).Value()
		case null.Byte:
			value, _ = nodes[i].(null.Byte).Value()
		case null.Bytes:
			value, _ = nodes[i].(null.Bytes).Value()
		case null.Float32:
			value, _ = nodes[i].(null.Float32).Value()
		case null.Float64:
			value, _ = nodes[i].(null.Float64).Value()
		case null.Int:
			value, _ = nodes[i].(null.Int).Value()
		case null.Int8:
			value, _ = nodes[i].(null.Int8).Value()
		case null.Int16:
			value, _ = nodes[i].(null.Int16).Value()
		case null.Int32:
			value, _ = nodes[i].(null.Int32).Value()
		case null.Int64:
			value, _ = nodes[i].(null.Int64).Value()
		case null.Uint:
			value, _ = nodes[i].(null.Uint).Value()
		case null.Uint8:
			value, _ = nodes[i].(null.Uint8).Value()
		case null.Uint16:
			value, _ = nodes[i].(null.Uint16).Value()
		case null.Uint32:
			value, _ = nodes[i].(null.Uint32).Value()
		case null.Uint64:
			value, _ = nodes[i].(null.Uint64).Value()
		default:
			value = nodes[i]
		}
		ret[i] = EdgeData{
			Node:   value,
			Cursor: cursorFunc(nodes[i]),
		}
	}

	return ret
}
