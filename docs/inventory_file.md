# Inventory file

The inventory file is written in **JSON** format and contains all your machines.

In it's simplest form an _inventory_ file is an array of object with a `name`
property, that's it. But, if you define additional keys on those objects you
will be able to configure how `hotomata` communicates with your machine over SSH
and more (set machine specific vars that can later be used in templates).

## Simplest form

```json
[
  {
    "name": "do-nyc3-web1.example.com"
  }
]
```

According to `hotomata` defaults, when provisioning via SSH, it will try to using:

```yaml
hostname: do-nyc3-web1.example.com
username: root
port: 22
key: ~/.ssh/id_rsa
```

## Machine fields

You can put any keys in the machine object but few of them have a special signification:

| Key          | Default | Description   |
|--------------|---------|---------------|
| name         |         | Pretty name for machine, used in logs, also used as default hostname |
| ssh_hostname |         | Hostname to reach machine, dns or ip |
| ssh_username | root    | Username to login as |
| ssh_port     | 22      | Port to connect to |
| ssh_password |         | (Optional) Password for authentication |
| ssh_key      | ~/.ssh/id_rsa | path on local computer to ssh key to use |
| other fields |         | All other fields will be passed in to templates as global vars |

## Groups

Now, everybody ends up with multiple web instances, 2 load balances, 3 db replicas
and copy and pasting ssh configuration pains the eyes, so, groups exist!

Groups are a JSON objects very similar to machines, but, they support two special
keys:

| Key        | Description |
|------------|-------------|
| group_name | The name to give to the group, has no real use for now but is needed to discern `machines` from `groups` |
| machines   | An array of the machines that will inherit from the groups vars |

```json
[
  {
    "group_name": "web",
    "ssh_port": 2222,
    "ssh_key": "/home/bob/.ssh/aws_infra",
    "provider": "aws",
    "machines": [
      {
        "name": "web1",
        "ssh_hostname": "web1.us-east.aws.example.com"
      },
      {
        "name": "web2",
        "ssh_hostname": "web2.us-east.aws.example.com"
      }
    ]
  },
  {
    "name": "lb",
    "ssh_hostname": "lb.nyc3.do.example.com",
    "ssh_key": "/home/bob/.ssh/do_infra"
  }
]
```

An it goes deeper:

```json
[
  {
    "group_name": "aws",
    "ssh_port": 2222,
    "ssh_key": "/home/bob/.ssh/aws_infra",
    "provider": "aws",
    "machines": [
      {
        "group_name": "db",
        "tags": ["8gb", "db", "ssd", "mem"],
        "role": "db",
        "ssh_username": "ubuntu",
        "machines": [
          {
            "name": "web1",
            "ssh_hostname": "web1.us-east.aws.example.com"
          },
          {
            "name": "web2",
            "ssh_hostname": "web2.us-east.aws.example.com"
          },
          {
            "name": "web3",
            "ssh_hostname": "web3.us-east.aws.example.com"
          }
        ]
      },
      {
        "group_name": "lb",
        "tags": ["2gb", "lb", "spinning", "cpu"],
        "role": "lb",
        "ssh_username": "core",
        "machines": [
          {
            "name": "lb1",
            "ssh_hostname": "lb1.us-east.aws.example.com"
          },
          {
            "name": "lb2",
            "ssh_hostname": "lb2.us-east.aws.example.com"
          }
        ]
      }
    ]
  }
]
```
