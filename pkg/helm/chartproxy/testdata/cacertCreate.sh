#!/bin/bash
<<<<<<< HEAD
echo quit | openssl s_client -showcerts -servername localhost -connect localhost:9553 > cacert.pem
=======
echo quit | openssl s_client -showcerts -servername localhost -connect localhost:9553 > cacert.pem
>>>>>>> 4205782c2b (Tls bug fix changes)
