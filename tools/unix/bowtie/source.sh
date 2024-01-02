cd {{toolDir}} #cd into the toolDir, basically where the tool was cloned into
./configure --prefix={{toolDir}}
make
make install