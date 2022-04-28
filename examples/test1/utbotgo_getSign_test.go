package main

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"
)

type utbotgo_args_getSign struct {
	Ax int `json:"x"`
}

type utbotgo_result_getSign int // FIXME can't use when function return many values

func utbotgo_parseArgs_getSign() (res []utbotgo_args_getSign, err error) {
	var data []byte
	data, err = os.ReadFile("utbotgo/functions/getSign/args.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &res)
	return
}

func utbotgo_parseAnswers_getSign() (res []utbotgo_result_getSign, err error) {
	var data []byte
	data, err = os.ReadFile("utbotgo/functions/getSign/answers.json")
	if err != nil {
		return nil, nil
	}
	err = json.Unmarshal(data, &res)
	return
}

func utbotgo_parseTests_getSign() (args []utbotgo_args_getSign, answers []utbotgo_result_getSign, err error) {
	args, err = utbotgo_parseArgs_getSign()
	if err != nil {
		return
	}
	answers, err = utbotgo_parseAnswers_getSign()
	if err != nil {
		return
	}
	return
}

func utbotgo_writeResults_getSign(results []utbotgo_result_getSign) (err error) {
	var data []byte
	data, err = json.Marshal(results)
	if err != nil {
		return
	}
	err = os.WriteFile("utbotgo/functions/getSign/results.json", data, 0666)
	return
}

func Test_utbotgo_getSign(t *testing.T) {
	args, answers, err := utbotgo_parseTests_getSign()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	results := make([]utbotgo_result_getSign, len(args))
	for i := range args {
		t.Run("test_#" + strconv.Itoa(i+1), func(t_ *testing.T) {
			results[i] = utbotgo_result_getSign(getSign(args[i].Ax))
			if i >= len(answers) {
				t_.Logf("\nexpected is unknown\nactual: %#v\n", results[i])
			} else if answers[i] != results[i] {
				t_.Logf("\nexpected: %#v\nactual: %#v\n", answers[i], results[i])
				t_.Fail()
			}
		})
	}
	err = utbotgo_writeResults_getSign(results)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
}
