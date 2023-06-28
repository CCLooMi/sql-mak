package god

type EXP struct {
	exp  interface{}
	args []interface{}
}

func NewEXP(exp string, args ...interface{}) *EXP {
	return &EXP{
		exp:  exp,
		args: args,
	}
}

func NewEXPFromSQLSM(sm *SQLSM) *EXP {
	return &EXP{
		exp:  sm,
		args: sm.args,
	}
}

func (e *EXP) Exp() string {
	if s, ok := e.exp.(string); ok {
		return s
	}
	return e.exp.(*SQLSM).ExpSQL()
}

func (e *EXP) Args() []interface{} {
	if _, ok := e.exp.(string); ok {
		return e.args
	}
	return e.exp.(*SQLSM).Args()
}

func Now() *EXP {
	return &EXP{
		exp: "NOW()",
	}
}

func UUID() *EXP {
	return &EXP{
		exp: "REPLACE(UUID(),'-','')",
	}
}

func Exp(sm *SQLSM) *EXP {
	return &EXP{
		exp:  sm,
		args: sm.args,
	}
}

func ExpStr(exp string, args ...interface{}) *EXP {
	return &EXP{
		exp:  exp,
		args: args,
	}
}
