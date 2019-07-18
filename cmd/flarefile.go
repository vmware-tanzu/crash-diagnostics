package cmd

const flarefile = `FROM local
WORKDIR /tmp/flareout
COPY /var/log/messages
COPY /var/log/syslog
COPY /var/log/system.log
`
