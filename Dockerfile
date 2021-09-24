FROM photon:4.0
RUN tdnf install -y vim >> /dev/null
COPY ldaputils /usr/local/bin