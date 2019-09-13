package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	vkapi "github.com/Dimonchik0036/vk-api"
	discordgo "github.com/bwmarrin/discordgo"
)

var (
	chatID     		int64
	consoleMSG 		string
	vkToken  		string
	dg        	        *discordgo.Session
	botToken  		string
        discordChannelID 	string
	discordBotID 		string
 	myID			int64
	
	//CHANGE ALL!
func init() {
	botToken = "" //Get one in Public VK Settings in Work with API: https://vk.com/club12345?act=tokens
	vkToken = "" //https://vkhost.github.io
	//Check VK messages in Chat #1 (first chat of bot's conversation) You need to Enable Bot messages and Bot Ivbiting to chats
	//in VK public settings https://vk.com/club1234?act=messages&tab=bots
	//And Enable LongPoll API v.5.85
	//Invite Public to conversation
	chatID = 1
	//Set channel ID 
	discordChannelID = "12345" 
        discordBotID = "12345" 
	//Owner's ID of VKontakte public
 	myID  = 1234 //CHANGE IT
}

func main() {
	//Discord Part (Get message from chat #)
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
        // Get Messages from chat in VK.com
	getFromVK(dg, vkToken, chatID)
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()

}

//This Function will be called each time when need to send msg to VK
func sendToVK(token string, message string, ID int64) {
	//VK Part

	//client, err := vkapi.NewClientFromLogin("<username>", "<password>", vkapi.ScopeMessages)
	client, err := vkapi.NewClientFromToken(token)
	if err != nil {
		log.Panic(err)
	}

	client.Log(true)

	if err := client.InitLongPoll(0, 2); err != nil {
		log.Panic(err)
	}

	//Send one consoleMSG to chatID!
	client.SendMessage(vkapi.NewMessage(vkapi.NewDstFromChatID(ID), message))
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.ChannelID == discordChannelID {
		//Get new message
		var prefix = ""

		//If not messages from server (BOT) = from users in discord channel
		if m.Author.ID != discordBotID {
			prefix = "[" + m.Author.Username + "]"
			consoleMSG = prefix + ": " + m.Content
		} else { //if messages from Discord Minecraft Bot
			// Create replacer with pairs as arguments.
			//Fixing some EMOJIS to VK format
			r := strings.NewReplacer(":octagonal_sign:", "&#9940;",
				":white_check_mark:", "&#9989;",
				":heavy_plus_sign:", "&#10133;",
				":heavy_minus_sign:", "&#10134;",
				":skull:", "&#128128;",
				":tada:", "&#127881;",
				":medal:", "&#127942;")

			// Replace all pairs.
			result := r.Replace(m.Content)
			consoleMSG = result
		}

		//Call sendTOVk
		sendToVK(vkToken, consoleMSG, chatID)
	}
}

// //This Function will be called each time when new message in chat created
func getFromVK(s *discordgo.Session, token string, chID int64) {
	//VK Part
	client, err := vkapi.NewClientFromToken(token)
	if err != nil {
		log.Panic(err)
	}

	client.Log(true)

	if err := client.InitLongPoll(0, 2); err != nil {
		log.Panic(err)
	}

	updates, _, err := client.GetLPUpdatesChan(100, vkapi.LPConfig{25, vkapi.LPModeAttachments})
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil || !update.IsNewMessage() || update.Message.Outbox() {
			continue
		}

		log.Printf("%s", update.Message.String())

		//Make for all admins of public	
		//Send update.Message from chatID to Discord
		if update.Message.FromID == myID {

			var msgText = update.Message.Text

			//Send TO Discord
			messageToDiscordCreate(s, discordChannelID, msgText)

		}
	}

}

// This function will be called each time when new message appeared in VK
func messageToDiscordCreate(s *discordgo.Session, chID string, msg string) {
	s.ChannelMessageSend(
		chID,
		msg,
	)

}
