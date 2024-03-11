# net-cat | imazlin | bdesouza

## REQUIRED

- Before starting this project you will need of go and an IDE
### Installation Golang & IDE Advice
- Linux Terminal
```
sudo apt-get update
sudo apt-get install golang

sudo apt-get update
sudo apt-get install code
```
- Windows Mac  
Go -> [GOLANG](https://go.dev/doc/install)  
Go -> [VSCODE](https://code.visualstudio.com/download)

### Launching Program
After launching the project with vscode you can open
a new terminal in the top of page  
1. Write in the terminal 
```
go mod init NETCAT
```
2. On the same terminal execute the command
```
go get "github.com/jroimartin/gocui"
```
3. Compile the project doing
``` 
go build -o TCPChat
```
ENJOY ;)  

## PRESENTATION
netcat is a project who permit the communications between many clients like a chatbox passing by a server
- To launch 
```
.\TCPChat port
```
if port is not specified the port 8989 is launched by default
- If someone want to connect
```
nc ip port
```

### PROGRAM  
- When a user enter in server the others users are informed, the user received the precedents logs 
- When a user write a message the others users received his message with the time when the message has been sended  
the user who has send the message and the message ``` [2020-01-20 15:48:41][client.name]:[client.message] ```  
- the logs of each server are saved in a file
- the user can rename himself with ``` "/rename" ```
- the user can see the users in the server with ``` -users" ```
- the user can see the availables commands with ``` --help" ```
- the program accept at most 10 persons

### ALGORITHM
- We launch a server with net.Listen, we listen on the connexions with a limit of 10  
- We manage with gofuncs the actions from each connexions in asynchronous  
- The program request to the conn an username to interact with others users
- for each activity, (client, entry, exit, rename) message a message is send by the channel, who send it to the others connexions  

### VIEWS
-When the project started a view of the server is displayed on the terminal, Its the main view of the server
- The main view of the server display the users activities : entrance, exit, rename
- When a user logs on the server a view is made for him, the view display the messages send by the user
- You can kick a user from the serveur by clicking in the button click
- You can start another server in writing is port in the designated area


## FILEST
-To launch the filetest In terminal root of project
```
go test
```

- The filetest test the starting of a server and his limit connexions and  
if the connexions reach the server

## DOCKER

- To launch the docker In terminal root of project
```
cd docker
sh netcat.sh
```

The docker is launched in port 8989, you can open other servers in the range 8990...9000