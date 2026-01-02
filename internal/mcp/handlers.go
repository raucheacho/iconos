package mcpserver

import (
"context"
"fmt"
"strconv"
"strings"

"iconos/internal/icon"
"iconos/internal/presets"
"iconos/internal/utils"

"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GenerateIconsInput struct {
	Input      string `json:"input" jsonschema:"Chemin vers l'image source (PNG ou JPG)"`
	Output     string `json:"output,omitempty" jsonschema:"Répertoire de sortie (défaut: icons)"`
	Prefix     string `json:"prefix,omitempty" jsonschema:"Préfixe des fichiers (défaut: icon)"`
	Format     string `json:"format,omitempty" jsonschema:"Format de sortie: png ou jpg (défaut: png)"`
	Sizes      string `json:"sizes,omitempty" jsonschema:"Tailles séparées par virgules (ex: 16,32,48,128)"`
	Preset     string `json:"preset,omitempty" jsonschema:"Preset: chrome-extension ou web ou pwa ou favicon"`
	KeepRatio  *bool  `json:"keep_ratio,omitempty" jsonschema:"Conserver le ratio (défaut: true)"`
	Background string `json:"background,omitempty" jsonschema:"Couleur de fond: transparent ou #RRGGBB"`
}

type GenerateIconsOutput struct {
	Message string   `json:"message"`
	Files   []string `json:"files"`
}

func handleGenerateIcons(ctx context.Context, req *mcp.CallToolRequest, input GenerateIconsInput) (*mcp.CallToolResult, GenerateIconsOutput, error) {
	if input.Input == "" {
		return nil, GenerateIconsOutput{}, fmt.Errorf("paramètre 'input' requis")
	}
	if !utils.FileExists(input.Input) {
		return nil, GenerateIconsOutput{}, fmt.Errorf("fichier introuvable: %s", input.Input)
	}
	output := defaultString(input.Output, "icons")
	prefix := defaultString(input.Prefix, "icon")
	format := strings.ToLower(defaultString(input.Format, "png"))
	background := defaultString(input.Background, "transparent")
	keepRatio := true
	if input.KeepRatio != nil {
		keepRatio = *input.KeepRatio
	}
	if format != "png" && format != "jpg" && format != "jpeg" {
		return nil, GenerateIconsOutput{}, fmt.Errorf("format non supporté: %s", format)
	}
	sizes, err := determineSizes(input.Sizes, input.Preset)
	if err != nil {
		return nil, GenerateIconsOutput{}, err
	}
	if err := utils.EnsureDir(output); err != nil {
		return nil, GenerateIconsOutput{}, err
	}
	bgColor := icon.ParseBackgroundColor(background)
	if err := icon.Generate(input.Input, output, prefix, format, sizes, keepRatio, bgColor); err != nil {
		return nil, GenerateIconsOutput{}, err
	}
	var files []string
	for _, size := range sizes {
		files = append(files, fmt.Sprintf("%s/%s%d.%s", output, prefix, size, format))
	}
	return nil, GenerateIconsOutput{Message: "Icônes générées avec succès", Files: files}, nil
}

type GenerateFaviconsInput struct {
	Input      string `json:"input" jsonschema:"Chemin vers l'image source"`
	Output     string `json:"output,omitempty" jsonschema:"Répertoire de sortie (défaut: icons)"`
	KeepRatio  *bool  `json:"keep_ratio,omitempty" jsonschema:"Conserver le ratio (défaut: true)"`
	Background string `json:"background,omitempty" jsonschema:"Couleur de fond (défaut: transparent)"`
	HTML       bool   `json:"html,omitempty" jsonschema:"Générer favicon.html avec les balises link"`
}

type GenerateFaviconsOutput struct {
	Message string   `json:"message"`
	Files   []string `json:"files"`
}

func handleGenerateFavicons(ctx context.Context, req *mcp.CallToolRequest, input GenerateFaviconsInput) (*mcp.CallToolResult, GenerateFaviconsOutput, error) {
	if input.Input == "" {
		return nil, GenerateFaviconsOutput{}, fmt.Errorf("paramètre 'input' requis")
	}
	if !utils.FileExists(input.Input) {
		return nil, GenerateFaviconsOutput{}, fmt.Errorf("fichier introuvable: %s", input.Input)
	}
	output := defaultString(input.Output, "icons")
	background := defaultString(input.Background, "transparent")
	keepRatio := true
	if input.KeepRatio != nil {
		keepRatio = *input.KeepRatio
	}
	if err := utils.EnsureDir(output); err != nil {
		return nil, GenerateFaviconsOutput{}, err
	}
	bgColor := icon.ParseBackgroundColor(background)
	faviconSizes := presets.FaviconSizes()
	if err := icon.GenerateFavicons(input.Input, output, faviconSizes, keepRatio, bgColor); err != nil {
		return nil, GenerateFaviconsOutput{}, err
	}
	var files []string
	for _, size := range faviconSizes {
		files = append(files, fmt.Sprintf("%s/favicon-%dx%d.png", output, size, size))
	}
	files = append(files, fmt.Sprintf("%s/favicon.ico", output))
	if input.HTML {
		if err := icon.GenerateFaviconHTML(output, faviconSizes); err != nil {
			return nil, GenerateFaviconsOutput{}, err
		}
		files = append(files, fmt.Sprintf("%s/favicon.html", output))
	}
	return nil, GenerateFaviconsOutput{Message: "Favicons générés avec succès", Files: files}, nil
}

type ListPresetsInput struct{}
type ListPresetsOutput struct {
	Presets map[string][]int `json:"presets"`
}

func handleListPresets(ctx context.Context, req *mcp.CallToolRequest, input ListPresetsInput) (*mcp.CallToolResult, ListPresetsOutput, error) {
	result := make(map[string][]int)
	for _, name := range presets.ListPresets() {
		sizes, _ := presets.GetPreset(name)
		result[name] = sizes
	}
	return nil, ListPresetsOutput{Presets: result}, nil
}

func determineSizes(sizesStr, presetName string) ([]int, error) {
	if sizesStr != "" {
		return parseSizes(sizesStr)
	}
	if presetName != "" {
		sizes, ok := presets.GetPreset(presetName)
		if !ok {
			return nil, fmt.Errorf("preset inconnu: %s", presetName)
		}
		return sizes, nil
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
	return result, nil
}

func defaultString(val, def string) string {
	if val == "" {
		return def
	}
	return val
}
