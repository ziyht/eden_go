package elog

// InitFromYml this will init loggers from a yaml file, note that the logger with same name will be replaced with new one
// you can get the sample cfg content from SampleCfg
func InitFromYml(file string) {

	cfgs := readCfgFromYaml(file)

	dfcfg := cfgs.Cfgs[cfgDefaultKey]

	initDfLogger(dfcfg)

	for name, cfg := range cfgs.Cfgs {
		if name == cfgDefaultKey{
			continue
		}

		NewLogger(name, cfg)
	}
}

func init() {
	initSyslog()
	initDfLogger(nil)
}