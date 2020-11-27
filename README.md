![GitHub language count](https://img.shields.io/github/languages/count/rhiskey/spigot2vk_admin) ![Discord](https://img.shields.io/discord/224962875716796418) ![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/rhiskey/spigot2vk_admin) ![Lines of code](https://img.shields.io/tokei/lines/github/rhiskey/spigot2vk_admin) ![GitHub](https://img.shields.io/github/license/rhiskey/spigot2vk_admin)
![GitHub followers](https://img.shields.io/github/followers/rhiskey?label=Follow%20me&style=social)
# spigot2vk_admin
Resend messages from TCP port to VK Conversation vise versa
1. Build exec file for: `gox`

You can specify package or platform
```
gox -os="linux"
gox -osarch="linux/amd64"
```
2. Place `config.json` nearby your `spigot2vk_admin` executable and run;

By default listens on ports:
8338, 8336

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
https://blog.golang.org/docker
`docker build -t spigot2vk .`
`docker run --name sp2vk --rm spigot2vk`
