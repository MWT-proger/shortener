# Рассчет покрытия тестами

```bash
cd cmd/shortener
 ```

``` bash
go test -v -coverprofile ../../docs/cover.out ../../...
go tool cover -html ../../docs/cover.out -o ../../docs/cover.html
open ../../docs/cover.html
```