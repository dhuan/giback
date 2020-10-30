cat /giback_test/tmp/id_rsa.pub >> /home/git/.ssh/authorized_keys

mkdir /srv/git/test.git
cd /srv/git/test.git
git init --bare

mkdir /srv/git/test2.git
cd /srv/git/test2.git
git init --bare

chown -R git /srv/git

service ssh start

tail -f /dev/null
