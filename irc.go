// vim: ft=go
package main

import (
   "net"
   "fmt"
   "os"
   "bufio"
   "strings"
   "regexp"
)

type IRCConnection struct {
   socket         net.Conn
   pread, pwrite  chan string
   Error          chan os.Error
   nick           string
   user           string
   registered     bool
   server         string
}

func (i *IRCConnection) Connect(server string) os.Error {
   i.server = server
   fmt.Printf("I am now connecting to the '%s' server.\n", i.server)
   var err os.Error
   i.socket, err = net.Dial("tcp", "", i.server)

   if err != nil {
      return err
   }

   println("I have successfully made a connection to the aforementioned server.")
   i.pread = make(chan string, 100)
   i.pwrite = make(chan string, 100)
   i.Error = make(chan os.Error, 10)

   return nil

}

func (irc *IRCConnection) Join(channel string) {
   irc.socket.Write(strings.Bytes("JOIN " + channel + "\r\n"))
}

func IRC(nick string, user string) *IRCConnection {
   irc := new(IRCConnection)
   irc.registered = false
   irc.pread = make(chan string, 100)
   irc.pwrite = make(chan string, 100)
   irc.Error = make(chan os.Error)
   irc.nick = nick
   irc.user = user
   return irc
}


func main() {
   con := IRC("Go", "Go")
   err := con.Connect("eighthbit.net:6667")
   if err != nil {
      fmt.Printf("%s\n", err)
      fmt.Printf("%#v\n", con)
      os.Exit(1)
   }
   b := strings.Bytes("USER hi hi hi hi\r\nNICK go\r\n")
   con.socket.Write(b)
   br := bufio.NewReader(con.socket)
   var source, nick, user, host, printmsg string
   con.Join("#bots")

   for {
      // Live.
      msg, err := br.ReadString('\n')
      if err != nil {
         println(err)
      }
      println(msg)
      msg = msg[0 : len(msg)-2] // kill \r\n
      if msg[0] == ':' {
         if i := strings.Index(msg, " "); i > -1 {
            source = msg[1:i]
            msg = msg[i+1 : len(msg)]
         } else {
            fmt.Printf("Misformed msg from server: %#s\n", msg)
         }
         if i, j := strings.Index(source, "!"), strings.Index(source, "@"); i > -1 && j > -1 {
            nick = source[0:i]
            user = source[i+1 : j]
            host = source[j+1 : len(source)]
         }
      }
      ping,err := regexp.MatchString("PING", msg)
      if ping {
         con.socket.Write(strings.Bytes("PONG :irc.eighthbit.net"))
         println("*** POKE *** I have just ***PONG***'d the server!")
      } else { // Not a ping, so show the message. We also handle commands here.
         args := strings.Split(msg, " :", 2)
         // args = strings.Split(args[0], " ", 0)
         if len(args) > 1 {
           printmsg = fmt.Sprintf("%s (%s@%s) said: %s", nick, user, host, args[1])
           println(printmsg)
         }
      }

   }
}
