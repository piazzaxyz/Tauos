# τ tauos

TUI minimalista para notas e tarefas, inspirado no Notion. Feito em Go com [Bubbletea](https://github.com/charmbracelet/bubbletea) e [Lipgloss](https://github.com/charmbracelet/lipgloss). Suporta **markdown inline**, integração com **Notion** e **GitHub Issues**.

## Instalação

```bash
git clone <repo>
cd tauos
go install .
```

Ou compilar direto:

```bash
go build -o tauos .
```

## Uso

```bash
tauos
```

## Estrutura do projeto

```
tauos/
├── integrations/
│   ├── notion.go       # cliente Notion API
│   └── github.go       # cliente GitHub API
├── config.go           # carregamento de config
├── model.go            # lógica da TUI (Bubbletea)
├── view.go             # renderização (Lipgloss + markdown)
├── styles.go           # tema Catppuccin Mocha
├── storage.go          # tipos de dados e persistência JSON
└── main.go             # entrypoint
```

## Configuração de integrações

Crie o arquivo `~/.config/tauos/config.json` (gerado automaticamente vazio no primeiro uso):

```json
{
  "notion_token": "secret_xxx",
  "notion_database_id": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
  "github_token": "ghp_xxx",
  "github_owner": "seu-usuario",
  "github_repo": "seu-repositorio"
}
```

### Notion
1. Crie uma integração em [notion.so/my-integrations](https://www.notion.so/my-integrations)
2. Compartilhe o banco de dados com a integração
3. Copie o token e o ID do banco (32 chars da URL da página)

### GitHub
1. Gere um token em Settings → Developer Settings → Personal Access Tokens
2. O token precisa da permissão `repo` (ou `public_repo` para repos públicos)

## Atalhos

### Sidebar (páginas)

| Tecla | Ação |
|-------|------|
| `↑` / `k` | Página anterior |
| `↓` / `j` | Próxima página |
| `enter` | Abrir página |
| `n` | Nova página |
| `r` | Renomear página |
| `d` | Deletar página |
| `tab` | Ir para conteúdo |
| `q` | Sair |

### Conteúdo (blocos)

| Tecla | Ação |
|-------|------|
| `↑` / `k` | Bloco anterior |
| `↓` / `j` | Próximo bloco |
| `enter` | Toggle checkbox |
| `e` | Editar bloco (textarea multiline) |
| `a` | Adicionar bloco |
| `K` / `J` | Mover bloco para cima/baixo |
| `d` | Deletar bloco |
| `ctrl+n` | Enviar bloco ao Notion |
| `ctrl+g` | Criar issue no GitHub |
| `tab` | Voltar para sidebar |

### Ao adicionar um bloco (`a`)

| Tecla | Tipo |
|-------|------|
| `t` | Texto |
| `h` | Título |
| `c` | Checkbox |
| `C` | Código |
| `q` | Citação |
| `k` | Card (markdown completo) |
| `/` | Divisor |
| `esc` | Cancelar |

### Ao editar um bloco (`e` ou `a` + tipo)

| Tecla | Ação |
|-------|------|
| `ctrl+s` | Salvar |
| `esc` | Cancelar |

## Markdown suportado

Dentro de qualquer bloco de texto, card ou título, o conteúdo renderiza com:

| Sintaxe | Resultado |
|---------|-----------|
| `**texto**` | **negrito** |
| `*texto*` | *itálico* |
| `` `código` `` | código inline |
| `~~texto~~` | ~~tachado~~ |

Dentro de **cards** (`k`), suporte completo a:
- `# Título`, `## Subtítulo`
- `- item` ou `* item` (listas)
- `> citação`
- `---` (divisor)

## Tipos de bloco

| Tipo | Descrição |
|------|-----------|
| **Texto** | Parágrafo com markdown inline |
| **Título** | Destaque em azul com `##` |
| **Checkbox** | Tarefa com toggle (○ / ●) |
| **Código** | Bloco com fundo escuro e barra lateral |
| **Citação** | Texto itálico com barra lateral |
| **Card** | Caixa com borda arredondada, markdown completo |
| **Divisor** | Linha horizontal |

## Dados

Notas salvas automaticamente em:

```
~/.local/share/tauos/pages.json
```

Config em:

```
~/.config/tauos/config.json
```

## Tecnologias

- [Go](https://golang.org)
- [Bubbletea](https://github.com/charmbracelet/bubbletea) — framework TUI
- [Lipgloss](https://github.com/charmbracelet/lipgloss) — estilização
- [Bubbles](https://github.com/charmbracelet/bubbles) — textinput, textarea
- Tema: [Catppuccin Mocha](https://github.com/catppuccin/catppuccin)
