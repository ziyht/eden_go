package elog

import (
	"sync"
	"time"

	"go.uber.org/zap/zapcore"
)


var (
	_dfTimeFmt = "2006-01-02_15:04:05.000"

	_timeEncoders = map[string]zapcore.TimeEncoder {}

	_itemEncoderMu sync.Mutex
)

// getTimeEncoder get a timeEncoder by specified fmt, if fmt not set, using _dfTimeFmt
func getTimeEncoder(fmt ...string) zapcore.TimeEncoder {

	_itemEncoderMu.Lock()
	defer _itemEncoderMu.Unlock()

	_fmt := _dfTimeFmt
	if len(fmt) > 0 {
		_fmt = fmt[0]
	}

	enc, exist := _timeEncoders[_fmt]

	// create a new one if not exist
	if !exist {
		enc = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(_fmt))
		}
		_timeEncoders[_fmt] = enc
	}

	return enc
}


func _coloredLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder){
	enc.AppendString(getColoredTagByLevel(Level(l)))
}
func _noColoredLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder){
	enc.AppendString(getTagByLevel(Level(l)))
}

func getLevelEncoder(color colorSwitch) zapcore.LevelEncoder {

	if color == ColorOn{
		return _coloredLevelEncoder
	}
	if color == ColorOff {
		return _noColoredLevelEncoder
	}
  
	// todo: add auto color funcs
	return _coloredLevelEncoder
}