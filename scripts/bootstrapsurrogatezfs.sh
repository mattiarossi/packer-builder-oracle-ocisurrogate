#!/bin/bash
yum install -y oracle-epel-release-el7.x86_64
THIS_KERNEL_PACKAGE=$(for i in $(yum -v list kernel-uek-devel --show-duplicates); do echo $(uname -r | sed 's/\.x86_64//') | grep -o $i; done)
yum install -y kernel-uek-devel-$THIS_KERNEL_PACKAGE
yum install -y http://download.zfsonlinux.org/epel/zfs-release.el7_6.noarch.rpm
yum install -y zfs zfs-dracut grub2-efi-modules.x86_64
modprobe zfs
dmesg
systemctl preset zfs-import-cache zfs-import-scan zfs-mount zfs-share zfs-zed zfs.target
systemctl enable zfs-import-scan
systemctl list-unit-files | grep zfs

sgdisk -Zg -n1:2048:+210M -t1:EF00 -c1:EFI -n2:0:+1G -t2:EF02 -c2:GRUB -n3:0:0 -t3:BF01 -c3:ZFS /dev/sdb
echo "Rebooting Now ... see you in a short while"
reboot
