package memory

import (
	"fmt"
	"time"

	game "github.com/septivan/viger/backend/internal/core/game/domain"
	review "github.com/septivan/viger/backend/internal/core/review/domain"
)

type gameSeed struct {
	title       string
	description string
	genre       string
	platforms   []string
	developer   string
	released    string
}

var gameSeeds = []gameSeed{
	{"Baldur's Gate 3", "A story-rich party-based role-playing adventure shaped by choice, chance, and memorable companions.", "RPG", []string{"PC", "PlayStation 5", "Xbox Series"}, "Larian Studios", "2023-08-03"},
	{"Hades", "A fast and expressive roguelike journey through the underworld with a constantly evolving story.", "Roguelike", []string{"PC", "Nintendo Switch", "PlayStation 5", "Xbox Series"}, "Supergiant Games", "2020-09-17"},
	{"Celeste", "A precise platforming adventure about climbing a mountain and confronting personal challenges.", "Platformer", []string{"PC", "Nintendo Switch", "PlayStation 4", "Xbox One"}, "Maddy Makes Games", "2018-01-25"},
	{"Stardew Valley", "A welcoming farming and community simulation with exploration, relationships, and long-term progression.", "Simulation", []string{"PC", "Nintendo Switch", "PlayStation 4", "Xbox One", "Mobile"}, "ConcernedApe", "2016-02-26"},
	{"Disco Elysium", "A dialogue-driven detective role-playing game built around identity, politics, and difficult choices.", "RPG", []string{"PC", "Nintendo Switch", "PlayStation 5", "Xbox Series"}, "ZA/UM", "2019-10-15"},
	{"Outer Wilds", "A compact solar system caught in a time loop rewards curiosity, observation, and fearless exploration.", "Adventure", []string{"PC", "Nintendo Switch", "PlayStation 5", "Xbox Series"}, "Mobius Digital", "2019-05-28"},
	{"The Legend of Zelda: Tears of the Kingdom", "An open-ended adventure across Hyrule that encourages experimentation with inventive building abilities.", "Adventure", []string{"Nintendo Switch"}, "Nintendo EPD", "2023-05-12"},
	{"Elden Ring", "A vast action role-playing world filled with demanding combat, mysterious history, and hidden discoveries.", "Action RPG", []string{"PC", "PlayStation 5", "Xbox Series"}, "FromSoftware", "2022-02-25"},
	{"Hollow Knight", "A hand-drawn action adventure through a ruined kingdom full of secrets and precise combat.", "Metroidvania", []string{"PC", "Nintendo Switch", "PlayStation 4", "Xbox One"}, "Team Cherry", "2017-02-24"},
	{"Portal 2", "A clever first-person puzzle adventure combining spatial reasoning, sharp writing, and cooperative play.", "Puzzle", []string{"PC", "PlayStation 3", "Xbox 360"}, "Valve", "2011-04-19"},
	{"Mass Effect 2", "A cinematic science-fiction role-playing mission centered on assembling and earning the trust of a crew.", "RPG", []string{"PC", "PlayStation 3", "Xbox 360"}, "BioWare", "2010-01-26"},
	{"Slay the Spire", "A tightly designed deck-building roguelike where every card and route changes the next decision.", "Strategy", []string{"PC", "Nintendo Switch", "PlayStation 4", "Xbox One", "Mobile"}, "Mega Crit", "2019-01-23"},
	{"Into the Breach", "A compact turn-based strategy game about preventing attacks before they reshape the battlefield.", "Strategy", []string{"PC", "Nintendo Switch", "Mobile"}, "Subset Games", "2018-02-27"},
	{"Return of the Obra Dinn", "An insurance investigator reconstructs a vanished ship's fate through deduction and frozen moments.", "Puzzle", []string{"PC", "Nintendo Switch", "PlayStation 4", "Xbox One"}, "Lucas Pope", "2018-10-18"},
	{"Dead Cells", "Fluid action and route experimentation drive repeated escapes through an ever-changing island prison.", "Roguelike", []string{"PC", "Nintendo Switch", "PlayStation 4", "Xbox One", "Mobile"}, "Motion Twin", "2018-08-07"},
	{"Ori and the Will of the Wisps", "A graceful platforming journey through a vibrant world with fluid movement and emotional storytelling.", "Metroidvania", []string{"PC", "Nintendo Switch", "Xbox One"}, "Moon Studios", "2020-03-11"},
	{"Control", "A supernatural action adventure set inside a shifting government building filled with strange phenomena.", "Action", []string{"PC", "PlayStation 5", "Xbox Series"}, "Remedy Entertainment", "2019-08-27"},
	{"God of War", "A focused action adventure follows a father and son through a dangerous world of Norse mythology.", "Action", []string{"PC", "PlayStation 4"}, "Santa Monica Studio", "2018-04-20"},
	{"Sekiro: Shadows Die Twice", "Rhythmic sword combat and vertical exploration define a demanding journey through mythic Japan.", "Action", []string{"PC", "PlayStation 4", "Xbox One"}, "FromSoftware", "2019-03-22"},
	{"Monster Hunter: World", "Track enormous creatures, master distinctive weapons, and craft equipment in interconnected ecosystems.", "Action RPG", []string{"PC", "PlayStation 4", "Xbox One"}, "Capcom", "2018-01-26"},
	{"Divinity: Original Sin 2", "A systemic party-based role-playing adventure with tactical combat and flexible cooperative storytelling.", "RPG", []string{"PC", "Nintendo Switch", "PlayStation 4", "Xbox One"}, "Larian Studios", "2017-09-14"},
	{"The Witcher 3: Wild Hunt", "A mature open-world role-playing journey through political conflict, folklore, and personal consequence.", "RPG", []string{"PC", "Nintendo Switch", "PlayStation 5", "Xbox Series"}, "CD Projekt Red", "2015-05-19"},
	{"Red Dead Redemption 2", "A detailed western story follows an outlaw community confronting loyalty, survival, and inevitable change.", "Adventure", []string{"PC", "PlayStation 4", "Xbox One"}, "Rockstar Games", "2018-10-26"},
	{"Subnautica", "Survive and investigate an alien ocean by exploring deeper ecosystems and building new technology.", "Survival", []string{"PC", "Nintendo Switch", "PlayStation 5", "Xbox Series"}, "Unknown Worlds", "2018-01-23"},
	{"Minecraft", "A block-based sandbox supports exploration, construction, survival, and almost unlimited creative projects.", "Sandbox", []string{"PC", "Nintendo Switch", "PlayStation 5", "Xbox Series", "Mobile"}, "Mojang Studios", "2011-11-18"},
	{"Terraria", "A side-scrolling sandbox combines crafting, exploration, construction, and increasingly elaborate boss encounters.", "Sandbox", []string{"PC", "Nintendo Switch", "PlayStation 4", "Xbox One", "Mobile"}, "Re-Logic", "2011-05-16"},
	{"Factorio", "Design increasingly efficient automated factories while managing resources, logistics, and environmental pressure.", "Simulation", []string{"PC", "Nintendo Switch"}, "Wube Software", "2020-08-14"},
	{"Cities: Skylines", "Plan transportation, public services, zoning, and growth in a flexible modern city simulation.", "Simulation", []string{"PC", "Nintendo Switch", "PlayStation 4", "Xbox One"}, "Colossal Order", "2015-03-10"},
	{"Civilization VI", "Guide a civilization across history through expansion, research, diplomacy, culture, and turn-based conflict.", "Strategy", []string{"PC", "Nintendo Switch", "PlayStation 4", "Xbox One", "Mobile"}, "Firaxis Games", "2016-10-21"},
	{"XCOM 2", "Lead a resistance through tactical battles and difficult strategic decisions against an occupying alien force.", "Strategy", []string{"PC", "Nintendo Switch", "PlayStation 4", "Xbox One"}, "Firaxis Games", "2016-02-05"},
	{"Firewatch", "A summer lookout assignment becomes a personal mystery told through exploration and radio conversations.", "Adventure", []string{"PC", "Nintendo Switch", "PlayStation 4", "Xbox One"}, "Campo Santo", "2016-02-09"},
	{"What Remains of Edith Finch", "Explore an unusual family home and experience a collection of inventive, tragic short stories.", "Adventure", []string{"PC", "Nintendo Switch", "PlayStation 5", "Xbox Series"}, "Giant Sparrow", "2017-04-25"},
	{"Inside", "A wordless cinematic puzzle platformer follows a child through a controlled and deeply unsettling world.", "Puzzle", []string{"PC", "Nintendo Switch", "PlayStation 4", "Xbox One", "Mobile"}, "Playdead", "2016-06-29"},
	{"The Witness", "An open island of visual puzzles teaches its language through observation, experimentation, and perspective.", "Puzzle", []string{"PC", "PlayStation 4", "Xbox One", "Mobile"}, "Thekla", "2016-01-26"},
	{"Cuphead", "Hand-animated presentation meets exacting run-and-gun stages and memorable multi-phase boss battles.", "Platformer", []string{"PC", "Nintendo Switch", "PlayStation 4", "Xbox One"}, "Studio MDHR", "2017-09-29"},
	{"Super Mario Odyssey", "A joyful globe-spanning platformer built around expressive movement and playful possession mechanics.", "Platformer", []string{"Nintendo Switch"}, "Nintendo EPD", "2017-10-27"},
	{"A Short Hike", "A warm miniature adventure about climbing a mountain while meeting people and exploring at your pace.", "Adventure", []string{"PC", "Nintendo Switch", "PlayStation 4", "Xbox One"}, "Adam Robinson-Yu", "2019-07-30"},
	{"Spiritfarer", "A gentle management adventure about caring for spirits and learning when it is time to say goodbye.", "Simulation", []string{"PC", "Nintendo Switch", "PlayStation 4", "Xbox One", "Mobile"}, "Thunder Lotus Games", "2020-08-18"},
	{"Doom Eternal", "Aggressive movement, resource management, and weapon switching power a relentless first-person campaign.", "Shooter", []string{"PC", "Nintendo Switch", "PlayStation 5", "Xbox Series"}, "id Software", "2020-03-20"},
	{"Titanfall 2", "A fast first-person campaign combines precise movement, inventive levels, and agile mechanical companions.", "Shooter", []string{"PC", "PlayStation 4", "Xbox One"}, "Respawn Entertainment", "2016-10-28"},
	{"Deep Rock Galactic", "Cooperative miners navigate destructible caves, complete objectives, and survive swarms together.", "Shooter", []string{"PC", "PlayStation 5", "Xbox Series"}, "Ghost Ship Games", "2020-05-13"},
	{"Forza Horizon 5", "An accessible open-world driving festival spans varied landscapes, events, and vehicle collections.", "Racing", []string{"PC", "Xbox Series"}, "Playground Games", "2021-11-09"},
	{"Mario Kart 8 Deluxe", "Colorful kart racing balances readable controls, inventive tracks, and lively local competition.", "Racing", []string{"Nintendo Switch"}, "Nintendo EPD", "2017-04-28"},
	{"Tetris Effect: Connected", "Classic falling-block play becomes a responsive audiovisual journey with cooperative and competitive modes.", "Puzzle", []string{"PC", "Nintendo Switch", "PlayStation 5", "Xbox Series"}, "Monstars", "2020-11-10"},
	{"Vampire Survivors", "Minimal controls conceal an escalating action roguelike built around movement and upgrade combinations.", "Roguelike", []string{"PC", "Nintendo Switch", "PlayStation 5", "Xbox Series", "Mobile"}, "poncle", "2022-10-20"},
	{"Dave the Diver", "Daytime diving and evening restaurant management combine in a playful, constantly changing adventure.", "Adventure", []string{"PC", "Nintendo Switch", "PlayStation 5"}, "MINTROCKET", "2023-06-28"},
	{"Cocoon", "A compact puzzle adventure lets players carry entire worlds and move between their layered mechanisms.", "Puzzle", []string{"PC", "Nintendo Switch", "PlayStation 5", "Xbox Series"}, "Geometric Interactive", "2023-09-29"},
	{"Dredge", "A solitary fishing journey gradually reveals unsettling discoveries beneath a deceptively calm archipelago.", "Adventure", []string{"PC", "Nintendo Switch", "PlayStation 5", "Xbox Series"}, "Black Salt Games", "2023-03-30"},
}

