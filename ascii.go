package main

import (
	"fmt"

	"github.com/fatih/color"
)

const (
	asciiTitle1 = "   _         _             _                          _       _       "
	asciiTitle2 = "  | |_ ___ _| |___ ___ ___| |_    _ _ _ ___ ___ _____| |_ ___| |___   "
	asciiTitle3 = "  | . | -_| . |  _| . |  _| '_|  | | | | . |  _|     |   | . | | -_|  "
	asciiTitle4 = "  |___|___|___|_| |___|___|_,_|  |_____|___|_| |_|_|_|_|_|___|_|___|  "

	asciiConfig1 = "             ___ _     "
	asciiConfig2 = " ___ ___ ___|  _|_|___ "
	asciiConfig3 = "|  _| . |   |  _| | . |"
	asciiConfig4 = "|___|___|_|_|_| |_|_  |"
	asciiConfig5 = "                  |___|"

	asciiDiagram1  = " ┌───────────────────────────────────────────────────────────────────┐ "
	asciiDiagram2  = " │                       Your Home Network                           │ "
	asciiDiagram3  = " │                                                                   │ "
	asciiDiagram4  = " │                                                                   │ "
	asciiDiagram5  = " │              'Yeah this is totally a    ┌───────────────────────┐ │ "
	asciiDiagram6  = " │               local Bedrock server;     │    Your Computer      │ │ "
	asciiDiagram7  = " │ ┌──────────┐  connect to me!'           │ ┌───────────────────┐ │ │ "
	asciiDiagram8  = " │ │ Xbox   ◄─┼────────────────────────────┼─┼─►                 │ │ │ "
	asciiDiagram9  = " │ └──────────┘                            │ │ Bedrock Wormhole  │ │ │ "
	asciiDiagram10 = " │ ┌──────────┐                            │ │         ▲         │ │ │ "
	asciiDiagram11 = " │ │ Switch ◄─┼────────────────────────────┼─┼─►       │         │ │ │ "
	asciiDiagram12 = " │ └──────────┘                            │ └─────────┼─────────┘ │ │ "
	asciiDiagram13 = " │                                         └───────────┼───────────┘ │ "
	asciiDiagram14 = " └─────────────────────────────────────────────────────┼─────────────┘ "
	asciiDiagram15 = "                                                       │               "
	asciiDiagram16 = "                                                       ▼               "
	asciiDiagram17 = "                                               ┌─────────────────┐     "
	asciiDiagram18 = "                                               │ Some Internet   │     "
	asciiDiagram19 = "                                               │ Bedrock Server  │     "
	asciiDiagram20 = "                                               └─────────────────┘     "
)

func printRainbowTildes(length int) {
	colors := []*color.Color{
		color.New(color.FgRed),
		color.New(color.FgYellow),
		color.New(color.FgGreen),
		color.New(color.FgCyan),
		color.New(color.FgBlue),
		color.New(color.FgMagenta),
	}
	result := ""
	for i := 0; i < length; i++ {
		result += colors[i%len(colors)].Sprint("~")
	}
	fmt.Println(result)
}

func PrintTitle() {
	color.New(color.FgRed, color.BgBlack).Println()
	color.New(color.FgRed, color.BgBlack).Println(asciiTitle1)
	color.New(color.FgYellow, color.BgBlack).Println(asciiTitle2)
	color.New(color.FgGreen, color.BgBlack).Println(asciiTitle3)
	color.New(color.FgCyan, color.BgBlack).Println(asciiTitle4)
	color.New(color.FgRed, color.BgBlack).Println("                                                                      ")
	fmt.Println()
	printRainbowTildes(80)
	fmt.Println()
}

func PrintDiagram() {
	diagramStyle := color.New(color.FgBlue, color.BgHiBlack)
	diagramStyle.Println(asciiDiagram1)
	diagramStyle.Println(asciiDiagram2)
	diagramStyle.Println(asciiDiagram3)
	diagramStyle.Println(asciiDiagram4)
	diagramStyle.Println(asciiDiagram5)
	diagramStyle.Println(asciiDiagram6)
	diagramStyle.Println(asciiDiagram7)
	diagramStyle.Println(asciiDiagram8)
	diagramStyle.Println(asciiDiagram9)
	diagramStyle.Println(asciiDiagram10)
	diagramStyle.Println(asciiDiagram11)
	diagramStyle.Println(asciiDiagram12)
	diagramStyle.Println(asciiDiagram13)
	diagramStyle.Println(asciiDiagram14)
	diagramStyle.Println(asciiDiagram15)
	diagramStyle.Println(asciiDiagram16)
	diagramStyle.Println(asciiDiagram17)
	diagramStyle.Println(asciiDiagram18)
	diagramStyle.Println(asciiDiagram19)
	diagramStyle.Println(asciiDiagram20)
	fmt.Println()
}

func PrintDescription() {
	white := color.New(color.FgWhite)
	white.Println("This proxy forwards Minecraft Bedrock Edition traffic from this machine")
	white.Println("to a remote server, making this computer appear as the remote server on")
	white.Println("your local network. This allows Xbox, PlayStation, Switch, and mobile")
	white.Println("devices to connect to remote servers as if they were local.")
	fmt.Println()
	white.Println("Both this machine and your gaming devices must be on the same local network.")
	fmt.Println()
}

func PrintWarning() {
	warningStyle := color.New(color.FgBlack, color.BgYellow)
	warningStyle.Println("Warning: This will NOT allow your friends on other networks to play on your")
	warningStyle.Println("local worlds. This will only allow local game consoles to play on remote/internet")
	warningStyle.Println("servers as if it is a LAN world.")
	fmt.Println()
	printRainbowTildes(80)
	fmt.Println()
}

func PrintConfigPrompt() {
	blueBold := color.New(color.FgBlue, color.Bold)
	blueBold.Println(asciiConfig1)
	blueBold.Println(asciiConfig2)
	blueBold.Println(asciiConfig3)
	blueBold.Println(asciiConfig4)
	blueBold.Println(asciiConfig5)
	blueBold.Println("")
}
