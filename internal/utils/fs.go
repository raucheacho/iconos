package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// EnsureDir crée un répertoire s'il n'existe pas
func EnsureDir(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("impossible de créer le répertoire %s: %w", path, err)
	}
	return nil
}

// FileExists vérifie si un fichier existe
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// GetExtension retourne l'extension d'un fichier en minuscules
func GetExtension(path string) string {
	ext := filepath.Ext(path)
	if len(ext) > 0 {
		return ext[1:] // Enlève le point
	}
	return ""
}

// BuildOutputPath construit le chemin de sortie pour une icône
func BuildOutputPath(outDir, prefix, format string, size int) string {
	filename := fmt.Sprintf("%s%d.%s", prefix, size, format)
	return filepath.Join(outDir, filename)
}

// BuildFaviconPath construit le chemin de sortie pour un favicon
func BuildFaviconPath(outDir string, size int) string {
	filename := fmt.Sprintf("favicon-%dx%d.png", size, size)
	return filepath.Join(outDir, filename)
}

// BuildICOPath construit le chemin pour le fichier ICO
func BuildICOPath(outDir string) string {
	return filepath.Join(outDir, "favicon.ico")
}
