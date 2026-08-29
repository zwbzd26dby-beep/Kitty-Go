// Package gui implements a minimal web chat UI backed by the shared Core.
package gui

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"

	"github.com/zwbzd26dby-beep/Kitty-Go/internal/agent"
)

const page = `<!DOCTYPE html>
<html lang="pl">
<head>
<meta charset="utf-8">
<title>Kitty GUI</title>
<style>
body{font-family:sans-serif;max-width:700px;margin:2em auto;padding:0 1em}
#chat{border:1px solid #ccc;border-radius:8px;padding:1em;min-height:300px;margin-bottom:1em}
.msg{margin:.4em 0;padding:.4em;border-radius:6px}
.user{background:#eef6ff}
.kitty{background:#f3f3f3}
input{width:100%;padding:.6em;box-sizing:border-box}
</style>
</head>
<body>
<h1>Kitty</h1>
<div id="chat"></div>
<input id="in" placeholder="Wpisz wiadomość…">
<script>
const chat=document.getElementById('chat');
function add(cls,text){
 const d=document.createElement('div');d.className='msg '+cls;d.textContent=text;chat.appendChild(d);
}
document.getElementById('in').addEventListener('keydown',async e=>{
 if(e.key!=='Enter')return;
 const inp=e.target;const text=inp.value.trim();if(!text)return;
 inp.value='';add('user',text);
 const r=await fetch('/agent/respond',{method:'POST',headers:{'Content-Type':'application/json'},
   body:JSON.stringify({message:text})});
 const j=await r.json();add('kitty',j.content||('Błąd: '+r.status));
});
</script>
</body>
</html>
`

var tmpl = template.Must(template.New("gui").Parse(page))

// Server renders the chat page and proxies responses through the agent.
type Server struct {
	agent agent.Agent
}

// NewServer builds a GUI server over the agent.
func NewServer(a agent.Agent) *Server {
	return &Server{agent: a}
}

// Handler returns the request router.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handlePage)
	mux.HandleFunc("/agent/respond", s.handleRespond)
	return mux
}

func (s *Server) handlePage(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = tmpl.Execute(w, nil)
}

func (s *Server) handleRespond(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Message) == "" {
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}
	resp, err := s.agent.Process(r.Context(), req.Message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"content": resp.Content})
}
