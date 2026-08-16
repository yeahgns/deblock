[🇺🇸 English](README.md) | 🇧🇷 Português | [🇪🇸 Español](README.es.md)

# Deblock

Um assistente de terminal que sobe um servidor Minecraft Java Edition em poucos toques de teclado. Sem caçar link de download, sem editar `server.properties` na mão, sem precisar pesquisar "como aceitar a EULA do Minecraft" às 11 horas da noite.

```
██████╗ ███████╗██████╗ ██╗      ██████╗  ██████╗██╗  ██╗
██╔══██╗██╔════╝██╔══██╗██║     ██╔═══██╗██╔════╝██║ ██╔╝
██║  ██║█████╗  ██████╔╝██║     ██║   ██║██║     █████╔╝
██║  ██║██╔══╝  ██╔══██╗██║     ██║   ██║██║     ██╔═██╗
██████╔╝███████╗██████╔╝███████╗╚██████╔╝╚██████╗██║  ██╗
╚═════╝ ╚══════╝╚═════╝ ╚══════╝ ╚═════╝  ╚═════╝╚═╝  ╚═╝
```

## Por que isso existe

Configurar um servidor Minecraft do jeito manual significa achar o link certo de download, lembrar quais flags `-Xmx`/`-Xms` passar pro Java, editar o `server.properties` num editor de texto, e aceitar um arquivo de EULA na mão antes de qualquer coisa sequer ligar.

Nada disso é difícil, na verdade. Só chato o suficiente pra atrapalhar a parte que você realmente quer: jogar com os amigos.

O Deblock transforma tudo isso numa conversinha rápida no terminal.

## Como funciona

```
$ deblock
→ escolhe um nome pro servidor (o Deblock cria a pasta pra ele aqui mesmo)
→ busca a versão mais recente do Minecraft direto na API da Mojang
   (ou deixa você digitar uma versão específica)
→ configura MOTD, dificuldade, modo de jogo, whitelist, porta, memória
→ pede pra aceitar a EULA da Mojang
→ baixa o server.jar oficial
→ gera o start.sh / start.bat
→ pergunta se quer já iniciar o servidor ali mesmo, no mesmo terminal
```

## O que ele faz

- **Assistente interativo de instalação** — responde umas perguntas e pronto, sem editar arquivo na mão
- **Fala direto com as APIs oficiais** (Mojang) pra sempre resolver o `server.jar` certo e atualizado
- **Autocontido por servidor** — tudo fica numa pasta só, com o nome que você escolher
- **Detecta instalações existentes** — rode o Deblock de novo na mesma pasta pra reconfigurar, reinstalar ou só iniciar
- **Multiplataforma** — um binário só em Go, sem dependência própria (você ainda precisa do Java pra rodar o servidor em si, só não pra rodar o Deblock)

## Compatibilidade

| Sistema | Rodar o Deblock | Rodar o servidor Minecraft |
|---|---|---|
| Linux | ✅ | precisa de um JDK (21+ recomendado) |
| macOS | ✅ | precisa de um JDK (21+ recomendado) |
| Windows | ✅ | precisa de um JDK (21+ recomendado) |

Versões mais antigas do Minecraft podem exigir uma versão mais antiga do Java. Se o servidor não subir, é a primeira coisa a checar.

Só o **Vanilla** aparece no menu por enquanto. O suporte a Paper e Fabric já está implementado e testado por baixo do capô (veja [Limitações conhecidas](#limitações-conhecidas)).

## Instalação

Ainda não tem binário pronto pra baixar, então por enquanto você precisa do [Go](https://go.dev/dl/) 1.24+ instalado:

```bash
git clone https://github.com/yeahgns/deblock.git
cd deblock
go run .
```

Ou compila um binário próprio pra reusar sem precisar recompilar toda vez:

```bash
go build -o deblock .
./deblock        # Linux/Mac
deblock.exe      # Windows
```

## Uso

É só rodar e responder as perguntas — o Deblock cria a pasta do servidor bem onde você rodou o comando, sem precisar digitar caminho nenhum:

```bash
go run .
# ou, se já compilou:
./deblock
```

Ele vai pedir um nome pro servidor, uma versão do Minecraft (padrão é a última release), e o básico do `server.properties` (MOTD, máximo de jogadores, porta, dificuldade, modo de jogo, whitelist, online-mode, memória). No final, pergunta se você quer já iniciar o servidor.

### Gerenciando um servidor que você já configurou

Roda o Deblock de novo dentro dessa mesma pasta e ele reconhece a instalação existente:

- **Editar as configurações** — muda o `server.properties` sem mexer no jar
- **Reinstalar do zero** — apaga o jar e baixa de novo
- **Só iniciar** — pula a configuração inteira

## Comandos úteis

```bash
# roda o assistente
go run .

# roda a suíte de testes
go test ./...

# inicia um servidor já configurado sem passar pelo assistente
cd pasta-do-meu-servidor
./start.sh       # Linux/Mac
start.bat        # Windows
```

## Estrutura

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
    ├── loaders/       # fala com as APIs da Mojang/PaperMC/Fabric
    ├── download/       # baixa o server.jar com barra de progresso
    ├── props/          # le e escreve o server.properties
    └── startscript/    # gera o start.sh / start.bat
```

## Limitações conhecidas

- Só o Vanilla aparece no menu por enquanto. Paper e Fabric já estão totalmente implementados (`internal/loaders`), só ficam escondidos até passarem por mais teste no mundo real.
- Um servidor por pasta — não tem um painel de múltiplos servidores, é só rodar o Deblock de novo numa pasta diferente pra ter um segundo servidor.
- O Deblock não expõe seu servidor pra internet (port forward, túneis, etc.) — isso fica por sua conta, de propósito, já que depende muito da sua própria configuração de rede.
- Ainda não tem binário pronto pra baixar — compilar a partir do código-fonte com Go é o único jeito de instalar por enquanto.
- Cada loader fala com uma API pública de terceiro (Mojang, PaperMC, Fabric). Se algum deles mudar o formato da resposta, aquele loader especificamente pode quebrar até ser corrigido — eles ficam isolados em arquivos próprios dentro de `internal/loaders` exatamente pra isso ser um ajuste pequeno, não uma reescrita.

## Licença

MIT