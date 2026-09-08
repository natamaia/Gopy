# Publicação e Deploy (V1)

## Instalação Nativa
Para empacotar e publicar a versão 1.0, criamos um script `install.sh`. 
Esse script executa as seguintes etapas:
1. Compila o Go de forma otimizada (removendo tabelas de debug para deixar o binário ultra leve usando `-ldflags="-s -w"`).
2. Instala o binário em `~/.local/bin/gopy` (padrão de binários do usuário no Linux).
3. Cria um arquivo `.desktop` no diretório `~/.config/autostart/`. Isso garante que o sistema de DE (Desktop Environment) do Linux chame o nosso Daemon assim que o usuário fizer login.

## Atalhos (Global Hotkeys)
Em vez de implementarmos um "Keylogger" no código em Go lendo eventos diretos do teclado via JSON (o que pode ser bloqueado pelo Wayland ou sistemas de segurança modernos do Linux), optamos pela integração **nativa da DE**:
- O sistema operacional gerencia o atalho de teclado global (ex: `Super + V`).
- Quando ativado, ele apenas invoca `/home/natanael/.local/bin/gopy --ui`.
- Isso evita que o nosso Daemon consuma recursos de CPU lendo todos os cliques de tecla globais.

## Preparação para a V2 (Visão de Futuro)
Para resolver a dependência estrita ao GTK3 (que causa a necessidade de instalar `libgtk-3-dev` e pode dar crash em sistemas sem as mesmas bibliotecas gráficas), o planejamento para a V2 será:
- Abandonar a UI nativa `gotk3`.
- Migrar para um framework Go multiplataforma e com renderização própria (Ex: **Fyne**, **Wails** ou **Lorca**).
- Isso garantirá que o binário do Gopy contenha sua própria janela de renderização embarcada, permitindo que ele rode até no Windows ou distribuições Linux base que não usem GTK sem necessitar de novas dependências CGO!
