package presets

// Presets définit les tailles d'icônes pour différents cas d'usage
var Presets = map[string][]int{
	"chrome-extension": {16, 32, 48, 128},
	"web":              {16, 32, 48, 64, 128, 256},
	"pwa":              {192, 512},
	"favicon":          {16, 32, 48},
}

// DefaultSizes retourne les tailles par défaut
func DefaultSizes() []int {
	return []int{16, 32, 48, 128}
}

// FaviconSizes retourne les tailles pour les favicons
func FaviconSizes() []int {
	return []int{16, 32, 48}
}

// GetPreset retourne les tailles pour un preset donné
func GetPreset(name string) ([]int, bool) {
	sizes, ok := Presets[name]
	return sizes, ok
}

// ListPresets retourne la liste des presets disponibles
func ListPresets() []string {
	keys := make([]string, 0, len(Presets))
	for k := range Presets {
		keys = append(keys, k)
	}
	return keys
}
