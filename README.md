# ORVIT

<p>
    BOT para o Whatsapp utilizando Whatsmeow
</p>

## Funções
- **Criação de Figurinhas**: Converte Imagens|Video|Gifs para figurinhas
- **Inteligência Artifical**: Integração com Google Gemini e Hugging Face
- **Banco de Dados**: Persistência com SQLite

## Pré **requisitos**
- [Go](https://go.dev/)
- [Git](https://git-scm.com/)
- **ffmpeg**: Necessário para conversão de mídia (figurinhas)
- **webpmux**: Necessário para manipulação de imagens WebP

### Instalando Dependencias De Mídias
```bash
sudo apt update
sudo apt install ffmpeg webp
```

## Instalação/Configuração

#### Clone o repositório
```bash
git clone https://github.com/viitorags/orvit
cd orvit
```

#### Configure as variáveis de ambiente
```base
cp .env_example .env
```
#### Edite o arquivo .env

|Váriavel|Descrição|
|:--------|:--------|
|BOT_NAME|Nome do Bot/Binário|
|DB_PATH |Caminho para o banco SQLite (ex: file:data/orvit.db?_foreign_keys=on)|
|GEMINI_API_KEY|Chave de API do Google Gemini|
|HUGGING_KEY|Chave de API do Hugging Face|
|BOT_PREFIX|Prefixo para acionar comandos (ex: !, .)|

#### Baixando as Dependencias do Go
```bash
go mod download
```

## Executar
O projeto tem um [Makefile](./Makefile) para facilitar a execução

compilar:
```bash
make build
```
executar:
```bash
make run
```
limpar binario:
```bash
make clean
```
## Exemplos de Uso
- !menu: Exibe a lista de comandos disponíveis.
- !ping: Verifica se o bot está online e envia um !pong.
- !fig (ou !s): Responda a uma imagem ou vídeo para criar uma figurinha.
- !info: Exibe informações sobre o grupo.

## Estrutura
```bash
.
├── cmd/bot/           # Ponto de entrada (main.go)
├── data/              # Armazenamento do banco de dados SQLite
├── internal/
│   ├── bot/           # Lógica do cliente WhatsApp e Handlers
│   ├── commands/      # Implementação de cada comando (!ping, ...)
│   ├── helpers/       # Funções auxiliares (processamento de mídia)
│   └── services/      # Integrações externas (Gemini, Hugging Face)
└── Makefile           # Scripts de automação
```

## Licença
Este projeto está licenciado sob a licença MIT - veja o arquivo [LICENSE](./LICENSE) para mais detalhes.
