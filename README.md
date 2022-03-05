# DILEMA

## О проекте

Dilema - это DI-контейнер для серверных приложений на языке Go, позволяющий довольно легко управлять зависимостями. На данный момент проект находится на стадии развития и улучшения, любые советы и комментарии приветствуются. В будущих версиях планируется добавить поддержку графа зависимостей для определения циклических зависимостей

## Основные типы

1. Dicon - непосредственно DI-контейнер. Предоставляет API для управления зависимостями
2. Destroyable - интерфейс с единственным методом Destroy() error. В случае очистки приложения вызывается метод Destroy у всех контейнеров, которые реализуют данный интерфейс, что позволяет в случае остановки работы приложения очистить соединение с базой данных, очистить кэш или выполнить другие действия
3. CallResults - специальный тип данных, возвращаемый при запуске функции. Введен для удобства подстановки значений результатов в переменные

Отдельно хотелось бы отметить, что возможны два типа контейнеров - временные (Temporary) и постоянные (Singletone). Временный контейнер создаётся каждый раз заново при получении, постоянный создаётся на этапе регистрации и в дальнейшем возвращает одну и ту же структуру.

## API

Существует два типа методов: обычные и с префиксом Must - в случае внутренней ошибки вызывают panic(err)

* ***RegisterTemporary(alias string, serviceInit interface{}) error*** - региструет временный контейнер, где serviceInit - функция, возвращающая либо интерфейсное значение, либо интерфейсное значение и ошибку. Метод возвращает ошибку несоответствие типов, если таковая была найдена. При попытке зарегистрировать новый контейнер с ранее использованным алиасом так же будет возвращена ошибка. **Существует версия MustRegisterTemporary**
* ***RegisterSingletone(alias string, serviceInit interface{}, args ...interface{}) error*** - региструет постоянный контейнер, где serviceInit - функция, возвращающая либо интерфейсное значение, либо интерфейсное значение и ошибку, agrs - необходимые для создания сервиса аргументы, переданные в том порядке, в каком их принимает функция-конструктор. Метод возвращает ошибку несоответствие типов, если таковая была найдена, и ошибки, возникающие при создании контейнера При попытке зарегистрировать новый контейнер с ранее использованным алиасом так же будет возвращена ошибка. **Существует версия MustRegisterSingletone**
* ***GetSingletone(alias string) (interface{}, error)*** - возвращает зарегестрированный ранее постоянный контейнер по переданному алиасу. В случае, если контейнер не найден, будет возвращена ошибка. **Существует версия MustGetSingletone**
* ***ProcessSingletone(alias string, container interface{}) error*** - позволяет подставить зарегестрированный ранее по переданному алиасу постоянный контейнер в переменную, переданную по ссылке как аргумент *container*. В случае, если контейнер не найден либо подстановка невозможна, будет возвращена ошибка. **Существует версия MustProcessSingletone**
* ***GetTemporary(alias string, args ...interface{}) (interface{}, error)*** - возвращает зарегестрированный ранее временный контейнер по переданному алиасу. Для создания контейнера необходимо передать аргументы в том порядке, в котором их принимает функция-конструктор. В случае, если контейнер не найден либо создание контейнера завершилось неудачей, будет возвращена ошибка. **Существует версия MustGetTemporary**
* ***ProcessTemporary(alias string, container interface{}, args ...interface{}) error*** - позволяет подставить зарегестрированный ранее по переданному алиасу постоянный контейнер в переменную, переданную по ссылке как аргумент *container*. Для создания контейнера необходимо передать аргументы в том порядке, в котором их принимает функция-конструктор. В случае, если контейнер не найден, подстановка невозможна или создание контейнера завершилось неудачей, будет возвращена ошибка. **Существует версия MustProcessTemporary**
* ***ProcessStruct(str interface{}) error*** - "собирает" структуру, переданную в метод в качестве аргумента, при условии, что поля публичные и имеют типы загеристрированных ранее постоянных контейнеров. Также у полей структуры могут быть указан тег *dilema:"имя_котейнера"* для надёжности подстановки. Если сборка структуры или подстановка завершились неудачей, будет возвращена ошибка. **Существует версия MustProcessStruct**
* ***Run(function interface{}, args ...interface{}) (CallResults, error)*** - вызывает функцию, переданную как первый аргумент. Агрументы функции необходимо передавать в том порядке, в котором они требуются для вызова функции. Возвращает CallResults и ошибку, если таковая возникла при запуске функции. **Существует версия MustRun**
* ***Recover(function interface{}, args ...interface{}) (cr CallResults, err error)*** - вызывает функцию, переданную как первый аргумент. Агрументы функции необходимо передавать в том порядке, в котором они требуются для вызова функции. Возвращает CallResults и ошибку, если таковая возникла при запуске функции. В случае возникновения паники внутри функции данный метод обрабатывает и возвращает возникшую ошибку
* ***RecoverAndClean(function interface{}, args ...interface{}) (cr CallResults, err error)*** - вызывает функцию, переданную как первый аргумент. Агрументы функции необходимо передавать в том порядке, в котором они требуются для вызова функции. Возвращает CallResults и ошибку, если таковая возникла при запуске функции. В случае возникновения паники внутри функции данный метод обрабатывает и возвращает возникшую ошибку, а также вызывает метод Destroy() у зарегистрированных постоянных контейнеров
* ***(СallResults) Process(values ...interface{}) error*** - позволяет установить переданным по ссылке переменным значения результаты выполнения функции. Переменные необходимо передавать в том порядке, как их возвращает функция. Если подстановка завершится неудачей, будет возвращена ошибка. **Существует версия MustProcess**
* ***Ctx() context.Context*** - позволяет получить контекст DI-контейнера
* ***SetCtx(ctx context.Context)*** - позволяет установить новый контекст в DI-контейнер
* ***AddToCtx(alias string, value interface{})*** - позволяет добавить новое значение в контекст с указанным алиасом
* ***GetFromCtx(alias string) interface{}*** - позволяет получить значение из контекста по переданному алиасу

## Пример

Примеры использования библиотеки можно найти в пакете **example**