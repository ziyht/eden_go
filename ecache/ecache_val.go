package ecache

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/cockroachdb/errors"
	"github.com/ziyht/eden_go/utils/ptr"
)

type ValType byte

const (
  Nil      ValType = 0x00
  BOOL     ValType = 0x01
  I8       ValType = 0x02
  I16      ValType = 0x03
  I32      ValType = 0x04
  I64      ValType = 0x05
  U8       ValType = 0x06
  U16      ValType = 0x07
  U32      ValType = 0x08
  U64      ValType = 0x09
  F32      ValType = 0x0a
  F64      ValType = 0x0b
  TIME     ValType = 0x0c
  DURATION ValType = 0x0d
  BYTES    ValType = 0x0e   // including string
  ITEM     ValType = 0x0f
  VT_MAX   ValType = 0x10
  VT_ERR   ValType = 0xff
)

var (
  type_strs = []string{
    "none",
    "bool",
    "int8",
    "int16",
    "int32",
    "int64",
    "uint8",
    "uint16",
    "uint32",
    "uint64",
    "float32",
    "float64",
    "time",
    "duration",
    "[]byte",
    "item",
  }
)

type Val struct {
  meta   [4]byte
  d      []byte
}

func (d *Val)unmarshal(b []byte){
  if len(b) < 4 {
    d.setErr("invalid input data, must be at least 4 bytes")
    return
  }

  *d = Val{meta: [4]byte{b[0], b[1], b[2], b[3]}, d: b[4:]}
}

func (d *Val)marshal()[]byte{
  return append(d.meta[:], d.d...)
}

// note: string type will be considered as []byte
func NewVal(d any) (out *Val, err error){
  out = new(Val)
  out.Reset(d)
  return 
}

// note: this will reset type and data by d 
func (v *Val)Reset(d any) error {
  if d == nil {
    v.__reset(Nil)
  }

  switch d_ := d.(type) {
    case bool         : v.setBool(d_)
    case int8         : v.setI8( d_)
    case int16        : v.setI16(d_)  
    case int32        : v.setI32(d_)  
    case int64        : v.setI64(d_)
    case uint8        : v.setU8( d_)
    case uint16       : v.setU16(d_)  
    case uint32       : v.setU32(d_)  
    case uint64       : v.setU64(d_)
    case float32      : v.setF32(d_)
    case float64      : v.setF64(d_)
    case time.Time    : v.setTime(d_)
    case time.Duration: v.setDuration(d_)
    case []byte       : v.setBytes(d_); 
    case string       : v.setBytes(ptr.StringToBytes(d_));
    case Val          : *v = d_
    case *Val         : *v = *d_
    case ValType      : if d_ >= VT_MAX {
                          v.setErr(fmt.Sprintf("input type %d overflow", d_))
                        }
                        v.__reset(d_)

    default: v.setErr(fmt.Sprintf("type %T not supported", d_))
             return v.Error()
  }

  return nil
}

func (r *Val)Clone()*Val{
  return &Val{meta: r.meta, d: r.d[:]}
}

func (r *Val)AppendInt32(v int32)error{
  if r.__isType(I32) {
    r.__appendUint32(uint32(v))
    return nil
  } else if r.__isType(Nil){
    r.__reset(I32)
    r.__appendUint32(uint32(v))
    return nil
  }
  return fmt.Errorf("invald operation, you can not append a %v value to a Raw with type '%s'", v, r.__typeStr())
}

func (r *Val)AppendInt64(v int64)error{
  if r.__isType(I32) {
    r.__appendUint64(uint64(v))
    return nil
  } else if r.__isType(Nil){
    r.__reset(I32)
    r.__appendUint64(uint64(v))
    return nil
  }
  return fmt.Errorf("invald operation, you can not append a %v value to a Raw with type '%s'", v, r.__typeStr())
}

func (r *Val)String()string{
  return "todo"
}

func (r *Val)setBool (v bool ){ r.__reset(BOOL ); r.__appendBool  (v) }
func (d *Val)setI8 (v int8 ){ d.__reset(I8 ); d.__appendUint8 (uint8(v )) }
func (d *Val)setI16(v int16){ d.__reset(I16); d.__appendUint16(uint16(v)) }
func (d *Val)setI32(v int32){ d.__reset(I32); d.__appendUint32(uint32(v)) }
func (d *Val)setI64(v int64){ d.__reset(I64); d.__appendUint64(uint64(v)) }

