package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

const (
	uploadDir      = "./uploads"
	maxUploadSize  = 10 << 20 // 10 MB
	allowedFormKey = "myFile"
)

func init() {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	absPath, _ := filepath.Abs("./../index.html")
	http.ServeFile(w, r, absPath)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		if strings.Contains(err.Error(), "request body too large") {
			http.Error(w, "Файл слишком большой (максимум 10MB)", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "Неверный запрос", http.StatusBadRequest)
		}
		return
	}
	file, header, err := r.FormFile(allowedFormKey)
	if err != nil {
		http.Error(w, "Неверный файл", http.StatusBadRequest)
		return
	}
	defer file.Close()
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
		return
	}
	content := string(fileBytes)
	if strings.TrimSpace(content) == "" {
		http.Error(w, "Файл пуст", http.StatusBadRequest)
		return
	}

	converted, err := service.AutoDetectAndConvert(content)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка конвертации: %v", err), http.StatusInternalServerError)
		return
	}

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".txt"
	}
	outputFilename := time.Now().UTC().Format("20060102150405") + "_converted" + ext
	outputPath := filepath.Join(uploadDir, outputFilename)

	if err := os.WriteFile(outputPath, []byte(converted), 0644); err != nil {
		http.Error(w, "Ошибка сохранения файла", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8; text/html")
	w.Write([]byte(converted))

}
