package icon

import (
	"image"
	"image/color"

	"github.com/disintegration/imaging"
)

// ResizeWithRatio redimensionne une image en conservant le ratio
// L'image est centrée dans un carré de la taille cible avec un fond configurable
func ResizeWithRatio(img image.Image, targetSize int, bgColor color.Color) image.Image {
	bounds := img.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()

	// Calcul du ratio pour que l'image tienne dans le carré cible
	ratio := float64(targetSize) / float64(max(srcWidth, srcHeight))
	newWidth := int(float64(srcWidth) * ratio)
	newHeight := int(float64(srcHeight) * ratio)

	// Redimensionnement avec Lanczos (haute qualité)
	resized := imaging.Resize(img, newWidth, newHeight, imaging.Lanczos)

	// Création du canvas carré avec le fond
	canvas := imaging.New(targetSize, targetSize, bgColor)

	// Calcul de la position pour centrer l'image
	offsetX := (targetSize - newWidth) / 2
	offsetY := (targetSize - newHeight) / 2

	// Collage de l'image redimensionnée sur le canvas
	return imaging.Paste(canvas, resized, image.Pt(offsetX, offsetY))
}

// ResizeExact redimensionne une image en forçant les dimensions exactes (sans conserver le ratio)
func ResizeExact(img image.Image, targetSize int) image.Image {
	return imaging.Resize(img, targetSize, targetSize, imaging.Lanczos)
}

// ParseBackgroundColor parse une couleur de fond
func ParseBackgroundColor(bg string) color.Color {
	if bg == "transparent" || bg == "" {
		return color.Transparent
	}

	// Parse hex color (#RRGGBB ou #RGB)
	if len(bg) > 0 && bg[0] == '#' {
		bg = bg[1:]
	}

	var r, g, b uint8
	switch len(bg) {
	case 6:
		r = hexToByte(bg[0:2])
		g = hexToByte(bg[2:4])
		b = hexToByte(bg[4:6])
	case 3:
		r = hexToByte(string(bg[0]) + string(bg[0]))
		g = hexToByte(string(bg[1]) + string(bg[1]))
		b = hexToByte(string(bg[2]) + string(bg[2]))
	default:
		return color.Transparent
	}

	return color.RGBA{R: r, G: g, B: b, A: 255}
}

func hexToByte(hex string) uint8 {
	var val uint8
	for _, c := range hex {
		val *= 16
		switch {
		case c >= '0' && c <= '9':
			val += uint8(c - '0')
		case c >= 'a' && c <= 'f':
			val += uint8(c - 'a' + 10)
		case c >= 'A' && c <= 'F':
			val += uint8(c - 'A' + 10)
		}
	}
	return val
}
