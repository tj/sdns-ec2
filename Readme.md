
# sdns-ec2

 EC2 resolver for [SDNS](https://github.com/tj/sdns).

## Installation

 Via binary [releases](https://github.com/tj/sdns-ec2/releases) or:

```
$ go get github.com/tj/sdns-ec2
```

## Example

Config:

```
bind: ":5000"
upstream:
  - 8.8.8.8
  - 8.8.4.4
domains:
  - name: ec2.local.
    command: sdns-ec2
  - name: ec2.
    command: sdns-ec2 --region us-west-1 --ttl 600
```

dig:

```
$ dig @127.0.0.1 -p 5000 +short site-02.ec2
10.30.41.228
```

# License

MIT