func (d *Val)setU8 (v uint8 ){ d.__reset(U8 ); d.__appendUint8 (v) }
func (d *Val)setU16(v uint16){ d.__reset(U16); d.__appendUint16(v) }
func (d *Val)setU32(v uint32){ d.__reset(U32); d.__appendUint32(v) }
func (d *Val)setU64(v uint64){ d.__reset(U64); d.__appendUint64(v) }

func (d *Val)setF32(v float32){ d.__reset(F32); d.__appendUint32(*(*uint32)(unsafe.Pointer(&v))) }
func (d *Val)setF64(v float64){ d.__reset(F64); d.__appendUint64(*(*uint64)(unsafe.Pointer(&v))) }

func (d *Val)setBytes(v []byte){ d.__reset(BYTES ); d.d = append(d.d, v...) }
func (d *Val)setItem(v Item) error { d.__reset(ITEM ); b, e := v.Marshal(); if e != nil {return e}; d.d = append(d.d, b...); return nil}

func (d *Val)setTime(t time.Time)error{
  bytes, err := t.MarshalBinary()
  if err != nil {
    return err
  }
  d.__reset(TIME)
  d.__appendBytes(bytes)
  return nil
}

func (d *Val)setDuration(du time.Duration){
  d.__reset(DURATION)
  d.__appendUint64(uint64(du))
}

func (d *Val)setErr(err string){
  d.__reset(VT_ERR); d.__appendBytes(ptr.StringToBytes(err))
}

func (d *Val)Bool()(bool ){ if d.__isType(BOOL ) {return d.__Uint8 () == 1  }; return false }
func (d *Val)I8 ()(int8 ){ if d.__isType(I8 ) {return int8 (d.__Uint8 ())}; return 0}
func (d *Val)I16()(int16){ if d.__isType(I16) {return int16(d.__Uint16())}; return 0}
func (d *Val)I32()(int32){ if d.__isType(I32) {return int32(d.__Uint32())}; return 0}
func (d *Val)I64()(int64){ if d.__isType(I64) {return int64(d.__Uint64())}; return 0}
func (d *Val)U8 ()(uint8 ){ if d.__isType(U8 ) {return uint8 (d.__Uint8 ())}; return 0}
func (d *Val)U16()(uint16){ if d.__isType(U16) {return uint16(d.__Uint16())}; return 0}
func (d *Val)U32()(uint32){ if d.__isType(U32) {return uint32(d.__Uint32())}; return 0}
func (d *Val)U64()(uint64){ if d.__isType(U64) {return uint64(d.__Uint64())}; return 0}
func (d *Val)F32()(float32){ if !d.__isType(F32) {return 0}; val := d.__Uint32(); return *(*float32)(unsafe.Pointer(&val)) }
func (d *Val)F64()(float64){ if !d.__isType(F64) {return 0}; val := d.__Uint64(); return *(*float64)(unsafe.Pointer(&val)) }
func (d *Val)Str()(string){ if d.__isType(BYTES ) { return ptr.BytesToString(d.d) }; return "" }
func (d *Val)Bytes()([]byte){ if d.__isType(BYTES ) { return d.d }; return nil }
func (d *Val)Time()(out time.Time){ if d.__isType(TIME) { out.UnmarshalBinary(d.d) }; return}
func (d *Val)Duration()(time.Duration){ if d.__isType(DURATION) { return time.Duration(d.__Uint64()) }; return 0}
func (d *Val)Error()error{ if d.__isType(VT_ERR) { return errors.New(ptr.BytesToString(d.d)) }; return nil }

