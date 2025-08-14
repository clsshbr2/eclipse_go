# eclipse_go

Modulos do painel eclipse feito em Go.

---

## Requisitos

- Linux (64-bit)
- Go 1.25.0
- Git

---

## Instalação do Go 1.25.0

Siga os passos abaixo para instalar ou atualizar o Go:

1. **Baixar o arquivo binário do Go 1.25.0:**

```bash
wget https://go.dev/dl/go1.25.0.linux-amd64.tar.gz
```

2. **Remover versões antigas do Go (se houver):**
```bash
sudo rm -rf /usr/local/go
```

3. **Extrair o Go para /usr/local:**
```bash
sudo tar -C /usr/local -xzf go1.25.0.linux-amd64.tar.gz
```

4. **Adicionar o Go ao PATH:**
```bash
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
source ~/.profile
```

5. **Verificar a instalação:**
```bash
go version
```

6. **Clonar o projeto**
```bash
git clone https://github.com/clsshbr2/eclipse_go.git
cd eclipse_go
```

7. **Instalar dependências**
```bash
go mod tidy
```

8. **Criar config.json**
```bash
{
    "authToken": "senha_da_vps_base64",
    "url": "https://painel_eclipse.com.br",
    "porta": 8989
}
```

9. **Executar o modulo**
```bash
go run main.go
```