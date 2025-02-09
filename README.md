# Fake DNS
A stupid simple DNS server which answers as authoritative nameserver for each zone requested.  
> Attention: Should only be used for testing and mocking. Don't allow access from public to this server!

## How it works
This server will respond to A, NS and SOA record requests with a static response regardless of the zone name.  
Returned nameservers can be configured via json config file.
At the moment an IP can be specified, but it does not have an effect. (Maybe using this later for glue records)

## Why would I do this
Mocking zone availability in a test environment where a zone should be "reachable".  
You can also do this by running a fully fledged DNS server and do some preparation work.
But this is quite an overkill, especially if you don't want to bring a huge dependency to the test environment.

## Build
Binaries are saved within bin directory after calling make command.
```bash
make build # build for current OS
make build_all # build for Linux, Mac OS and Windows
```

## Usage
```bash
./fake-dns -c ./path/to/config.json 
```
Default location is: ``./config.json``  
An example configuration file can be found in ``config.json.dist``.