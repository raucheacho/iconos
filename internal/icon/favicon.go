package icon

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

// GenerateFavicons génère les favicons PNG et ICO
func GenerateFavicons(srcPath, outDir string, sizes []int, keepRatio bool, bgColor color.Color) error {
	img, err := Load(srcPath)
	if err != nil {
		return err
	}

	var images []image.Image

	// Génération des PNG
	for _, size := range sizes {
		resized := Resize(img, size, keepRatio, bgColor)
		images = append(images, resized)

		outPath := filepath.Join(outDir, fmt.Sprintf("favicon-%dx%d.png", size, size))
		if err := Save(resized, outPath); err != nil {
			return err
		}
		fmt.Printf("✓ Généré: %s (%dx%d)\n", outPath, size, size)
	}

	// Génération du ICO
	icoPath := filepath.Join(outDir, "favicon.ico")
	if err := CreateICO(images, icoPath); err != nil {
		return err
	}
	fmt.Printf("✓ Généré: %s (multi-résolutions)\n", icoPath)

	return nil
}

// CreateICO crée un fichier ICO à partir de plusieurs images
// Format ICO: https://en.wikipedia.org/wiki/ICO_(file_format)
func CreateICO(images []image.Image, path string) error {
	if len(images) == 0 {
		return fmt.Errorf("aucune image fournie pour le fichier ICO")
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("impossible de créer %s: %w", path, err)
	}
	defer file.Close()

	// Encode chaque image en PNG
	var pngData [][]byte
	for _, img := range images {
		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			return fmt.Errorf("erreur d'encodage PNG: %w", err)
		}
		pngData = append(pngData, buf.Bytes())
	}

	// ICO Header (6 bytes)
	header := make([]byte, 6)
	binary.LittleEndian.PutUint16(header[0:2], 0)                    // Reserved
	binary.LittleEndian.PutUint16(header[2:4], 1)                    // Type: 1 = ICO
	binary.LittleEndian.PutUint16(header[4:6], uint16(len(images))) // Number of images

	if _, err := file.Write(header); err != nil {
		return err
	}

	// Calcul de l'offset de départ des données d'image
	dataOffset := 6 + 16*len(images)

	// Directory entries (16 bytes chacune)
	for i, img := range images {
		bounds := img.Bounds()
		width := bounds.Dx()
		height := bounds.Dy()

		// Taille 256 est encodée comme 0
		w := uint8(width)
		h := uint8(height)
		if width >= 256 {
			w = 0
		}
		if height >= 256 {
			h = 0
		}

		entry := make([]byte, 16)
		entry[0] = w
		entry[1] = h
		entry[2] = 0
		entry[3] = 0
		binary.LittleEndian.PutUint16(entry[4:6], 1)
		binary.LittleEndian.PutUint16(entry[6:8], 32)
		binary.LittleEndian.PutUint32(entry[8:12], uint32(len(pngData[i])))
		binary.LittleEndian.PutUint32(entry[12:16], uint32(dataOffset))

		if _, err := file.Write(entry); err != nil {
			return err
		}

		dataOffset += len(pngData[i])
	}

	// Image data
	for _, data := range pngData {
		if _, err := file.Write(data); err != nil {
			return err
		}
	}

	return nil
}

// GenerateFaviconHTML génère un fichier HTML avec les balises link pour les favicons
func GenerateFaviconHTML(outDir string, sizes []int) error {
	var b strings.Builder
	b.WriteString(`<!DOCTYPE html>
<html>
<head>
    <title>Favicon Preview</title>
    <!-- Favicons générés par iconos -->
    <link rel="icon" type="image/x-icon" href="favicon.ico">
`)
	for _, size := range sizes {
		fmt.Fprintf(&b, `    <link rel="icon" type="image/png" sizes="%dx%d" href="favicon-%dx%d.png">
`, size, size, size, size)
	}

	b.WriteString(`</head>
<body>
    <h1>Favicon Preview</h1>
    <p>Les favicons ont été générés avec succès.</p>
</body>
</html>`)

	path := filepath.Join(outDir, "favicon.html")
	if err := os.WriteFile(path, []byte(b.String()), 0644); err != nil {
		return fmt.Errorf("impossible de créer %s: %w", path, err)
	}
	fmt.Printf("✓ Généré: %s\n", path)
	return nil
}
