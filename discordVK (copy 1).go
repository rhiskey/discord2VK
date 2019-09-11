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
	//dg *discordgo.Session
)

var (
	chatID     int64
	consoleMSG string
	vkToken    string
	dg         *discordgo.Session
	botToken   string
)

func init() {
	botToken = ""
	//flag.StringVar(&botToken, "t", "", "Bot Token")
	//flag.Parse()
	vkToken = "" //dublicate

	//Check VK messages in Chat #1
	chatID = 1
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

	////Set channel ID 
	var discordChannelID = "12345"
	//// s.State.ChannelID
	if m.ChannelID == discordChannelID {
		//Get new message
		var kolovanjaPublic2ID = "12345"
		var prefix = ""
		//consoleMSG = ""

		//If not messages from server (BOT) = from users in discord channel
		if m.Author.ID != kolovanjaPublic2ID {
			prefix = "[" + m.Author.Username + "]"
			consoleMSG = prefix + ": " + m.Content
		} else { //if messages from Discord Minecraft Bot
			// Create replacer with pairs as arguments.
			r := strings.NewReplacer(":octagonal_sign:", "&#9940;",
				":white_check_mark:", "&#9989;",
				":heavy_plus_sign:", "&#10133;",
				":heavy_minus_sign:", "&#10134;",
				":skull:", "&#128128;",
				":tada:", "&#127881;",
				":medal:", "&#127942;")

			// Replace all pairs.
			result := r.Replace(m.Content)

			/* 			m.Content = strings.Replace(m.Content, ":octagonal_sign:", "&#9940;", -1)    //stop server
			   			m.Content = strings.Replace(m.Content, ":white_check_mark:", "&#9989;", -1)  //start server
			   			m.Content = strings.Replace(m.Content, ":heavy_plus_sign:", "&#10133;", -1)  //player join
			   			m.Content = strings.Replace(m.Content, ":heavy_minus_sign:", "&#10134;", -1) //player leave
			   			m.Content = strings.Replace(m.Content, ":skull:", "&#128128;", -1)           //dead
			   			m.Content = strings.Replace(m.Content, ":tada:", "&#127881;", -1)            //first time
			   			m.Content = strings.Replace(m.Content, ":medal:", "&#127942;", -1)           //achievement */

			//consoleMSG = formatted
			consoleMSG = result
		}

		//Call sendTOVk
		chatID = 1
		
		sendToVK(vkToken, consoleMSG, chatID)
	}

	/* 	if strings.HasPrefix(m.Content, "/start") {
		go StartCommand(m)
	} */

	// Here you can make other commands

	/* 	// If the message is "ping" reply with "Pong!"
	   	if m.Content == "ping" {
	   		s.ChannelMessageSend(m.ChannelID, "Pong!")
	   	}

	   	// If the message is "pong" reply with "Ping!"
	   	if m.Content == "pong" {
	   		s.ChannelMessageSend(m.ChannelID, "Ping!")
	   	} */
}

// //This Function will be called each time when new message in chat created
func getFromVK(s *discordgo.Session, token string, chID int64) {
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

	updates, _, err := client.GetLPUpdatesChan(100, vkapi.LPConfig{25, vkapi.LPModeAttachments})
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil || !update.IsNewMessage() || update.Message.Outbox() {
			continue
		}

		log.Printf("%s", update.Message.String())
		//update.Message.
		//var id = client.

		//Make for all admins of public
		var myID int64 = 1234
		var iliyaID int64 = 1234
		//Send update.Message from chatID to Discord
		if update.Message.FromID == myID || update.Message.FromID == iliyaID {

			var msgText = update.Message.Text
			var discordChannelID = "1234"

			//Send TO Discord
			messageToDiscordCreate(s, discordChannelID, msgText)

		}

		//if update.Message.Text == "/start" {

		//client.SendMessage(vkapi.NewMessage(vkapi.NewDstFromUserID(update.Message.FromID), "Hello!"))
		//}

	}

	/* 	updates, _, err := client.GetLPUpdatesChan(100, vkapi.LPConfig{25, vkapi.LPModeAttachments})
	   if err != nil {
		   log.Panic(err)
	   } */

	//Send one consoleMSG to chatID!
	//client.SendMessage(vkapi.NewMessage(vkapi.NewDstFromChatID(ID), message))
}

// This function will be called each time when new message appeared in VK
func messageToDiscordCreate(s *discordgo.Session, chID string, msg string) {
	//var msg = message
	// Send text message
	s.ChannelMessageSend(
		chID,
		msg,
		//fmt.Sprintf(message)
	)

	//discordgo.MessageSend(discordgo.MessageCreate(*messageText),)
	//discordgo.MessageCreate()
	//discordgo.Message

}

/* func StartCommand(update vkapi.LPUpdate) {
	user, err := GetUser(update.Message.From.ID)
	if err != nil {
		log.Fatal(err)
		return
	}

	msg := vkapi.NewMessage(vkapi.NewDstFromUserID(update.Message.FromID), fmt.Sprintf("Hello, %s!", user.Name))
	client.SendMessage(msg)
} */
