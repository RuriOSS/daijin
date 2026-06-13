#!/bin/sh
#
# SPDX-License-Identifier: Apache-2.0
# This file is part of termux-container.
#
# Copyright (c) 2023 Moe-hacker
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
#
# Fix `command not found`
PATH=$PATH:"/bin:/sbin:/usr/bin"
# Fix networking and other permission issues on Android.
# From Tmoe: https://github.com/2moe/Tmoe
groupadd aid_system -g 1000 || groupadd aid_system -g 1074
groupadd aid_radio -g 1001
groupadd aid_bluetooth -g 1002
groupadd aid_graphics -g 1003
groupadd aid_input -g 1004
groupadd aid_audio -g 1005
groupadd aid_camera -g 1006
groupadd aid_log -g 1007
groupadd aid_compass -g 1008
groupadd aid_mount -g 1009
groupadd aid_wifi -g 1010
groupadd aid_adb -g 1011
groupadd aid_install -g 1012
groupadd aid_media -g 1013
groupadd aid_dhcp -g 1014
groupadd aid_sdcard_rw -g 1015
groupadd aid_vpn -g 1016
groupadd aid_keystore -g 1017
groupadd aid_usb -g 1018
groupadd aid_drm -g 1019
groupadd aid_mdnsr -g 1020
groupadd aid_gps -g 1021
groupadd aid_media_rw -g 1023
groupadd aid_mtp -g 1024
groupadd aid_drmrpc -g 1026
groupadd aid_nfc -g 1027
groupadd aid_sdcard_r -g 1028
groupadd aid_clat -g 1029
groupadd aid_loop_radio -g 1030
groupadd aid_media_drm -g 1031
groupadd aid_package_info -g 1032
groupadd aid_sdcard_pics -g 1033
groupadd aid_sdcard_av -g 1034
groupadd aid_sdcard_all -g 1035
groupadd aid_logd -g 1036
groupadd aid_shared_relro -g 1037
groupadd aid_dbus -g 1038
groupadd aid_tlsdate -g 1039
groupadd aid_media_ex -g 1040
groupadd aid_audioserver -g 1041
groupadd aid_metrics_coll -g 1042
groupadd aid_metricsd -g 1043
groupadd aid_webserv -g 1044
groupadd aid_debuggerd -g 1045
groupadd aid_media_codec -g 1046
groupadd aid_cameraserver -g 1047
groupadd aid_firewall -g 1048
groupadd aid_trunks -g 1049
groupadd aid_nvram -g 1050
groupadd aid_dns -g 1051
groupadd aid_dns_tether -g 1052
groupadd aid_webview_zygote -g 1053
groupadd aid_vehicle_network -g 1054
groupadd aid_media_audio -g 1055
groupadd aid_media_video -g 1056
groupadd aid_media_image -g 1057
groupadd aid_tombstoned -g 1058
groupadd aid_media_obb -g 1059
groupadd aid_ese -g 1060
groupadd aid_ota_update -g 1061
groupadd aid_automotive_evs -g 1062
groupadd aid_lowpan -g 1063
groupadd aid_hsm -g 1064
groupadd aid_reserved_disk -g 1065
groupadd aid_statsd -g 1066
groupadd aid_incidentd -g 1067
groupadd aid_secure_element -g 1068
groupadd aid_lmkd -g 1069
groupadd aid_llkd -g 1070
groupadd aid_iorapd -g 1071
groupadd aid_gpu_service -g 1072
groupadd aid_network_stack -g 1073
groupadd aid_shell -g 2000
groupadd aid_cache -g 2001
groupadd aid_diag -g 2002
groupadd aid_oem_reserved_start -g 2900
groupadd aid_oem_reserved_end -g 2999
groupadd aid_net_bt_admin -g 3001
groupadd aid_net_bt -g 3002
groupadd aid_inet -g 3003
groupadd aid_net_raw -g 3004
groupadd aid_net_admin -g 3005
groupadd aid_net_bw_stats -g 3006
groupadd aid_net_bw_acct -g 3007
groupadd aid_readproc -g 3009
groupadd aid_wakelock -g 3010
groupadd aid_uhid -g 3011
groupadd aid_everybody -g 9997
groupadd aid_misc -g 9998
groupadd aid_nobody -g 9999
groupadd aid_app_start -g 10000
groupadd aid_app_end -g 19999
groupadd aid_cache_gid_start -g 20000
groupadd aid_cache_gid_end -g 29999
groupadd aid_ext_gid_start -g 30000
groupadd aid_ext_gid_end -g 39999
groupadd aid_ext_cache_gid_start -g 40000
groupadd aid_ext_cache_gid_end -g 49999
groupadd aid_shared_gid_start -g 50000
groupadd aid_shared_gid_end -g 59999
groupadd aid_isolated_start -g 99000
groupadd aid_isolated_end -g 99999
groupadd aid_user_offset -g 100000
# Fix root user's permissions.
usermod -a -G aid_system,aid_radio,aid_bluetooth,aid_graphics,aid_input,aid_audio,aid_camera,aid_log,aid_compass,aid_mount,aid_wifi,aid_adb,aid_install,aid_media,aid_dhcp,aid_sdcard_rw,aid_vpn,aid_keystore,aid_usb,aid_drm,aid_mdnsr,aid_gps,aid_media_rw,aid_mtp,aid_drmrpc,aid_nfc,aid_sdcard_r,aid_clat,aid_loop_radio,aid_media_drm,aid_package_info,aid_sdcard_pics,aid_sdcard_av,aid_sdcard_all,aid_logd,aid_shared_relro,aid_dbus,aid_tlsdate,aid_media_ex,aid_audioserver,aid_metrics_coll,aid_metricsd,aid_webserv,aid_debuggerd,aid_media_codec,aid_cameraserver,aid_firewall,aid_trunks,aid_nvram,aid_dns,aid_dns_tether,aid_webview_zygote,aid_vehicle_network,aid_media_audio,aid_media_video,aid_media_image,aid_tombstoned,aid_media_obb,aid_ese,aid_ota_update,aid_automotive_evs,aid_lowpan,aid_hsm,aid_reserved_disk,aid_statsd,aid_incidentd,aid_secure_element,aid_lmkd,aid_llkd,aid_iorapd,aid_gpu_service,aid_network_stack,aid_shell,aid_cache,aid_diag,aid_oem_reserved_start,aid_oem_reserved_end,aid_net_bt_admin,aid_net_bt,aid_inet,aid_net_raw,aid_net_admin,aid_net_bw_stats,aid_net_bw_acct,aid_readproc,aid_wakelock,aid_uhid,aid_everybody,aid_misc,aid_nobody,aid_app_start,aid_app_end,aid_cache_gid_start,aid_cache_gid_end,aid_ext_gid_start,aid_ext_gid_end,aid_ext_cache_gid_start,aid_ext_cache_gid_end,aid_shared_gid_start,aid_shared_gid_end,aid_isolated_start,aid_isolated_end,aid_user_offset root
# Fix apt.
usermod -g aid_inet _apt 2>/dev/null
# Fix gentoo emerge.
usermod -a -G aid_inet,aid_net_raw portage 2>/dev/null
# Fix permission of su.
# proot do not need this
chown root:root /bin/su
chmod 777 /bin/su
# Create mountpoints.
[ -e /dev ] || mkdir /dev
[ -e /proc ] || mkdir /proc
[ -e /sys ] || mkdir /sys
# Fix dns problem.
rm -f /etc/resolv.conf >>/dev/null 2>&1
echo nameserver 8.8.8.8 >>/etc/resolv.conf
echo nameserver 114.114.114.114 >>/etc/resolv.conf
printf "\033[33mFor alpine linux, there might be errors, do not care\n"
