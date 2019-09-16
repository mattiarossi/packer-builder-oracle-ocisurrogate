#!/bin/bash

#Cleanup script, unmounts the chroot environment and exports the pool

umount -lf /mnt/dev/
umount -lf /mnt/sys/
umount -lf /mnt/proc/
umount /mnt/boot/efi
umount /mnt/boot
zfs umount -a
zpool export rpool