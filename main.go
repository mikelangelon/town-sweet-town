package main

import (
	"bytes"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/mikelangelon/town-sweet-town/logic"
	"github.com/mikelangelon/town-sweet-town/world/npc"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/joelschutz/stagehand"
	"github.com/mikelangelon/town-sweet-town/assets"
	"github.com/mikelangelon/town-sweet-town/common"
	"github.com/mikelangelon/town-sweet-town/graphics"
	"github.com/mikelangelon/town-sweet-town/scenes"
	"log/slog"
)

func main() {

	town1, err := graphics.NewMapScene(assets.TileMapPacked, assets.Town1, assets.TileMapPackedTSX, common.ScreenWidth, common.ScreenHeight, common.Scale)
	if err != nil {
		slog.Error("crash parseTileSet", "error", err)
		return
	}
	houseFactory, err := graphics.NewHouseFactory()
	people1, err := graphics.NewMapScene(assets.TileMapPacked, assets.People1, assets.TileMapPackedTSX, common.ScreenWidth, common.ScreenHeight, common.Scale)
	if err != nil {
		slog.Error("crash parseTileSet", "error", err)
		return
	}
	charFactory, err := graphics.NewCharFactory(assets.Characters, assets.CharactersTSX, common.Scale)
	if err != nil {
		slog.Error("crash parseTileSet", "error", err)
		return
	}
	fancyTownFactory, err := graphics.NewCharFactory(assets.FancyTown, assets.FancyTownTSX, common.Scale)
	if err != nil {
		slog.Error("crash parseTileSet", "error", err)
		return
	}
	tinyTownFactory, err := graphics.NewCharFactory(assets.TileMapPacked, assets.TileMapPackedTSX, common.Scale)
	if err != nil {
		slog.Error("crash parseTileSet", "error", err)
		return
	}
	npcFactory := npc.NPCFactory{
		CharFactory: charFactory,
		Scale:       common.Scale,
	}
	if err != nil {
		slog.Error("crash parseTileSet", "error", err)
		return
	}
	common.BigFont, _ = loadFont(30)
	common.NormalFont, _ = loadFont(20)
	common.TinyFont, _ = loadFont(16)
	common.MegaTinyFont, _ = loadFont(12)
	gameLogic := logic.GameLogic{
		NPCFactory:       npcFactory,
		CharFactory:      charFactory,
		FancyTownFactory: fancyTownFactory,
		TinyTownFactory:  tinyTownFactory,
		HouseFactory:     houseFactory,
		RulesApplier: &npc.RuleApplier{Rules: []npc.Rule{
			npc.RentRule,
			npc.EatingRule,
			npc.RoofRule,
			npc.CompleteTownRule,
		}},
	}
	townAudio, menuAudio := audios()
	state := gameLogic.NextDay(scenes.State{})
	state = scenes.State{
		Status:        scenes.Menu,
		MenuSong:      menuAudio,
		TownSillySong: townAudio,
	}
	people1Scene := scenes.NewEntrance("people", people1)
	town1Scene := scenes.NewTown("town1", town1)
	town1Scene.TransitionPoints = scenes.Transition{
		Position:  common.Position{X: 24 * 16, Y: 6 * 16},
		Scene:     people1Scene,
		Direction: stagehand.RightToLeft,
	}
	people1Scene.TransitionPoints = scenes.Transition{
		Position:  common.Position{X: 0 * 16, Y: 6 * 16},
		Scene:     town1Scene,
		Direction: stagehand.LeftToRight,
	}
	menuScene := scenes.NewMenu(town1Scene, gameLogic)
	town1Scene.MenuScene = scenes.Transition{
		Scene:     menuScene,
		Direction: stagehand.TopToBottom,
	}
	sm := stagehand.NewSceneManager[scenes.State](menuScene, state)
	ebiten.SetWindowSize(common.ScreenWidth, common.ScreenHeight)
	if err := ebiten.RunGame(sm); err != nil {
		slog.Error("something went wrong", "err", err)
	}
}

func loadFont(size float64) (font.Face, error) {
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	}), nil
}

func audios() (*audio.Player, *audio.Player) {
	var player, menuPlayer *audio.Player
	{
		audioContext := audio.NewContext(48000)
		decoded, err := mp3.DecodeWithSampleRate(48000, bytes.NewReader(assets.TownSong))
		if err != nil {
			slog.With("error", err).Error("weird audio issue")
		}

		loop := audio.NewInfiniteLoop(decoded, decoded.Length())
		player, err = audioContext.NewPlayer(loop)
		if err != nil {
			slog.With("error", err).Error("problem creating player")
		}

		decodedMenu, err := mp3.DecodeWithSampleRate(48000, bytes.NewReader(assets.MenuLoop))
		if err != nil {
			slog.With("error", err).Error("weird audio issue")
		}

		loopMenu := audio.NewInfiniteLoop(decodedMenu, decodedMenu.Length())
		menuPlayer, err = audioContext.NewPlayer(loopMenu)
		if err != nil {
			slog.With("error", err).Error("problem creating player")
		}
	}
	return player, menuPlayer
}
