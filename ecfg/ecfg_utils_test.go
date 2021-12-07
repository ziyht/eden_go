package ecfg

import "testing"


type testCfg struct {

	Key1 string
	Key2 int
}

func Test_parsingCfgFromStr(t *testing.T){
	{
		var d testCfg
		content := `
Key1: val1
Key2: 2
`
		err := parsingCfgFromStr(content, "yml", "", &d)
		if err != nil {
			t.Fatalf("err occured: %s", err)
		}
		if d.Key1 != "val1" { t.Errorf("val not match") }
		if d.Key2 != 2      { t.Errorf("val not match") }
	}
	
}