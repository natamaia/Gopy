# Gopy - Histórico de Área de Transferência

Gopy é um gerenciador de Área de Transferência (Clipboard) ultraleve criado em **Go** nativamente para ambientes **Linux**. Ele foi projetado para rodar silenciosamente consumindo o mínimo de memória possível em distribuições modernas e antigas.

Ele conta com dois modos num binário único:
- **Daemon:** Fica invisível no fundo, monitorando todas as cópias que você faz de textos, links e imagens, salvando as informações.
- **Janela (UI):** Uma interface minimalista estilo notificação, que aparece sob demanda para você colar o que precisa de forma rápida.

## 🚀 Funcionalidades Principais
- 📝 **Cópia de Textos & Links:** Encurta visualmente os links e cria menus do tipo *dropdown* (expansores) para não poluir sua tela ao copiar textos longos.
- 🖼️ **Cópia de Imagens:** Suporte total a captura de screenshots ou imagens copiadas do navegador/sistema, criando miniaturas automáticas.
- 📌 **Fixação Inteligente:** Você pode fixar (★) cópias importantes para que elas não sejam apagadas nunca. Os arquivos físicos são realocados em segurança.
- 🧹 **Volátil por Padrão:** O histórico de cópias não-fixadas se autodestrói magicamente a cada reinício da máquina para preservar o armazenamento do seu PC!
- 🪟 **Janela Fantasma:** A janela não possui os botões normais de fechar e não ocupa espaço na barra de tarefas. Se você clicar fora dela, ela some sozinha.

## 💻 Requisitos do Sistema
- **Sistema Operacional:** Linux (com ambiente gráfico)
- **Linguagem / Compilação:** Go 1.20+
- **Bibliotecas Gráficas:** GTK3 (através da biblioteca `gotk3`)

### Dependência para Compilação:
Para que o compilador do Go consiga entender a interface gráfica nativa do seu Linux, você deve garantir que os cabeçalhos de desenvolvimento do GTK3 estejam instalados no seu PC:
```bash
sudo apt update && sudo apt install -y libgtk-3-dev
```

## 🛠️ Como Instalar
Clonou o repositório? Basta rodar o script de instalação nativo:
```bash
chmod +x install.sh
./install.sh
```

**O que o Instalador faz?**
1. Compila o arquivo cortando sobras de debug para ele ficar bem pequeno.
2. Copia para sua pasta oficial de usuário (`~/.local/bin/gopy`).
3. Adiciona na inicialização do sistema (`~/.config/autostart`) para que o "modo escuta" (Daemon) ligue automaticamente quando você acessar o computador.

## ⌨️ Como Usar e Configurar
Para imitar o atalho do Windows (`Windows + V`):
1. No seu Linux, abra as configurações de **Atalhos de Teclado**.
2. Vá em **Atalhos Personalizados**.
3. Adicione o comando: `/home/SEU_USUARIO/.local/bin/gopy --ui`
4. Pressione as teclas Super (Botão Windows) + V para vincular.

Pronto! Agora toda vez que você copiar algo, basta apertar o atalho que a janelinha minimalista irá aparecer perto de você!

---
*(Desenvolvido em Go + GTK3)*
