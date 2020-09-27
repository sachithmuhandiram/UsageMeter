## Useful Links

[JSON response](https://medium.com/@vivek_syngh/http-response-in-golang-4ca1b3688d6)

## Creating user

`CREATE USER 'usagemeter'@'mysqlcontainer' IDENTIFIED BY '7890';`
`GRANT ALL ON usagemeter.* TO 'usagemeter'@'mysqlcontainer';`
`CREATE USER 'usagemeter'@'%' IDENTIFIED BY '7890';`
`GRANT ALL PRIVILEGES ON *.* TO 'usagemeter'@'%' WITH GRANT OPTION;`


## Stored Procedures

[Samples](https://www.mysqltutorial.org/stored-procedures-parameters.aspx)
```
DELIMITER $$

CREATE PROCEDURE GetUserDetails( IN  device_ip VARCHAR(16) ) BEGIN SELECT userChain,email,isManager,defaultQuota from users WHERE userChain= (SELECT userChain FROM userDevices WHERE deviceIP=device_ip); END
$$

DELIMITER ;
```

```
CREATE PROCEDURE GetManagersEmail(IN user_chain varchar(50)) begin select email from users where userChain IN (select managerChain from userManagers where userChain=user_chain); END$$
```

Int to String convert

```
import (
    "strconv"
    "fmt"
)

func main() {
    t := strconv.Itoa(123)
    fmt.Println(t)
}
```

String to Int

`strconv.Atoi(string)`

### Links

[Github Markdown](https://towardsdatascience.com/build-a-stunning-readme-for-your-github-profile-9b80434fe5d7)