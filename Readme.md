
# sdns

 Little nameserver resolving via arbitrary command(s).

 __Warning__: This is a work-in-progress, you have been warned!

## Installation

 Via binary [releases](https://github.com/tj/sdns/releases) or:

```
$ go get github.com/tj/sdns
```

## Usage

 Run with config:

```
$ sdns < domains.yml
```

 Configuration example:

```yml
bind: ":5000"
upstream:
  - 8.8.8.8
  - 8.8.4.4
domains:
  - name: foo
    command: |
      echo '{ "type": "A", "value": "1.1.1.1", "ttl": 60 }'
  - name: bar
    command: |
      echo '{ "type": "A", "value": "1.1.1.2", "ttl": 60 }'
  - name: foo.bar
    command: |
      echo '{ "type": "A", "value": "1.1.1.3", "ttl": 300 }'
  - name: boom
    command: |
      echo 'something goes boom' && exit 1
```

 Dig it:

```
$ dig @127.0.0.1 -p 5000 something.foo +short
1.1.1.1

$ dig @127.0.0.1 -p 5000 something.bar +short
1.1.1.2

$ dig @127.0.0.1 -p 5000 something.foo.bar +short
1.1.1.3

$ dig @127.0.0.1 -p 5000 segment.com +short
54.213.169.105
```

# License

 MIT