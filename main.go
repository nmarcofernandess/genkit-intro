package main

import (
    "context"
    "errors"
    "fmt"
    "log"
    "net/http"

    "github.com/firebase/genkit/go/ai"
    "github.com/firebase/genkit/go/genkit"
    "github.com/firebase/genkit/go/plugins/googleai"
)

func main() {
    ctx := context.Background()

    if err := googleai.Init(ctx, nil); err != nil {
        log.Fatal(err)
    }

    genkit.DefineFlow("menuSuggestionFlow", func(ctx context.Context, input string) (string, error) {
        m := googleai.Model("gemini-1.5-flash")
        if m == nil {
            return "", errors.New("menuSuggestionFlow: failed to find model")
        }

        resp, err := ai.Generate(ctx, m,
            ai.WithConfig(&ai.GenerationCommonConfig{Temperature: 1}),
            ai.WithTextPrompt(fmt.Sprintf(`Suggest an item for the menu of a %s themed restaurant`, input)))
        if err != nil {
            return "", err
        }

        text := resp.Text()
        return text, nil
    })

    if err := genkit.Init(ctx, nil); err != nil {
        log.Fatal(err)
    }

    // Adicionando um handler HTTP
    http.HandleFunc("/suggest", func(w http.ResponseWriter, r *http.Request) {
        input := r.URL.Query().Get("input")
        if input == "" {
            http.Error(w, "Missing input parameter", http.StatusBadRequest)
            return
        }

        result, err := genkit.RunFlow(ctx, "menuSuggestionFlow", input)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        fmt.Fprintf(w, "Suggestion: %s", result)
    })

    // Iniciando o servidor HTTP
    log.Println("Server listening on http://127.0.0.1:3400")
    log.Fatal(http.ListenAndServe(":3400", nil))
}
