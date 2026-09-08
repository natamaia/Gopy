# Gopy - Conclusão do Desenvolvimento

## 08 de Setembro de 2026 - Versão Final

O projeto **Gopy**, um gerenciador de Área de Transferência ultraleve focado em ambientes Linux GTK (sem Docker), foi desenvolvido e refinado com sucesso, atingindo todas as exigências propostas.

### Arquitetura Consolidada
O projeto roda em um binário único sob dois modos distintos:
1. **Modo Daemon (`./gopy`):** Fica em background utilizando `glib.TimeoutAdd` (integrado ao main thread do GTK) para interceptar modificações na clipboard. Ele lida com a prioridade de Imagens sobre Textos, evitando duplicações, e também interpreta ativamente as cópias de arquivos locais feitas via Nautilus/Nemo, convertendo links diretos para exibição visual.
2. **Modo UI (`./gopy --ui`):** Chamado sob demanda para exibir a notificação.

### Componentização
A Interface de Usuário (UI) foi componentizada utilizando o padrão de caixas flexíveis (`gtk.Box`), dividindo a responsabilidade da interface entre textos e imagens:

* **Componente de Texto:** Lida com links (fazendo o parser da URL e as marcações GTK de estilo azul/sublinhado) e com textos extensos (que são minimizados através de um `gtk.Expander` modificado para mostrar o resumo nas 37 primeiras letras).
* **Componente de Imagem:** Renderiza a miniatura dimensionada da imagem (`PixbufNewFromFileAtScale`). Conta com as funções adicionais de visualização em janela própria e de abrir o diretório local via `xdg-open`.

### Soluções Técnicas Encontradas & Resolvidas
- **Bug da Biblioteca (gotk3):** Atualização forçada para a `branch master` para contornar falhas do gerador CGO na versão antiga lançada em 2021.
- **Bug da EventBox Obscura:** Resolvido desabilitando a renderização visual (`SetVisibleWindow(false)`) da caixa que captura os cliques.
- **Persistência Dinâmica:** Separação do histórico regular volátil (que fica guardado apenas na duração da sessão e apagado pelo Daemon no startup) para o persistente `pinned.json` (que não tem prazo de validade).
- **Problema de Cross-Device Link:** A conversão da função `os.Rename` nativa para um processo de `Cópia e Deleção` (`moveFile`) foi essencial para que imagens salvas em `/tmp` (Virtual/RAM) pudessem ser salvas na persistência dentro do HD em `~/.config`.

### Dependências para Rodar em Máquinas Novas
Para rodar este código a partir do zero em outras distribuições Linux baseadas em Debian, será sempre exigido:
`sudo apt update && sudo apt install -y libgtk-3-dev`
