# Diário de Construção (Build Log)

## 08 de Setembro de 2026 - Primeiros Passos e Dependências

1. **Dependências do Sistema:** 
   O pacote `gotk3` necessita dos headers C do GTK para realizar a ponte (CGO) entre Go e as bibliotecas nativas do Linux.
   - Instalação necessária: `sudo apt install -y libgtk-3-dev`.

2. **Gerenciamento de Pacotes Go:**
   - Comandos executados: `go mod init` e `go mod tidy`.

3. **O "Congelamento" do Build Inicial:**
   - Ao executar `go build` pela primeira vez em um projeto que utiliza `gotk3`, o processo pode parecer travado. Isso ocorre porque o compilador do Go está gerando e compilando milhares de linhas de código C para realizar a conversão (bindings).

4. **Erro "undefined: callback" no gotk3:**
   - A compilação falhou no pacote gdk (`gdk_since_3_22.go:140:180: undefined: callback`). Esse é um problema conhecido da versão v0.6.4 (que é de 2021) em contato com versões mais recentes da linguagem Go e CGO.
   - **Solução Aplicada:** Atualizamos a dependência do `gotk3` para o commit mais recente (`master`) que já contém a correção para esse bug usando `go get github.com/gotk3/gotk3/gtk@master`.
