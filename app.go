package main

import (
	"fmt"
	"github.com/xlab/closer"
	"log"
	"lyag/bot"
	"lyag/core"
	"lyag/data"
	"lyag/storage"
	"os"
	"strconv"
	"time"
)

func getRequiredVar(name string) string {
	value, isSet := os.LookupEnv(name)

	if isSet == false {
		fmt.Printf("Required %s environment variable\n", name)
		os.Exit(1)
	}

	return value
}

func main() {
	const tokenVar = "LYAG_TG_BOT_TOKEN"
	const clientNameVar = "LYAG_CLIENT_NAME"
	const clientPhoneVar = "LYAG_CLIENT_PHONE"
	const clientEmailVar = "LYAG_CLIENT_EMAIL"
	const chatIdVar = "LYAG_TG_BOT_CHAT_ID"

	token := getRequiredVar(tokenVar)
	clientName := getRequiredVar(clientNameVar)
	clientPhone := getRequiredVar(clientPhoneVar)
	clientEmail := getRequiredVar(clientEmailVar)
	chatId, err := strconv.ParseInt(getRequiredVar(chatIdVar), 10, 64)

	if err != nil {
		log.Panicf("%s must be int!\n%s", chatIdVar, err)
	}

	staffIds := make([]string, 0, len(data.StaffMap))
	for k := range data.StaffMap {
		staffIds = append(staffIds, k)
	}

	duration := time.Minute * 3

	settings := core.Settings{
		StaffIds:      staffIds,
		DaysDepth:     7,
		Duration:      duration,
		TelegramToken: token,
		Name:          clientName,
		Phone:         clientPhone,
		Email:         clientEmail,
		ChatID:        chatId,
	}

	go run(settings)
	closer.Hold()
}

func run(settings core.Settings) {
	token, err := core.FetchToken("https://n18056.yclients.com/company:104300")

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Token found: %s", token)

	yClients := core.NewYClients(token, "", &settings)
	redis := storage.NewRedis()

	lyagBot, err := bot.NewLyagBot(&settings, yClients, redis)

	if err != nil {
		log.Panicf("Can't create LyagBot: %s", err)
		return
	}

	for {
		log.Println("Begin")

		lyagBot.Run(true)

		log.Printf("Waiting %0.0f sec", settings.Duration.Seconds())
		time.Sleep(settings.Duration)
		log.Println("One more time...")
	}
}
