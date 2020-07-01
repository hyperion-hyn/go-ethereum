package storage

import (
    "fmt"
    "reflect"
    "unsafe"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/state"
)

type setterFunc func(value interface{})

type Storage struct {
    db *state.StateDB
    addr common.Address
    slot int
    object interface{}
    dirty map[string]interface{}
    setter setterFunc
}


func NewStorage(db *state.StateDB, addr common.Address, slot int, obj interface{}, setter setterFunc) (*Storage){
    return &Storage{
        db: db,
        addr: addr,
        slot: slot,
        object: obj,
        dirty: make(map[string]interface{}, 0),
        setter: setter,
    }
}

func GetUnexportedField(field reflect.Value) interface{} {
    return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
}

func (s* Storage) Get(nameOrIndex interface{} ) (*Storage) {
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

func (s* Storage) GetByName(name string) (*Storage) {
    v := reflect.ValueOf(s.object)
    indirectVal := reflect.Indirect(v)
    switch indirectVal.Kind() {
    case reflect.Struct:
        if sf, ok := indirectVal.Type().FieldByName(name); ok {
            fv := indirectVal.FieldByName(name)
            var val reflect.Value
            if fv.CanAddr() {
                val =  fv.Addr()
            } else {
                val = fv
            }

            isUnexported := sf.PkgPath != ""
            if !isUnexported {
                obj := val.Interface()
                return NewStorage(s.db, s.addr, 0, obj, func(value interface{}) {
                    target := reflect.Indirect(reflect.ValueOf(obj))
                    target.Set(reflect.ValueOf(value))
                })
            } else {
                obj := reflect.NewAt(reflect.Indirect(val).Type(), unsafe.Pointer(val.Pointer())).Interface()
                return NewStorage(s.db, s.addr, 0, obj, func(value interface{}) {
                    target := reflect.Indirect(reflect.ValueOf(obj))
                    target.Set(reflect.ValueOf(value))
                })
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
        return NewStorage(s.db, s.addr, 0, obj.Interface(), func(value interface{}) {
            indirectVal.SetMapIndex(reflect.ValueOf(name), reflect.ValueOf(value))
        })
    default:
        panic(fmt.Sprintf("%s not supported", indirectVal.Kind()))
    }
}

func (s* Storage) GetByIndex(index int) (*Storage) {
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
            val =  fv.Addr()
        } else {
            val = fv
        }

        obj := val.Interface()
        return NewStorage(s.db, s.addr, 0, obj, func(value interface{}) {
            target := reflect.Indirect(reflect.ValueOf(obj))
            target.Set(reflect.ValueOf(value))
        })

    default:
        panic(fmt.Sprintf("%s not supported", indirectVal.Kind()))
    }
}

func (s *Storage) Value() (interface{}){
    v := reflect.ValueOf(s.object)
    indirectVal := reflect.Indirect(v)
    return indirectVal.Interface()
}

func (s *Storage) Value2(isLazy bool) interface{} {
    return nil
}

func (s *Storage) SetValue(data interface{}) (*Storage) {
    s.setter(data)
    return s
}

func (s *Storage) Contain(data interface{}) bool {
    return true
}

func (s *Storage) DeleteByName(name string) bool {
    return true
}

func (s *Storage) DeleteByIndex(index int) error {
    return nil
}

func (s *Storage) Len() int {
    // support array and map
    return 0
}

func (s *Storage) Keys() interface{} {
    // only support map
    return 0
}