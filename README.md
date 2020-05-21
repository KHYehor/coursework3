# Курсовая Работа
## Тема: Параллельное Программирование. Реализация Мьютекса.
### Про курсовую
* Выполнено по материалу из [лекции](https://youtu.be/JNLrITevhRI).
* Язык Реализации - golang 1.14.
* Успешные [тесты](https://travis-ci.com/github/KHYehor/coursework3) на TravisCI.

Не смотря на тот факт, что в языке уже есть реализация Мьютекса ( ["sync/mutex"](https://golang.org/pkg/sync/#Mutex) ) из коробки - 
реализовать его собственноручно не составляет больших проблем, в виду того,
что в языке есть модуль для выполнения атомарных операций ["sync/atomic"](https://golang.org/pkg/sync/atomic/)
(что хочется сказать, довольно удобный и понятый).

В моей реализации - мьютекс выглядит следующим [образом](https://github.com/KHYehor/coursework3/blob/37fa7de8685d77648e618621084cab8f530e4cd0/mutex.go#L19):
```go
// Mutex with timeout
type Mutex struct {
	state int32
	timeout time.Duration
}
``` 
Переменная `state` служит для отслеживания текущего состояния 
мьютекса - захвачен ли он каким-либо другим потоком. Ее размер мог бы быть и меньше
но атомарные операции языка позволяет работать только с `int32` и `int64` числами.
Для дополнения, и контроля - я добавил поле `timeout` - дабы было ограничение 
по времени на контроль мьютекса процессом.

Далее метод для [получение](https://github.com/KHYehor/coursework3/blob/37fa7de8685d77648e618621084cab8f530e4cd0/mutex.go#L25) мьютекса:
```go
// Getting mutex for control
func (m *Mutex) getMutex() bool {
	if atomic.CompareAndSwapInt32(&m.state, UNLOCKED, LOCKED) {
		return true
	}
	// Waiting for getting mutex
	start := time.Now()
	for {
		if atomic.CompareAndSwapInt32(&m.state, UNLOCKED, LOCKED) {
			// Finish stopwatch
			total := time.Now().Sub(start)
			// Printing total time of waiting for the mutex
			log.Printf("Mutex has been holding: %d ms", total.Microseconds())
			return true
		}
		total := time.Now().Sub(start)
		if total > m.timeout {
			panic(fmt.Sprintf(ERROR_GET, m.timeout.Microseconds()))
		}
	}
}
```
Состоит из двух логических частей - первая эта проверка есть ли возможность перевести
мьютекс в состояние "Захвачено", если же не вышло, тогда он уходит в синхронное ожидание,
пока мьютекс - либо не освободится другим процессом, либо не истечет предельно допустим время
на удержание мьютекса другими потоками.

Последние метод - для [освобождения](https://github.com/KHYehor/coursework3/blob/37fa7de8685d77648e618621084cab8f530e4cd0/mutex.go#L47) контроля над мьютексом.
```go
// Releasing mutex to let deal with it
func (m *Mutex) releaseMutex() bool {
	if atomic.CompareAndSwapInt32(&m.state, LOCKED, UNLOCKED) {
		return true
	}
	panic(fmt.Sprintf(ERROR_RELEASE))
}
```
Состоит из проверки - возможно ли освободить мьютекс, тоесть захватил ли его кто-либо. Если же этого 
не получилось сделать - значит вышла логическая ошибка в использовании, и какой-то из потоков
не сделал захват, и освобождать нечего.


