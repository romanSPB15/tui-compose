# Changelog

## [3.0.0] – 2026-06-29

### Added
- Контейнеры `HBox` и `VBox` с автоматическим расчётом позиций (интерфейс `Container`).
- Виджет `InputField` – однострочное поле ввода с курсором.
- Поддержка мыши: SGR (1006), события кликов.
- Виджеты `Canvas` (16 цветов) и `CanvasRGB` (True Color) с двумя режимами отображения (`PixelTwoSymbol`, `PixelOneSymbol`).
- Методы `Width()` и `Height()` для `Canvas` и `CanvasRGB`.
- Методы `SetContent`, `SetOverlay`, `ShowOverlay`, `HideOverlay` для управления содержимым окна.
- Структура `Page` с методом `Open()`.
- Методы `SetTitle`, `CopyToClipboard`.
- Интерфейсы `TextInput`, `Container`, `ClickableAt`.
- Собственный парсер ANSI-последовательностей.
- `RegisterClickHandler` – глобальные обработчики мыши.
- Методы `AddWidgets`, `Clear` – заменены на `SetContent`.


### Changed
- Удалена зависимость от `github.com/eiannone/keyboard` – теперь используется собственный парсер ввода.
- Система фокуса переключения по Tab/Shift+Tab.
- Перерисовка оптимизирована (перерисовываются только изменённые строки).
- `RegisterKeyHandler` – регистрация обработчиков клавиатуры теперь с `KeyboardEvent`.

### Removed
- Зависимость от `github.com/eiannone/keyboard`, `github.com/charmbracelet/x/term` и `github.com/acarl005/stripansi`.


### Fixed
- Исправлены ошибки Canvas\[RGB\].

### Dependencies
- Единственная внешняя зависимость – `golang.org/x/term`.