#!/bin/bash
sed -i 's/GRUB_CMDLINE_LINUX.*/GRUB_CMDLINE_LINUX="crashkernel=auto LANG=en_US.UTF-8 console=tty0 console=ttyS0,9600 rd.luks=0 rd.lvm=0 rd.md=0 rd.dm=0 netroot=iscsi:169.254.0.2:::1:iqn.2015-02.oracle.boot:uefi iscsi_param=node.session.timeo.replacement_timeout=6000 net.ifnames=1 nvme_core.shutdown_timeout=10 ipmi_si.tryacpi=0 ipmi_si.trydmi=0 ipmi_si.trydefaults=0 libiscsi.debug_libiscsi_eh=1 network-config=e2NvbmZpZzogZGlzYWJsZWR9Cg== loglevel=4 boot=zfs root=ZFS=rpool\/ROOT\/oel"/'  /etc/default/grub

echo 'GRUB_PRELOAD_MODULES="part_gpt zfs"' >> /etc/default/grub
grub2-mkconfig -o /boot/efi/EFI/redhat/grub.cfg

sed -i "s/.*set root.*/        set root='hd0,gpt2'/" /boot/efi/EFI/redhat/grub.cfg
sed -i 's/root=ZFS=(null)\/ROOT\/oel//' /boot/efi/EFI/redhat/grub.cfg


dracut -f -v
sed -i '/\/.*xfs/s/^/# /' /etc/fstab
sed -i 's/.*swap.*/\/dev\/zvol\/rpool\/swap\tswap\tswap\t defaults,_netdev,x-initrd.mount 0 0/' /etc/fstab
fs_uuid=$(blkid -o value -s UUID /dev/sdb2); echo $fs_uuid
echo "UUID=${fs_uuid} /boot                       ext4     defaults,_netdev,_netdev,x-initrd.mount 0 0" >> /etc/fstab
grub2-probe /
grub2-install -d /usr/lib/grub/x86_64-efi /dev/sdb