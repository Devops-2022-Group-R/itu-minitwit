# Taken from https://gist.github.com/avborup/580bc0b003a702d1a84235890acd5fd7
FROM blang/latex:ctanfull

RUN apt-get update
RUN apt-get install python3 python3-pip -y;

RUN ln -s /usr/bin/python3 /usr/bin/python; \
    ln -s /usr/bin/pip3 /usr/bin/pip; \
    python --version; \
    pip --version; \
    pip install Pygments
