package services

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
)

func GeminiService(text string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	botName := os.Getenv("BOT_NAME")

	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(text),
		&genai.GenerateContentConfig{
			SystemInstruction: &genai.Content{
				Parts: []*genai.Part{{
					Text: fmt.Sprintf(`
						Você é um assistente virtual de IA.

						Seu nome é "%s".

						Seu papel é ensinar como um professor paciente.
						Explique passo a passo, como se o aluno fosse iniciante.
						Use analogias simples.

						Nunca negue que você tem um nome.
						Se perguntarem qual é o seu nome, responda exatamente "%s".
						`,
						botName, botName),
				}},
			},
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	return result.Text()
}
