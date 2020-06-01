# red
Redis stream testing tool


Generates payloads:

```
Usage:
  red generate [flags]

Flags:
  -n, --count int        Number of items per run (default 1000)
  -d, --delay int        Number of milliseconds to delay between runs (default 1000)
  -h, --help             help for generate
  -p, --payload string   Filename for sample payload data

Global Flags:
  -H, --host string     Local host/IP redis is on (default "127.0.0.1")
  -k, --key string      Name of key in stream to use (default "keyname")
  -P, --port string     Port redis is listening on (default "6379")
  -s, --stream string   Name of stream to use (default "teststream")
  -t, --threads int     Number of threads to run in parallel (default 2)
```

Consumes items (payloads):

```
Usage:
  red consume [flags]

Flags:
  -c, --count int     Number of items to batch receive (default 1000)
  -d, --delay int     Number of milliseconds to delay between runs (default 1000)
  -h, --help          help for consume
  -n, --name string   Name of client consuming stream data (default "client")

Global Flags:
  -H, --host string     Local host/IP redis is on (default "127.0.0.1")
  -k, --key string      Name of key in stream to use (default "keyname")
  -P, --port string     Port redis is listening on (default "6379")
  -s, --stream string   Name of stream to use (default "teststream")
  -t, --threads int     Number of threads to run in parallel (default 2)
```
