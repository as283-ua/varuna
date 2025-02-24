#!/bin/bash
spki=$("$(dirname "$0")/echo_spki.sh" $1)

google-chrome --enable-quic --allow-insecure-localhost --origin-to-force-quic-on=127.0.0.1:4433 --user-data-dir=/tmp/temp-chrome --ignore-certificate-errors-spki-list="$spki"