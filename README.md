# τ tauos

Um TUI minimalista para anotações e tarefas, inspirado no Notion. Feito em Go com [Bubbletea](https://github.com/charmbracelet/bubbletea) e [Lipgloss](https://github.com/charmbracelet/lipgloss).

## Instalação

```bash
git clone <repo>
cd tauos
go install .
```

Ou para instalar manualmente em `~/.local/bin`:

```bash
go build -o ~/.local/bin/tauos .
```

## Uso

```bash
tauos
```

## Atalhos

### Sidebar (páginas)

| Tecla | Ação |
|-------|------|
| `↑` / `k` | Página anterior |
| `↓` / `j` | Próxima página |
| `enter` | Abrir página |
| `n` | Nova página |
| `d` | Deletar página |
| `tab` | Ir para o conteúdo |

### Conteúdo (blocos)

| Tecla | Ação |
|-------|------|
| `↑` / `k` | Bloco anterior |
| `↓` / `j` | Próximo bloco |
| `enter` | Editar texto / toggle checkbox |
| `e` | Editar bloco selecionado |
| `a` | Adicionar bloco |
| `d` | Deletar bloco |
| `tab` | Voltar para sidebar |

### Ao adicionar um bloco (`a`)

| Tecla | Tipo |
|-------|------|
| `t` | Texto |
| `h` | Título |
| `c` | Checkbox |
| `esc` | Cancelar |

### Geral

| Tecla | Ação |
|-------|------|
| `q` | Sair |

## Tipos de bloco

- **Texto** — parágrafo simples
- **Título** — destaque em azul com prefixo `##`
- **Checkbox** — tarefa com toggle (☐ / ☑)

## Dados

As anotações são salvas automaticamente em:

```
~/.local/share/tauos/pages.json
```

## Tecnologias

- [Go](https://golang.org)
- [Bubbletea](https://github.com/charmbracelet/bubbletea) — framework TUI
- [Lipgloss](https://github.com/charmbracelet/lipgloss) — estilização
- [Bubbles](https://github.com/charmbracelet/bubbles) — componentes (textinput)
- Tema: [Catppuccin Mocha](https://github.com/catppuccin/catppuccin)
