package slashcommands

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/haashi/omega-strikers-bot/internal/discord"
	log "github.com/sirupsen/logrus"
)

type SlashCommand interface {
	Name() string
	Description() string
	Run(s *discordgo.Session, i *discordgo.InteractionCreate)
	Options() []*discordgo.ApplicationCommandOption
	RequiredPerm() *int64
}

var registeredCommands []*discordgo.ApplicationCommand
var commands = []SlashCommand{Join{}, Leave{}, Result{}, Who{}, Link{}, Unlink{}, Update{}, Cancel{}, Credits{}, Predict{}}

// This doesn't perfectly compare options, but I can't be bothered deep checking literally everything.
func compareApplicationCommandOption(o1 *discordgo.ApplicationCommandOption, o2 *discordgo.ApplicationCommandOption) bool {
	return o1.Type == o2.Type &&
		o1.Name == o2.Name &&
		o1.Description == o2.Description &&
		o1.Required == o2.Required
}

func compareApplicationCommandOptions(o1 []*discordgo.ApplicationCommandOption, o2 []*discordgo.ApplicationCommandOption) bool {
	if len(o1) != len(o2) {
		return false
	}
	for i := range o1 {
		if !compareApplicationCommandOption(o1[i], o2[i]) {
			return false
		}
	}
	return true
}

func compareCommands(slashcommand SlashCommand, appcommand *discordgo.ApplicationCommand) bool {
	return appcommand.Name == slashcommand.Name() &&
		appcommand.Description == slashcommand.Description() &&
		compareApplicationCommandOptions(appcommand.Options, slashcommand.Options()) &&
		*appcommand.DefaultMemberPermissions == *slashcommand.RequiredPerm()
}

func Init() {
	session := discord.GetSession()

	commandHandlers := make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
	for _, command := range commands {
		commandHandlers[command.Name()] = command.Run
	}
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	registeredCommands = make([]*discordgo.ApplicationCommand, len(commands))
	previouslyRegisteredCommands, err := session.ApplicationCommands(session.State.User.ID, discord.GuildID)
	if err != nil {
		log.Errorf("cannot get previously registered commands.")
	}
	for i, command := range commands {
		// I don't care about O(n^2) complexity, we won't have that many commands.
		skip := false
		for _, prevCommand := range previouslyRegisteredCommands {
			if compareCommands(command, prevCommand) && os.Getenv("mode") != "prod" {
				registeredCommands[i] = prevCommand
				skip = true
				break
			}
		}
		if skip {
			log.Debugf("skipped registering command %s, as it was a duplicate of a previously declared one.", command.Name())
			continue
		}
		appCommand := &discordgo.ApplicationCommand{
			Name:                     command.Name(),
			Description:              command.Description(),
			Options:                  command.Options(),
			DefaultMemberPermissions: command.RequiredPerm(),
		}
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, discord.GuildID, appCommand)
		if err != nil {
			log.Fatalf("Cannot create '%v' command: %v", command.Name(), err)
		}
		registeredCommands[i] = cmd
	}
}

func Stop() {
	/*session := discord.GetSession()

	log.Println("removing commands...")
	// We need to fetch the commands, since deleting requires the command ID.
	// We are doing this from the returned commands on line 375, because using
	// this will delete all the commands, which might not be desirable, so we
	// are deleting only the commands that we added.
	registeredCommands, err := session.ApplicationCommands(session.State.User.ID, discord.GuildID)
	if err != nil {
		log.Errorf("Could not fetch registered commands: %v", err)
	}

	for _, v := range registeredCommands {
		err := session.ApplicationCommandDelete(session.State.User.ID, discord.GuildID, v.ID)
		if err != nil {
			log.Errorf("cannot delete '%v' command: %v", v.Name, err)
		}
	}*/
}
