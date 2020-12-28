set -x
set -e

cat /giback/test/tmp/id_rsa.pub >> /home/git/.ssh/authorized_keys

rm -rf /srv/git/test*

mkdir /srv/git/test.git
cd /srv/git/test.git
git init --bare

mkdir /srv/git/test2.git
cd /srv/git/test2.git
git init --bare

chown -R git /srv/git

service ssh start

cd /giback

go clean -testcache
GIT_SSH_COMMAND="ssh -o StrictHostKeyChecking=no -i $(pwd)/test/tmp/id_rsa" go test -v ./test
