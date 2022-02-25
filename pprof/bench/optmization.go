package profiling

import (
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
	"strconv"
	"strings"
)

func createRequest() string {
	payload := make([]int, 100, 100)
	for i := 0; i < 100; i++ {
		payload[i] = i
	}
	req := Request{"demo_transaction", payload}
	v, err := json.Marshal(&req)
	if err != nil {
		panic(err)
	}
	return string(v)
}

func processRequest(reqs []string) []string {
	reps := []string{}
	for _, req := range reqs {
		reqObj := &Request{}
		jsoniter.Unmarshal([]byte(req), reqObj)

		var buf strings.Builder
		for _, e := range reqObj.PayLoad {
			buf.WriteString(strconv.Itoa(e))
			buf.WriteString(",")
		}
		repObj := &Response{reqObj.TransactionID, buf.String()}
		repJson, err := json.Marshal(&repObj)
		if err != nil {
			panic(err)
		}
		reps = append(reps, string(repJson))
	}
	return reps
}

func processRequestOld(reqs []string) []string {
	reps := []string{}
	for _, req := range reqs {
		reqObj := &Request{}
		jsoniter.Unmarshal([]byte(req), reqObj)
		//ret := ""
		var buf strings.Builder
		for _, e := range reqObj.PayLoad {
			//ret += strconv.Itoa(e) + ","
			buf.WriteString(strconv.Itoa(e))
			buf.WriteString(",")
		}
		repObj := &Response{reqObj.TransactionID, buf.String()}
		repJson, err := json.Marshal(&repObj)
		if err != nil {
			panic(err)
		}
		reps = append(reps, string(repJson))
	}
	return reps
}
