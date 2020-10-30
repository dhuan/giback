cat /giback_test/tmp/id_rsa.pub >> /home/git/.ssh/authorized_keys

mkdir /srv/git/test.git
cd /srv/git/test.git
git init --bare
chown -R git /srv/git

service ssh start

tail -f /dev/null
