package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"tatria/langchain"
	"tatria/response"

	"github.com/tmc/langchaingo/chains"
)

type Res struct {
	MSG string `json:"msg"`
	Err error  `json:"err"`
}

type Body struct {
	SysLog string `json:"syslog"`
}

func Process(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body Body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Println("Handling:", body.SysLog)
	input := map[string]any{
		"input": body.SysLog,
	}
	_, err := chains.Call(context.Background(), langchain.Executor, input)
	if err != nil {
		fmt.Println(err)
		response.Error(w, "LLM Internal Error", 500)
		return
	}

	fmt.Println("Done!")
	response.JSON(w, Res{MSG: "done", Err: nil}, 200)
}
