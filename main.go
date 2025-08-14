package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

type Config struct {
	Porta     int    `json:"porta"`
	AuthToken string `json:"authToken"`
	URL       string `json:"url"`
}

var (
	config        Config
	caminhoDelete = "./usersToDelete.json"
	caminhoAddSSH = "./usersToaddssh.json"
	mutexDelete   sync.Mutex
	mutexAddSSH   sync.Mutex
	urlOnline     string
	urlBk         string
)

func main() {
	// Carrega config
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Arquivo config.json não encontrado: %v", err)
	}
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("Erro ao ler config: %v", err)
	}

	urlOnline = fmt.Sprintf("%s/onlines.php", config.URL)
	urlBk = fmt.Sprintf("%s/bk.php", config.URL)

	mux := http.NewServeMux()
	mux.HandleFunc("/", authMiddleware(handleRoot))

	// Cron jobs
	s := gocron.NewScheduler(time.Local)
	// s.Every(1).Minute().Do(cronOnlines)
	s.Every(3).Seconds().Do(cronDeleteUsers)
	s.Every(3).Seconds().Do(cronAddUsersSSH)
	s.Every(15).Minutes().Do(cronBackup)
	s.StartAsync()

	addr := fmt.Sprintf(":%d", config.Porta)
	log.Printf("Servidor Go iniciado na porta %d", config.Porta)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != fmt.Sprintf("Bearer %s", config.AuthToken) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Autenticação necessária"})
			return
		}
		next(w, r)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Comando string                 `json:"comando"`
		ExecCmd string                 `json:"exec"`
		Dados   map[string]interface{} `json:"dados"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch req.Comando {
	case "exec":
		out, err := execCommand(req.ExecCmd)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"icon": "error", "mensagem": "Erro ao executar comando",
			})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"icon": "success", "mensagem": "comando executado", "saida": out,
		})

	// Aqui você chamaria suas funções equivalentes do pacote modulos
	case "criarTestssh":
		// exemplo: modulos.CriarTeste(...)
	case "criaruserSsh":
	case "criarUserv2":
	case "criarUserxray":
	case "deleteUsers":
	case "userDeleteALL":
	case "userSinc":
	default:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Comando inválido"})
	}
}

func execCommand(cmd string) (string, error) {
	command := exec.Command("bash", "-c", cmd)
	var out bytes.Buffer
	command.Stdout = &out
	command.Stderr = &out
	err := command.Run()
	return out.String(), err
}

// func cronOnlines() {
// 	log.Println("⏰ Cron onlines rodando")
// 	onV2, errV2 := modulos.GetOnlinesV2()
// 	if errV2 != nil {
// 		log.Println("Erro getonlinesV2:", errV2)
// 	}
// 	onSSH, errSSH := modulos.GetOnlineUsers()
// 	if errSSH != nil {
// 		log.Println("Erro getOnlineUsers:", errSSH)
// 	}
// 	all := append(onV2, onSSH...)
// 	ip, _ := getPublicIP()
// 	data := map[string]interface{}{
// 		"ip":      ip,
// 		"onlines": all,
// 	}
// 	body, _ := json.Marshal(data)
// 	http.Post(urlOnline, "application/json", bytes.NewBuffer(body))
// }

func cronDeleteUsers() {
	mutexDelete.Lock()
	defer mutexDelete.Unlock()
	// Aqui você faz a lógica igual Node.js lendo JSON, deletando e salvando de novo
}

func cronAddUsersSSH() {
	mutexAddSSH.Lock()
	defer mutexAddSSH.Unlock()
	// Aqui você faz a lógica igual Node.js para criar usuários
}

func cronBackup() {
	http.Get(urlBk)
}

func getPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var ipResp struct {
		IP string `json:"ip"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&ipResp); err != nil {
		return "", err
	}
	return ipResp.IP, nil
}
