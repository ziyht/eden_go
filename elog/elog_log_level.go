package elog

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

type LevelEnabler interface {
	Enabled(Level) bool
}

type Level zapcore.Level

func (l Level) Enabled(lvl Level) bool {
	return lvl >= l
}

const (
	LEVEL_DEBUG  = Level(zapcore.DebugLevel)
	LEVEL_INFO   = Level(zapcore.InfoLevel)
	LEVEL_WARN   = Level(zapcore.WarnLevel)
	LEVEL_ERROR  = Level(zapcore.ErrorLevel)
	LEVEL_FATAL  = Level(zapcore.FatalLevel)
	LEVEL_DPANIC = Level(zapcore.DPanicLevel)
	LEVEL_PANIC  = Level(zapcore.PanicLevel)
	LEVEL_ALL    = LEVEL_DEBUG
	LEVEL_NONE   = LEVEL_FATAL
	LEVEL_OFF    = LEVEL_FATAL + 1

	LEVELS_DEBUG  = "debug"
	LEVELS_INFO   = "info"
	LEVELS_WARN   = "warn"
	LEVELS_ERROR  = "error"
	LEVELS_FATAL  = "fatal"
	LEVELS_DPANIC = "dpanic"
	LEVELS_PANIC  = "panic"

	LEVELT_DEBUG  = "DBG"
	LEVELT_INFO   = "INF"
	LEVELT_WARN   = "WRN"
	LEVELT_ERROR  = "ERR"
	LEVELT_FATAL  = "FTL"
	LEVELT_DPANIC = "DPN"
	LEVELT_PANIC  = "PNC"	
)

var (
	supportLevelStrs = []string{LEVELS_DEBUG, LEVELS_INFO, LEVELS_WARN, LEVELS_ERROR, LEVELS_FATAL, LEVELS_DPANIC, LEVELS_PANIC}
)

var _levelTags = [...]string {
	LEVELT_DEBUG, 
	LEVELT_INFO, 
	LEVELT_WARN,
	LEVELT_ERROR,
	LEVELT_FATAL,
	LEVELT_DPANIC,
	LEVELT_PANIC,
}

var _levelTagsColored = [...]string {
	getColoredStr(getColorByLevel(LEVEL_DEBUG ), LEVELT_DEBUG),
  getColoredStr(getColorByLevel(LEVEL_INFO  ), LEVELT_INFO),
  getColoredStr(getColorByLevel(LEVEL_WARN  ), LEVELT_WARN),
  getColoredStr(getColorByLevel(LEVEL_ERROR ), LEVELT_ERROR),
  getColoredStr(getColorByLevel(LEVEL_FATAL ), LEVELT_FATAL),
  getColoredStr(getColorByLevel(LEVEL_DPANIC), LEVELT_DPANIC),
  getColoredStr(getColorByLevel(LEVEL_PANIC ), LEVELT_PANIC),
}

func getLevelByStr(levelStr string) (Level, error) {

	switch levelStr {
	case LEVELS_DEBUG : return LEVEL_DEBUG , nil
	case LEVELS_INFO  : return LEVEL_INFO  , nil
	case LEVELS_WARN  : return LEVEL_WARN  , nil
	case LEVELS_ERROR : return LEVEL_ERROR , nil
	case LEVELS_FATAL : return LEVEL_FATAL , nil
	case LEVELS_DPANIC: return LEVEL_DPANIC, nil
	case LEVELS_PANIC : return LEVEL_PANIC , nil
	}

	return LEVEL_ALL, fmt.Errorf("unsupported Level str '%s', you can use: %s", levelStr, supportLevelStrs)
}

func getTagByLevel(l Level) string {
	return _levelTags[l + 1]
}

func getColoredTagByLevel(l Level) string {
	return _levelTagsColored[l + 1]
}