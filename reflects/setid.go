package reflects

import (
	"errors"
	"reflect"
)

func SetID(i interface{}, newID string) error {

	r := reflect.ValueOf(i)

	if r.Kind() != reflect.Ptr {
		return errors.New("Ptr should be given, else Pass By Value prevent setting struct ID field remotely")
	}

	val, ok := idReflectValue(r)

	if !ok {
		return errors.New("could not locate ID field in the given structure")
	}

	val.Set(reflect.ValueOf(newID))

	return nil
}
