Sample shell script:

sudo /home/rad/passenger_worker_killer -mode live -limit 400 -passenger_version 5 -passenger_memory_stat_path /usr/sbin >> /home/rad/passenger_worker_killer.log

Sample crontab entry (run every 10 minutes):

*/10 * * * * /home/rad/passenger_worker_killer.sh

Run touch /home/rad/passenger_worker_killer.log first so you don't get the "file does not exist" error
