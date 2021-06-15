## fabric-explorer 接口

### 获取区块信息


  请求方式|请求URL
  ---|---
  GET| /v1/jzsg/block?page=1&count=10

**请求参数说明**

|参数名|类型|说明|
|:-----  |:-----|-----                           |
|page |integer   |页数  |
|count |integer   |  |

**返回示例**
``` 
{
    "errCode": 0,
    "errMsg": "",
    "data": {
        "total": 5,
        "list": [
            {
                "block_height": 4,
                "block_hash": "c0a5143e04ecc48763f0806a88139a0b986e5a1f7959843210a1a351b97edf08",
                "block_pre_hash": "a00e1b4470c7b6d458deb1ebaa06743c6babf1f8718b227da50034d67e6a53dd",
                "block_time_stamp": "2021-06-15T08:50:04.204357133Z",
                "block_confirm_time": 0,
                "block_tx_count": 1,
                "block_size": 7741,
                "CreatedAt": "2021-06-15T08:50:04.208427124Z"
            },
            {
                "block_height": 3,
                "block_hash": "1c33f565c88ea73c264a9bc5d2b3dd20e387a67623d90bea28d569102f6c66cf",
                "block_pre_hash": "fc2c59109183fa3e65599548b6749d2b37195639604a1dfa6716094b1de6c9d3",
                "block_time_stamp": "2021-06-15T08:50:04.199284398Z",
                "block_confirm_time": 0,
                "block_tx_count": 1,
                "block_size": 9957,
                "CreatedAt": "2021-06-15T08:50:04.203991431Z"
            },
            {
                "block_height": 2,
                "block_hash": "5ad01abfa4c43c794c7e0cb8a48b0326a96cf1124759c3eaf147c0f3669af4b4",
                "block_pre_hash": "f552f2d0aaba538557824362565ddbb6205c0eb95542e9f4ff71d9ff74ce36cc",
                "block_time_stamp": "2021-06-15T08:50:04.194590916Z",
                "block_confirm_time": 0,
                "block_tx_count": 1,
                "block_size": 7277,
                "CreatedAt": "2021-06-15T08:50:04.198896315Z"
            },
            {
                "block_height": 1,
                "block_hash": "28981775ad6a0244e3467b30a6511aef833f50243aa5a86541d366d303bf8d80",
                "block_pre_hash": "f8ca5c432052f248c87313b724f021e86f524c4309ce26fdf1d30c0ae1fd77b4",
                "block_time_stamp": "2021-06-15T08:50:04.18779821Z",
                "block_confirm_time": 0,
                "block_tx_count": 1,
                "block_size": 8039,
                "CreatedAt": "2021-06-15T08:50:04.194336709Z"
            },
            {
                "block_height": 0,
                "block_hash": "aa5788f2cd7db5b77172d0f758fa0366f972235f52903a02e77504fedf5ac73a",
                "block_pre_hash": "",
                "block_time_stamp": "2021-06-15T08:50:04.180028631Z",
                "block_confirm_time": 0,
                "block_tx_count": 1,
                "block_size": 15185,
                "CreatedAt": "2021-06-15T08:50:04.187695394Z"
            }
        ]
    }
}
```

**返回参数说明**

|参数名|类型|说明|
|:-----  |:-----|-----                           |
|total |integer   |总记录数  |
|block_height |integer   |区块高度  |
|block_hash |string   |区块hash|
|block_pre_hash |string   |上一区块hash|
|block_tx_count |integer   |区块内交易数量  |
|block_time_stamp |datetime   |区块生成时间戳  |

### 根据区块高度获取指定区块交易信息

**请求**

请求方式|请求URL
---|---
GET| /v1/jzsg/block/:height 

**请求参数说明**

|参数名|类型|说明|
|:-----  |:-----|-----                           |
|height |integer   |区块高度  |

