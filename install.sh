#!/bin/bash

if [ -z "$1" -o -z "$2" -o -z "$3" ]; then
    read -p "Digite A porta que deseja Usar: " PORTA
    read -p "Digite o dominio do seu painel: " URL
    read -p "Digite a senha da sua vps: " TOKEN
else
    PORTA=$1
    URL=$2
    TOKEN=$3
fi

senha_base64=$(echo -n "$TOKEN" | base64)

wget https://go.dev/dl/go1.25.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
source ~/.profile

cd /usr/local
git clone https://github.com/clsshbr2/eclipse_go.git
cd eclipse_go
go mod tidy

cat > config.json <<EOF
{
    "authToken": "$senha_base64",
    "url": "https://$URL",
    "porta": $PORTA
}
EOF

SERVICE_NAME="eclipse_go"
SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"
PROJECT_DIR="/usr/local/eclipse_go"
USER="root"

# Criar arquivo de serviço
echo "Criando serviço systemd em $SERVICE_FILE..."
sudo bash -c "cat > $SERVICE_FILE <<EOF
[Unit]
Description=Eclipse Go Service
After=network.target

[Service]
Type=simple
User=$USER
WorkingDirectory=$PROJECT_DIR
ExecStart=/usr/local/go/bin/go run $PROJECT_DIR/main.go
Restart=always

[Install]
WantedBy=multi-user.target
EOF"

# Recarregar systemd
echo "Recarregando systemd..."
sudo systemctl daemon-reload

# Habilitar serviço para iniciar junto com a VPS
echo "Habilitando serviço para iniciar junto com a VPS..."
sudo systemctl enable $SERVICE_NAME

# Iniciar serviço agora
echo "Iniciando o serviço..."
sudo systemctl start $SERVICE_NAME

# Mostrar status
sudo systemctl status $SERVICE_NAME --no-pager
