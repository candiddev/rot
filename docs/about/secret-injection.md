---
categories:
- feature
description: Rot makes injecting secrets into tools easy
title: Secret Injection
type: docs
---

Users can inject secrets from Rot into their favorite tools like [Ansible]({{< ref "/docs/guides/integration/ansible" >}}) and [OpenTofu/Terraform]({{< ref "/docs/guides/integration/opentofu-terraform" >}}) using Environment Variables:

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
