#!/bin/bash
umount -lf /mnt/dev/
umount -lf /mnt/sys/
umount -lf /mnt/proc/
umount /mnt/boot/efi
umount /mnt/boot
zfs umount -a
zpool export rpool