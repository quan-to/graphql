package graphql

import (
	"encoding/base64"
	"github.com/satori/go.uuid"
	"github.com/jinzhu/gorm"
)

type Model interface {
	GetID() *uuid.UUID
}

// Not working
func GetNodes(o interface{}) []interface{} {
	obj := o.([]interface{})
	gNodes := make([]interface{}, len(obj))
	for i, v := range obj {
		gNodes[i] = v
	}

	return gNodes
}

// Not working
func GetCursors(o interface{}) (startCursor, endCursor string) {

	obj := o.([]interface{})
	if len(obj) > 0 {
		v := obj[0]
		startCursor = FromGormToCursor(v.(Model).GetID())
		v = obj[len(obj)-1]
		endCursor = FromGormToCursor(v.(Model).GetID())
	}

	return
}

func ApplyGormFilters(e *gorm.DB, args map[string]interface{}) (*gorm.DB) {
	first, last, after, before := ParseArgs(args)

	if after != nil {
		e = e.Where("id > ?", after)
	}

	if before != nil {
		e = e.Where("id < ?", before)
	}

	if last != -1 {
		panic("Last not implemented")
	}

	if first > 1000 {
		first = 1000
	}

	e = e.Limit(first)

	return e

}

func ParseArgs(args map[string]interface{}) (first int, last int, after *uuid.UUID, before *uuid.UUID) {
	first = 10
	last = -1

	if val, ok := args["After"]; ok {
		after = FromCursorToGorm(val.(string))
	}

	if val, ok := args["Before"]; ok {
		before = FromCursorToGorm(val.(string))
	}

	if val, ok := args["First"]; ok {
		first = val.(int)
	}

	if val, ok := args["Last"]; ok {
		last = val.(int)
	}

	return
}

func GetOutput(totalCount int64, startCursor string, endCursor string, edges []EdgeData) (output *ConnectionData) {

	output = &ConnectionData{
		TotalCount: totalCount,
		PageInfo: PageInfo{
			StartCursor: startCursor,
			EndCursor:   endCursor,
		},
		Edges: edges,
	}

	return
}

func FromGormToCursor(id *uuid.UUID) string {
	return base64.StdEncoding.EncodeToString(id.Bytes())
}

func FromCursorToGorm(cursor string) *uuid.UUID {
	var sid, err = base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		panic(err)
	}

	id, err := uuid.FromBytes(sid)
	if err != nil {
		panic(err)
	}

	return &id
}
