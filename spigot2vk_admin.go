package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	 "container/list"

	// structs "structs/structs"
	"syscall"
	"time"

	vkapi "github.com/Dimonchik0036/vk-api"
)

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

//INITIAL VARS
var (
	consoleChatID int64

	vkCommunityToken       string
	vkUserToken            string
	portTCPChatUplink      string
	portTCPChatDownlink    string
	portTCPConsoleUplink   string
	portTCPConsoleDownlink string

	portTCPConsoleJSONUplink   string
	portTCPConsoleJSONDownlink string

	IDList []int64 //Init Slice typeof int64 (analog of List in GOlang)

	cfg map[string]interface{}
	//msg map[string]interface{}
	str     strings.Builder // Collect one big string -> send to VK
	msgJSON struct{}        //structure of msg

	// Linked list
	queue := list.New()
)

// TODO: TCPListener2 (MC->VK) one more create read JSON on different port
// TODO: TCPListener2 (VK ->VK) one more create read JSON on different port

// Config
func init() {
	// Read JSON file as []byte
	jsonByte, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	fmt.Println("Contents of config file:", string(jsonByte))

	if json.Unmarshal(jsonByte, &cfg); err != nil {
		panic(err)
	}
	//fmt.Println(cfg)

	consoleChatID = int64(cfg["consoleChatID"].(float64)) // in order to get from the properties num number, 2000000001
	strs := cfg["IDList"].([]interface{})                 // in order to get an array of interfaces ...
	var myID = int64(strs[0].(float64))                   // ... and then get the string from it
	var grishaID = int64(strs[1].(float64))

	vkCommunityToken = cfg["vkCommunityToken"].(string)
	vkUserToken = cfg["vkUserToken"].(string)

	portTCPChatUplink = cfg["portTCPChatUplink"].(string)
	portTCPChatDownlink = cfg["portTCPChatDownlink"].(string)
	portTCPConsoleUplink = cfg["portTCPConsoleUplink"].(string)
	portTCPConsoleDownlink = cfg["portTCPConsoleDownlink"].(string)
	portTCPConsoleJSONUplink = cfg["portTCPConsoleJsonUplink"].(string)
	portTCPConsoleJSONDownlink = cfg["portTCPConsoleJsonDownlink"].(string)

	IDList = append(IDList, myID, grishaID)
	//msgJSON = new(structs.MessageJSON)
}

func main() {
	var isComm bool //isCommunity?
	isComm = false

	// Goroutine
	//JavaPlugin Socket TCP Part (Get message from Java #)
	if isComm == true {
		go TCPServer(portTCPChatDownlink, true) //Read Chat - Send BOT VK
		////Check VK messages in Public Group
		go getFromVK(vkCommunityToken, true) //Read VK BOT
	} else {
		go TCPServer(portTCPConsoleDownlink, false) //Read Console - Send ADMIN CONFA CONSOLE CHAT
		go getFromVK(vkUserToken, false)            //Read CONSOLE CHAT

		go TCPListenerJSON(portTCPConsoleJSONDownlink) //Read Console through JSON -> Send to VK
		//go getFromVK

	}
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

//This Function will be called each time when need to send msg to VK
func sendToVK(token string, message string, IDs []int64, consoleID int64, isCommunity bool) {
	//IDs := []int{1,2,3}
	//VK Part
	client, err := vkapi.NewClientFromToken(token)
	if err != nil {
		log.Panic(err)
	}

	client.Log(false)

	//Do trycatch
	defer func() {
		if err := client.InitLongPoll(0, 2); err != nil {
			log.Panic(err)
			return
		}
	}()
	if isCommunity == true {
		//Send All users
		for _, id := range IDs {
			client.SendMessage(vkapi.NewMessage(vkapi.NewDstFromUserID(id), message))
			//time.Sleep(1500 * time.Millisecond)
		}
	} else {
		//send to ADMIN CHAT
		_, err := client.SendMessage(vkapi.NewMessage(vkapi.NewDstFromChatID(consoleID), message))
		time.Sleep(2 * time.Second)
		if err != nil {
			print("Error Code: \n")
			println(err)
			//time.Sleep(2 * time.Second)
		}
	}
}

// //This Function will be called each time when new message in chat created
func getFromVK(token string, isCommunity bool) { //isCommunity == true => messages from BOT API
	client, err := vkapi.NewClientFromToken(token)
	if err != nil {
		log.Panic(err)
	}

	client.Log(false)

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
		// logs
		log.Printf("%s", update.Message.String())

		//Send update.Message from chatID to Java
		//if update.Message.FromID != myID || update.Message.FromID != iliyaID || update.Message.FromID != grishaID {
		//	continue
		//} else {
		var msgText = update.Message.Text

		client.MarkMessageAsRead(update.Message.ID)

		//Add new user ID from Community Messages
		//Send to TCP Socket Java
		if isCommunity == true {
			var newID = update.Message.FromID //get user ID
			//update.Message.FromChat
			// TODO: Save IDLIST to file AND MYSQL
			// Add new ID TO Chat newsletter
			for _, oneID := range IDList {
				if oneID == newID {
					break
				}
				IDList = append(IDList, newID) //Add new user ID from Community Messages
			}
			go TCPClient(msgText, portTCPChatUplink)
		} else { //Admin console
			//Check if beginning of msgText is '/'  if(strings.Index(msgText, "/")==1)
			if strings.HasPrefix(msgText, "/") || (strings.Index(msgText, "/") == 1) {
				//var prefixIndex = strings.Index(msgText, '/')
				formattedMsg := strings.Replace(msgText, "/", "", 1)
				go TCPClient(formattedMsg, portTCPConsoleUplink)
				go TCPClientJSON(formattedMsg, portTCPConsoleJSONUplink)
			}
		}

	}
}

