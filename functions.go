package dilema

// // GetFromStruct returnStruct which is implementation of provided interface.
// // Provided interface must contain fields typed with provided earlier services
// // Now it must contain only singleTones
// func (di *dicon) GetFromStruct(unionStruct interface{}, args ...interface{}) interface{} {
// 	t, _, err := checkProvidedTypeIsCorrectStruct(unionStruct)
// 	if err != nil {
// 		return nil
// 	}

// 	fields := make([]reflect.StructField, 0)
// 	fieldsValues := make(map[reflect.Type]reflect.Value)
// 	argsIndex := 0
// 	for i := 0; i < t.NumMethod(); i++ {
// 		t, sType, ok := di.findServiceByMethod(t.Method(i))
// 		if !ok {
// 			return nil
// 		}

// 		var value reflect.Value
// 		switch sType {
// 		case singleToneType:
// 			value = di.singletonesByAlias[t]
// 		case temporalType:
// 			constr := di.temporalByAlias[t]
// 			value, err = di.createService(constr.Creator(), &argsIndex, args...)
// 			if err != nil {
// 				return nil
// 			}
// 		}

// 		field := reflect.StructField{
// 			Name:      t.Name(),
// 			Type:      t,
// 			Anonymous: false,
// 		}
// 		fields = append(fields, field)
// 		fieldsValues[t] = value
// 	}

// 	structPtr := reflect.New(reflect.StructOf(fields))
// 	structValue := structPtr.Elem()

// 	for i := 0; i < structValue.NumField(); i += 1 {
// 		field := structValue.Field(i)
// 		value := fieldsValues[field.Type()]
// 		structValue.Field(i).Set(value)
// 	}

// 	return structValue.Interface()
// }

// func (di *dicon) findServiceByMethod(method reflect.Method) (reflect.Type, containerType, bool) {
// 	method.Index = 0
// 	di.mutex.Lock()
// 	action, ok := di.destroyers[method]
// 	di.mutex.Unlock()
// 	return action.actionType, action.serviceType, ok
// }

// func (di *dicon) addToCache(t reflect.Type, v reflect.Value) {
// 	di.mutex.Lock()
// 	di.cache[t] = v
// 	di.mutex.Unlock()
// }
