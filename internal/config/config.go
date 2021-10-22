package config

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config contains the various values that are configurable for the bot.
type Config struct {
	Bot struct {
		Lang        string
		GuildID     string
		CityChanID  string
		RegCommands bool
		DelCommands bool
	}
	Discord struct {
		Token string
	}
	Repository struct {
		Sqlite struct {
			Location string
		}
	}
}

// GetConfig reads the config file and flags, then applies environment variable overrides.
func GetConfig() (*Config, error) {
	cfg := &Config{}

	// initialize config variables
	viper.SetEnvPrefix("BOT")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/bot")
	viper.AddConfigPath("./configs")

	// read config file
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg.Bot.Lang = viper.GetString("bot.lang")
	cfg.Bot.GuildID = viper.GetString("bot.guild-id")
	cfg.Bot.CityChanID = viper.GetString("bot.city-channel")

	cfg.Repository.Sqlite.Location = viper.GetString("repository.sqlite.location")

	// define flags
	pflag.StringVarP(&cfg.Discord.Token, "token", "t", "", "Bot Token")
	pflag.BoolVarP(&cfg.Bot.RegCommands, "reg-commands", "r", false, "Create Discord slash commands")
	pflag.BoolVarP(&cfg.Bot.DelCommands, "del-commands", "d", false, "Remove Discord slash commands")

	// parse and bind flags
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return nil, fmt.Errorf("failed to bind command line flags: %w", err)
	}

	return cfg, nil
}
