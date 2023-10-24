module main

import net.websocket

const (
    oauth = "gp762nuuoqcoxypju8c569th9wz7q5"
    user = "deki_helper"
    channel = "deki_senpai_tm"
)

fn main() {
	mut ws := websocket.new_client('wss://irc-ws.chat.twitch.tv:443', ) or {
		panic(err)
	}

	println("asdasd")

	ws.on_open(fn (mut c websocket.Client) ! {
		println("Starting")
		c.write_string('PASS oauth:$oauth') or { panic(err) }
		println("Auth done")
		c.write_string('NICK $user') or { panic(err) }
		println("Nick done")
		c.write_string('JOIN #$channel') or { panic(err) }
		println("Join done")
	})

	ws.on_message(fn (mut c websocket.Client, m &websocket.Message) ! {
		if m.payload.str().contains('Hello World') {
			c.write_string('PRIVMSG #$channel :cringe') or { panic(err) }
		}
		if m.payload.str().contains('PING') {
			c.write_string('PONG') or { panic(err) }
		}
	})

	ws.listen() or { panic(err) }
}
