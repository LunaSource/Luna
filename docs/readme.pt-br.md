# 🌙 Luna — Gerador de Commits AI

Luna gera mensagens de commit usando a API do OpenRouter, permitindo usar sua própria chave e escolher qualquer modelo suportado.

## ✨ Funcionalidades

* **Seletor de arquivos interativo**: escolha quais arquivos alterados colocar em staged, direto na TUI
* **Commits por arquivo**: um commit para cada arquivo selecionado
* **Alimentado pelo OpenRouter**: resumos AI baseados nas diffs dos arquivos, usando o modelo que você escolher
* **Prefixos convencionais**: adiciona prefixo se estiver ausente
* **Emojis escolhidos pela IA**: com `-e`, o modelo escolhe um emoji estilo gitmoji que combina com a intenção do commit
* **Controle de tamanho**: alvo < 60 caracteres, máximo configurável (padrão 72)
* **Filtragem inteligente**: ignora binários e imagens comuns

## Como funciona

1. Lista todos os arquivos alterados via `git status --porcelain` (modificados, staged, não rastreados), exceto os ignorados
2. Você escolhe quais incluir: `↑`/`↓` para mover, `espaço` para selecionar, `a` para selecionar todos, `enter` para confirmar
3. Coloca os arquivos selecionados em staged (`git add`) e processa um por um
4. Envia o diff de cada arquivo para o modelo configurado no OpenRouter
5. Se `-e` estiver ativo, o modelo prefixa a mensagem com um emoji correspondente; se a resposta não contiver prefixo conhecido, um é selecionado aleatoriamente de: `chore:`, `refactor:`, `feat:`, `fix:`, `docs:`, `test:`, etc.
6. Trunca para `maxCommitLength` e realiza commit com `git commit -m <mensagem> -- <arquivo>`

## Requisitos

* Windows
* Git instalado e disponível no PATH
* Chave API do OpenRouter (`https://openrouter.ai/keys`)
* Uma [Nerd Font](https://www.nerdfonts.com/) definida como fonte do terminal, para que os ícones da interface sejam exibidos corretamente

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

* Chave API e modelo: Global → Projeto → Padrão
* Outras configurações: Projeto → Padrão

Configurações padrão (do código):

* `ignoredPatterns`: `*.exe`, `*.dll`, `*.png`, `*.jpg`, `*.jpeg`, `*.gif`, `*.bin`
* `commitPrefixes`: `chore:`, `refactor:`, `feat:`, `fix:`, `docs:`, `test:`
* `maxCommitLength`: `72`
* `defaultEmoji`: `false`

### Defina sua chave API e o modelo

```bash
$ Luna apikey SUA_CHAVE_OPENROUTER
$ Luna model openai/gpt-4o-mini
```

Isso salva a chave e o modelo no `.lunarc` global. Reabra o terminal após definir. Veja os modelos disponíveis em [openrouter.ai/models](https://openrouter.ai/models).

## Uso

Execute Luna dentro de um repositório Git com qualquer alteração (staged, não staged ou não rastreada) — você escolhe quais arquivos colocar em staged direto na TUI.

### Comandos e aliases

* `help` | `h`: Mostra ajuda
* `commit` | `c`: Gera e commita mensagens por arquivo
* `apikey <SUA_CHAVE>` | `k <SUA_CHAVE>`: Define a chave API
* `model <MODEL_ID>` | `m <MODEL_ID>`: Define o modelo do OpenRouter
* `config` | `cfg` com subcomandos:

  * `init`: Cria `.lunacfg` no diretório atual
  * `show`: Mostra configuração mesclada
  * `edit`: Placeholder (não implementado ainda)

### Fluxo típico

```bash
$ Luna commit # ou Luna c
```

Use `↑`/`↓` para mover, `espaço` para selecionar arquivos, `a` para selecionar todos, `enter` para colocar em staged e commitar.

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

* Erro: `Set API key using 'luna apikey' first`

  * Execute `Luna apikey SUA_CHAVE` e reabra o terminal
* Erro: `Set model using 'luna model' first`

  * Execute `Luna model MODEL_ID` (veja [openrouter.ai/models](https://openrouter.ai/models))
* Erro ao rodar comandos Git

  * Verifique se está em um repositório Git e se o Git está instalado
* Nenhuma alteração encontrada

  * Nada aparece no `git status`; modifique ou crie um arquivo primeiro
* Chave API não funciona

  * Verifique se a chave é válida e se o modelo escolhido está acessível na sua conta OpenRouter

---

Feito com ❤️ por hax — versão 1.3 (Beta)
