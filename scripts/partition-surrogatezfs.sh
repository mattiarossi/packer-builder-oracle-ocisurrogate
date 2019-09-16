#!/bin/bash

## Second stage of the bootstrap script

#Create ZFS pool and datasets

partprobe
zpool create   -f  -o altroot=/mnt     -o ashift=12     -o cachefile=/etc/zfs/zpool.cache     -O canmount=off     -O compression=lz4     -O atime=off     -O normalization=formD   -o feature@hole_birth=disabled -o feature@embedded_data=disabled -m none     rpool     /dev/sdb3
zfs create     -o canmount=off     -o mountpoint=none     rpool/ROOT
zfs create     -o canmount=noauto     -o mountpoint=/     rpool/ROOT/oel
zpool status
zfs mount rpool/ROOT/oel
zfs create     -o setuid=off     -o mountpoint=/home     rpool/home
zfs create     -o mountpoint=/root     rpool/home/root
zfs create     -o setuid=off     -o overlay=on     -o mountpoint=/var     rpool/var
zfs create     -o com.sun:auto-snapshot=false     -o mountpoint=/var/cache     rpool/var/cache
zfs create     -o com.sun:auto-snapshot=false     -o mountpoint=/var/tmp     rpool/var/tmp
zfs create     -o mountpoint=/var/spool     rpool/var/spool
zfs create     -o exec=on     -o mountpoint=/var/lib     rpool/var/lib
zfs create     -o mountpoint=/var/log     rpool/var/log
zfs create     -o mountpoint=/tmp     rpool/tmp
zfs create -V 4G -b $(getconf PAGESIZE) -o compression=zle       -o logbias=throughput -o sync=always       -o primarycache=metadata -o secondarycache=none       -o com.sun:auto-snapshot=false rpool/swap
zfs set quota=8G rpool/tmp
zfs set quota=8G rpool/var
zfs set quota=8G rpool/home
zfs list
mkswap /dev/zvol/rpool/swap

#Setup EFI and Boot partitions
fs_uuid=$(blkid -o value -s UUID /dev/sda1| tr -d "-"); echo $fs_uuid
mkfs.msdos -i $fs_uuid /dev/sdb1
mkfs.ext4 /dev/sdb2

#Setup chroot environment

mkdir -p /mnt/boot
mount /dev/sdb2 /mnt/boot
mkdir -p /mnt/boot/efi
mount /dev/sdb1 /mnt/boot/efi
rsync -ax  / /mnt/
touch /mnt/.autorelabel
mount --bind /dev /mnt/dev/
mount --bind /proc /mnt/proc/
mount --bind /sys/ /mnt/sys/
rsync -ax --delete  /boot/ /mnt/boot/
rsync -ax --delete  /boot/efi/ /mnt/boot/efi/
mkdir -p /mnt/etc/zfs
/bin/cp -p /etc/zfs/zpool.cache /mnt/etc/zfs/zpool.cache 
mkdir -p /mnt/boot/efi/EFI/redhat/x86_64-efi
cp -a /usr/lib/grub/x86_64-efi/zfs* /mnt/boot/efi/EFI/redhat/x86_64-efi
cp -a /usr/lib/grub/x86_64-efi/*dos* /mnt/boot/efi/EFI/redhat/x86_64-efi

