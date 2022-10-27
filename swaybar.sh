#!/bin/sh

## program state
#
previous_network_time=$(date "+%s.%N")
previous_network_tx_bytes=$(cat /sys/class/net/enp3s0/statistics/tx_bytes)
previous_network_rx_bytes=$(cat /sys/class/net/enp3s0/statistics/rx_bytes)

## header
#
# man -s7 swaybar-protocol
#
echo '{ "version": 1 }'
echo -n '['

while /bin/true; do
	## date
	#
	# Input:
	#   Wed Oct 26 01:33:28 PM PDT 2022
	#
	# Output:
	#   Wed Oct 26 13:33
	date_formatted="$(date "+%a %b %d %H:%M")"

	## uptime
	#
	# Input:
	#   13:32:15 up  6:27,  1 user,  load average: 0.79, 0.69, 0.61
	#
	# Output:
	#   6:27
	#
	uptime_formatted="U: $(uptime | awk '{ print $3 }' | tr -d ',')"

	## cput temperature
	#
	# Input:
	#   25700
	#
	# Output:
	#   25.7 C
	#
	temperature_value=$(cat /sys/class/hwmon/hwmon3/temp1_input)
	temperature_celcius=$(echo "$temperature_value/1000" | bc -l)
	printf -v temperature_formatted "CPU: %0.1f C" $temperature_celcius

	## memory
	#
	# Input:
	#                  total        used        free      shared  buff/cache   available
	#   Mem:              31           2          24           0           4          27
	#   Swap:              7           0           7
	#
	# Output:
	#   T: 31 Gb U: 2 Gb F: 24 Gb
	memory_total=$(free -g | grep "Mem" | awk '{ print $2 }')
	memory_used=$(free -g | grep "Mem" | awk '{ print $3 }')
	memory_free=$(free -g | grep "Mem" | awk '{ print $4 }')
	printf -v memory_formatted "T: %d Gb U: %d Gb F: %d Gb" $memory_total $memory_used $memory_free

	## network stats
	#
	network_time=$(date "+%s.%N")
	network_tx_bytes=$(cat /sys/class/net/enp3s0/statistics/tx_bytes)
	network_rx_bytes=$(cat /sys/class/net/enp3s0/statistics/rx_bytes)

	network_tx_mbps=$(echo "($network_tx_bytes-$previous_network_tx_bytes)/($network_time-$previous_network_time)*8/1000000" | bc -l)
	network_rx_mbps=$(echo "($network_rx_bytes-$previous_network_rx_bytes)/($network_time-$previous_network_time)*8/1000000" | bc -l)

	previous_network_time=$network_time
	previous_network_tx_bytes=$network_tx_bytes
	previous_network_rx_bytes=$network_rx_bytes

	printf -v network_formatted "Rx: %0.2f Tx: %0.2f (Mbit/s)" $network_tx_mbps $network_rx_mbps

	## network info
	#
	network_ip_address=$(ip address | grep "enp3s0" | grep "inet" | awk '{ print $2 }')
	network_speed=$(cat /sys/class/net/enp3s0/speed)
	printf -v network_info_formatted "E: %s %d (Mbit/s)" $network_ip_address $network_speed

	# generate the output for swaybar
	echo "[ { 'full_text': '$network_formatted' },{ 'full_text': '$network_info_formatted' },{ 'full_text': '$memory_formatted' },{ 'full_text': '$temperature_formatted' },{ 'full_text': '$uptime_formatted' },{ 'full_text': '$date_formatted' } ],"

	sleep 1
done
