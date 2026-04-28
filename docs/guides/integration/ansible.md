---
categories:
- guide
description: How to use Ansible with Rot
title: Ansible
---

Rot can integrate with [Ansible](https://www.ansible.com/) by wrapping commands using {{% cli run %}}.

## Prerequisites

- An existing Rot configuration setup with secret values
- Ansible installed

## Using Rot via Environment Lookup

Ansible can access environment variables using the [`ansible.builtin.env`](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/env_lookup.html) lookup command:

{{< highlight bash >}}
ca:
  crt: "{{ lookup('ansible.builtin.env', 'ca_crt' )}}"
{{< /highlight >}}

This lookup provider can access any environment variable provided to Ansible.  Lets add some Values to Rot and test it out:

{{< highlight bash >}}
$ echo secret1 | rot value-add rot/ansible/SECRET1
$ echo secret2 | rot value-add rot/ansible/SECRET2
$ echo secret3 | rot value-add rot/ansible/SECRET3
{{< /highlight >}}

We'll create a really basic Ansible example playbook that simply debug logs these values to localhost:

{{< highlight bash >}}
$ cat > playbook.yaml << EOF
- hosts: 127.0.0.1
  connection: local
  tasks:
    - name: Print values from Rot
      ansible.builtin.debug:
        msg: "{{ lookup('ansible.builtin.env', 'SECRET1') }} {{ lookup('ansible.builtin.env', 'SECRET2') }} {{ lookup('ansible.builtin.env', 'SECRET3') }}"
EOF
{{< /highlight >}}

Lets execute this playbook by wrapping `ansible-playbook` with {{% cli run %}}:

{{< highlight bash >}}
$ rot run rot/ansible ansible-playbook playbook.yaml
[WARNING]: No inventory was parsed, only implicit localhost is available
[WARNING]: provided hosts list is empty, only localhost is available. Note that
the implicit localhost does not match 'all'

PLAY [127.0.0.1] ***************************************************************

TASK [Gathering Facts] *********************************************************
ok: [127.0.0.1]

TASK [Print values from Rot] ***************************************************
ok: [127.0.0.1] => {
    "msg": "*** *** ***"
}

PLAY RECAP *********************************************************************
127.0.0.1                  : ok=2    changed=0    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0
{{< /highlight >}}

{{% cli run %}} will mask our secrets if they print to stdout, so we should just see *** *** ***.  Lets unmask these just to make sure they're printing OK (don't try this in production!):

{{< highlight bash >}}
$ rot -x unmask='["SECRET1", "SECRET2"]' run rot/ansible ansible-playbook playbook.yaml
[WARNING]: No inventory was parsed, only implicit localhost is available
[WARNING]: provided hosts list is empty, only localhost is available. Note that
the implicit localhost does not match 'all'

PLAY [127.0.0.1] ***************************************************************

TASK [Gathering Facts] *********************************************************
ok: [127.0.0.1]

TASK [Print values from Rot] ***************************************************
ok: [127.0.0.1] => {
    "msg": "secret1 secret2 ***"
}

PLAY RECAP *********************************************************************
127.0.0.1                  : ok=2    changed=0    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0
{{< /highlight >}}

This time we unmasked two of the values and successfully printed them: secret1 and secret2.  You should never need to unmask values as this only applies to stdout--Ansible will use the environment variables which are unmasked.
