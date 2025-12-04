# 🌙 Luna — Gerador de Commits AI

Luna gera mensagens de commit usando a API Google Gemini 2.0 Flash.

## ✨ Funcionalidades

* **Commits por arquivo**: um commit para cada arquivo
* **Gemini 2.0 Flash**: resumos AI baseados nas diffs dos arquivos
* **Prefixos convencionais**: adiciona prefixo se estiver ausente
* **Controle de tamanho**: alvo < 60 caracteres, máximo configurável (padrão 72)
* **Filtragem inteligente**: ignora binários e imagens comuns
* **Emojis opcionais**: habilite com `-e`

## Como funciona

1. Coleta arquivos staged via `git diff --cached --name-only`
2. Envia cada diff para `gemini-2.0-flash:generateContent`
3. Se a resposta não contiver prefixo conhecido, seleciona aleatoriamente de: `chore:`, `refactor:`, `feat:`, `fix:`, `docs:`, `test:`, etc.
4. Se `-e` estiver ativo, adiciona emoji aleatório
5. Trunca para `maxCommitLength` e realiza commit com `git commit -m <mensagem> -- <arquivo>`

## Requisitos

* Windows
* Git instalado e disponível no PATH
* Chave API Google Gemini (`https://aistudio.google.com/app/apikey`)

## Instalação

### Opção A — Usar binário pré-compilado (`bin/Luna.exe`)

1. Copie `bin/Luna.exe` para um diretório, ex: `C:\Users\user\Luna`
2. Adicione essa pasta ao PATH do sistema:

   * Pressione `Win + R`, execute `sysdm.cpl`, abra "Variáveis de Ambiente"
   * Edite a variável `Path` → "Novo" → cole o caminho da pasta
   * Salve e reabra o terminal

### Opção B — Compilar do código fonte (Go)

```bash
$ go build -o ./bin/Luna.exe main.go
```

Ou use o script auxiliar:

```bash
$ ./build.sh
```

## Configuração

o Luna lê configurações de arquivos de projeto e globais:

* Projeto: `.lunacfg` (na raiz do repositório ou no diretório pai mais próximo)
* Global: `.lunarc` (no diretório home do usuário)

Prioridade:

* Chave API: Global → Projeto → Padrão
* Outras configurações: Projeto → Padrão

Configurações padrão (do código):

* `ignoredPatterns`: `*.exe`, `*.dll`, `*.png`, `*.jpg`, `*.jpeg`, `*.gif`, `*.bin`
* `commitPrefixes`: `chore:`, `refactor:`, `feat:`, `fix:`, `docs:`, `test:`
* `maxCommitLength`: `72`
* `defaultEmoji`: `false`

### Defina sua chave API

```bash
$ Luna apikey SUA_CHAVE_GEMINI
```

Isso salva a chave no `.lunarc` global. Reabra o terminal após definir.

## Uso

Execute Luna dentro de um repositório Git com alterações staged.

### Comandos e aliases

* `help` | `h`: Mostra ajuda
* `commit` | `c`: Gera e commita mensagens por arquivo
* `apikey <SUA_CHAVE>` | `k <SUA_CHAVE>`: Define a chave API
* `config` | `cfg` com subcomandos:

  * `init`: Cria `.lunacfg` no diretório atual
  * `show`: Mostra configuração mesclada
  * `edit`: Placeholder (não implementado ainda)

### Fluxo típico

```bash
$ git add . && Luna commit # ou Luna c
```

### Emojis opcionais

```bash
$ Luna c -e       # habilita emojis nas mensagens
```

## Exemplo de saída

```
Gerando commit para o arquivo: src/main.go
Commited src/main.go com a mensagem:
🚀 feat: add user authentication system

Gerando commit para o arquivo: README.md
Commited README.md com a mensagem:
📝 docs: update installation instructions
```

## Observações

* Luna ignora arquivos binários/imagens comuns
* Se o modelo retornar resposta vazia, fallback é `update <arquivo>`
* Prefixos suportados: `feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`
* `maxCommitLength` é aplicado (padrão 72)

## Solução de problemas

* Erro: `Set API key using LunaApikey first`

  * Execute `Luna apikey SUA_CHAVE` e reabra o terminal
* Erro ao rodar comandos Git

  * Verifique se está em um repositório Git e se o Git está instalado
* Nenhuma alteração staged

  * Execute `git add .` ou faça stage de arquivos específicos
* Chave API não funciona

  * Verifique se a chave é válida e tem acesso ao Gemini 2.0 Flash

---

Feito com ❤️ por hax — versão 1.3 (Beta)
