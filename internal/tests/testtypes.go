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
