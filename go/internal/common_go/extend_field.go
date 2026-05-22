package common_go

import (
	"errors"
	"fmt"
	"slices"

	"gitee.com/cruvie/kk_go_kit/kk_id"
	"github.com/cruvie/kk-scheduler/go/kk_scheduler"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func CheckFields(msg proto.Message) (err error) {
	reflectMsg := msg.ProtoReflect()
	desc := reflectMsg.Descriptor()

	fields := desc.Fields()

	var requiredFieldNames []protoreflect.Name
	var uuidFieldNames []protoreflect.Name
	var ipFieldNames []protoreflect.Name
	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		behaviors := fieldBehaviors(field)
		{
			if slices.Contains(behaviors, kk_scheduler.FieldBehavior_REQUIRED) {
				// 检查字段是否设置
				if !reflectMsg.Has(field) {
					requiredFieldNames = append(requiredFieldNames, field.Name())
				}
			}
			if reflectMsg.Has(field) {
				// 设置了值
				if slices.Contains(behaviors, kk_scheduler.FieldBehavior_UUID7) {
					uuidFieldNames = checkUUID7(reflectMsg, field)
				}
			}
		}
	}
	if len(requiredFieldNames) != 0 {
		err = errors.Join(fmt.Errorf("field(s) required: %s; ", requiredFieldNames))
	}
	if len(uuidFieldNames) != 0 {
		err = errors.Join(err, fmt.Errorf("field(s) uuid7 invalid: %s; ", uuidFieldNames))
	}
	if len(ipFieldNames) != 0 {
		err = errors.Join(err, fmt.Errorf("field(s) ip invalid: %s; ", ipFieldNames))
	}
	return err
}

func checkUUID7(reflectMsg protoreflect.Message, field protoreflect.FieldDescriptor) (fieldNames []protoreflect.Name) {
	if field.IsList() {
		for index := range reflectMsg.Get(field).List().Len() {
			if !(kk_id.ValidateUUID7(reflectMsg.Get(field).List().Get(index).String())) {
				fieldNames = append(fieldNames, field.Name())
			}
		}
	} else {
		if !(kk_id.ValidateUUID7(reflectMsg.Get(field).String())) {
			fieldNames = append(fieldNames, field.Name())
		}
	}
	return fieldNames
}

func fieldBehaviors(field protoreflect.FieldDescriptor) []kk_scheduler.FieldBehavior {
	// 获取字段选项
	opts := field.Options()
	if opts == nil {
		return nil
	}

	// 从字段选项中获取 field_behavior 扩展
	v := proto.GetExtension(opts, kk_scheduler.E_FieldBehavior)

	// 检查是否能转换为 FieldBehavior 切片
	if behaviors, ok := v.([]kk_scheduler.FieldBehavior); ok {
		return behaviors
	}
	return nil
}
