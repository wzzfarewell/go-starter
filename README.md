# go-starter

This is a code generator of Go. Using this tool you can generate a project by database configuration, it will generate a project architecture based on ```bxcodec/go-clean-arch```. The generated project contains packages like domain, repository, service and delivery module.

## 1. Quick Start

1. Using this command to install go-starter.
```shell
go install github.com/wzzfarewell/go-starter@latest
```

2. Write a configuration file for generator by using a toml or a yaml format. This file should contains your project's path, project's module name and your database information. For example:
```toml
project-path = "./go-starter-example"
module-name = "github.com/wzzfarewell/go-starter-example"

[db]
db-name = "merak_example"
host = "localhost"
password = "farewell"
port = 3306
user = "root"

[[tables]]
name = "t_user"
struct-name = "User"
package-name = "user"

[[tables]]
name = "t_user_info"
struct-name = "UserInfo"
package-name = "user"
```

3. Generate your project code by a configuration file.
```shell
go-starter -c application.toml
```

4. Run your project code.
```Go
cd /path/to/yourproject
go mod tidy
go run main.go
```
