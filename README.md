https://gorm.io/docs/
https://gin-gonic.com/docs/quickstart/
https://pkg.go.dev/golang.org/x/crypto#section-readme
https://github.com/golang-jwt/jwt
https://github.com/joho/godotenv
https://github.com/githubnemo/CompileDaemon


Crear la carpeta del nuevo proyecto en C:\Users\marco\go\src\github.com\calleros
entrar con Git Bash a la carpeta

go mod init github.com/calleros/sich
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
go get -u github.com/gin-gonic/gin
go get -u golang.org/x/crypto/bcrypt
go get -u github.com/dgrijalva/jwt-go
go get github.com/joho/godotenv
go get github.com/githubnemo/CompileDaemon
go install github.com/githubnemo/CompileDaemon

iniciar con (code .) visual studio

iniciar el compilerdaemon:
compiledaemon --command="./sich"