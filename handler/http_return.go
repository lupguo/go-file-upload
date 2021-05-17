package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// HTTPJson json响应
func HTTPJson(w http.ResponseWriter, ret map[string][]string) {
	retJson, _ := json.Marshal(ret)
	w.Header().Set("Content-Type", "app/json")
	w.Write(retJson)
}

// HTTPServerError http错误响应
func HTTPServerError(w http.ResponseWriter, format string, v ...interface{}) {
	errorMsg := fmt.Sprintf(format, v)
	log.Println(errorMsg)
	http.Error(w, errorMsg, http.StatusInternalServerError)
}
