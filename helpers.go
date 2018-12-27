package graphql

import (
	"encoding/base64"
	"github.com/satori/go.uuid"
	. "github.com/volatiletech/sqlboiler/queries/qm"
)

type Model interface {
	GetID() *uuid.UUID
}

func ApplyBoilerFilters(args map[string]interface{}) (mods []QueryMod) {
	first, last, after, before := ParseArgs(args)

	if after != nil {
		mods = append(mods, Where("id > ?", after))
	}

	if before != nil {
		mods = append(mods, Where("id < ?", before))
	}

	if last != -1 {
		panic("Last not implemented")
	}

	if first > 1000 {
		first = 1000
	}

	mods = append(mods, Limit(first))

	return

}

func ParseArgs(args map[string]interface{}) (first int, last int, after *uuid.UUID, before *uuid.UUID) {
	first = 10
	last = -1

	if val, ok := args["After"]; ok {
		after = FromCursorToUUID(val.(string))
	}

	if val, ok := args["Before"]; ok {
		before = FromCursorToUUID(val.(string))
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

func FromBoilerToCursor(id string) string {
	return base64.StdEncoding.EncodeToString([]byte(id))
}

func FromCursorToBoiler(cursor string) string {
	var sid, err = base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		panic(err)
	}

	return string(sid)
}

func FromUUIDToCursor(id *uuid.UUID) string {
	return base64.StdEncoding.EncodeToString(id.Bytes())
}

func FromCursorToUUID(cursor string) *uuid.UUID {
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