# Hotomata

_Lightweight configuration management tool written in Go (yaml+ssh)_

Hotomata is hopefully easy to use and candidate for serious projects but in all
cases a great learning experience and fun project.

For simpler use cases it is more approachable than bigger players in the field.

The biggest differentiator with other CM tools out there is that _Hotomata_ really
embodies simplicity by promoting small commands doing one thing well that can
be composed in complex operations by composing them together. Unix and Golang
concepts yes, yes. All that without writing custom plugins once and still using
your beloved YAML popular for Ansible or SaltStack folks, no funky DSL to learn.

**Hotomata** draws inspiration from _Ansible_ and _SaltStack_ but also from few
open source, not widely popular projects like [dynport/urknall](http://github.com/dynport/urknall)
or [sudharsh/henchman](http://github.com/sudharsh/henchman) and even few SaaSes
out there like [commando.io](https://commando.io)

- or [devo.ps](https://devo.ps)
- or [bigpanda.io](https://bigpanda.io/)
- or [stackstorm](http://stackstorm.com/)
- Cacti
- Splunk
- VictorOps
- or ...

Hotomata masterplans end up being conceptually close to this:

```
masterplan db
  machines db-master
  vars [...]
  plans
    - common
        vars [...]
        plans
          - hostname
              run: "echo '{{.varx}}' >> /etc/hostname"
          - lang
              vars [...]
              plans
                - upload
                    vars [...]
                    run "scp ..."
                    local true
                - service_reboot
                    vars [...]
                    run "sudo service {{.service}} restart"
    - db
        vars [...]
        plans
          - ...
    - upload
        vars [...]
        run "rsync {{.localDir}} {{.remoteDir}}"
        local true
```

## Getting started

```bash
go get github.com/merd/hotomata/cmd/hotomata
hotomata -h
```

## Documentation

There is plenty of documentation being written, even before code sometimes.
Reading them is a good introduction to what `hotomata` does but mostly how it
solves it:

- [Overview](https://github.com/merd/hotomata/blob/master/docs/overview.md)
- [Masterplan file](https://github.com/merd/hotomata/blob/master/docs/masterplan_file.md)
- [Inventory file](https://github.com/merd/hotomata/blob/master/docs/inventory_file.md)

## CLI Tools

### `hotomata`

Main tool used for execution of **masterplans** against remotes and remote execution
facilities.

### `hotomata-vault`

Tool used to **decrypt**, **encrypt**, **view**, **edit**, **create** and **rekey**
vaults of var files. Those are the place you can store your secrets safely and
commit them source control with less worry than with plaintext.

### `hotomata-inventory`

A simple tool to inspect (**print**) an inventory file's contents as seen and
parsed by `hotomata` and to validate (**check**) that an inventory file is syntactically
valid.

## Contributions

... are welcomed, hoping the docs are giving you a good idea of what's going on,
if not I hope the code is clear enough to understand what's happening, in any case
drop me a line if you want to help or know more

## License

MIT
