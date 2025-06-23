package services

import (
	"backend/interfaces"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type LocalStorage struct {
	baseDir string
	baseURL string
}

// NewLocalStorage creates a new local storage instance
func NewLocalStorage() interfaces.StorageService {
	baseDir := "local_storage"
	
	// Ensure storage directory exists
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		log.Printf("Warning: Failed to create storage directory %s: %v", baseDir, err)
	}
	
	// Base URL for serving local files (this would need to be configured based on your setup)
	baseURL := "http://localhost:8080/storage"
	
	log.Println("Local storage initialized with directory:", baseDir)
	return &LocalStorage{
		baseDir: baseDir,
		baseURL: baseURL,
	}
}

func (s *LocalStorage) UploadFromURL(sourceURL, destinationPath string) (string, error) {
	log.Printf("Local storage: Downloading from URL %s to path %s", sourceURL, destinationPath)
	
	// Download the file from the source URL
	resp, err := http.Get(sourceURL)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download file, status code: %d", resp.StatusCode)
	}
	
	// Create the full local path
	fullPath := filepath.Join(s.baseDir, destinationPath)
	
	// Ensure the directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}
	
	// Create the destination file
	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()
	
	// Copy the content
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}
	
	// Return the local URL
	localURL := fmt.Sprintf("%s/%s", s.baseURL, destinationPath)
	log.Printf("Local storage: File saved successfully, URL: %s", localURL)
	return localURL, nil
}

func (s *LocalStorage) UploadBase64Image(base64Image, destinationPath string) (string, error) {
	log.Printf("Local storage: Saving base64 image to path %s", destinationPath)
	
	// Extract the base64 data
	var base64Data string
	var contentType string
	
	if strings.HasPrefix(base64Image, "data:") {
		// Extract the base64 part from the data URL
		parts := strings.Split(base64Image, ",")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid base64 image format: wrong number of parts")
		}
		
		// Extract content type from the data URL
		metaParts := strings.Split(parts[0], ";")
		if len(metaParts) < 1 {
			return "", fmt.Errorf("invalid base64 image format: missing content type")
		}
		
		contentTypeParts := strings.Split(metaParts[0], ":")
		if len(contentTypeParts) != 2 {
			return "", fmt.Errorf("invalid base64 image format: invalid content type")
		}
		
		contentType = contentTypeParts[1]
		base64Data = parts[1]
		
		log.Printf("Local storage: Extracted content type: %s", contentType)
	} else {
		// Assume it's already base64 data without the data URL prefix
		base64Data = base64Image
		contentType = "image/png" // Default to PNG
	}
	
	// Validate base64 data
	if len(base64Data) == 0 {
		return "", fmt.Errorf("empty base64 data")
	}
	
	// Decode the base64 data
	imageData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 image: %v", err)
	}
	
	// Create the full local path
	fullPath := filepath.Join(s.baseDir, destinationPath)
	
	// Ensure the directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}
	
	// Create the destination file
	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()
	
	// Write the image data
	_, err = file.Write(imageData)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}
	
	// Return the local URL
	localURL := fmt.Sprintf("%s/%s", s.baseURL, destinationPath)
	log.Printf("Local storage: Base64 image saved successfully, URL: %s", localURL)
	return localURL, nil
}

func (s *LocalStorage) DeleteFile(filePath string) error {
	fullPath := filepath.Join(s.baseDir, filePath)
	
	if err := os.Remove(fullPath); err != nil {
		// If file doesn't exist, we consider it a success
		if os.IsNotExist(err) {
			log.Printf("Local storage: File %s doesn't exist, considering delete successful", filePath)
			return nil
		}
		return fmt.Errorf("failed to delete file: %v", err)
	}
	
	log.Printf("Local storage: File %s deleted successfully", filePath)
	return nil
} 