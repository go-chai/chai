package tests

type TestStruct struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`
}

type TestInnerResponse struct {
	FooFoo int `json:"foo_foo"`
	BarBar int `json:"bar_bar"`
}
type TestResponse struct {
	Foo string `json:"foo"`
	Bar string `json:"bar"`

	TestInnerResponse TestInnerResponse `json:"test_inner_response"`
}

type TestRequest struct {
	Foo string `json:"foob"`
	Bar string `json:"barb"`

	TestInnerResponse TestInnerResponse `json:"test_inner_responseb"`
}

type TestError struct {
	Message string `json:"message"`
}

func (e TestError) Error() string {
	return e.Message
}

type TestErrorPtr struct {
	Message string `json:"message"`
}

func (e *TestErrorPtr) Error() string {
	return e.Message
}

type TestErrorMap map[string]string

func (e TestErrorMap) Error() string {
	return "test error map"
}

type TestErrorMapPtr map[string]string

func (e *TestErrorMapPtr) Error() string {
	return "test error map ptr"
}
