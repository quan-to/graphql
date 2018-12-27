package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/volatiletech/null"
)

func BoilerFieldResolveFn(p graphql.ResolveParams) (interface{}, error) {
	return GetBoilerNullableValue(p.Source)
}

func GetBoilerNullableValue(v interface{}) (interface{}, error) {
	switch v.(type) {
	case null.Bool:
		return v.(null.Bool).Value()
	case null.String:
		return v.(null.String).Value()
	case null.Time:
		return v.(null.Time).Value()
	case null.Byte:
		return v.(null.Byte).Value()
	case null.Bytes:
		return v.(null.Bytes).Value()
	case null.Float32:
		return v.(null.Float32).Value()
	case null.Float64:
		return v.(null.Float64).Value()
	case null.Int:
		return v.(null.Int).Value()
	case null.Int8:
		return v.(null.Int8).Value()
	case null.Int16:
		return v.(null.Int16).Value()
	case null.Int32:
		return v.(null.Int32).Value()
	case null.Int64:
		return v.(null.Int64).Value()
	case null.Uint:
		return v.(null.Uint).Value()
	case null.Uint8:
		return v.(null.Uint8).Value()
	case null.Uint16:
		return v.(null.Uint16).Value()
	case null.Uint32:
		return v.(null.Uint32).Value()
	case null.Uint64:
		return v.(null.Uint64).Value()
	default:
		return v, nil
	}
}
