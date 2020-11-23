package main

var key string = "demo key keykeykeykeykeykeykey"
var salt string = "demo salt saltsaltsaltsaltsaltsalt"

func main() {
	lg := getLogger()

	args := argparser("test kcp")
	side := args.get("", "side", "s", "\"c\" for client, \"s\" for server")
	addr := args.get("", "addr", "127.0.0.1", "address for listen or connect to")
	port := args.getInt("", "port", "12345", "port for listen or connect to")
	sleepSecond := args.getInt("", "start", "1", "sleep second")
	step := args.getInt("", "step", "1", "number of second to increase for every loop")
	args.parseArgs()

	if side == "c" {
		lg.trace("connect", addr, "on", port)
		c := kcpConnect(addr, port, key, salt)
		lg.trace("connect success")
		c.send("ping")
		lg.trace("send ping done")
		lasttime := now()
		for {
			buf := c.recv(4096)
			lg.debug("recv:", buf, fmtTimeDuration(now()-lasttime))
			c.send(buf)
			lasttime = now()
		}
	} else if side == "s" {
		lg.trace("listen", addr, "on", port)
		k := kcpListen(addr, port, key, salt)
		lg.trace("waiting for connection")
		c := <-k.accept()
		for {
			c.send(toString(now()))
			print(c.recv(4096))
			sleep(sleepSecond)
			sleepSecond += step
		}
	}
}