func (d *Val)GetBool()(bool,  error){ if e := d.__checkType(BOOL ); e != nil {return false, e}; return d.__Uint8 () == 1, nil }
func (d *Val)GetI8 ()(int8,  error){ if e := d.__checkType(I8 ); e != nil {return 0, e}; return int8 (d.__Uint8 ()), nil }
func (d *Val)GetI16()(int16, error){ if e := d.__checkType(I16); e != nil {return 0, e}; return int16(d.__Uint16()), nil }
func (d *Val)GetI32()(int32, error){ if e := d.__checkType(I32); e != nil {return 0, e}; return int32(d.__Uint32()), nil }
func (d *Val)GetI64()(int64, error){ if e := d.__checkType(I64); e != nil {return 0, e}; return int64(d.__Uint64()), nil }
func (d *Val)GetU8 ()(uint8,  error){ if e := d.__checkType(U8 ); e != nil {return 0, e}; return d.__Uint8 (), nil }
func (d *Val)GetU16()(uint16, error){ if e := d.__checkType(U16); e != nil {return 0, e}; return d.__Uint16(), nil }
func (d *Val)GetU32()(uint32, error){ if e := d.__checkType(U32); e != nil {return 0, e}; return d.__Uint32(), nil }
func (d *Val)GetU64()(uint64, error){ if e := d.__checkType(U64); e != nil {return 0, e}; return d.__Uint64(), nil }
func (d *Val)GetF32()(float32, error){ if e := d.__checkType(F32); e != nil {return 0, e}; val := d.__Uint32(); return *(*float32)(unsafe.Pointer(&val)), nil }
func (d *Val)GetF64()(float64, error){ if e := d.__checkType(F64); e != nil {return 0, e}; val := d.__Uint64(); return *(*float64)(unsafe.Pointer(&val)), nil }
func (d *Val)GetBytes()([]byte, error){ if e := d.__checkType(BYTES ); e != nil {return nil, e}; return d.d, nil }
func (d *Val)GetTime()(out time.Time, e error){ if e = d.__checkType(TIME ); e != nil {return }; if e = out.UnmarshalBinary(d.d); e != nil {return }; return }
func (d *Val)GetDuration()(out time.Duration, e error){ if e = d.__checkType(DURATION ); e != nil {return }; return time.Duration(d.__Uint64()), nil }

func (d *Val)Type()ValType{
  return ValType(d.meta[0])
}

func (d *Val)__set_type(t ValType){
  d.meta[0] = byte(t)
}

func (d *Val)__clear(){
  d.d = d.d[:0]
}

func (d *Val)__isType(t ValType)bool{
  return d.meta[0] == byte(t)
}

func (d *Val)__reset(t ValType){
  d.__set_type(t)
  d.__clear()
}

func (d *Val)__checkType(t ValType) error {
  if !d.__isType(t){
    return fmt.Errorf("invalid type %s", d.__typeStr())
  }
  return nil
}

func (d *Val)__typeStr() string {
  if d.Type() >= VT_MAX{
    return "(TypeOverload)"
  }
  return type_strs[d.Type()]
}

func (d *Val)__appendUint64(v uint64) {
	d.d = append(d.d,
		byte(v),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24),
		byte(v>>32),
		byte(v>>40),
		byte(v>>48),
		byte(v>>56),
	)
}

func (d *Val)__appendUint32(v uint32) {
	d.d = append(d.d,
		byte(v),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24),
	)
}

func (d *Val)__appendUint16(v uint16) {
	d.d = append(d.d,
		byte(v),
		byte(v>>8),
	)
}

func (d *Val)__appendUint8(v uint8) {
	d.d = append(d.d,
		byte(v),
	)
}

func (d *Val)__appendBool(v bool) {
  if v { d.d = append(d.d, byte(1))
  } else { d.d = append(d.d, byte(0))}
}

func (d *Val)__appendBytes(v []byte) {
  d.d = append(d.d, v...)
}

func (d *Val)__Uint64()uint64{
  b := d.d
	_ = b[7] // bounds check hint to compiler; see golang.org/issue/14808
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

func (d *Val)__Uint32()uint32{
  b := d.d
	_ = b[3] // bounds check hint to compiler; see golang.org/issue/14808
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24 
}

func (d *Val)__Uint16()uint16{
  b := d.d
	_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
	return uint16(b[0]) | uint16(b[1])<<8
}

func (d *Val)__Uint8()uint8{
  b := d.d
	_ = b[0] // bounds check hint to compiler; see golang.org/issue/14808
	return uint8(b[0])
}