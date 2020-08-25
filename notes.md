## Useful Links

[JSON response](https://medium.com/@vivek_syngh/http-response-in-golang-4ca1b3688d6)

## Stored Procedures

[Samples](https://www.mysqltutorial.org/stored-procedures-parameters.aspx)
```
DELIMITER $$

CREATE PROCEDURE SELECT userChain,email,isManager,defaultQuota from users WHERE userChain= (SELECT userChain FROM userDevices WHERE deviceIP=device_ip)
END
$$

DELIMITER ;
```