// TCPServer - get message FROM MC and send to VK
func TCPServer(port string, isComm bool) {
	listener, err := net.Listen("tcp", port)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()
	fmt.Printf("Server is listening..: %v\n", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			//conn.Close()
		}
		go handleConnection(conn, isComm) // goroutine to handle request
		time.Sleep(2000 * time.Millisecond)
	}
}

// TCPListenerJSON - get Json message FROM MC and send to VK
func TCPListenerJSON(port string) {
	listener, err := net.Listen("tcp", port)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()
	fmt.Printf("Server is listening..: %v\n", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			//conn.Close()
		}
		go handleConnectionJSON(conn) // goroutine to handle request
		time.Sleep(2000 * time.Millisecond)
	}
}

// TCP handler
func handleConnection(conn net.Conn, isCommunity bool) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		message := scanner.Text()
		fmt.Println("Message Received:", message)

		// TODO:
		// Collect message in big string, with /n on each EOL

		// if length >2000, then send, refresh string
		//Stack or Queue
		if isCommunity == true {
			sendToVK(vkCommunityToken, message, IDList, consoleChatID, true)
			time.Sleep(2000 * time.Millisecond)
		} else {

			if queue.Len() > 2 {
				// Send message 
				message2send := queue.Front() //First Element
				message2send.Value // -> to VK

				sendToVK(vkUserToken, message2send, IDList, consoleChatID, false)
				time.Sleep(2000 * time.Millisecond)

				queue.Remove(message2send) // Dequee
			} else{
				queue.PushBack(message)
			}

			//queue.PushBack(message)
			
			//sendToVK(vkUserToken, message, IDList, consoleChatID, false)
			//time.Sleep(2000 * time.Millisecond)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error:", err)
	}
}

// TCP handler from console to VK
func handleConnectionJSON(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		message := scanner.Text()
		fmt.Println("Json received:", message)

		jsonBytes := []byte(message)

		msgStruct := MessageJSON{}
		// Json Deserialize to Class type,payload,ip,port -> send only json.payload
		json.Unmarshal(jsonBytes, &msgStruct)

		msgToSend := msgStruct.payload
		// typeOfMessage := msgStruct.typeValue
		// ip := msgStruct.ip
		// port := msgStruct.port

		/// TODO: Queue
		/// Collect message in big string, with /n on each EOL
		// str.WriteString(msgToSend)
		// if str.Len()+len(msgToSend)+1 > 2000 {
		// 	sendToVK(vkUserToken, str.String(), IDList, consoleChatID, false)
		// 	time.Sleep(2000 * time.Millisecond)
		// 	str.Reset()
		// } else {
		// 	sendToVK(vkUserToken, str.String(), IDList, consoleChatID, false)
		// 	time.Sleep(2000 * time.Millisecond)
		// }

		sendToVK(vkUserToken, msgToSend, IDList, consoleChatID, false)
		time.Sleep(2000 * time.Millisecond)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error:", err)
	}
}

// TCPClient - send msg to MC
func TCPClient(message string, port string) {
	conn, err := net.Dial("tcp", port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// send message to spigot server
	if n, err := conn.Write([]byte(message)); n == 0 || err != nil {
		fmt.Println(err)
		return
	}

}

// TCPClientJSON - send msg from VK to MC in JSON
func TCPClientJSON(message string, port string) {
	conn, err := net.Dial("tcp", port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	// TODO: TEST conversation to INT
	portInt, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}
	messageVar := &MessageJSON{
		TypeValue: "console",
		payload:   message,
		ip:        "",
		port:      portInt,
	}
	//Create JSON Serialize -> send
	jsonString, _ := json.Marshal(messageVar)

	// send JSON to spigot server
	if n, err := conn.Write([]byte(jsonString)); n == 0 || err != nil {
		fmt.Println(err)
		return
	}
}
