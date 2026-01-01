package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"iconos/internal/icon"
	"iconos/internal/presets"
	"iconos/internal/utils"

	"github.com/spf13/cobra"
)

var (
	sizes      string
	outDir     string
	prefix     string
	format     string
	preset     string
	favicon    bool
	noRatio    bool
	background string
	genHTML    bool
	genManifest bool

	// Version info
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// SetVersion sets the version info from main
func SetVersion(v, c, d string) {
	version = v
	commit = c
	date = d
}

var rootCmd = &cobra.Command{
	Use:     "iconos [image]",
	Short:   "Générateur d'icônes pour extensions, web, PWA et mobile",
	Version: version,
	Long: `iconos génère automatiquement des icônes redimensionnées et des favicons
à partir d'une image source (PNG ou SVG).

Exemples:
  iconos input.png                          # Génère les icônes par défaut
  iconos input.png --favicon                # Génère aussi les favicons
  iconos input.png --preset pwa             # Utilise le preset PWA
  iconos input.png --sizes 64,128,256       # Tailles personnalisées
  iconos input.png --no-ratio               # Force le carré exact
  iconos input.png --bg "#ffffff"           # Fond blanc au lieu de transparent`,
	Args: cobra.ExactArgs(1),
	RunE: run,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&sizes, "sizes", "s", "", "Tailles des icônes (ex: 16,32,48,128)")
	rootCmd.Flags().StringVarP(&outDir, "out", "o", "icons", "Répertoire de sortie")
	rootCmd.Flags().StringVarP(&prefix, "prefix", "p", "icon", "Préfixe des fichiers")
	rootCmd.Flags().StringVarP(&format, "format", "f", "png", "Format de sortie (png, jpg)")
	rootCmd.Flags().StringVar(&preset, "preset", "", "Preset de tailles (chrome-extension, web, pwa, favicon)")
	rootCmd.Flags().BoolVar(&favicon, "favicon", false, "Générer les favicons (PNG + ICO)")
	rootCmd.Flags().BoolVar(&noRatio, "no-ratio", false, "Forcer le carré exact (ne pas conserver le ratio)")
	rootCmd.Flags().StringVar(&background, "bg", "transparent", "Couleur de fond (#RRGGBB ou transparent)")
	rootCmd.Flags().BoolVar(&genHTML, "html", false, "Générer favicon.html")
	rootCmd.Flags().BoolVar(&genManifest, "manifest", false, "Générer manifest.json pour PWA")
}

func run(cmd *cobra.Command, args []string) error {
	inputPath := args[0]

	// Vérification du fichier source
	if !utils.FileExists(inputPath) {
		return fmt.Errorf("fichier introuvable: %s", inputPath)
	}

	// Validation du format
	format = strings.ToLower(format)
	if format != "png" && format != "jpg" && format != "jpeg" {
		return fmt.Errorf("format non supporté: %s (utilisez png ou jpg)", format)
	}

	// Détermination des tailles
	iconSizes, err := determineSizes()
	if err != nil {
		return err
	}

	// Création du répertoire de sortie
	if err := utils.EnsureDir(outDir); err != nil {
		return err
	}

	// Parse de la couleur de fond
	bgColor := icon.ParseBackgroundColor(background)
	keepRatio := !noRatio

	fmt.Printf("📁 Source: %s\n", inputPath)
	fmt.Printf("📂 Sortie: %s/\n", outDir)
	if keepRatio {
		fmt.Printf("📐 Mode: conservation du ratio (fond: %s)\n", background)
	} else {
		fmt.Printf("📐 Mode: carré exact (ratio ignoré)\n")
	}
	fmt.Println()

	// Génération des icônes
	if err := icon.Generate(inputPath, outDir, prefix, format, iconSizes, keepRatio, bgColor); err != nil {
		return err
	}

	// Génération des favicons si demandé
	if favicon {
		fmt.Println()
		fmt.Println("🌐 Génération des favicons...")
		faviconSizes := presets.FaviconSizes()
		if err := icon.GenerateFavicons(inputPath, outDir, faviconSizes, keepRatio, bgColor); err != nil {
			return err
		}

		if genHTML {
			if err := icon.GenerateFaviconHTML(outDir, faviconSizes); err != nil {
				return err
			}
		}
	}

	// Génération du manifest PWA si demandé
	if genManifest {
		if err := generateManifest(iconSizes); err != nil {
			return err
		}
	}

	fmt.Println()
	fmt.Println("✅ Génération terminée!")
	return nil
}

func determineSizes() ([]int, error) {
	// Priorité: --sizes > --preset > défaut
	if sizes != "" {
		return parseSizes(sizes)
	}

	if preset != "" {
		presetSizes, ok := presets.GetPreset(preset)
		if !ok {
			available := strings.Join(presets.ListPresets(), ", ")
			return nil, fmt.Errorf("preset inconnu: %s (disponibles: %s)", preset, available)
		}
		return presetSizes, nil
	}

	return presets.DefaultSizes(), nil
}

func parseSizes(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	result := make([]int, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		size, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("taille invalide: %s", p)
		}
		if size <= 0 || size > 1024 {
			return nil, fmt.Errorf("taille hors limites (1-1024): %d", size)
		}
		result = append(result, size)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("aucune taille spécifiée")
	}

	return result, nil
}

func generateManifest(sizes []int) error {
	manifest := `{
  "name": "My App",
  "short_name": "App",
  "icons": [
`
	for i, size := range sizes {
		manifest += fmt.Sprintf(`    {
      "src": "%s%d.%s",
      "sizes": "%dx%d",
      "type": "image/%s"
    }`, prefix, size, format, size, size, format)
		if i < len(sizes)-1 {
			manifest += ","
		}
		manifest += "\n"
	}
	manifest += `  ],
  "theme_color": "#ffffff",
  "background_color": "#ffffff",
  "display": "standalone"
}`

	path := outDir + "/manifest.json"
	if err := os.WriteFile(path, []byte(manifest), 0644); err != nil {
		return fmt.Errorf("impossible de créer %s: %w", path, err)
	}
	fmt.Printf("✓ Généré: %s\n", path)
	return nil
}
