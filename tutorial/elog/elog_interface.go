package elog

import "github.com/ziyht/eden_go/elog"

func runInterfaceTutorial(){

	// interfaceTagTutorail()
	// interfaceLevelTutorail()
	// interfaceMultiFileTutorail()
	// interfaceSLevelTutorail()
	interfaceColorTutorail()
}

func interfaceTagTutorail(){

	logger := elog.LoggerFromInterface(map[string]string{
		"tag": "testtag",
	})
	log := logger.Log()
  log.Debugf("interfaceTagTutorail: tag should be testtag")  // in default setting, it will not output
	log.Infof( "interfaceTagTutorail: tag should be testtag")
	log.Warnf( "interfaceTagTutorail: tag should be testtag")
	log.Errorf("interfaceTagTutorail: tag should be testtag")

	log = logger.Log(elog.Opt().Tags("level1", "level2"))
  log.Debugf("interfaceTagTutorail: tag should be testtag.level1.level2")  // in default setting, it will not output
	log.Infof( "interfaceTagTutorail: tag should be testtag.level1.level2")
	log.Warnf( "interfaceTagTutorail: tag should be testtag.level1.level2")
	log.Errorf("interfaceTagTutorail: tag should be testtag.level1.level2")

	log = logger.Log(elog.Opt().ReTags("level1", "level2"))
  log.Debugf("interfaceTagTutorail: tag should be level1.level2")  // in default setting, it will not output
	log.Infof( "interfaceTagTutorail: tag should be level1.level2")
	log.Warnf( "interfaceTagTutorail: tag should be level1.level2")
	log.Errorf("interfaceTagTutorail: tag should be level1.level2")
}

func interfaceLevelTutorail(){

	logger := elog.LoggerFromInterface(map[string]string{
		"tag": "interfaceLevelTutorail",
		"f_level": "debug",
		"c_level": "debug",
	})
	log := logger.Log(elog.Opt().Tags("debug"))
  log.Debugf("output debug")  // in default setting, it will not output
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")

	logger = elog.LoggerFromInterface(map[string]interface{}{
		"tag": "interfaceLevelTutorail",
		"f_level": "info",
		"c_level": "info",
	})
	log = logger.Log(elog.Opt().Tags("info"))
  log.Debugf("output debug")  // in default setting, it will not output
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")

	logger = elog.LoggerFromInterface(map[string]interface{}{
		"tag": "interfaceLevelTutorail",
		"f_level": []string{"debug", "debug"},
		"c_level": []string{"debug", "debug"},
	})
	log = logger.Log(elog.Opt().Tags("debug-debug"))
  log.Debugf("output debug")  // in default setting, it will not output
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")

	logger = elog.LoggerFromInterface(map[string]interface{}{
		"tag": "interfaceLevelTutorail",
		"f_level": []string{"info", "error"},
		"c_level": []string{"info", "error"},
	})
	log = logger.Log(elog.Opt().Tags("info-error"))
  log.Debugf("output debug")  // in default setting, it will not output
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")

	logger = elog.LoggerFromInterface(map[string]interface{}{
		"tag": "interfaceLevelTutorail",
		"f_level": []string{"debug", "info"},
		"c_level": []string{"warn", "error"},
	})
	log = logger.Log(elog.Opt().Tags("f:debug-info|c:warn-error"))
  log.Debugf("output debug")  // in default setting, it will not output
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")
}

func interfaceMultiFileTutorail(){

	logger := elog.LoggerFromInterface([]interface{}{
		map[string]interface{}{
			"tag": "interfaceMultiFileTutorail.file1.debug-info",
			"filename": "<APP>_debug",
			"level": []string{"debug", "info"},
		},
		map[string]interface{}{
			"tag": "interfaceMultiFileTutorail.file2.warn-error",
			"filename": "<APP>_warn",
			"level": []string{"warn", "error"},
		},	
	})
	log := logger.Log()
  log.Debugf("output debug")  // in default setting, it will not output
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")
}

func interfaceSLevelTutorail(){

	logger := elog.LoggerFromInterface([]interface{}{
		map[string]interface{}{
			"tag": "interfaceSLevelTutorail.debug",
			"console": "stdout",
			"level": "debug",
			"slevel": "debug",
		},
	})
	log := logger.Log()
  log.Debugf("output debug")  // in default setting, it will not output
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")


	logger = elog.LoggerFromInterface([]interface{}{
		map[string]interface{}{
			"tag": "interfaceSLevelTutorail.info",
			"console": "stdout",
			"level": "debug",
			"slevel": "info",
		},	
	})
	log = logger.Log()
  log.Debugf("output debug")  // in default setting, it will not output
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")

	logger = elog.LoggerFromInterface([]interface{}{
		map[string]interface{}{
			"tag": "interfaceSLevelTutorail.info-info",
			"console": "stdout",
			"level": "debug",
			"slevel": []string{"info", "warn"},
		},	
	})
	log = logger.Log()
  log.Debugf("output debug")  // in default setting, it will not output
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")
}

func interfaceColorTutorail(){

	logger := elog.LoggerFromInterface([]interface{}{
		map[string]interface{}{
			"tag"    : "interfaceColorTutorail.color=false",
			"console": "stdout",
			"level"  : "debug",
			"color"  : "false",
			"slevel" : "fatal",
		},
	})
	log := logger.Log()
  log.Debugf("output debug")  // in default setting, it will not output
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")

	logger = elog.LoggerFromInterface([]interface{}{
		map[string]interface{}{
			"tag"    : "interfaceColorTutorail.color=true",
			"console": "stdout",
			"level"  : "debug",
			"color"  : "true",
			"slevel" : "fatal",
		},
	})
	log = logger.Log()
  log.Debugf("output debug")  // in default setting, it will not output
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")

	logger = elog.LoggerFromInterface([]interface{}{
		map[string]interface{}{
			"tag"    : "interfaceColorTutorail.color=auto",
			"console": "stdout",
			"level"  : "debug",
			"color"  : "auto",
			"slevel" : "fatal",
		},
	})
	log = logger.Log()
  log.Debugf("output debug")  // in default setting, it will not output
	log.Infof( "output info")
	log.Warnf( "output warn")
	log.Errorf("output error")
}