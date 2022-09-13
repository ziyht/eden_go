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
  UNKNOWN  ValType = 0x00
  BOOL     ValType = 0x01
  INT8     ValType = 0x02
  INT16    ValType = 0x03
  INT32    ValType = 0x04
  INT64    ValType = 0x05
  UINT8    ValType = 0x06
  UINT16   ValType = 0x07
  UINT32   ValType = 0x08
  UINT64   ValType = 0x09
  FLOAT32  ValType = 0x0a
  FLOAT64  ValType = 0x0b
  TIME     ValType = 0x0c
  DURATION ValType = 0x0d
  BYTES    ValType = 0x0e   // including string
  ITEM     ValType = 0x0f
  VT_MAX   ValType = 0x10
  VT_ERR   ValType = 0x11
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

// note: string type will be considered as []byte
func NewVal(a any) (out *Val, err error){
  out = &Val{}
  
  switch v := a.(type) {
    case bool         : out.setBool(v)
    case int8         : out.setInt8(v)
    case int16        : out.setInt16(v)  
    case int32        : out.setInt32(v)  
    case int64        : out.setInt64(v)
    case uint8        : out.setUInt8(v)
    case uint16       : out.setUInt16(v)  
    case uint32       : out.setUInt32(v)  
    case uint64       : out.setUInt64(v)
    case float32      : out.setFloat32(v)
    case float64      : out.setFloat64(v)
    case time.Time    : out.setTime(v)
    case time.Duration: out.setDuration(v)
    case []byte       : out.setBytes(v); 
    case string       : out.setBytes(ptr.StringToBytes(v));

    default: err = fmt.Errorf("type %t not supported", v)
  }

  return 
}

func (d *Val)unmarshal(b []byte){
  if len(b) < 4 {
    *d = Val{meta: [4]byte{byte(VT_ERR), 0, 0, 0}, d: ptr.StringToBytes("invalid input data, must be at least 4 bytes")}
    return
  }
  
  if len(b) <= 32{
    *d = Val{meta: [4]byte{b[0], b[1], b[2], b[3]}, d: b[4:]}
    return
  }

  m := b[len(b)-4:]
  *d = Val{meta: [4]byte{m[0], m[1], m[2], m[3]}, d: b[:len(b)-4]}
}

func (d *Val)marshal()[]byte{
  total := len(d.meta) + len(d.d)

  if total <= 32 {
    return append(d.meta[:], d.d...)
  }

  return append(d.d, d.meta[:]...)
}

// this will reset the type of Raw and clear the internal data
// we assuming that the input type not equal or bigger than
func (r *Val)Reset(t ValType) {
  if t >= VT_MAX{
    return 
  }
  r.__reset(t)
}

func (r *Val)Clone()*Val{
  return &Val{meta: r.meta, d: r.d}
}

func (r *Val)AppendInt32(v int32)error{
  if r.__isType(INT32) {
    r.__appendUint32(uint32(v))
    return nil
  } else if r.__isType(UNKNOWN){
    r.__reset(INT32)
    r.__appendUint32(uint32(v))
    return nil
  }
  return fmt.Errorf("invald operation, you can not append a %v value to a Raw with type '%s'", v, r.__typeStr())
}

func (r *Val)AppendInt64(v int64)error{
  if r.__isType(INT32) {
    r.__appendUint64(uint64(v))
    return nil
  } else if r.__isType(UNKNOWN){
    r.__reset(INT32)
    r.__appendUint64(uint64(v))
    return nil
  }
  return fmt.Errorf("invald operation, you can not append a %v value to a Raw with type '%s'", v, r.__typeStr())
}

func (r *Val)String()string{
  return "todo"
}

func (r *Val)setBool (v bool ){ r.__reset(BOOL ); r.__appendBool  (v) }
func (d *Val)setInt8 (v int8 ){ d.__reset(INT8 ); d.__appendUint8 (uint8(v )) }
func (d *Val)setInt16(v int16){ d.__reset(INT16); d.__appendUint16(uint16(v)) }
func (d *Val)setInt32(v int32){ d.__reset(INT32); d.__appendUint32(uint32(v)) }
func (d *Val)setInt64(v int64){ d.__reset(INT64); d.__appendUint64(uint64(v)) }

func (d *Val)setUInt8 (v uint8 ){ d.__reset(UINT8 ); d.__appendUint8 (v) }
func (d *Val)setUInt16(v uint16){ d.__reset(UINT16); d.__appendUint16(v) }
func (d *Val)setUInt32(v uint32){ d.__reset(UINT32); d.__appendUint32(v) }
func (d *Val)setUInt64(v uint64){ d.__reset(UINT64); d.__appendUint64(v) }

