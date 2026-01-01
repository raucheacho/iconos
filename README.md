# iconos 🎨

Générateur d'icônes CLI pour extensions, web, PWA et mobile.

## Installation

### Homebrew (macOS/Linux)

```bash
brew install raucheacho/tap/iconos
```

### Scoop (Windows)

```powershell
scoop bucket add raucheacho https://github.com/raucheacho/scoop-bucket
scoop install iconos
```

### Go install

```bash
go install github.com/raucheacho/iconos@latest
```

### Build from source

```bash
git clone https://github.com/raucheacho/iconos.git
cd iconos
go build -o iconos .
```

## Utilisation

### Commande de base

```bash
iconos input.png
```

Génère dans `icons/` :

- `icon16.png` (16x16)
- `icon32.png` (32x32)
- `icon48.png` (48x48)
- `icon128.png` (128x128)

### Avec favicons

```bash
iconos input.png --favicon
```

Génère en plus :

- `favicon-16x16.png`
- `favicon-32x32.png`
- `favicon-48x48.png`
- `favicon.ico` (multi-résolutions)

### Options

| Option           | Description            | Défaut         |
| ---------------- | ---------------------- | -------------- |
| `--sizes`, `-s`  | Tailles personnalisées | `16,32,48,128` |
| `--out`, `-o`    | Répertoire de sortie   | `icons`        |
| `--prefix`, `-p` | Préfixe des fichiers   | `icon`         |
| `--format`, `-f` | Format (png, jpg)      | `png`          |
| `--preset`       | Preset de tailles      | -              |
| `--favicon`      | Générer les favicons   | `false`        |
| `--no-ratio`     | Forcer le carré exact  | `false`        |
| `--bg`           | Couleur de fond        | `transparent`  |
| `--html`         | Générer favicon.html   | `false`        |
| `--manifest`     | Générer manifest.json  | `false`        |

### Presets disponibles

| Preset             | Tailles                  |
| ------------------ | ------------------------ |
| `chrome-extension` | 16, 32, 48, 128          |
| `web`              | 16, 32, 48, 64, 128, 256 |
| `pwa`              | 192, 512                 |
| `favicon`          | 16, 32, 48               |

### Exemples

```bash
# Preset Chrome Extension
iconos logo.png --preset chrome-extension

# Preset PWA avec manifest
iconos logo.png --preset pwa --manifest

# Tailles personnalisées
iconos logo.png --sizes 64,128,256,512

# Forcer le carré (sans conserver le ratio)
iconos logo.png --no-ratio

# Fond blanc au lieu de transparent
iconos logo.png --bg "#ffffff"

# Tout générer
iconos logo.png --favicon --html --manifest --preset web
```

## Conservation du ratio

Par défaut, iconos conserve le ratio d'aspect de l'image source :

- L'image est redimensionnée pour tenir dans le carré cible
- Les zones vides sont remplies par un fond transparent (ou la couleur `--bg`)

Avec `--no-ratio`, l'image est étirée pour remplir exactement le carré.

## Support SVG

Pour convertir des SVG en PNG, installez [resvg](https://github.com/RazrFalcon/resvg) :

```bash
# macOS
brew install resvg

# Conversion
resvg input.svg input.png
iconos input.png
```

## Licence

MIT
