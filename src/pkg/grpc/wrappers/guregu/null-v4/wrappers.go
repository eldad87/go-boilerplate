package null_v4

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
	"gopkg.in/guregu/null.v4"
)

// Convert From TO

func StringValueToNull(sv *wrappers.StringValue, ns *null.String) {
	if sv != nil {
		ns.SetValid(sv.Value)
	} else {
		ns.Valid = false
	}
}
func NullToStringValue(ns *null.String) *wrappers.StringValue {
	if ns != nil && ns.Valid {
		return &wrappers.StringValue{Value: ns.String}
	}
	return nil
}

func NullToTimestamp(ns *null.Time) *timestamp.Timestamp {
	if ns != nil && ns.Valid {
		t, _ := ptypes.TimestampProto(ns.Time)
		return t
	}
	return nil

}
