
# anydns

 Little nameserver resolving via arbitrary command(s).

 __Warning__: This is a work-in-progress, you have been warned!

## Usage

 Run with config:

```
$ anydns < domains.yml
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

# License

 MIT