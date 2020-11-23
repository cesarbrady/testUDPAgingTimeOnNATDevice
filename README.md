We will often use network address translation on IPv4 network .

After the data packet is forwarded by the NAT gateway , other machine on the public network cannot obtain the IP address of the machine on the private network behind the NAT gateway , but can Only obtain the ip address and port of the NAT gateway .

The NAT device dynamically maintains a table of port mapping for packet forwading ,because there are many machine on the private network that use the same NAT device , but the number of port on NAT  device is limited , so it needs to be automatically delete the port mapping rules within a period of time after no data transmission.

when the port mapping rule takes effect , machines on public network can send data to the machines which on the  internal network, and if the rule is delete, machines on internal network  is unable to receive data from machines on public network .

This period of time is call aging-time .

So , how to test the aging-time for UDP protocol on NAT gateway ?

I wrote a golang program to perform the test , it use kcp protocol which is a reliable UDP protocol .

when the program running in server mode , it will listen a port and wating for clients to send data to , after the connection established , it will send some data to the client at an interval , the interval will increase in each cycle .

when the program running in client mode , it will connect to the server, send some  data to ensure a port mapping is deploy on NAT device , and then wait for the data from the server and print the interval between every two received data .

The lastest time that printedis the aging-time of the UDP protocol on the NAT  device .

```go
package main

var key string = "demo key keykeykeykeykeykeykey"
var salt string = "demo salt saltsaltsaltsaltsaltsalt"

func main() {
    lg := getLogger()

    args := argparser("test kcp")
    side := args.get("", "side", "s", "\"c\" for client, \"s\" for server")
    addr := args.get("", "addr", "127.0.0.1", "address for listen or connect to")
    port := args.getInt("", "port", "12345", "port for listen or connect to")
    sleepSecond := args.getInt("", "start", "30", "sleep second")
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
```

the program runs in server mode that looks like this

```bash
root@server:~# ./kcp --side s --addr 0.0.0.0 --port 12345 --step 1   
11-23 00:44:42   1 [TRAC] listen 0.0.0.0 on 12345   
11-23 00:44:42   1 [TRAC] waiting for connection  
"ping"  
"1606063489"
"1606063519"
...
```

the program runs in client mode that looks like this

```bash
root@xserver:~# ./run --side c --addr j.googleapies.com --port 12345   
11-23 00:44:49   1 [TRAC] connect j.googleapies.com on 12345   
11-23 00:44:49   1 [TRAC] connect success  
11-23 00:44:49   1 [TRAC] send ping done  
11-23 00:45:20   1 [DEBU] "recv:" "1606063519" "30 seconds"  
11-23 00:45:51   1 [DEBU] "recv:" "1606063550" "31 seconds"  
11-23 00:46:23   1 [DEBU] "recv:" "1606063582" "32 seconds"  
11-23 00:46:56   1 [DEBU] "recv:" "1606063615" "33 seconds"
...
11-23 02:34:36   1 [DEBU] "recv:" "1606070075" "1 minute 58 seconds"
11-23 02:36:35   1 [DEBU] "recv:" "1606070194" "1 minute 59 seconds"
11-23 02:38:35   1 [DEBU] "recv:" "1606070314" "2 minutes 0 second"
^C
root@xserver:~#
```

from the output at above , you can seed that the aging-time here is 2 minutes .

Github: https://github.com/cesarbrady/testUDPAgingTimeOnNATDevice

