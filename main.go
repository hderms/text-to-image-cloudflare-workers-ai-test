//go:build js && wasm
package main

import (
	"encoding/json"
	"net/http"

	"github.com/syumai/workers"
	"github.com/syumai/workers/cloudflare/ai"
)

type PromptRequest struct {
	Prompt string `json:"prompt"`
}

type TextToImageRequest struct {
	Prompt string `json:"prompt"`
}


func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		ttir := &TextToImageRequest{}
		err := json.NewDecoder(req.Body).Decode(ttir)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ai, err := ai.NewAI("AI")
		if err != nil {
			panic(err)
		}
		imagePrompt := &PromptRequest{
			Prompt: ttir.Prompt,
		}

		imageBinary, err := ai.RunBytes("@cf/leonardo/lucid-origin", imagePrompt, "image")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "image/jpg")
		w.Write(imageBinary)
	})
	workers.Serve(nil) // use http.DefaultServeMux
}
