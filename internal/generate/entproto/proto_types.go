package entproto

import (
	"entgo.io/ent/schema/field"
	"github.com/jhump/protoreflect/v2/protobuilder"
	entpb "github.com/syralon/entc-gen-go/proto/syralon/entproto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TypeMapping interface {
	Mapping(field.Type) *protobuilder.FieldType
}

var (
	EntityTypeMapping typeMapping = map[field.Type]*protobuilder.FieldType{
		field.TypeBool:    protobuilder.FieldTypeBool(),
		field.TypeTime:    protobuilder.FieldTypeImportedMessage((&timestamppb.Timestamp{}).ProtoReflect().Descriptor()),
		field.TypeJSON:    protobuilder.FieldTypeBytes(),
		field.TypeUUID:    protobuilder.FieldTypeString(),
		field.TypeBytes:   protobuilder.FieldTypeBytes(),
		field.TypeString:  protobuilder.FieldTypeString(),
		field.TypeInt8:    protobuilder.FieldTypeInt32(),
		field.TypeInt16:   protobuilder.FieldTypeInt32(),
		field.TypeInt32:   protobuilder.FieldTypeInt32(),
		field.TypeInt:     protobuilder.FieldTypeInt64(),
		field.TypeInt64:   protobuilder.FieldTypeInt64(),
		field.TypeUint8:   protobuilder.FieldTypeUint32(),
		field.TypeUint16:  protobuilder.FieldTypeUint32(),
		field.TypeUint32:  protobuilder.FieldTypeUint32(),
		field.TypeUint:    protobuilder.FieldTypeUint64(),
		field.TypeUint64:  protobuilder.FieldTypeUint64(),
		field.TypeFloat32: protobuilder.FieldTypeFloat(),
		field.TypeFloat64: protobuilder.FieldTypeDouble(),
		// field.TypeEnum:    nil, // TODO
		// field.TypeOther:   nil, // TODO
	}

	OperationTypeMapping typeMapping = map[field.Type]*protobuilder.FieldType{
		field.TypeBool:    protobuilder.FieldTypeImportedMessage((&entpb.BoolField{}).ProtoReflect().Descriptor()),
		field.TypeTime:    protobuilder.FieldTypeImportedMessage((&entpb.TimestampField{}).ProtoReflect().Descriptor()),
		field.TypeJSON:    protobuilder.FieldTypeImportedMessage((&entpb.BytesField{}).ProtoReflect().Descriptor()),
		field.TypeUUID:    protobuilder.FieldTypeImportedMessage((&entpb.StringField{}).ProtoReflect().Descriptor()),
		field.TypeBytes:   protobuilder.FieldTypeImportedMessage((&entpb.BytesField{}).ProtoReflect().Descriptor()),
		field.TypeEnum:    protobuilder.FieldTypeImportedMessage((&entpb.Int32Field{}).ProtoReflect().Descriptor()),
		field.TypeString:  protobuilder.FieldTypeImportedMessage((&entpb.StringField{}).ProtoReflect().Descriptor()),
		field.TypeInt8:    protobuilder.FieldTypeImportedMessage((&entpb.Int32Field{}).ProtoReflect().Descriptor()),
		field.TypeInt16:   protobuilder.FieldTypeImportedMessage((&entpb.Int32Field{}).ProtoReflect().Descriptor()),
		field.TypeInt32:   protobuilder.FieldTypeImportedMessage((&entpb.Int32Field{}).ProtoReflect().Descriptor()),
		field.TypeInt:     protobuilder.FieldTypeImportedMessage((&entpb.Int64Field{}).ProtoReflect().Descriptor()),
		field.TypeInt64:   protobuilder.FieldTypeImportedMessage((&entpb.Int64Field{}).ProtoReflect().Descriptor()),
		field.TypeUint8:   protobuilder.FieldTypeImportedMessage((&entpb.Uint32Field{}).ProtoReflect().Descriptor()),
		field.TypeUint16:  protobuilder.FieldTypeImportedMessage((&entpb.Uint32Field{}).ProtoReflect().Descriptor()),
		field.TypeUint32:  protobuilder.FieldTypeImportedMessage((&entpb.Uint32Field{}).ProtoReflect().Descriptor()),
		field.TypeUint:    protobuilder.FieldTypeImportedMessage((&entpb.Uint64Field{}).ProtoReflect().Descriptor()),
		field.TypeUint64:  protobuilder.FieldTypeImportedMessage((&entpb.Uint64Field{}).ProtoReflect().Descriptor()),
		field.TypeFloat32: protobuilder.FieldTypeImportedMessage((&entpb.FloatField{}).ProtoReflect().Descriptor()),
		field.TypeFloat64: protobuilder.FieldTypeImportedMessage((&entpb.DoubleField{}).ProtoReflect().Descriptor()),
		// field.TypeOther:   nil,
	}

	TypePaginator = protobuilder.FieldTypeImportedMessage((&entpb.Paginator{}).ProtoReflect().Descriptor())
)

type typeMapping map[field.Type]*protobuilder.FieldType

func (m typeMapping) Mapping(t field.Type) *protobuilder.FieldType {
	return m[t]
}
