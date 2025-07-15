package api

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gnojus/wedl/transfer"
)

type Server struct {
	port string
}

type DownloadRequest struct {
	WetransferURL string `json:"wetransfer_url"`
	Password      string `json:"password,omitempty"`
}

type WetransferResponse struct {
	FileName   string `json:"fileName"`
	FileBinary string `json:"fileBinary"`
}

type InfoResponse struct {
	Success   bool   `json:"success"`
	Filename  string `json:"filename"`
	Size      int    `json:"size"`
	URL       string `json:"dl_url"`
	Error     string `json:"error,omitempty"`
}

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

func NewServer(port string) *Server {
	return &Server{port: port}
}

func (s *Server) Start() error {
	http.HandleFunc("/health", s.healthHandler)
	http.HandleFunc("/wetransfer", s.wetransferHandler)
	http.HandleFunc("/info", s.infoHandler)

	log.Printf("Starting API server on port %s", s.port)
	return http.ListenAndServe(":"+s.port, nil)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) wetransferHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.WetransferURL == "" {
		http.Error(w, "wetransfer_url is required", http.StatusBadRequest)
		return
	}

	resp, dlResp, err := transfer.GetDlResponse(req.WetransferURL, req.Password)
	if err != nil {
		http.Error(w, "Failed to get download response: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	filename := dlResp.DlFilename
	if filename == "" {
		filename = transfer.FilenameFromUrl(resp.Request.URL.String())
	}
	if filename == "" {
		filename = "download"
	}

	// Read the entire file content into memory
	fileData, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode to base64
	fileBinary := base64.StdEncoding.EncodeToString(fileData)

	response := WetransferResponse{
		FileName:   filename,
		FileBinary: fileBinary,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) infoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.WetransferURL == "" {
		http.Error(w, "wetransfer_url is required", http.StatusBadRequest)
		return
	}

	resp, dlResp, err := transfer.GetDlResponse(req.WetransferURL, req.Password)
	if err != nil {
		response := InfoResponse{
			Success: false,
			Error:   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	resp.Body.Close()

	response := InfoResponse{
		Success:  true,
		Filename: dlResp.DlFilename,
		Size:     dlResp.DlSize,
		URL:      dlResp.DlUrl,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

