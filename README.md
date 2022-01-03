# qsync
an experimental protocol for safe file synchronization within Qubes OS

see https://zek.manedwolf.website/syncing-files-across-qubes/ for more information

### how-to
- the VM managing the `qsyncd` daemon needs access to the VM list using the Qubes Admin API.  
  to allow access to it, add the following to `/etc/qubes-rpc/policy/admin.vm.List` (assuming aforementioned VM is named `sync`)
  ```
  sync $adminvm allow,target=$adminvm
  ```
