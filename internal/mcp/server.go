package mcpserver

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Serve démarre le serveur MCP sur stdio
func Serve(version string) error {
	server := mcp.NewServer(
		&mcp.Implementation{Name: "iconos", Version: version},
		nil,
	)

	// Tool: generate_icons
	mcp.AddTool(server, &mcp.Tool{
		Name:        "generate_icons",
		Description: "Génère des icônes redimensionnées à partir d'une image source (PNG ou JPG)",
	}, handleGenerateIcons)

	// Tool: generate_favicons
	mcp.AddTool(server, &mcp.Tool{
		Name:        "generate_favicons",
		Description: "Génère des favicons PNG et un fichier ICO multi-résolutions",
	}, handleGenerateFavicons)

	// Tool: list_presets
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_presets",
		Description: "Liste les presets de tailles disponibles",
	}, handleListPresets)

	return server.Run(context.Background(), &mcp.StdioTransport{})
}
