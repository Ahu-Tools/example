package security

import (
	"fmt"
	"reflect"

	"github.com/Ahu-Tools/example/crypto"
	"gorm.io/gorm"
)

// BlindIndexCallback is a generic hook for GORM
func BlindIndexCallback(db *gorm.DB) {
	if db.Statement.Schema == nil {
		return
	}

	for _, field := range db.Statement.Schema.Fields {
		targetFieldName := field.Tag.Get("blind")
		if targetFieldName == "" {
			continue
		}

		// 1. Get Source Value (Pointer Safe)
		fieldValue, isZero := field.ValueOf(db.Statement.Context, db.Statement.ReflectValue)

		// Handle Pointers
		v := reflect.ValueOf(fieldValue)
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				continue
			}
			v = v.Elem()
		}
		if isZero && v.IsZero() {
			continue
		}

		// 2. Compute Hash
		sourceString := fmt.Sprintf("%v", v.Interface())
		hashedValue := crypto.GlobalEncrypter.ComputeBlindIndex(sourceString)

		// 3. METHOD A: Update the Struct directly (Reflection)
		// This ensures the Go object in memory has the new value
		dest := db.Statement.ReflectValue
		if dest.Kind() == reflect.Ptr {
			dest = dest.Elem()
		}

		if targetField := dest.FieldByName(targetFieldName); targetField.IsValid() && targetField.CanSet() {
			// Check if target is a pointer (*string) or value (string)
			if targetField.Kind() == reflect.Ptr {
				// Handle *string
				strPtr := &hashedValue
				targetField.Set(reflect.ValueOf(strPtr))
			} else {
				// Handle string
				targetField.SetString(hashedValue)
			}
		}

		// 4. METHOD B: Tell GORM explicitly (The Statement)
		// This ensures it gets into the SQL INSERT/UPDATE map
		db.Statement.SetColumn(targetFieldName, hashedValue)
	}
}
