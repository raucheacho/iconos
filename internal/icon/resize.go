package icon

import (
	"fmt"
	"image"
	"image/color"
	"path/filepath"

	"github.com/disintegration/imaging"
)

// Load charge une image depuis un fichier (PNG ou JPEG)
func Load(path string) (image.Image, error) {
	img, err := imaging.Open(path)
	if err != nil {
		return nil, fmt.Errorf("impossible de charger l'image %s: %w", path, err)
	}
	return img, nil
}

// Resize redimensionne une image selon les options
func Resize(img image.Image, size int, keepRatio bool, bgColor color.Color) image.Image {
	if keepRatio {
		return ResizeWithRatio(img, size, bgColor)
	}
	return ResizeExact(img, size)
}

// Save sauvegarde une image au format spécifié
func Save(img image.Image, path string) error {
	if err := imaging.Save(img, path); err != nil {
		return fmt.Errorf("impossible de sauvegarder %s: %w", path, err)
	}
	return nil
}

// Generate génère toutes les icônes pour les tailles données
func Generate(srcPath, outDir, prefix, format string, sizes []int, keepRatio bool, bgColor color.Color) error {
	img, err := Load(srcPath)
	if err != nil {
		return err
	}

	for _, size := range sizes {
		resized := Resize(img, size, keepRatio, bgColor)
		outPath := filepath.Join(outDir, fmt.Sprintf("%s%d.%s", prefix, size, format))

		if err := Save(resized, outPath); err != nil {
			return err
		}
		fmt.Printf("✓ Généré: %s (%dx%d)\n", outPath, size, size)
	}

	return nil
}
