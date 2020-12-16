package moego

type Keyword string

const (
	Break       Keyword = "break"
	Default             = "default"
	Func                = "func"
	Interface           = "interface"
	Select              = "select"
	Case                = "case"
	Defer               = "defer"
	Go                  = "go"
	Map                 = "map"
	Struct              = "struct"
	Chan                = "chan"
	Else                = "else"
	Goto                = "goto"
	Package             = "package"
	Switch              = "switch"
	Const               = "const"
	Fallthrough         = "fallthrough"
	If                  = "if"
	Range               = "range"
	Type                = "type"
	Continue            = "continue"
	For                 = "for"
	Import              = "import"
	Return              = "return"
	Var                 = "var"
)

var KeywordColor = map[Keyword]Color{
	Break:       FG_CYAN,
	Default:     FG_CYAN,
	Interface:   FG_CYAN,
	Select:      FG_CYAN,
	Case:        FG_CYAN,
	Defer:       FG_CYAN,
	Go:          FG_CYAN,
	Map:         FG_CYAN,
	Struct:      FG_CYAN,
	Chan:        FG_CYAN,
	Else:        FG_CYAN,
	Goto:        FG_CYAN,
	Switch:      FG_CYAN,
	Const:       FG_CYAN,
	Fallthrough: FG_CYAN,
	Return:      FG_CYAN,
	Range:       FG_CYAN,
	Type:        FG_CYAN,
	Continue:    FG_CYAN,
	For:         FG_CYAN,
	If:          FG_CYAN,
	Package:     FG_CYAN,
	Import:      FG_CYAN,
	Func:        FG_CYAN,
	Var:         FG_CYAN,
}
