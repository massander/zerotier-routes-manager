# ZeroTier Managed Routes CLI

A command-line tool to manage ZeroTier network routes. This CLI allows to set routes to apps or websites in a convinient way, providing control over network routing without needing to access the ZeroTier web interface.

## Configuration

Ensure you have a valid ZeroTier API token and have completed vpn `exit-node` setup as suggested in Zero Tier Docs [docs](https://docs.zerotier.com/exitnode) and have at least 2 authorized device in the network.

Network ID, Exit Node address and LAN config can be passed directly as config value:  

```
{
    networkId: 0000000000000000
    exitNode: 192.168.196.x
    lan: 192.168.196.0/24
    apps: [
        ...
    ]
}
```

ZeroTier API token must be set as environment variable:

```bash
export ZT_TOKEN=zeroTierApiTokenToChange
```

Full config example (zt_route.hjson):

```
    networkId: 0000000000000000
    exitNode: 192.168.196.x
    lan: 192.168.196.0/24
    apps: [
        {
            name: chatgpt
            domains: [
                chatgpt.com
                chat.openai.com
            ]
            // next value is not neccessary
            ips: [
                188.114.99.229
                188.114.98.229
            ]
        }
        {
            name: soundcloud
            domains: [
                soundcloud.com
            ]
        }
    ]
```

After runnig command config will be updated with Zero Tier routing ruls for each group. Example:

```
{
    ...
    apps: [
        {
            name: xxxxxxxxx
            ...
            routes: [
                {
                    // domain addreess 
                    target: xxx.xxx.xxx.xxx/32
                    // exit node address defined earlier
                    via: xxx.xxx.xxx.xxx 
                }
            ]
        }
    ]
    
}
```

## Usage

```
zt-routes.exe --help
zt-routes is a CLI tool for managing ZeroTier Managed Routes

Usage:
  zt-routes [flags]

Flags:
  -c, --config string   config file (default "./zt_routes.hjson")
      --debug           option to update local config without updating ZT config
  -h, --help            help for zt-routes
  -v, --version         version for zt-routes
```