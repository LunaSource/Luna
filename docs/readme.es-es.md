# 🌙 Luna — Generador de Commits AI

Luna genera mensajes de commit concisos usando la API de OpenRouter, permitiéndote usar tu propia clave y elegir cualquier modelo compatible.

## ✨ Funcionalidades

- **Selector de archivos interactivo**: elige qué archivos modificados poner en staging, directo desde la TUI
- **Commits por archivo**: un commit por cada archivo seleccionado
- **Impulsado por OpenRouter**: resúmenes AI basados en los diffs de los archivos, usando el modelo que elijas
- **Prefijos convencionales**: agrega el prefijo si falta
- **Emojis elegidos por la IA**: con `-e`, el modelo elige un emoji estilo gitmoji que coincide con la intención del commit
- **Control de longitud**: objetivo < 60 caracteres, máximo configurable (por defecto 72)
- **Filtrado inteligente**: ignora binarios e imágenes comunes

## Cómo funciona

1. Lista todos los archivos modificados vía `git status --porcelain` (modificados, en staging, sin seguimiento), excepto los ignorados
2. Eliges qué archivos incluir: `↑`/`↓` para moverte, `espacio` para seleccionar, `a` para seleccionar todos, `enter` para confirmar
3. Pone en staging los archivos seleccionados (`git add`) y los procesa uno por uno
4. Envía el diff de cada archivo al modelo de OpenRouter configurado
5. Si `-e` está activo, el modelo antepone un emoji acorde a la mensaje; si la respuesta no tiene un prefijo conocido, se elige uno al azar entre: `chore:`, `refactor:`, `feat:`, `fix:`, `docs:`, `test:`, etc.
6. Trunca a `maxCommitLength` y realiza el commit con `git commit -m <mensaje> -- <archivo>`

## Requisitos

- Windows
- Git instalado y disponible en el PATH
- Clave API de OpenRouter (`https://openrouter.ai/keys`)
- Una [Nerd Font](https://www.nerdfonts.com/) configurada como fuente de la terminal, para que los íconos de la interfaz se vean correctamente

## Instalación

### Opción A — Usar el binario precompilado (`bin/Luna.exe`)

1. Copia `bin/Luna.exe` a un directorio, por ejemplo `C:\Users\usuario\Luna`
2. Agrega esa carpeta al PATH del sistema:
   - Presiona `Win + R`, ejecuta `sysdm.cpl`, abre "Variables de entorno"
   - Edita la variable `Path` → "Nuevo" → pega la ruta de la carpeta
   - Guarda y reabre la terminal

### Opción B — Compilar desde el código fuente (Go)

```bash
$ go build -o ./bin/Luna.exe main.go
```

O usa el script auxiliar:

```bash
$ ./build.sh
```

## Configuración

Luna lee la configuración desde archivos de proyecto y globales:

- Proyecto: `.lunacfg` (en la raíz del repositorio o en el directorio padre más cercano)
- Global: `.lunarc` (en el directorio home del usuario)

Prioridad:

- Clave API y modelo: Global → Proyecto → Predeterminado
- Otras configuraciones: Proyecto → Predeterminado

Configuración predeterminada (desde el código):

- `ignoredPatterns`: `*.exe`, `*.dll`, `*.png`, `*.jpg`, `*.jpeg`, `*.gif`, `*.bin`
- `commitPrefixes`: `chore:`, `refactor:`, `feat:`, `fix:`, `docs:`, `test:`
- `maxCommitLength`: `72`
- `defaultEmoji`: `false`

### Configura tu clave API y modelo

```bash
$ Luna apikey TU_CLAVE_OPENROUTER
$ Luna model openai/gpt-4o-mini
```

Esto guarda la clave y el modelo en tu `.lunarc` global. Reabre la terminal después de configurarlos. Explora los modelos disponibles en [openrouter.ai/models](https://openrouter.ai/models).

## Uso

Ejecuta Luna dentro de un repositorio Git con cualquier cambio (en staging, sin staging o sin seguimiento) — elegirás qué archivos poner en staging desde la TUI.

### Comandos y alias

- `help` | `h`: Muestra la ayuda
- `commit` | `c`: Genera y commitea mensajes por archivo
- `apikey <TU_CLAVE>` | `k <TU_CLAVE>`: Configura la clave API
- `model <MODEL_ID>` | `m <MODEL_ID>`: Configura el modelo de OpenRouter
- `config` | `cfg` con subcomandos:
  - `init`: Crea `.lunacfg` en el directorio actual
  - `show`: Muestra la configuración combinada
  - `edit`: Placeholder (aún no implementado)

### Flujo típico

```bash
$ Luna commit # o Luna c
```

Usa `↑`/`↓` para moverte, `espacio` para seleccionar archivos, `a` para seleccionar todos, `enter` para ponerlos en staging y commitearlos.

### Emojis opcionales

```bash
$ Luna c -e       # habilita emojis en los mensajes
```

## Ejemplo de salida

```
Generating commit for file: src/main.go
Committed src/main.go with message:
🚀 feat: add user authentication system

Generating commit for file: README.md
Committed README.md with message:
📝 docs: update installation instructions
```

## Notas

- Luna omite archivos binarios/imágenes comunes
- Si el modelo devuelve una respuesta vacía, el fallback es `update <archivo>`
- Prefijos soportados: `feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`
- `maxCommitLength` se aplica siempre (por defecto 72)

## Solución de problemas

- Error: `Set your API key first: luna apikey <KEY>`
  - Ejecuta `Luna apikey TU_CLAVE` y reabre la terminal
- Error: `Set your model first: luna model <MODEL_ID>`
  - Ejecuta `Luna model MODEL_ID` (ver [openrouter.ai/models](https://openrouter.ai/models))
- Error al ejecutar comandos Git
  - Verifica que estés dentro de un repositorio Git y que Git esté instalado
- No se encontraron cambios
  - No aparece nada en `git status`; modifica o crea un archivo primero
- La clave API no funciona
  - Verifica que la clave sea válida y que el modelo elegido esté disponible en tu cuenta de OpenRouter

---

Hecho con ❤️ por hax — versión 1.3 (Beta)