**返回示例**
``` 
{
    "errCode": 0,
    "errMsg": "",
    "data": [
        {
            "tx_id": "fba9ad52abfecbb89d55cb8c2b86f59f5e244150286a4c9ff9413c68f34e450f",
            "tx_hash": "3a5adb932e6acf2de8564db8612f4a2ff73688424ee70c2f00def41cbdc4546d",
            "tx_index": 0,
            "tx_block_height": 4,
            "tx_block_hash": "c0a5143e04ecc48763f0806a88139a0b986e5a1f7959843210a1a351b97edf08",
            "tx_start_time": "2021-06-15T08:32:18.065251683Z",
            "tx_finish_time": "2021-06-15T08:50:04.204357133Z",
            "tx_confirm_time": 1066,
            "tx_type": "ENDORSER_TRANSACTION",
            "tx_valid_status": true,
            "tx_size": 4229,
            "last_config_block_num": 0,
            "function": "updateEdu",
            "args": "[\"{\\\"docType\\\":\\\"\\\",\\\"EntityID\\\":\\\"10002\\\",\\\"EnterpriseName\\\":\\\"南京可信区块链与算法经济研究院有限公司\\\",\\\"EnterpriseType\\\":\\\"有限责任公司\\\",\\\"EnterpriseCorporation\\\":\\\"石宁\\\",\\\"EnterpriseCode\\\":\\\"91320191MA2006WN9B\\\",\\\"EnterpriseAddress\\\":\\\"中国(江苏)自由贸易试验区南京片区研创园团结路99号孵鹰大厦1515室 \\\",\\\"RegisteredCapital\\\":\\\"1000万\\\",\\\"EstablishmentDate\\\":\\\"2014年7月\\\",\\\"BusinessTermFrom\\\":\\\"中国人民大学\\\",\\\"BusinessTermTo\\\":\\\"行政管理\\\",\\\"RegistrationAuthority\\\":\\\"南京市江北新区管理委员会行政审批局\\\",\\\"ApprovalDate\\\":\\\"四年\\\",\\\"RegistrationStatus\\\":\\\"普通全日制\\\",\\\"BusinessScope\\\":\\\"本科\\\",\\\"Historys\\\":null}\",\"eventModifyEdu\"]",
            "proposal_response": "CMgBGhLkv6Hmga/mm7TmlrDmiJDlip8=",
            "chaincode_id": "educc",
            "chaincode_name": "",
            "chaincode_version": "",
            "chaincode_type": "UNDEFINED",
            "chaincode_path": "",
            "endorse_policy": "",
            "init_param": "",
            "signature": "MEUCIQDHb+Z6PYjB0cHNcJgRw/L8gN/piv8hz5H0rZJ8JRulywIgJZtHdV9wSv3ZMTVo+dk8oSXiqhNnALOz2d3uA849FA8=",
            "endorse_signature": "MEQCIHpSc5aCmvFD6sJaIIj8Z7iDcqTS/VuiIfBPCXXtxWqsAiBaGE8yFGRrLeX7Lftao9q34Wbvj+4ni69er8YPG6dn0g==",
            "anchor_info": "",
            "policy": "",
            "msp_id": "org1.kevin.kongyixueyuan.com",
            "org_name": "",
            "ca_cert": "",
            "tx_action_cert": "{\"Raw\":\"MIICTzCCAfagAwIBAgIRANbzg+NSqCtBRbuY3qJcGh8wCgYIKoZIzj0EAwIwgYsxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1TYW4gRnJhbmNpc2NvMSUwIwYDVQQKExxvcmcxLmtldmluLmtvbmd5aXh1ZXl1YW4uY29tMSgwJgYDVQQDEx9jYS5vcmcxLmtldmluLmtvbmd5aXh1ZXl1YW4uY29tMB4XDTE4MDkyNjA1MzAxMFoXDTI4MDkyMzA1MzAxMFoweDELMAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBGcmFuY2lzY28xDzANBgNVBAsTBmNsaWVudDErMCkGA1UEAwwiVXNlcjFAb3JnMS5rZXZpbi5rb25neWl4dWV5dWFuLmNvbTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABD0VBT8P9fLqMtAeO9xbOOCp1R1Vp0xRKGd8VY9lq04nldbrw1ymc8cAWdA/bggp/Ejs9x07YPsrE9y54WPv/dajTTBLMA4GA1UdDwEB/wQEAwIHgDAMBgNVHRMBAf8EAjAAMCsGA1UdIwQkMCKAIHJ+ae1KAaIEzVO/SpfCwcuUdBlQT4KFH2rlY8PJbeo6MAoGCCqGSM49BAMCA0cAMEQCIBQxT2c6U7UbNzoVtrEMzGupz3l4y8/LgenFHElcnZaPAiBvjfLXTOoLnupMbvYpxWJH6sk1nAaTIMNi3HHxXzAyDA==\",\"RawTBSCertificate\":\"MIIB9qADAgECAhEA1vOD41KoK0FFu5jeolwaHzAKBggqhkjOPQQDAjCBizELMAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBGcmFuY2lzY28xJTAjBgNVBAoTHG9yZzEua2V2aW4ua29uZ3lpeHVleXVhbi5jb20xKDAmBgNVBAMTH2NhLm9yZzEua2V2aW4ua29uZ3lpeHVleXVhbi5jb20wHhcNMTgwOTI2MDUzMDEwWhcNMjgwOTIzMDUzMDEwWjB4MQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZyYW5jaXNjbzEPMA0GA1UECxMGY2xpZW50MSswKQYDVQQDDCJVc2VyMUBvcmcxLmtldmluLmtvbmd5aXh1ZXl1YW4uY29tMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEPRUFPw/18uoy0B473Fs44KnVHVWnTFEoZ3xVj2WrTieV1uvDXKZzxwBZ0D9uCCn8SOz3HTtg+ysT3LnhY+/91qNNMEswDgYDVR0PAQH/BAQDAgeAMAwGA1UdEwEB/wQCMAAwKwYDVR0jBCQwIoAgcn5p7UoBogTNU79Kl8LBy5R0GVBPgoUfauVjw8lt6jo=\",\"RawSubjectPublicKeyInfo\":\"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEPRUFPw/18uoy0B473Fs44KnVHVWnTFEoZ3xVj2WrTieV1uvDXKZzxwBZ0D9uCCn8SOz3HTtg+ysT3LnhY+/91g==\",\"RawSubject\":\"MHgxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1TYW4gRnJhbmNpc2NvMQ8wDQYDVQQLEwZjbGllbnQxKzApBgNVBAMMIlVzZXIxQG9yZzEua2V2aW4ua29uZ3lpeHVleXVhbi5jb20=\",\"RawIssuer\":\"MIGLMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZyYW5jaXNjbzElMCMGA1UEChMcb3JnMS5rZXZpbi5rb25neWl4dWV5dWFuLmNvbTEoMCYGA1UEAxMfY2Eub3JnMS5rZXZpbi5rb25neWl4dWV5dWFuLmNvbQ==\",\"Signature\":\"MEQCIBQxT2c6U7UbNzoVtrEMzGupz3l4y8/LgenFHElcnZaPAiBvjfLXTOoLnupMbvYpxWJH6sk1nAaTIMNi3HHxXzAyDA==\",\"SignatureAlgorithm\":10,\"PublicKeyAlgorithm\":3,\"PublicKey\":{\"Curve\":{\"P\":115792089210356248762697446949407573530086143415290314195533631308867097853951,\"N\":115792089210356248762697446949407573529996955224135760342422259061068512044369,\"B\":41058363725152142129326129780047268409114441015993725554835256314039467401291,\"Gx\":48439561293906451759052585252797914202762949526041747995844080717082404635286,\"Gy\":36134250956749795798585127919587881956611106672985015071877198253568414405109,\"BitSize\":256,\"Name\":\"P-256\"},\"X\":27628223760828831623261995870142172107280465733072272330475334640380863532583,\"Y\":67774346888107567310683961435169519899600574688230971389732616376878552317398},\"Version\":3,\"SerialNumber\":285719194240628003630845251746273171999,\"Issuer\":{\"Country\":[\"US\"],\"Organization\":[\"org1.kevin.kongyixueyuan.com\"],\"OrganizationalUnit\":null,\"Locality\":[\"San Francisco\"],\"Province\":[\"California\"],\"StreetAddress\":null,\"PostalCode\":null,\"SerialNumber\":\"\",\"CommonName\":\"ca.org1.kevin.kongyixueyuan.com\",\"Names\":[{\"Type\":[2,5,4,6],\"Value\":\"US\"},{\"Type\":[2,5,4,8],\"Value\":\"California\"},{\"Type\":[2,5,4,7],\"Value\":\"San Francisco\"},{\"Type\":[2,5,4,10],\"Value\":\"org1.kevin.kongyixueyuan.com\"},{\"Type\":[2,5,4,3],\"Value\":\"ca.org1.kevin.kongyixueyuan.com\"}],\"ExtraNames\":null},\"Subject\":{\"Country\":[\"US\"],\"Organization\":null,\"OrganizationalUnit\":[\"client\"],\"Locality\":[\"San Francisco\"],\"Province\":[\"California\"],\"StreetAddress\":null,\"PostalCode\":null,\"SerialNumber\":\"\",\"CommonName\":\"User1@org1.kevin.kongyixueyuan.com\",\"Names\":[{\"Type\":[2,5,4,6],\"Value\":\"US\"},{\"Type\":[2,5,4,8],\"Value\":\"California\"},{\"Type\":[2,5,4,7],\"Value\":\"San Francisco\"},{\"Type\":[2,5,4,11],\"Value\":\"client\"},{\"Type\":[2,5,4,3],\"Value\":\"User1@org1.kevin.kongyixueyuan.com\"}],\"ExtraNames\":null},\"NotBefore\":\"2018-09-26T05:30:10Z\",\"NotAfter\":\"2028-09-23T05:30:10Z\",\"KeyUsage\":1,\"Extensions\":[{\"Id\":[2,5,29,15],\"Critical\":true,\"Value\":\"AwIHgA==\"},{\"Id\":[2,5,29,19],\"Critical\":true,\"Value\":\"MAA=\"},{\"Id\":[2,5,29,35],\"Critical\":false,\"Value\":\"MCKAIHJ+ae1KAaIEzVO/SpfCwcuUdBlQT4KFH2rlY8PJbeo6\"}],\"ExtraExtensions\":null,\"UnhandledCriticalExtensions\":null,\"ExtKeyUsage\":null,\"UnknownExtKeyUsage\":null,\"BasicConstraintsValid\":true,\"IsCA\":false,\"MaxPathLen\":-1,\"MaxPathLenZero\":false,\"SubjectKeyId\":null,\"AuthorityKeyId\":\"cn5p7UoBogTNU79Kl8LBy5R0GVBPgoUfauVjw8lt6jo=\",\"OCSPServer\":null,\"IssuingCertificateURL\":null,\"DNSNames\":null,\"EmailAddresses\":null,\"IPAddresses\":null,\"URIs\":null,\"PermittedDNSDomainsCritical\":false,\"PermittedDNSDomains\":null,\"ExcludedDNSDomains\":null,\"PermittedIPRanges\":null,\"ExcludedIPRanges\":null,\"PermittedEmailAddresses\":null,\"ExcludedEmailAddresses\":null,\"PermittedURIDomains\":null,\"ExcludedURIDomains\":null,\"CRLDistributionPoints\":null,\"PolicyIdentifiers\":null}",
            "endorse_certs": "[{\"Raw\":\"MIICTTCCAfOgAwIBAgIQDSShRDCopD+wCdMqq8klyzAKBggqhkjOPQQDAjCBizELMAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBGcmFuY2lzY28xJTAjBgNVBAoTHG9yZzEua2V2aW4ua29uZ3lpeHVleXVhbi5jb20xKDAmBgNVBAMTH2NhLm9yZzEua2V2aW4ua29uZ3lpeHVleXVhbi5jb20wHhcNMTgwOTI2MDUzMDEwWhcNMjgwOTIzMDUzMDEwWjB2MQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZyYW5jaXNjbzENMAsGA1UECxMEcGVlcjErMCkGA1UEAxMicGVlcjAub3JnMS5rZXZpbi5rb25neWl4dWV5dWFuLmNvbTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABAo5bSlqhqjM8Xl4hb0cj+a6dpU/XOvTMvahkTSzK9Bga6NzuAkC69MUExsfrL1zzZzm9bLAiRnUavKj6CxvA52jTTBLMA4GA1UdDwEB/wQEAwIHgDAMBgNVHRMBAf8EAjAAMCsGA1UdIwQkMCKAIHJ+ae1KAaIEzVO/SpfCwcuUdBlQT4KFH2rlY8PJbeo6MAoGCCqGSM49BAMCA0gAMEUCIQCLKXhkPe/ZUPtnwOTe/sU92KLYDQBo0LVXI88m6T2fnAIgKkUEK2AoxRoJfrXAglRgVV6VwEOLph0s4K4yEJJCejk=\",\"RawTBSCertificate\":\"MIIB86ADAgECAhANJKFEMKikP7AJ0yqrySXLMAoGCCqGSM49BAMCMIGLMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZyYW5jaXNjbzElMCMGA1UEChMcb3JnMS5rZXZpbi5rb25neWl4dWV5dWFuLmNvbTEoMCYGA1UEAxMfY2Eub3JnMS5rZXZpbi5rb25neWl4dWV5dWFuLmNvbTAeFw0xODA5MjYwNTMwMTBaFw0yODA5MjMwNTMwMTBaMHYxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1TYW4gRnJhbmNpc2NvMQ0wCwYDVQQLEwRwZWVyMSswKQYDVQQDEyJwZWVyMC5vcmcxLmtldmluLmtvbmd5aXh1ZXl1YW4uY29tMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAECjltKWqGqMzxeXiFvRyP5rp2lT9c69My9qGRNLMr0GBro3O4CQLr0xQTGx+svXPNnOb1ssCJGdRq8qPoLG8DnaNNMEswDgYDVR0PAQH/BAQDAgeAMAwGA1UdEwEB/wQCMAAwKwYDVR0jBCQwIoAgcn5p7UoBogTNU79Kl8LBy5R0GVBPgoUfauVjw8lt6jo=\",\"RawSubjectPublicKeyInfo\":\"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAECjltKWqGqMzxeXiFvRyP5rp2lT9c69My9qGRNLMr0GBro3O4CQLr0xQTGx+svXPNnOb1ssCJGdRq8qPoLG8DnQ==\",\"RawSubject\":\"MHYxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1TYW4gRnJhbmNpc2NvMQ0wCwYDVQQLEwRwZWVyMSswKQYDVQQDEyJwZWVyMC5vcmcxLmtldmluLmtvbmd5aXh1ZXl1YW4uY29t\",\"RawIssuer\":\"MIGLMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZyYW5jaXNjbzElMCMGA1UEChMcb3JnMS5rZXZpbi5rb25neWl4dWV5dWFuLmNvbTEoMCYGA1UEAxMfY2Eub3JnMS5rZXZpbi5rb25neWl4dWV5dWFuLmNvbQ==\",\"Signature\":\"MEUCIQCLKXhkPe/ZUPtnwOTe/sU92KLYDQBo0LVXI88m6T2fnAIgKkUEK2AoxRoJfrXAglRgVV6VwEOLph0s4K4yEJJCejk=\",\"SignatureAlgorithm\":10,\"PublicKeyAlgorithm\":3,\"PublicKey\":{\"Curve\":{\"P\":115792089210356248762697446949407573530086143415290314195533631308867097853951,\"N\":115792089210356248762697446949407573529996955224135760342422259061068512044369,\"B\":41058363725152142129326129780047268409114441015993725554835256314039467401291,\"Gx\":48439561293906451759052585252797914202762949526041747995844080717082404635286,\"Gy\":36134250956749795798585127919587881956611106672985015071877198253568414405109,\"BitSize\":256,\"Name\":\"P-256\"},\"X\":4624592175453143886727000646717812298506798626301525926515277097285010051168,\"Y\":48686269532377460204917810720956692844809457151427503899766330920463240201117},\"Version\":3,\"SerialNumber\":17470157502631554683342591791288886731,\"Issuer\":{\"Country\":[\"US\"],\"Organization\":[\"org1.kevin.kongyixueyuan.com\"],\"OrganizationalUnit\":null,\"Locality\":[\"San Francisco\"],\"Province\":[\"California\"],\"StreetAddress\":null,\"PostalCode\":null,\"SerialNumber\":\"\",\"CommonName\":\"ca.org1.kevin.kongyixueyuan.com\",\"Names\":[{\"Type\":[2,5,4,6],\"Value\":\"US\"},{\"Type\":[2,5,4,8],\"Value\":\"California\"},{\"Type\":[2,5,4,7],\"Value\":\"San Francisco\"},{\"Type\":[2,5,4,10],\"Value\":\"org1.kevin.kongyixueyuan.com\"},{\"Type\":[2,5,4,3],\"Value\":\"ca.org1.kevin.kongyixueyuan.com\"}],\"ExtraNames\":null},\"Subject\":{\"Country\":[\"US\"],\"Organization\":null,\"OrganizationalUnit\":[\"peer\"],\"Locality\":[\"San Francisco\"],\"Province\":[\"California\"],\"StreetAddress\":null,\"PostalCode\":null,\"SerialNumber\":\"\",\"CommonName\":\"peer0.org1.kevin.kongyixueyuan.com\",\"Names\":[{\"Type\":[2,5,4,6],\"Value\":\"US\"},{\"Type\":[2,5,4,8],\"Value\":\"California\"},{\"Type\":[2,5,4,7],\"Value\":\"San Francisco\"},{\"Type\":[2,5,4,11],\"Value\":\"peer\"},{\"Type\":[2,5,4,3],\"Value\":\"peer0.org1.kevin.kongyixueyuan.com\"}],\"ExtraNames\":null},\"NotBefore\":\"2018-09-26T05:30:10Z\",\"NotAfter\":\"2028-09-23T05:30:10Z\",\"KeyUsage\":1,\"Extensions\":[{\"Id\":[2,5,29,15],\"Critical\":true,\"Value\":\"AwIHgA==\"},{\"Id\":[2,5,29,19],\"Critical\":true,\"Value\":\"MAA=\"},{\"Id\":[2,5,29,35],\"Critical\":false,\"Value\":\"MCKAIHJ+ae1KAaIEzVO/SpfCwcuUdBlQT4KFH2rlY8PJbeo6\"}],\"ExtraExtensions\":null,\"UnhandledCriticalExtensions\":null,\"ExtKeyUsage\":null,\"UnknownExtKeyUsage\":null,\"BasicConstraintsValid\":true,\"IsCA\":false,\"MaxPathLen\":-1,\"MaxPathLenZero\":false,\"SubjectKeyId\":null,\"AuthorityKeyId\":\"cn5p7UoBogTNU79Kl8LBy5R0GVBPgoUfauVjw8lt6jo=\",\"OCSPServer\":null,\"IssuingCertificateURL\":null,\"DNSNames\":null,\"EmailAddresses\":null,\"IPAddresses\":null,\"URIs\":null,\"PermittedDNSDomainsCritical\":false,\"PermittedDNSDomains\":null,\"ExcludedDNSDomains\":null,\"PermittedIPRanges\":null,\"ExcludedIPRanges\":null,\"PermittedEmailAddresses\":null,\"ExcludedEmailAddresses\":null,\"PermittedURIDomains\":null,\"ExcludedURIDomains\":null,\"CRLDistributionPoints\":null,\"PolicyIdentifiers\":null}]",
            "block_creator_cert": "{\"Raw\":\"MIICMDCCAdegAwIBAgIQUEaAxYYwDp6/JPn2sCqECzAKBggqhkjOPQQDAjCBgTELMAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBGcmFuY2lzY28xIDAeBgNVBAoTF2tldmluLmtvbmd5aXh1ZXl1YW4uY29tMSMwIQYDVQQDExpjYS5rZXZpbi5rb25neWl4dWV5dWFuLmNvbTAeFw0xODA5MjYwNTMwMTBaFw0yODA5MjMwNTMwMTBaMGQxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1TYW4gRnJhbmNpc2NvMSgwJgYDVQQDEx9vcmRlcmVyLmtldmluLmtvbmd5aXh1ZXl1YW4uY29tMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE+qV2XDyxyd0WyZSUqtyZUTi8A2pMgL2XWXUqU+AvcN7/U9rfNVx/gfdLybbXx7XOFmU7AWwmhCSb9cHu/K1SJaNNMEswDgYDVR0PAQH/BAQDAgeAMAwGA1UdEwEB/wQCMAAwKwYDVR0jBCQwIoAgAdW8HVVSDcVI1pVMKY8F6MvrxLfkqlybOQ0AbnBdLrEwCgYIKoZIzj0EAwIDRwAwRAIgOs3qP4F2I4Je0NqKClHxVhh3hW8fim10S6CRB3Ab+fcCIDcHupW1CckrNdghPSI8BuOEMIY9QR4vXoG8hI9+6wFs\",\"RawTBSCertificate\":\"MIIB16ADAgECAhBQRoDFhjAOnr8k+fawKoQLMAoGCCqGSM49BAMCMIGBMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZyYW5jaXNjbzEgMB4GA1UEChMXa2V2aW4ua29uZ3lpeHVleXVhbi5jb20xIzAhBgNVBAMTGmNhLmtldmluLmtvbmd5aXh1ZXl1YW4uY29tMB4XDTE4MDkyNjA1MzAxMFoXDTI4MDkyMzA1MzAxMFowZDELMAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBGcmFuY2lzY28xKDAmBgNVBAMTH29yZGVyZXIua2V2aW4ua29uZ3lpeHVleXVhbi5jb20wWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAT6pXZcPLHJ3RbJlJSq3JlROLwDakyAvZdZdSpT4C9w3v9T2t81XH+B90vJttfHtc4WZTsBbCaEJJv1we78rVIlo00wSzAOBgNVHQ8BAf8EBAMCB4AwDAYDVR0TAQH/BAIwADArBgNVHSMEJDAigCAB1bwdVVINxUjWlUwpjwXoy+vEt+SqXJs5DQBucF0usQ==\",\"RawSubjectPublicKeyInfo\":\"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE+qV2XDyxyd0WyZSUqtyZUTi8A2pMgL2XWXUqU+AvcN7/U9rfNVx/gfdLybbXx7XOFmU7AWwmhCSb9cHu/K1SJQ==\",\"RawSubject\":\"MGQxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1TYW4gRnJhbmNpc2NvMSgwJgYDVQQDEx9vcmRlcmVyLmtldmluLmtvbmd5aXh1ZXl1YW4uY29t\",\"RawIssuer\":\"MIGBMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZyYW5jaXNjbzEgMB4GA1UEChMXa2V2aW4ua29uZ3lpeHVleXVhbi5jb20xIzAhBgNVBAMTGmNhLmtldmluLmtvbmd5aXh1ZXl1YW4uY29t\",\"Signature\":\"MEQCIDrN6j+BdiOCXtDaigpR8VYYd4VvH4ptdEugkQdwG/n3AiA3B7qVtQnJKzXYIT0iPAbjhDCGPUEeL16BvISPfusBbA==\",\"SignatureAlgorithm\":10,\"PublicKeyAlgorithm\":3,\"PublicKey\":{\"Curve\":{\"P\":115792089210356248762697446949407573530086143415290314195533631308867097853951,\"N\":115792089210356248762697446949407573529996955224135760342422259061068512044369,\"B\":41058363725152142129326129780047268409114441015993725554835256314039467401291,\"Gx\":48439561293906451759052585252797914202762949526041747995844080717082404635286,\"Gy\":36134250956749795798585127919587881956611106672985015071877198253568414405109,\"BitSize\":256,\"Name\":\"P-256\"},\"X\":113370558804280920631218082688737421128268750658465427086772434739153156927710,\"Y\":115487935293500845012244754339790992686287974775839037023094545545887625335333},\"Version\":3,\"SerialNumber\":106704312240797079287081944370763170827,\"Issuer\":{\"Country\":[\"US\"],\"Organization\":[\"kevin.kongyixueyuan.com\"],\"OrganizationalUnit\":null,\"Locality\":[\"San Francisco\"],\"Province\":[\"California\"],\"StreetAddress\":null,\"PostalCode\":null,\"SerialNumber\":\"\",\"CommonName\":\"ca.kevin.kongyixueyuan.com\",\"Names\":[{\"Type\":[2,5,4,6],\"Value\":\"US\"},{\"Type\":[2,5,4,8],\"Value\":\"California\"},{\"Type\":[2,5,4,7],\"Value\":\"San Francisco\"},{\"Type\":[2,5,4,10],\"Value\":\"kevin.kongyixueyuan.com\"},{\"Type\":[2,5,4,3],\"Value\":\"ca.kevin.kongyixueyuan.com\"}],\"ExtraNames\":null},\"Subject\":{\"Country\":[\"US\"],\"Organization\":null,\"OrganizationalUnit\":null,\"Locality\":[\"San Francisco\"],\"Province\":[\"California\"],\"StreetAddress\":null,\"PostalCode\":null,\"SerialNumber\":\"\",\"CommonName\":\"orderer.kevin.kongyixueyuan.com\",\"Names\":[{\"Type\":[2,5,4,6],\"Value\":\"US\"},{\"Type\":[2,5,4,8],\"Value\":\"California\"},{\"Type\":[2,5,4,7],\"Value\":\"San Francisco\"},{\"Type\":[2,5,4,3],\"Value\":\"orderer.kevin.kongyixueyuan.com\"}],\"ExtraNames\":null},\"NotBefore\":\"2018-09-26T05:30:10Z\",\"NotAfter\":\"2028-09-23T05:30:10Z\",\"KeyUsage\":1,\"Extensions\":[{\"Id\":[2,5,29,15],\"Critical\":true,\"Value\":\"AwIHgA==\"},{\"Id\":[2,5,29,19],\"Critical\":true,\"Value\":\"MAA=\"},{\"Id\":[2,5,29,35],\"Critical\":false,\"Value\":\"MCKAIAHVvB1VUg3FSNaVTCmPBejL68S35KpcmzkNAG5wXS6x\"}],\"ExtraExtensions\":null,\"UnhandledCriticalExtensions\":null,\"ExtKeyUsage\":null,\"UnknownExtKeyUsage\":null,\"BasicConstraintsValid\":true,\"IsCA\":false,\"MaxPathLen\":-1,\"MaxPathLenZero\":false,\"SubjectKeyId\":null,\"AuthorityKeyId\":\"AdW8HVVSDcVI1pVMKY8F6MvrxLfkqlybOQ0AbnBdLrE=\",\"OCSPServer\":null,\"IssuingCertificateURL\":null,\"DNSNames\":null,\"EmailAddresses\":null,\"IPAddresses\":null,\"URIs\":null,\"PermittedDNSDomainsCritical\":false,\"PermittedDNSDomains\":null,\"ExcludedDNSDomains\":null,\"PermittedIPRanges\":null,\"ExcludedIPRanges\":null,\"PermittedEmailAddresses\":null,\"ExcludedEmailAddresses\":null,\"PermittedURIDomains\":null,\"ExcludedURIDomains\":null,\"CRLDistributionPoints\":null,\"PolicyIdentifiers\":null}",
            "created_at": "2021-06-15T08:50:04.210444367Z"
        }
    ]
}
```

**返回参数说明**

|参数名|类型|说明|
|:-----  |:-----|-----                           |
|tx_hash |string   |tx hash  |
|msp_id |string   |交易发起人签名（组织级别）  |
|tx_block_hash |string   |区块hash|
|tx_valid_status |bool   |交易是否成功|
|tx_block_hash |integer   |交易所在区块高度  |
|args |[]string   |合约请求参数  |
|function |string   |合约请求方法  |

## postman 示例

[explorer.postman_collection.json](./explorer.postman_collection.json)
