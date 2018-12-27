package graphql

import (
	"fmt"
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
		value = GetBoilerNullableValue(nodes[i])
		ret[i] = EdgeData{
			Node:   value,
			Cursor: cursorFunc(nodes[i]),
		}
	}

	return ret
}
