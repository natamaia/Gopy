# Histórico de Área de Transferência (Gopy) - Plano de Desenvolvimento

Este documento descreve o plano de desenvolvimento da aplicação de histórico de área de transferência ultraleve solicitada.

## 1. Requisitos
- **Linguagem:** Go (sem uso de Docker).
- **Ambiente:** Linux Debian (ou similar com interface gráfica), rodando como um serviço/daemon.
- **Leveza:** Projetado para ser rápido e consumir pouca memória em versões antigas.
- **Interface Gráfica:** 
  - Estilo notificação (minimalista).
  - Sem botões de fechar, minimizar ou maximizar (`SetDecorated(false)`).
  - Fechar automaticamente ao clicar fora da janela (perda de foco).
  - Não deve ocupar espaço no painel/taskbar.
- **Funcionalidades de Conteúdo:**
  - **Links:** Encurtar exibição (remover http/https), exibir sublinhado. Clicar com `Ctrl` para abrir diretamente no navegador.
  - **Imagens:** Exibir miniatura, botão para expandir/visualizar.
  - **Textos:** Opção de scroll lateral ou botão dropdown para expandir textos muito grandes.
- **Gerenciamento de Arquivos:**
  - Botão para abrir o diretório onde o item está salvo.
  - Itens não fixados vão para o `/tmp/gopy_clipboard`.
  - Itens fixados vão para diretório persistente (ex: `~/.config/gopy/pinned`).

## 2. Arquitetura
A aplicação será dividida em dois modos dentro do mesmo binário, assim como o MVP de referência.

### A) Modo Daemon (`gopy --daemon`)
- Roda em background.
- Monitora a área de transferência usando `gotk3/gotk3/gtk` (para suportar imagens e textos nativamente no GTK).
- Salva o histórico (metadados) em `~/.config/gopy/history.json`.
- Salva dados binários (imagens) em `/tmp/gopy_clipboard`.

### B) Modo UI (`gopy --ui`)
- Janela instanciada sob demanda via atalho do sistema (ex: `Super + V`).
- Lê os dados locais (arquivos JSON) e monta a lista verticalmente.
- Configura os eventos de UI:
  - `focus-out-event`: Fecha a janela (`win.Destroy()`).
  - Tratamento para Links (identificação por Regex e exibição customizada).
  - Expansão de textos via `gtk.Expander`.

## 3. Próximos Passos
1. Inicializar o módulo Go e dependências (gotk3).
2. Escrever a estrutura de dados (JSON) suportando textos, links e imagens.
3. Desenvolver o Daemon que observa a Clipboard.
4. Desenvolver a UI com GTK, implementando estilo "sem decoração" e "fecha ao perder foco".
5. Testes no ambiente Linux local.

*Data da última atualização: 08/09/2026*
