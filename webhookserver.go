package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"

	"fmt"
	"net/http"

	"github.com/TylerBrock/colorjson"
)

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		token := req.URL.Query().Get("token")
		if token == "" {
			body, err := io.ReadAll(req.Body)
			if err != nil {
				panic(err)
			}
			var obj map[string]interface{}
			err = json.Unmarshal(body, &obj)

			if err != nil {
				log.Println(err)
				return
			}

			f := colorjson.NewFormatter()
			f.Indent = 2
			s, _ := f.Marshal(obj)

			fmt.Println(string(s))
			return
		}

		key := []byte("mysecret123")
		h := hmac.New(sha256.New, key)
		h.Write([]byte(token))

		hash := h.Sum(nil)
		encodedHash := base64.StdEncoding.EncodeToString(hash)
		response := "sha256=" + encodedHash
		_, _ = rw.Write([]byte(fmt.Sprintf(`{ "response_token": "%s"}`, response)))
	})

	http.ListenAndServe(":8082", http.DefaultServeMux)

}
