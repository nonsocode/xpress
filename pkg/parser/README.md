<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# parser

```go
import "github.com/nonsocode/xpress/pkg/parser"
```

## Index

- [Constants](<#constants>)
- [Variables](<#variables>)
- [type Array](<#Array>)
  - [func NewArray\(values \[\]Expr\) \*Array](<#NewArray>)
  - [func \(a \*Array\) Accept\(ctx context.Context, v Visitor\) EvaluationResult](<#Array.Accept>)
- [type Binary](<#Binary>)
  - [func NewBinary\(left Expr, operator Token, right Expr\) \*Binary](<#NewBinary>)
  - [func \(b \*Binary\) Accept\(ctx context.Context, v Visitor\) EvaluationResult](<#Binary.Accept>)
- [type Call](<#Call>)
  - [func NewCall\(callee Expr, arguments \[\]Expr\) \*Call](<#NewCall>)
  - [func \(c \*Call\) Accept\(ctx context.Context, v Visitor\) EvaluationResult](<#Call.Accept>)
- [type EvaluationError](<#EvaluationError>)
  - [func NewEvaluationError\(message string, args ...interface\{\}\) \*EvaluationError](<#NewEvaluationError>)
  - [func \(e \*EvaluationError\) Error\(\) string](<#EvaluationError.Error>)
- [type EvaluationResult](<#EvaluationResult>)
- [type Evaluator](<#Evaluator>)
  - [func NewInterpreter\(\) \*Evaluator](<#NewInterpreter>)
  - [func \(i \*Evaluator\) AddMember\(name string, member interface\{\}\) error](<#Evaluator.AddMember>)
  - [func \(i \*Evaluator\) Evaluate\(ctx context.Context, expr Expr\) \(interface\{\}, error\)](<#Evaluator.Evaluate>)
  - [func \(i \*Evaluator\) SetMembers\(members map\[string\]interface\{\}\) error](<#Evaluator.SetMembers>)
  - [func \(i \*Evaluator\) SetTimeout\(timeout time.Duration\)](<#Evaluator.SetTimeout>)
- [type Expr](<#Expr>)
- [type Get](<#Get>)
  - [func NewGet\(object Expr, name Token\) \*Get](<#NewGet>)
  - [func \(g \*Get\) Accept\(ctx context.Context, v Visitor\) EvaluationResult](<#Get.Accept>)
- [type Grouping](<#Grouping>)
  - [func NewGrouping\(expression Expr\) \*Grouping](<#NewGrouping>)
  - [func \(g \*Grouping\) Accept\(ctx context.Context, v Visitor\) EvaluationResult](<#Grouping.Accept>)
- [type Index](<#Index>)
  - [func NewIndex\(object Expr, index Expr\) \*Index](<#NewIndex>)
  - [func \(i \*Index\) Accept\(ctx context.Context, v Visitor\) EvaluationResult](<#Index.Accept>)
- [type Interpreter](<#Interpreter>)
- [type Lexer](<#Lexer>)
  - [func NewLexer\(source string\) \*Lexer](<#NewLexer>)
- [type Literal](<#Literal>)
  - [func NewLiteral\(value interface\{\}, raw string\) \*Literal](<#NewLiteral>)
  - [func \(l \*Literal\) Accept\(ctx context.Context, v Visitor\) EvaluationResult](<#Literal.Accept>)
- [type Map](<#Map>)
  - [func NewMap\(entries \[\]\*MapEntry\) \*Map](<#NewMap>)
  - [func \(m \*Map\) Accept\(ctx context.Context, v Visitor\) EvaluationResult](<#Map.Accept>)
- [type MapEntry](<#MapEntry>)
  - [func NewMapEntry\(key Expr, value Expr\) \*MapEntry](<#NewMapEntry>)
  - [func \(me \*MapEntry\) Accept\(ctx context.Context, v Visitor\) EvaluationResult](<#MapEntry.Accept>)
- [type Optional](<#Optional>)
  - [func NewOptional\(left Expr\) \*Optional](<#NewOptional>)
  - [func \(o \*Optional\) Accept\(ctx context.Context, v Visitor\) EvaluationResult](<#Optional.Accept>)
- [type ParseError](<#ParseError>)
  - [func \(p \*ParseError\) Accept\(ctx context.Context, v Visitor\) EvaluationResult](<#ParseError.Accept>)
  - [func \(pe \*ParseError\) Error\(\) string](<#ParseError.Error>)
- [type Parser](<#Parser>)
  - [func NewParser\(source string\) \*Parser](<#NewParser>)
  - [func \(p \*Parser\) Parse\(\) \(exp Expr\)](<#Parser.Parse>)
- [type Template](<#Template>)
  - [func NewTemplate\(expressions \[\]Expr\) \*Template](<#NewTemplate>)
  - [func \(t \*Template\) Accept\(ctx context.Context, v Visitor\) EvaluationResult](<#Template.Accept>)
- [type Ternary](<#Ternary>)
  - [func NewTernary\(condition Expr, trueExpr Expr, falseExpr Expr\) \*Ternary](<#NewTernary>)
  - [func \(t \*Ternary\) Accept\(ctx context.Context, v Visitor\) EvaluationResult](<#Ternary.Accept>)
- [type Token](<#Token>)
  - [func \(t Token\) String\(\) string](<#Token.String>)
- [type TokenType](<#TokenType>)
  - [func \(t TokenType\) String\(\) string](<#TokenType.String>)
- [type Unary](<#Unary>)
  - [func NewUnary\(operator Token, right Expr\) \*Unary](<#NewUnary>)
  - [func \(u \*Unary\) Accept\(ctx context.Context, v Visitor\) EvaluationResult](<#Unary.Accept>)
- [type Variable](<#Variable>)
  - [func NewVariable\(name Token\) \*Variable](<#NewVariable>)
  - [func \(v \*Variable\) Accept\(ctx context.Context, vis Visitor\) EvaluationResult](<#Variable.Accept>)
- [type Visitor](<#Visitor>)


## Constants

<a name="DefaultTimeout"></a>

```go
const (
    // DefaultTimeout is the default timeout for evaluating expressions.
    DefaultTimeout = 10 * time.Millisecond
)
```

## Variables

<a name="EvaluationCancelledErrror"></a>

```go
var (
    EvaluationCancelledErrror = NewEvaluationError("evaluation cancelled")
)
```

<a name="Array"></a>
## type [Array](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L63-L65>)



```go
type Array struct {
    // contains filtered or unexported fields
}
```

<a name="NewArray"></a>
### func [NewArray](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L166>)

```go
func NewArray(values []Expr) *Array
```



<a name="Array.Accept"></a>
### func \(\*Array\) [Accept](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L170>)

```go
func (a *Array) Accept(ctx context.Context, v Visitor) EvaluationResult
```



<a name="Binary"></a>
## type [Binary](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L9-L13>)



```go
type Binary struct {
    // contains filtered or unexported fields
}
```

<a name="NewBinary"></a>
### func [NewBinary](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L86>)

```go
func NewBinary(left Expr, operator Token, right Expr) *Binary
```



<a name="Binary.Accept"></a>
### func \(\*Binary\) [Accept](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L90>)

```go
func (b *Binary) Accept(ctx context.Context, v Visitor) EvaluationResult
```



<a name="Call"></a>
## type [Call](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L53-L56>)



```go
type Call struct {
    // contains filtered or unexported fields
}
```

<a name="NewCall"></a>
### func [NewCall](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L150>)

```go
func NewCall(callee Expr, arguments []Expr) *Call
```



<a name="Call.Accept"></a>
### func \(\*Call\) [Accept](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L154>)

```go
func (c *Call) Accept(ctx context.Context, v Visitor) EvaluationResult
```



<a name="EvaluationError"></a>
## type [EvaluationError](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/evalutator.go#L19-L21>)



```go
type EvaluationError struct {
    // contains filtered or unexported fields
}
```

<a name="NewEvaluationError"></a>
### func [NewEvaluationError](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/evalutator.go#L44>)

```go
func NewEvaluationError(message string, args ...interface{}) *EvaluationError
```



<a name="EvaluationError.Error"></a>
### func \(\*EvaluationError\) [Error](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/evalutator.go#L48>)

```go
func (e *EvaluationError) Error() string
```



<a name="EvaluationResult"></a>
## type [EvaluationResult](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/types.go#L31-L34>)



```go
type EvaluationResult interface {
    Get() interface{}
    Error() error
}
```

<a name="Evaluator"></a>
## type [Evaluator](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/evalutator.go#L13-L17>)



```go
type Evaluator struct {
    // contains filtered or unexported fields
}
```

<a name="NewInterpreter"></a>
### func [NewInterpreter](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/evalutator.go#L33>)

```go
func NewInterpreter() *Evaluator
```



<a name="Evaluator.AddMember"></a>
### func \(\*Evaluator\) [AddMember](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/evalutator.go#L52>)

```go
func (i *Evaluator) AddMember(name string, member interface{}) error
```



<a name="Evaluator.Evaluate"></a>
### func \(\*Evaluator\) [Evaluate](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/evalutator.go#L707>)

```go
func (i *Evaluator) Evaluate(ctx context.Context, expr Expr) (interface{}, error)
```



<a name="Evaluator.SetMembers"></a>
### func \(\*Evaluator\) [SetMembers](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/evalutator.go#L59>)

```go
func (i *Evaluator) SetMembers(members map[string]interface{}) error
```



<a name="Evaluator.SetTimeout"></a>
### func \(\*Evaluator\) [SetTimeout](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/evalutator.go#L40>)

```go
func (i *Evaluator) SetTimeout(timeout time.Duration)
```



<a name="Expr"></a>
## type [Expr](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/parser.go#L10-L12>)



```go
type Expr interface {
    Accept(context.Context, Visitor) EvaluationResult
}
```

<a name="Get"></a>
## type [Get](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L44-L47>)



```go
type Get struct {
    // contains filtered or unexported fields
}
```

<a name="NewGet"></a>
### func [NewGet](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L134>)

```go
func NewGet(object Expr, name Token) *Get
```



<a name="Get.Accept"></a>
### func \(\*Get\) [Accept](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L138>)

```go
func (g *Get) Accept(ctx context.Context, v Visitor) EvaluationResult
```



<a name="Grouping"></a>
## type [Grouping](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L15-L17>)



```go
type Grouping struct {
    // contains filtered or unexported fields
}
```

<a name="NewGrouping"></a>
### func [NewGrouping](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L94>)

```go
func NewGrouping(expression Expr) *Grouping
```



<a name="Grouping.Accept"></a>
### func \(\*Grouping\) [Accept](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L98>)

```go
func (g *Grouping) Accept(ctx context.Context, v Visitor) EvaluationResult
```



<a name="Index"></a>
## type [Index](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L58-L61>)



```go
type Index struct {
    // contains filtered or unexported fields
}
```

<a name="NewIndex"></a>
### func [NewIndex](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L158>)

```go
func NewIndex(object Expr, index Expr) *Index
```



<a name="Index.Accept"></a>
### func \(\*Index\) [Accept](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L162>)

```go
func (i *Index) Accept(ctx context.Context, v Visitor) EvaluationResult
```



<a name="Interpreter"></a>
## type [Interpreter](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/types.go#L27-L29>)



```go
type Interpreter interface {
    Evaluate(expr Expr) EvaluationResult
}
```

<a name="Lexer"></a>
## type [Lexer](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/lex.go#L80-L88>)



```go
type Lexer struct {
    // contains filtered or unexported fields
}
```

<a name="NewLexer"></a>
### func [NewLexer](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/lex.go#L100>)

```go
func NewLexer(source string) *Lexer
```



<a name="Literal"></a>
## type [Literal](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L19-L22>)



```go
type Literal struct {
    // contains filtered or unexported fields
}
```

<a name="NewLiteral"></a>
### func [NewLiteral](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L102>)

```go
func NewLiteral(value interface{}, raw string) *Literal
```



<a name="Literal.Accept"></a>
### func \(\*Literal\) [Accept](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L106>)

```go
func (l *Literal) Accept(ctx context.Context, v Visitor) EvaluationResult
```



<a name="Map"></a>
## type [Map](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L76-L78>)



```go
type Map struct {
    // contains filtered or unexported fields
}
```

<a name="NewMap"></a>
### func [NewMap](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L190>)

```go
func NewMap(entries []*MapEntry) *Map
```



<a name="Map.Accept"></a>
### func \(\*Map\) [Accept](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L194>)

```go
func (m *Map) Accept(ctx context.Context, v Visitor) EvaluationResult
```



<a name="MapEntry"></a>
## type [MapEntry](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L80-L83>)



```go
type MapEntry struct {
    // contains filtered or unexported fields
}
```

<a name="NewMapEntry"></a>
### func [NewMapEntry](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L198>)

```go
func NewMapEntry(key Expr, value Expr) *MapEntry
```



<a name="MapEntry.Accept"></a>
### func \(\*MapEntry\) [Accept](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L202>)

```go
func (me *MapEntry) Accept(ctx context.Context, v Visitor) EvaluationResult
```



<a name="Optional"></a>
## type [Optional](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L49-L51>)



```go
type Optional struct {
    // contains filtered or unexported fields
}
```

<a name="NewOptional"></a>
### func [NewOptional](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L142>)

```go
func NewOptional(left Expr) *Optional
```



<a name="Optional.Accept"></a>
### func \(\*Optional\) [Accept](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L146>)

```go
func (o *Optional) Accept(ctx context.Context, v Visitor) EvaluationResult
```



<a name="ParseError"></a>
## type [ParseError](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L71-L74>)



```go
type ParseError struct {
    // contains filtered or unexported fields
}
```

<a name="ParseError.Accept"></a>
### func \(\*ParseError\) [Accept](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L186>)

```go
func (p *ParseError) Accept(ctx context.Context, v Visitor) EvaluationResult
```



<a name="ParseError.Error"></a>
### func \(\*ParseError\) [Error](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L182>)

```go
func (pe *ParseError) Error() string
```



<a name="Parser"></a>
## type [Parser](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L33-L36>)



```go
type Parser struct {
    // contains filtered or unexported fields
}
```

<a name="NewParser"></a>
### func [NewParser](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/parser.go#L15>)

```go
func NewParser(source string) *Parser
```



<a name="Parser.Parse"></a>
### func \(\*Parser\) [Parse](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/parser.go#L21>)

```go
func (p *Parser) Parse() (exp Expr)
```



<a name="Template"></a>
## type [Template](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L29-L31>)



```go
type Template struct {
    // contains filtered or unexported fields
}
```

<a name="NewTemplate"></a>
### func [NewTemplate](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L118>)

```go
func NewTemplate(expressions []Expr) *Template
```



<a name="Template.Accept"></a>
### func \(\*Template\) [Accept](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L122>)

```go
func (t *Template) Accept(ctx context.Context, v Visitor) EvaluationResult
```



<a name="Ternary"></a>
## type [Ternary](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L38-L42>)



```go
type Ternary struct {
    // contains filtered or unexported fields
}
```

<a name="NewTernary"></a>
### func [NewTernary](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L126>)

```go
func NewTernary(condition Expr, trueExpr Expr, falseExpr Expr) *Ternary
```



<a name="Ternary.Accept"></a>
### func \(\*Ternary\) [Accept](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L130>)

```go
func (t *Ternary) Accept(ctx context.Context, v Visitor) EvaluationResult
```



<a name="Token"></a>
## type [Token](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/lex.go#L72-L77>)



```go
type Token struct {
    // contains filtered or unexported fields
}
```

<a name="Token.String"></a>
### func \(Token\) [String](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/lex.go#L92>)

```go
func (t Token) String() string
```



<a name="TokenType"></a>
## type [TokenType](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/lex.go#L78>)



```go
type TokenType int
```

<a name="EOF"></a>

```go
const (
    // Single-character tokens.
    EOF TokenType = iota
    ERROR
    LEFT_PAREN
    RIGHT_PAREN
    LEFT_BRACE
    RIGHT_BRACE
    LEFT_BRACKET
    RIGHT_BRACKET
    PERCENT
    COLON
    COMMA
    DOT
    MINUS
    PLUS
    SEMICOLON
    SLASH
    STAR
    QMARK
    // One or two character tokens.
    BANG
    BANG_EQUAL
    EQUAL
    EQUAL_EQUAL
    GREATER
    GREATER_EQUAL
    LESS
    LESS_EQUAL
    TEMPLATE_LEFT_BRACE
    TEMPLATE_RIGHT_BRACE
    OPTIONALCHAIN
    AND
    OR
    // Literals.
    IDENTIFIER
    STRING
    NUMBER
    TEXT

    FALSE
    TRUE
    NIL
)
```

<a name="TokenType.String"></a>
### func \(TokenType\) [String](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/lex.go#L96>)

```go
func (t TokenType) String() string
```



<a name="Unary"></a>
## type [Unary](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L24-L27>)



```go
type Unary struct {
    // contains filtered or unexported fields
}
```

<a name="NewUnary"></a>
### func [NewUnary](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L110>)

```go
func NewUnary(operator Token, right Expr) *Unary
```



<a name="Unary.Accept"></a>
### func \(\*Unary\) [Accept](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L114>)

```go
func (u *Unary) Accept(ctx context.Context, v Visitor) EvaluationResult
```



<a name="Variable"></a>
## type [Variable](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L67-L69>)



```go
type Variable struct {
    // contains filtered or unexported fields
}
```

<a name="NewVariable"></a>
### func [NewVariable](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L174>)

```go
func NewVariable(name Token) *Variable
```



<a name="Variable.Accept"></a>
### func \(\*Variable\) [Accept](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/nodes.go#L178>)

```go
func (v *Variable) Accept(ctx context.Context, vis Visitor) EvaluationResult
```



<a name="Visitor"></a>
## type [Visitor](<https://github.com/nonsocode/xpress/blob/main/pkg/parser/types.go#L8-L26>)



```go
type Visitor interface {
    // contains filtered or unexported methods
}
```

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
