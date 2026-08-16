[🇺🇸 English](README.md) | [🇧🇷 Português](README.pt-br.md) | 🇪🇸 Español

# Deblock

Un asistente de terminal que levanta un servidor de Minecraft Java Edition (vanilla) en un puñado de pasos — sin buscar el enlace de descarga correcto, sin editar `server.properties` a mano, sin tener que googlear "cómo aceptar el EULA de Minecraft" a las 11 de la noche.

```
██████╗ ███████╗██████╗ ██╗      ██████╗  ██████╗██╗  ██╗
██╔══██╗██╔════╝██╔══██╗██║     ██╔═══██╗██╔════╝██║ ██╔╝
██║  ██║█████╗  ██████╔╝██║     ██║   ██║██║     █████╔╝
██║  ██║██╔══╝  ██╔══██╗██║     ██║   ██║██║     ██╔═██╗
██████╔╝███████╗██████╔╝███████╗╚██████╔╝╚██████╗██║  ██╗
╚═════╝ ╚══════╝╚═════╝ ╚══════╝ ╚═════╝  ╚═════╝╚═╝  ╚═╝
```

## Por qué existe

Configurar un servidor de Minecraft a mano significa encontrar el enlace de descarga correcto, acordarte de qué flags `-Xmx`/`-Xms` pasarle a Java, editar el `server.properties` en un editor de texto, y aceptar un archivo de EULA a mano antes de que nada arranque siquiera. Nada de eso es difícil en sí, es solo lo bastante tedioso como para interponerse en lo que realmente te importa: jugar con tus amigos. Deblock convierte todo eso en una charla corta en la terminal.

## Cómo funciona

```
$ deblock
→ elegí un nombre para el servidor (Deblock crea la carpeta acá mismo)
→ busca la última versión de Minecraft directo desde la API de Mojang
   (o escribís una versión específica vos mismo)
→ configura MOTD, dificultad, modo de juego, whitelist, puerto, memoria
→ pide aceptar el EULA de Mojang
→ descarga el server.jar oficial
→ genera el start.sh / start.bat
→ opcionalmente inicia el servidor ahí mismo, en la misma terminal
```

## Qué hace

- **Asistente interactivo de instalación** — respondé unas preguntas y listo, sin editar archivos a mano
- **Habla directo con las APIs oficiales** (Mojang) para siempre resolver el `server.jar` correcto y actualizado
- **Autocontenido por servidor** — todo vive en una sola carpeta, con el nombre que elijas
- **Detecta instalaciones existentes** — volvé a correr Deblock en la misma carpeta para reconfigurar, reinstalar o simplemente iniciarlo
- **Multiplataforma** — un solo binario en Go, sin dependencias propias (igual necesitás Java para correr el servidor en sí, pero no para correr Deblock)

## Compatibilidad

| Sistema | Correr Deblock | Correr el servidor de Minecraft |
|---|---|---|
| Linux | ✅ | necesita un JDK (21+ recomendado) |
| macOS | ✅ | necesita un JDK (21+ recomendado) |
| Windows | ✅ | necesita un JDK (21+ recomendado) |

Versiones más antiguas de Minecraft pueden requerir una versión más vieja de Java — si el servidor no arranca, es lo primero que hay que revisar.

Por ahora solo aparece **Vanilla** en el menú. El soporte para Paper y Fabric ya está implementado y probado por debajo (mirá [Limitaciones conocidas](#limitaciones-conocidas)).

## Instalación

Todavía no hay binarios listos para descargar, así que por ahora necesitás [Go](https://go.dev/dl/) 1.24+ instalado:

```bash
git clone https://github.com/<tu-usuario>/deblock.git
cd deblock
go run .
```

O compilá un binario propio para reutilizar sin tener que recompilar cada vez:

```bash
go build -o deblock .
./deblock        # Linux/Mac
deblock.exe      # Windows
```

## Uso

Simplemente ejecutalo y respondé las preguntas — Deblock crea la carpeta del servidor justo donde corriste el comando, sin necesidad de escribir ninguna ruta:

```bash
go run .
# o, si ya lo compilaste:
./deblock
```

Te va a pedir un nombre para el servidor, una versión de Minecraft (por defecto la última release), y lo básico del `server.properties` (MOTD, máximo de jugadores, puerto, dificultad, modo de juego, whitelist, online-mode, memoria). Al final te pregunta si querés iniciar el servidor de una vez.

### Administrar un servidor que ya configuraste

Corré Deblock de nuevo dentro de esa misma carpeta y va a reconocer la instalación existente:

- **Editar la configuración** — cambiá el `server.properties` sin tocar el jar
- **Reinstalar desde cero** — borra el jar y lo descarga de nuevo
- **Solo iniciarlo** — se salta toda la configuración

## Comandos útiles

```bash
# corre el asistente
go run .

# corre la suite de tests
go test ./...

# inicia un servidor ya configurado sin pasar por el asistente
cd carpeta-de-mi-servidor
./start.sh       # Linux/Mac
start.bat        # Windows
```

## Estructura

```
deblock/
├── README.md
├── README.pt-br.md
├── README.es.md
├── LICENSE
├── go.mod
├── go.sum
├── main.go
├── main_test.go
└── internal/
    ├── loaders/       # habla con las APIs de Mojang/PaperMC/Fabric
    ├── download/       # descarga el server.jar con barra de progreso
    ├── props/          # lee y escribe el server.properties
    └── startscript/    # genera el start.sh / start.bat
```

## Limitaciones conocidas

- Por ahora solo Vanilla aparece en el menú. Paper y Fabric ya están completamente implementados (`internal/loaders`), pero se mantienen ocultos hasta pasar por más pruebas en el mundo real.
- Un servidor por carpeta — no hay un panel de múltiples servidores, simplemente corré Deblock de nuevo en otra carpeta para tener un segundo servidor.
- Deblock no expone tu servidor a internet (port forwarding, túneles, etc.) — eso queda de tu lado, a propósito, ya que depende mucho de tu propia configuración de red.
- Todavía no hay binarios listos para descargar — compilar desde el código fuente con Go es, por ahora, la única forma de instalarlo.
- Cada loader habla con una API pública de terceros (Mojang, PaperMC, Fabric). Si alguna cambia el formato de su respuesta, ese loader en particular puede romperse hasta que se corrija — están aislados en sus propios archivos dentro de `internal/loaders` justamente para que eso sea un arreglo chico, no una reescritura.

## Licencia

MIT
