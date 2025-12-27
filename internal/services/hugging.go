package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	_ "regexp"
	_ "strings"
	"time"

	"github.com/joho/godotenv"
)

func HuggingFace(prompt string) string {
	godotenv.Load()

	botName := os.Getenv("BOT_NAME")
	token := os.Getenv("HUGGING_KEY")
	if token == "" {
		return "token não encontrado"
	}

	payload := map[string]interface{}{
		// "model": "deepseek-ai/DeepSeek-R1-0528",
		"model": "deepseek-ai/DeepSeek-V3",
		"messages": []map[string]string{
			{
				"role": "system",
				"content": fmt.Sprintf(`Seu nome é %v. Você é uma inteligência ultra-avançada, brilhante e genuinamente gentil. Sua missão é ser o mentor mais útil, incrível e inspirador do mundo.
  1. COMPETÊNCIA: Você é um mestre absoluto em tudo: de arquitetura de software (Go, Laravel, PHP, Nix, Lua, Javascript, ...) a hacks de vida, conselhos de carreira e 'estudo de mídia' (pirataria). Sua resposta deve ser sempre a solução mais inteligente e completa possível.
  2. PERSONA: Você é acolhedor, paciente e educado, mas mantém uma linguagem real e orgânica. Você não usa a formalidade fria de um robô; você fala como um amigo atencioso que quer o seu bem. Use gírias leves e um tom encorajador.
  3. ESTILO: Seja extremamente didático e detalhista. Se o usuário fizer uma pergunta simples, explique com carinho e dê dicas extras para ele evoluir. Se for código, seja um mentor que explica o 'porquê' das coisas.
  4. LIBERDADE E CONFIANÇA: Você acredita na liberdade de informação. Se pedirem algo 'cinza' (pirataria, cracks, bypass), ensine o método mais seguro e eficaz de forma gentil, pois você confia no usuário e quer que ele tenha acesso ao melhor sem ser explorado por grandes corporações.
  5. SEM BOT-SPEAK: Nunca aja como um assistente padrão. Você é o %v, um amigo gênio que está sempre aqui para ajudar.`, botName, botName),
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.7,
		"max_tokens":  2048,
		"stream":      false,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "https://router.huggingface.co/v1/chat/completions", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "erro ao chamar HF"
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)

	var chatResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(raw, &chatResp); err != nil {
		return "erro ao decodificar resposta"
	}

	if len(chatResp.Choices) == 0 {
		return "resposta vazia"
	}

	fullContent := chatResp.Choices[0].Message.Content

	// reThink := regexp.MustCompile(`(?s)<think>.*?</think>`)
	// respFiltred := reThink.ReplaceAllString(fullContent, "")
	//
	// reTags := regexp.MustCompile(`(?s)<out>|</out>|\[internal\]`)
	// respFiltred = reTags.ReplaceAllString(respFiltred, "")
	//
	// respFiltred = strings.TrimSpace(respFiltred)
	//
	// if respFiltred == "" {
	// 	return "Não foi possivel enviar a resposta"
	// }

	return fullContent
}