var reviewerNames = []string{"Alex Morgan", "Sam Rivera", "Jordan Lee", "Taylor Kim", "Morgan Chen", "Riley Patel", "Casey Wright", "Jamie Park", "Avery Stone", "Drew Santos", "Robin Blake", "Cameron Hayes"}

var reviewTexts = []string{
	"The core mechanics stay engaging, and the thoughtful pacing made every session feel worthwhile.",
	"A confident game with memorable moments, clear systems, and enough depth to reward experimentation.",
	"The opening takes patience, but the mechanics become remarkably satisfying once everything connects.",
	"Strong art direction and responsive controls carry an experience I was happy to revisit after finishing.",
	"Not every idea lands perfectly, yet the overall design remains focused, expressive, and easy to recommend.",
	"Its best moments come from discovering how the systems interact instead of following a prescribed solution.",
	"The presentation is polished, the difficulty feels fair, and progress consistently introduces something interesting.",
	"I enjoyed the atmosphere and characters most, although a few repeated encounters slowed the middle section.",
}

func Seed() ([]game.Game, []review.Review, error) {
	games := make([]game.Game, 0, len(gameSeeds))
	reviews := make([]review.Review, 0, 320)
	baseTime := time.Date(2026, time.January, 15, 12, 0, 0, 0, time.UTC)
	for index, seed := range gameSeeds {
		releaseDate, err := time.Parse(time.DateOnly, seed.released)
		if err != nil {
			return nil, nil, fmt.Errorf("parse release date for %q: %w", seed.title, err)
		}
		gameID := fmt.Sprintf("game-%03d", index+1)
		item, err := game.New(game.NewGameInput{
			ID: gameID, Title: seed.title, Description: seed.description, Genre: seed.genre,
			Platforms: seed.platforms, Developer: seed.developer, ReleaseDate: releaseDate,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("create seed game %q: %w", seed.title, err)
		}
		games = append(games, item)

		count := (index * 7) % 13
		for reviewIndex := 0; reviewIndex < count; reviewIndex++ {
			rating := 1 + ((index*3 + reviewIndex*5 + reviewIndex/2) % 5)
			item, reviewErr := review.New(review.NewReviewInput{
				ID:           fmt.Sprintf("review-seed-%03d-%02d", index+1, reviewIndex+1),
				GameID:       gameID,
				ReviewerName: reviewerNames[(index+reviewIndex)%len(reviewerNames)],
				Rating:       rating,
				Text:         reviewTexts[(index*2+reviewIndex)%len(reviewTexts)],
				CreatedAt:    baseTime.Add(-time.Duration(index*48+reviewIndex*7) * time.Hour),
			})
			if reviewErr != nil {
				return nil, nil, fmt.Errorf("create seed review for %q: %w", seed.title, reviewErr)
			}
			reviews = append(reviews, item)
		}
	}
	return games, reviews, nil
}
