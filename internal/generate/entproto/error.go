package entproto

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"
)

//go:generate stringer -type ProtoType -trimprefix Proto
type ProtoType int

const (
	Protoenum ProtoType = iota + 1
	Protomessage
	Protoservice
	Protomethod
)

//go:generate stringer -type ErrorReason -trimprefix ErrorReason
type ErrorReason int

const (
	ErrorReasonUnknown ErrorReason = iota
	ErrorReasonNotFound
)

type ProtoError struct {
	Type   ProtoType
	Name   protoreflect.Name
	Reason ErrorReason
}

func (e *ProtoError) Error() string {
	return fmt.Sprintf("%s %s %s", e.Type, e.Name, e.Reason)
}

func ErrorMessageNotFound(name protoreflect.Name) error {
	return &ProtoError{
		Type:   Protomessage,
		Name:   name,
		Reason: ErrorReasonNotFound,
	}
}

func ErrorEnumNotFound(name protoreflect.Name) error {
	return &ProtoError{
		Type:   Protoenum,
		Name:   name,
		Reason: ErrorReasonNotFound,
	}
}
