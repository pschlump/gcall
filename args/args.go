package args

import "fmt"

// 0. Implements args -> map[name]string for args

type ArgRetType struct {
	Set       map[string]string
	Remainder []string
	Config    *ArgConfigType
}

type AnArgType struct {
	Name     string
	Abrev    string
	Default  string
	Required bool
	NoValue  bool // set to "1" if specifed
}

type ArgConfigType struct {
	Config map[string]AnArgType
}

func GetArgs(argConfig *ArgConfigType, off int, argc []string) (rv *ArgRetType, err error) {
	rv = &ArgRetType{
		Set:       make(map[string]string),
		Remainder: []string{},
	}
	// for ii, arg := range argc {
	for ii := off; ii < len(argc); ii++ {
		arg := argc[ii]
		if match, name, aa := argConfig.IsArg(arg); match {
			if aa.NoValue {
				rv.Set[name] = "1"
			} else {
				if ii+1 < len(argc) {
					rv.Set[name] = argc[ii+1]
					ii++
				} else {
					err = fmt.Errorf("Argument requring a value has no value, last argument")
					return
				}
			}
		} else {
			if arg[0:1] == "-" {
				rv.Remainder = append(rv.Remainder, arg)
				if ii+1 < len(argc) {
					ii++
					rv.Remainder = append(rv.Remainder, argc[ii])
				}
			} else {
				rv.Remainder = append(rv.Remainder, arg)
			}
		}
	}
	// Set Defaults
	for key, vv := range argConfig.Config {
		if vv.Default != "" {
			if _, ok := rv.Set[key]; !ok {
				rv.Set[key] = vv.Default
			}
		}
	}
	rv.Config = argConfig
	return
}

func (bb *ArgConfigType) IsArg(s string) (match bool, name string, aa AnArgType) {
	for key, vv := range bb.Config {
		if s == vv.Name || s == vv.Abrev {
			match, name, aa = true, key, vv
			return
		}
	}
	return
}

func (aa *ArgRetType) Usage(cmd string) {
	fmt.Printf("Usage: %s ", cmd)
	for _, vv := range aa.Config.Config {
		if vv.Name != "" && vv.Abrev != "" && vv.Required == false && vv.NoValue == false {
			fmt.Printf("[%s|%s] <value> ", vv.Name, vv.Abrev)
		} else if vv.Name != "" && vv.Abrev != "" && vv.Required == false && vv.NoValue == true {
			fmt.Printf("[%s|%s] ", vv.Name, vv.Abrev)
		} else if vv.Name != "" && vv.Required == false && vv.NoValue == false {
			fmt.Printf("[%s] <value> ", vv.Name)
		} else if vv.Name != "" && vv.Required == false && vv.NoValue == true {
			fmt.Printf("[%s] ", vv.Name)
		} else if vv.Abrev != "" && vv.Required == false && vv.NoValue == false {
			fmt.Printf("[%s] <value> ", vv.Abrev)
		} else if vv.Abrev != "" && vv.Required == false && vv.NoValue == true {
			fmt.Printf("[%s] ", vv.Abrev)
		} else if vv.Name != "" && vv.Abrev != "" && vv.Required == true && vv.NoValue == false {
			fmt.Printf("%s|%s <value> ", vv.Name, vv.Abrev)
		} else if vv.Name != "" && vv.Abrev != "" && vv.Required == true && vv.NoValue == true {
			fmt.Printf("%s|%s ", vv.Name, vv.Abrev)
		} else if vv.Name != "" && vv.Required == true && vv.NoValue == false {
			fmt.Printf("%s <value> ", vv.Name)
		} else if vv.Name != "" && vv.Required == true && vv.NoValue == true {
			fmt.Printf("%s ", vv.Name)
		} else if vv.Abrev != "" && vv.Required == true && vv.NoValue == false {
			fmt.Printf("%s <value> ", vv.Abrev)
		} else if vv.Abrev != "" && vv.Required == true && vv.NoValue == true {
			fmt.Printf("%s ", vv.Abrev)
		}
	}
	fmt.Printf("\n")
}
