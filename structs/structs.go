package structs

import "fmt"

// Config - data structure to work with json
type Config struct {
	consoleChatID int64

	vkCommunityToken       string //Community Token
	vkUserToken            string //Kolovanja Server
	portTCPChatUplink      string // outcoming messages FROM VK Community BOT to MC
	portTCPChatDownlink    string // incoming messages  FROM MC to VK Community BOT
	portTCPConsoleUplink   string // outcoming messages FROM VK Admin Chat to MC
	portTCPConsoleDownlink string // incoming messages  FROM MC to VK Admin Chat

	portTCPConsoleJSONUplink   string // outcoming messages FROM VK Admin Chat to MC
	portTCPConsoleJSONDownlink string // incoming messages  FROM MC to VK Admin Chat

	IDList []int64 //Init Slice typeof int64 (analog of List in GOlang), admin VK ID
}

// MessageJSON - data structure to work with json
type MessageJSON struct {
	TypeValue string `json:"type"`
	payload   string
	ip        string
	port      int64
}

func init() {
	fmt.Println("structs package initialized")
}

// func (s Config) Data(name string) struct {

// 	return s
// }
