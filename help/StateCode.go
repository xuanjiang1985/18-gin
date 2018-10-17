package help

type StateCode struct {
	Success        int
	ParamError     int
	ParamEmpty     int
	DbSelectError  int
	DbInsertError  int
	InterfaceError int
	UnAuthError    int
}

type BusinessCode struct {
	DbUserOk   int
	DbUserBand int
}

var Statecode *StateCode
var Bcode *BusinessCode

func init() {

	Statecode = &StateCode{
		Success:        0,
		ParamError:     1001,
		ParamEmpty:     1002,
		DbSelectError:  2001,
		DbInsertError:  2002,
		InterfaceError: 3001,
		UnAuthError:    4001,
	}

	Bcode = &BusinessCode{
		DbUserOk:   1,
		DbUserBand: 0,
	}

}
