package ecache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTypeData8(t *testing.T) {
	i8_1 := int8(1)
	i8_2 := int8(-1)
	u8_1 := uint8(i8_1)
	u8_2 := uint8(i8_2)

	b_i8_1, _ := NewVal(i8_1);
	b_i8_2, _ := NewVal(i8_2);
	b_u8_1, _ := NewVal(u8_1);
	b_u8_2, _ := NewVal(u8_2);

	v_i8_1, _ := b_i8_1.GetInt8()
	v_i8_2, _ := b_i8_2.GetInt8()
	v_u8_1, _ := b_u8_1.GetUInt8()
	v_u8_2, _ := b_u8_2.GetUInt8()

	assert.Equal(t, i8_1, v_i8_1)
	assert.Equal(t, i8_2, v_i8_2)
	assert.Equal(t, u8_1, v_u8_1)
	assert.Equal(t, u8_2, v_u8_2)
}

func TestTypeData16(t *testing.T) {
	i16_1 := int16(1)
	i16_2 := int16(-1)
	u16_1 := uint16(i16_1)
	u16_2 := uint16(i16_2)

	b_i16_1, _ := NewVal(i16_1);
	b_i16_2, _ := NewVal(i16_2);
	b_u16_1, _ := NewVal(u16_1);
	b_u16_2, _ := NewVal(u16_2);

	v_i16_1, _ := b_i16_1.GetInt16()
	v_i16_2, _ := b_i16_2.GetInt16()
	v_u16_1, _ := b_u16_1.GetUInt16()
	v_u16_2, _ := b_u16_2.GetUInt16()

	assert.Equal(t, i16_1, v_i16_1)
	assert.Equal(t, i16_2, v_i16_2)
	assert.Equal(t, u16_1, v_u16_1)
	assert.Equal(t, u16_2, v_u16_2)
}

func TestTypeData32(t *testing.T) {
	i32_1 := int32(1)
	i32_2 := int32(-1)
	u32_1 := uint32(i32_1)
	u32_2 := uint32(i32_2)
	f32_1 := float32(1.0)
	f32_2 := float32(-1.0)

	b_i32_1, _ := NewVal(i32_1)
	b_i32_2, _ := NewVal(i32_2)
	b_u32_1, _ := NewVal(u32_1)
	b_u32_2, _ := NewVal(u32_2)
	b_f32_1, _ := NewVal(f32_1)
	b_f32_2, _ := NewVal(f32_2)

	v_i32_1, _ := b_i32_1.GetInt32()
	v_i32_2, _ := b_i32_2.GetInt32()
	v_u32_1, _ := b_u32_1.GetUInt32()
	v_u32_2, _ := b_u32_2.GetUInt32()
  v_f32_1, _ := b_f32_1.GetFloat32()
	v_f32_2, _ := b_f32_2.GetFloat32()

	assert.Equal(t, i32_1, v_i32_1)
	assert.Equal(t, i32_2, v_i32_2)
	assert.Equal(t, u32_1, v_u32_1)
	assert.Equal(t, u32_2, v_u32_2)
	assert.Equal(t, f32_1, v_f32_1)
	assert.Equal(t, f32_2, v_f32_2)
}

func TestTypeData64(t *testing.T) {
	i64_1 := int64(1)
	i64_2 := int64(-1)
	u64_1 := uint64(i64_1)
	u64_2 := uint64(i64_2)
	f64_1 := 1.0
	f64_2 := -1.0

	b_i64_1, _ := NewVal(i64_1)
	b_i64_2, _ := NewVal(i64_2)
	b_u64_1, _ := NewVal(u64_1)
	b_u64_2, _ := NewVal(u64_2)
	b_f64_1, _ := NewVal(f64_1)
	b_f64_2, _ := NewVal(f64_2)

	v_i64_1, _ := b_i64_1.GetInt64()
	v_i64_2, _ := b_i64_2.GetInt64()
	v_u64_1, _ := b_u64_1.GetUInt64()
	v_u64_2, _ := b_u64_2.GetUInt64()
  v_f64_1, _ := b_f64_1.GetFloat64()
	v_f64_2, _ := b_f64_2.GetFloat64()

	assert.Equal(t, i64_1, v_i64_1)
	assert.Equal(t, i64_2, v_i64_2)
	assert.Equal(t, u64_1, v_u64_1)
	assert.Equal(t, u64_2, v_u64_2)
	assert.Equal(t, f64_1, v_f64_1)
	assert.Equal(t, f64_2, v_f64_2)

	// t.Logf("i64_1: %d  v_i64_1: %d\n", i64_1, v_i64_1)
	// t.Logf("i64_2: %d  v_i64_2: %d\n", i64_2, v_i64_2)
	// t.Logf("u64_1: %d  v_u64_1: %d\n", u64_1, v_u64_1)
	// t.Logf("u64_2: %d  v_u64_2: %d\n", u64_2, v_u64_2)
	// t.Logf("f64_1: %f  v_f64_1: %f\n", f64_1, v_f64_1)
	// t.Logf("f64_2: %f  v_f64_2: %f\n", f64_2, v_f64_2)
}

func TestTypeDataOther(t *testing.T) {

	t1 := time.Now().Local(); b_t1, _ := NewVal(t1);
	v_t1, _ := b_t1.GetTime()
	assert.Equal(t, t1, v_t1)
	//t.Logf("t1: %s  v_t1: %s\n", t1, v_t1)

	du1 := time.Since(t1); b_du1, _ := NewVal(du1);
	v_du1, _ := b_du1.GetDuration()
	assert.Equal(t, du1, v_du1)
	//t.Logf("t1: %s  v_du1: %s\n", du1, v_du1)

	b1 := []byte("hello world"); b_b1, _ := NewVal(b1);
	v_b1, _ := b_b1.GetBytes()
	assert.Equal(t, b1, v_b1)
}