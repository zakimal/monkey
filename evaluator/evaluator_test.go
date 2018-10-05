package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

// Integerを正しく評価できているかをテスト
func TestEvalIntegerExpression(t *testing.T) {

	// テストセット
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	// 各テストセットに対して
	for _, tt := range tests {

		// inputを評価して
		evaluated := testEval(tt.input)

		// 結果を確認
		testIntegerObject(t, evaluated, tt.expected)
	}
}

// 入力をレキサ・パーサに通して得られたASTをObjectに変換して返すヘルパー関数
func testEval(input string) object.Object {

	// 入力で初期化したレキサを生成
	l := lexer.New(input)

	// レキサをセットしたパーサを生成
	p := parser.New(l)

	// プログラムをパース
	program := p.ParseProgram()

	// 新たな環境を生成
	env := object.NewEnvironment()

	// パースした結果得られるASTを評価
	return Eval(program, env)
}

// 引数objがIntegerObject型で、かつ格納されている値が期待したものになっていることを確認するヘルパー関数
func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {

	// 引数がInteger型であることを確認
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T(%+v)", obj, obj)
		return false
	}

	// 格納してる値が期待したものになっていることを確認
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}

// Booleanを正しく評価できているかをテスト
func TestEvalBooleanExpression(t *testing.T) {

	// テストセット
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	// 各テストセットに対して
	for _, tt := range tests {

		// inputを評価して
		evaluated := testEval(tt.input)

		// 結果を確認
		testBooleanObject(t, evaluated, tt.expected)
	}
}

// 引数objが期待するBooleanObjectであることを確認するヘルパー関数
func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {

	// Boolean型であることを確認
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T(%+v)", obj, obj)
		return false
	}

	// 格納している値が期待したものであることを確認
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}

	return true
}

// !演算子の評価をテスト
func TestBangOperator(t *testing.T) {

	// テストセット
	// 5はtruthyに扱う
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	// 各テストセットに対して
	for _, tt := range tests {

		// inputを評価して
		evaluated := testEval(tt.input)

		// 結果を確認
		testBooleanObject(t, evaluated, tt.expected)
	}
}

// If-Else式の評価をテスト
func TestIfElseExpressions(t *testing.T) {

	// テストセット
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	// 各テストセットについて
	for _, tt := range tests {

		// inputを評価
		evaluated := testEval(tt.input)

		// 型アサーション
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}

	}
}

// 引数objがNullObjectであるかを確認するヘルパー関数
func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not Null. got=%T(%+v)", obj, obj)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {

	// テストセット
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`
if (10 > 1) {
	if (10 > 1) {
		return 10;
	}
	return 1;
}
`, 10},
	}

	// 各テストセットに対して
	for _, tt := range tests {

		// inputを評価
		evaluated := testEval(tt.input)

		// 正しいObjectが帰ってきていることを確認
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {

	// テストセット
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
if (10 > 1) {
  if (10 > 1) {
    return true + false;
  }

  return 1;
}
`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
	}

	// 各テストセットに対して
	for _, tt := range tests {

		// inputを評価
		evaluated := testEval(tt.input)

		// ErrorObjectが返っているはず
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

// Let文の評価をテストする
func TestLetStatements(t *testing.T) {

	// テストケース
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a;  b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	// 各テストケース対して
	for _, tt := range tests {

		// 各テストケースの結果返ってくるObjectはIntegerObjectであるはず
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

// 正しくFunction型のObjectを生成することができているかを確認するテスト
func TestFunctionObject(t *testing.T) {

	// 関数の入力
	input := "fn(x) { x + 2; };"

	// ここで評価
	evaluated := testEval(input)

	// 評価して得られたObjectの型を確認
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. gpt=%T(%+v)",
			evaluated, evaluated)
	}

	// 評価して得られたFunction型のObjectのParametersの個数を確認
	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong paramaters. Parameters=%+v",
			fn.Parameters)
	}

	// 評価して得られたFunction型のObjectのParametersのリテラルを確認
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	// 評価して得られたFunction型のObjectのBodyのリテラルを確認
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q",
			expectedBody, fn.Body.String())
	}
}

// 関数呼び出しのテスト
func TestFunctionApplication(t *testing.T) {

	// テストケース
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	// 各テストケースに対して
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}
