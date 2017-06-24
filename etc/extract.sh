#提取公钥和私钥密钥对,用于加密
openssl pkcs12 -in "$1" -nocerts -nodes -out "$1".key
#提取公钥证书，用于解密
openssl pkcs12 -in "$1" -clcerts -nokeys -out "$1".cer

