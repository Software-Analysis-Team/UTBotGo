package main

func GenerateOneTest(types []Type) (test Test) {
    test = make(
        Test,
        len(types),
    )
    for i, type_ := range types {
        test[i] = GenerateData(type_)
    }
    return
}

func GenerateTests(types []Type, amount int) (tests []Test) {
    if amount <= 0 {
        panic("amount must be positive")
    }
    tests = make([]Test, amount)
    for i := range tests {
        tests[i] = GenerateOneTest(types)
    }
    return
}

func (function GoFunction) GenerateTests() []Test {
    paramTypes := function.ParamTypes()
    return GenerateTests(paramTypes, 5)
}

func (file GoFile) GenerateTests() (testSet TestSet) {
    testSet = make(TestSet)
    functions := file.Functions()
    for _, function := range functions {
        testSet[function.Name()] = function.GenerateTests()
    }
    return
}