func (d *Val)setFloat32(v float32){ d.__reset(FLOAT32); d.__appendUint32(*(*uint32)(unsafe.Pointer(&v))) }
func (d *Val)setFloat64(v float64){ d.__reset(FLOAT64); d.__appendUint64(*(*uint64)(unsafe.Pointer(&v))) }

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

func (d *Val)Bool ()(bool ){ if d.__isType(BOOL ) {return d.__Uint8 () == 1  }; return false }

func (d *Val)Int8 ()(int8 ){ if d.__isType(INT8 ) {return int8 (d.__Uint8 ())}; return 0}
func (d *Val)Int16()(int16){ if d.__isType(INT16) {return int16(d.__Uint16())}; return 0}
func (d *Val)Int32()(int32){ if d.__isType(INT32) {return int32(d.__Uint32())}; return 0}
func (d *Val)Int64()(int64){ if d.__isType(INT64) {return int64(d.__Uint64())}; return 0}

func (d *Val)UInt8 ()(uint8 ){ if d.__isType(UINT8 ) {return uint8 (d.__Uint8 ())}; return 0}
func (d *Val)UInt16()(uint16){ if d.__isType(UINT16) {return uint16(d.__Uint16())}; return 0}
func (d *Val)UInt32()(uint32){ if d.__isType(UINT32) {return uint32(d.__Uint32())}; return 0}
func (d *Val)UInt64()(uint64){ if d.__isType(UINT64) {return uint64(d.__Uint64())}; return 0}

func (d *Val)Float32()(float32){ if !d.__isType(FLOAT32) {return 0}; val := d.__Uint32(); return *(*float32)(unsafe.Pointer(&val)) }
func (d *Val)Float64()(float64){ if !d.__isType(FLOAT64) {return 0}; val := d.__Uint64(); return *(*float64)(unsafe.Pointer(&val)) }

func (d *Val)Str()(string){ if d.__isType(BYTES ) { return ptr.BytesToString(d.d) }; return "" }
func (d *Val)Bytes()([]byte){ if d.__isType(BYTES ) { return d.d }; return nil }
func (d *Val)Time()(out time.Time){ if d.__isType(TIME) { out.UnmarshalBinary(d.d) }; return}
func (d *Val)Duration()(time.Duration){ if d.__isType(DURATION) { return time.Duration(d.__Uint64()) }; return 0}

func (d *Val)GetBool ()(bool,  error){ if e := d.__checkType(BOOL ); e != nil {return false, e}; return d.__Uint8 () == 1, nil }
func (d *Val)GetInt8 ()(int8,  error){ if e := d.__checkType(INT8 ); e != nil {return 0, e}; return int8 (d.__Uint8 ()), nil }
func (d *Val)GetInt16()(int16, error){ if e := d.__checkType(INT16); e != nil {return 0, e}; return int16(d.__Uint16()), nil }
func (d *Val)GetInt32()(int32, error){ if e := d.__checkType(INT32); e != nil {return 0, e}; return int32(d.__Uint32()), nil }
func (d *Val)GetInt64()(int64, error){ if e := d.__checkType(INT64); e != nil {return 0, e}; return int64(d.__Uint64()), nil }

func (d *Val)GetUInt8 ()(uint8,  error){ if e := d.__checkType(UINT8 ); e != nil {return 0, e}; return d.__Uint8 (), nil }
func (d *Val)GetUInt16()(uint16, error){ if e := d.__checkType(UINT16); e != nil {return 0, e}; return d.__Uint16(), nil }
func (d *Val)GetUInt32()(uint32, error){ if e := d.__checkType(UINT32); e != nil {return 0, e}; return d.__Uint32(), nil }
func (d *Val)GetUInt64()(uint64, error){ if e := d.__checkType(UINT64); e != nil {return 0, e}; return d.__Uint64(), nil }

func (d *Val)GetFloat32()(float32, error){ if e := d.__checkType(FLOAT32); e != nil {return 0, e}; val := d.__Uint32(); return *(*float32)(unsafe.Pointer(&val)), nil }
func (d *Val)GetFloat64()(float64, error){ if e := d.__checkType(FLOAT64); e != nil {return 0, e}; val := d.__Uint64(); return *(*float64)(unsafe.Pointer(&val)), nil }

func (d *Val)GetBytes()([]byte, error){ if e := d.__checkType(BYTES ); e != nil {return nil, e}; return d.d, nil }

func (d *Val)Error()error{
  if d.__isType(VT_ERR) {
    return errors.New(ptr.BytesToString(d.d))
  }
  return nil
}

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