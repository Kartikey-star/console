#!/bin/bash
cd ../testdata
curl --cert-type PEM   --cacert ../chartproxy/cacert.pem  --data-binary "@mychart-0.1.0.tgz" https://localhost:9553/api/charts
#push mariadb chart
<<<<<<< HEAD
curl --cert-type PEM   --cacert ../chartproxy/cacert.pem --data-binary "@mariadb-7.3.5.tgz" https://localhost:9553/api/charts
=======
curl --cert-type PEM   --cacert ../chartproxy/cacert.pem --data-binary "@mariadb-7.3.5.tgz" https://localhost:9553/api/charts
>>>>>>> 4205782c2b (Tls bug fix changes)
