# musicLib
store and manage music


### Install
1. Clone the repository:
```bash
git clone https://github.com/username/musicLib.git
```

2. Install dependencies:
```bash
cd musicLib
go mod tidy
go build
```

3. Run postgres:
```bash
make docker_build
```
or if builded:
```bash
make docker_run
```

4. Run the server:
```bash
make run
``` 

### Swagger
http://localhost:5000/swagger/index.html