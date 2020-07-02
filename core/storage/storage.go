package storage

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"unsafe"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/log"
)

type setterFunc func(value interface{})

type DirtyMapping map[common.Hash]common.Hash

type Storage struct {
	db     *state.StateDB
	addr   common.Address
	slot   int
	object interface{}
	dirty  DirtyMapping
	setter setterFunc
}

func NewStorage(db *state.StateDB, addr common.Address, slot int, obj interface{}, setter setterFunc) *Storage {
	return newStorage(db, addr, slot, obj, setter, nil)
}

func newStorage(db *state.StateDB, addr common.Address, slot int, obj interface{}, setter setterFunc, dirty DirtyMapping) *Storage {
	var val DirtyMapping
	if dirty == nil {
		val = make(DirtyMapping)
	} else {
		val = dirty
	}

	return &Storage{
		db:     db,
		addr:   addr,
		slot:   slot,
		object: obj,
		dirty:  val,
		setter: setter,
	}
}

func parseTag(tag string) (int, error) {
	var val int
	if _, err := fmt.Sscanf(tag, "slot%d", &val); err == nil {
		if val > 0 {
			return val, nil
		}
		return 0, errors.New(fmt.Sprintf("invalid tag: %s", tag))
	}

	return 0, errors.New(fmt.Sprintf("invalid tag: %s", tag))
}

func updateSlot(base int, slot int, typ reflect.Type, val reflect.Value, slots DirtyMapping) {
	s := common.BigToHash(big.NewInt(int64(base + slot)))
	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		slots[s] = common.BigToHash(big.NewInt(val.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	case reflect.Float32, reflect.Float64:
		// TODO:
	}
}

func GetUnexportedField(field reflect.Value) interface{} {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
}

func (s *Storage) Get(nameOrIndex interface{}) *Storage {
	val := reflect.ValueOf(nameOrIndex)
	switch val.Kind() {
	case reflect.String:
		return s.GetByName(nameOrIndex.(string))
	case reflect.Int:
		return s.GetByIndex(nameOrIndex.(int))
	default:
		panic(fmt.Sprintf("paramter of storage.Get() should be string or int, but %s got.", val.Kind()))
	}
}

func (s *Storage) GetByName(name string) *Storage {
	v := reflect.ValueOf(s.object)
	indirectVal := reflect.Indirect(v)
	switch indirectVal.Kind() {
	case reflect.Struct:
		if sf, ok := indirectVal.Type().FieldByName(name); ok {
			fv := indirectVal.FieldByName(name)
			var val reflect.Value
			if fv.CanAddr() {
				val = fv.Addr()
			} else {
				val = fv
			}

			slot, _ := sf.Tag.Lookup("storage")
			log.Info("Lookup", "slot", slot)

			isUnexported := sf.PkgPath != ""
			if !isUnexported {
				obj := val.Interface()
				return newStorage(s.db, s.addr, 0, obj, func(value interface{}) {
					target := reflect.Indirect(reflect.ValueOf(obj))
					target.Set(reflect.ValueOf(value))
					tag := sf.Tag.Get("storage")
					fmt.Println(":::", tag)
					if slot, err := parseTag(tag); err != nil {
						fmt.Println(">>>", tag)
						typ := sf.Type
						updateSlot(s.slot, slot, typ, reflect.ValueOf(value), s.dirty)
					}

				}, s.dirty)
			} else {
				obj := reflect.NewAt(reflect.Indirect(val).Type(), unsafe.Pointer(val.Pointer())).Interface()
				return newStorage(s.db, s.addr, 0, obj, func(value interface{}) {
					target := reflect.Indirect(reflect.ValueOf(obj))
					target.Set(reflect.ValueOf(value))
					tag := sf.Tag.Get("storage")
					if slot, err := parseTag(tag); err != nil {
						fmt.Println(">>>", tag)
						typ := sf.Type
						updateSlot(s.slot, slot, typ, reflect.ValueOf(value), s.dirty)
					}

				}, s.dirty)
			}
		} else {
			panic(fmt.Sprintf("type %s have no field %s", indirectVal.Type(), name))
		}

	case reflect.Map:
		if indirectVal.IsNil() {
			rv := reflect.MakeMap(indirectVal.Type())
			indirectVal.Set(rv)
		}
		obj := indirectVal.MapIndex(reflect.ValueOf(name))
		if !obj.IsValid() {
			indirectVal.SetMapIndex(reflect.ValueOf(name), reflect.Zero(indirectVal.Type().Elem()))
			obj = indirectVal.MapIndex(reflect.ValueOf(name))
		}
		return newStorage(s.db, s.addr, 0, obj.Interface(), func(value interface{}) {
			indirectVal.SetMapIndex(reflect.ValueOf(name), reflect.ValueOf(value))
		}, s.dirty)
	default:
		panic(fmt.Sprintf("%s not supported", indirectVal.Kind()))
	}
}

func (s *Storage) GetByIndex(index int) *Storage {
	v := reflect.ValueOf(s.object)
	indirectVal := reflect.Indirect(v)
	switch indirectVal.Kind() {
	case reflect.Array, reflect.Slice:
		if indirectVal.Len() <= index {
			typ := indirectVal.Type()
			cap := index - indirectVal.Len() + 1
			items := reflect.MakeSlice(typ, cap, cap)
			rv := reflect.AppendSlice(indirectVal, items)
			indirectVal.Set(rv)
		}

		fv := indirectVal.Index(index)
		var val reflect.Value
		if fv.CanAddr() {
			val = fv.Addr()
		} else {
			val = fv
		}

		obj := val.Interface()
		return newStorage(s.db, s.addr, 0, obj, func(value interface{}) {
			target := reflect.Indirect(reflect.ValueOf(obj))
			target.Set(reflect.ValueOf(value))
		}, s.dirty)

	default:
		panic(fmt.Sprintf("%s not supported", indirectVal.Kind()))
	}
}

func (s *Storage) Value() interface{} {
	v := reflect.ValueOf(s.object)
	indirectVal := reflect.Indirect(v)
	return indirectVal.Interface()
}

func (s *Storage) SetValue(data interface{}) *Storage {
	s.setter(data)
	return s
}

func (s *Storage) Flush() *Storage {
	log.Debug("Flush", "dirty", s.dirty)
	for k, v := range s.dirty {
		log.Debug("Flush", "addr", s.addr, "slot", k, "value", v)
		s.db.SetState(s.addr, k, v)
	}
	return s
}

func (s *Storage) StateDB()  *state.StateDB {
	return s.db
}