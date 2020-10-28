# spigot2vk_admin
Resend messages from TCP port to VK Conversation vise versa
1. Build exec file for: `gox`

You can specify package or platform
```
gox -os="linux"
gox -osarch="linux/amd64"
```
2. Place `config.json` nearby your `spigot2vk_admin` executable and run;
`config.json`:
```
{
    "consoleChatID": 0,
    "vkCommunityToken": "YOUR_VK_GROUP_TOKEN",
    "vkUserToken": "YOUR_USER_TOKEN",
    "portTCPChatUplink": "IP:PORT",
    "portTCPChatDownlink": ":PORT",
    "portTCPConsoleUplink": "IP:PORT",
    "portTCPConsoleDownlink": ":PORT",
    "portTCPConsoleJsonUplink": "IP:PORT",
    "portTCPConsoleJsonDownlink": ":PORT",
    "IDList": [VK_USER_ID_1, VK_USER_ID_2]
}
```
