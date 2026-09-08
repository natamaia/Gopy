#!/bin/bash

echo "Compilando Gopy..."
go build -ldflags="-s -w" -o gopy src/main.go

echo "Instalando binário em ~/.local/bin..."
mkdir -p ~/.local/bin
cp gopy ~/.local/bin/gopy

echo "Criando serviço de Inicialização Automática (Daemon)..."
mkdir -p ~/.config/autostart
cat <<DESKTOP > ~/.config/autostart/gopy-daemon.desktop
[Desktop Entry]
Type=Application
Exec=$HOME/.local/bin/gopy
Hidden=false
NoDisplay=false
X-GNOME-Autostart-enabled=true
Name=Gopy Daemon
Comment=Monitor de Área de Transferência (Clipboard)
DESKTOP

echo ""
echo "===================================================="
echo "✅ Gopy instalado com sucesso!"
echo "O Daemon já vai iniciar automaticamente quando você ligar o PC."
echo ""
echo "Para criar o atalho da Janela (--ui):"
echo "1. Abra Menu > Configurações do Sistema > Teclado > Atalhos"
echo "2. Vá em 'Atalhos Personalizados' e clique em 'Adicionar'"
echo "3. Nome: Gopy Histórico"
echo "4. Comando: $HOME/.local/bin/gopy --ui"
echo "5. Clique na linha não atribuída e pressione a tecla desejada (ex: Super + V)"
echo "===================================================="
