package irc

import "log"

import "net"
import "net/textproto"

import "strconv"
import "strings"
import "bufio"
import "errors"

const (
	MaxMessageLength = 510
)

const (
	End = "\r\n"
)

type IrcClient struct {
	Nick       string
	Pass       string
	Host       string
	Port       int
	Connection net.Conn
	CallBack   func(*IrcClient, string)
	Channel    string
}

func NewIrcClient() *IrcClient {
	return &IrcClient{}
}

func CheckPort(irc *IrcClient) *IrcClient {
	if irc.Port == 0 {
		irc.Port = 6667
		return irc
	} else {
		return irc
	}
}

func CheckHost(irc *IrcClient) (*IrcClient, error) {
	if irc.Host == "" {
		log.Fatal("[Error] Host can't be empty")
		return irc, errors.New("[Error] Host can't be empty")
	} else {
		return irc, nil
	}
}

func CheckChannel(irc *IrcClient) (*IrcClient, error) {
	if irc.Channel == "" {
		log.Fatal(("[Error] Channel can't be empty"))
		return irc, errors.New("[Error] Channel can't be empty")
	} else {
		return irc, nil
	}
}

func (i *IrcClient) SendMessage(message string) {
	i.Connection.Write([]byte("PRIVMSG " + i.Channel + " " + message + " " + " \r\n"))
}

func (i *IrcClient) Join() {
	conn, _ := net.Dial("tcp", i.Host+":"+strconv.Itoa(i.Port))
	i.Connection = conn

	i.Connection.Write([]byte("NICK " + i.Nick + " \r\n"))
	i.Connection.Write([]byte("USER " + i.Nick + " nohost noserver :golang\r\n"))
	i.Connection.Write([]byte("JOIN " + i.Channel + " \r\n"))

	start_connect(i)
}

func start_connect(client *IrcClient) {
	reader := bufio.NewReader(client.Connection)
	tp := textproto.NewReader(reader)

	for {

		line, _ := tp.ReadLine()

		resp := strings.Split(line, " ")

		if resp[0] == "PING" {
			client.Connection.Write([]byte("PONG \r\n"))
			continue
		}

		if resp[1] == "PRIVMSG" {
			mess := strings.Split(line, client.Channel)
			client.CallBack(client, mess[len(mess)-1])
			continue
		}

		if resp[0] == "QUIT" {
			client.Connection.Close()
			return
		}
	}
}
