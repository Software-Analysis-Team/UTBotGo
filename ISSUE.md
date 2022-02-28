# Проблема с неопределёнными константами в ll-файле

## Воспроизведение

1.  Собрать docker-образ
```
docker build -t gollvm_klee:5 .
```

2.  Запустить docker-контейнер
```
docker run -it --rm -v ${PWD}:/home gollvm_klee:5
```

3.  Транслировать Go-файл в ll-файл
```
cd /home
./go2ll simple.go
```

4.  Подправить полученный ll-файл
```
g++ -O3 -Wall $(llvm-config --cxxflags) -o modifier modifier.cpp $(llvm-config --ldflags --libs --system-libs)
./modifier simple.ll 2> simple_fixed.ll
```

5.  Сгенерировать тесты с помощью klee (не работает)
```
klee --disable-verify --external-calls=all simple_fixed.ll
```

## Замечания

  * `/go/lib64/libgo.so` -- стандартная библиотека для Go

  * `go build -n simple.go` -- команда для просмотра bash-скрипта сборки (не работает после запуска `go2ll`)
