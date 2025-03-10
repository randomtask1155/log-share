# log-share 

Host a basic auth web sever to share a temporary folder of files.  

usage

```
Usage of ./log-share:
  -cert string
        server cert
  -d string
        define a folder to share defaults to current working directory (default ".")
  -key string
        server key
  -p string
        define password or skip and program will generate one
  -port string
        set a static listening port (default "0")
  -t string
        How long to run before exiting: Valid time units are “ns”, “us” (or “µs”), “ms”, “s”, “m”, “h”. (default "24h")
  -u string
        define username defaults to admin (default "admin")
```

# examples

## share logs for 48 hours 

Remember valid time units are “ns”, “us” (or “µs”), “ms”, “s”, “m”, “h” where the default is "24h"

```
log-share -t 48h
```

## sepcify a specific user/password

```
log-share -u myuser -p mypassword
```

## start in TLS mode

```
log-share -cert cert.pem -key key.pem
```

## share a specific path

```
log-share -d /path/to/share
```

# nohup method

```
nohup log-share & tail -f nohup.out
```

# screen method
See this for a screen cheat-sheet https://gist.github.com/jctosta/af918e1618682638aa82

create a session

```
screen -S log-share-screen
```

Change into the directory you want to share 

```
cd /some/path/i/want/to/share
```

launch log share

```
$log-share &
[1] 108130
[#######@####### ########.bcm]$2025/03/07 14:32:34 staring server:  http://admin:XjJWvJLkKD1u@x.x.x.x:37891
2025/03/07 14:32:34 server will shutdown in  24h
```

detach screen 

```
crtl+a+d
```

jump back into screen session at a later time for more log sharing

```
screen -ls
screen -r log-share-screen
```


# creating a private self signed cert

```
openssl req -x509 \
   -nodes -newkey rsa:4096 -days 700 \
   -keyout key.pem \
   -out cert.pem \
   -subj "/C=US/ST=STATE/L=CITY/O=ORG/OU=ORG-UNIT/CN=logShare"
```

