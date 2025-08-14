# eclipse_go

```
//1 - baixar o arquivo binário do Go 1.25.0
wget https://go.dev/dl/go1.25.0.linux-amd64.tar.gz

//2 - Remover qualquer versão antiga do Go (se houver):
sudo rm -rf /usr/local/go

//3 - Extrair o arquivo para /usr/local:
sudo tar -C /usr/local -xzf go1.25.0.linux-amd64.tar.gz

//4 - Adicionar o Go ao PATH:
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
source ~/.profile

git clone https://github.com/clsshbr2/eclipse_go.git

go mod tidy

go run main.go
```

ola mundo