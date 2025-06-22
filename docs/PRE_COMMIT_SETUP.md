# Pre-commit Setup for go-locate

Этот документ описывает настройку и использование pre-commit хуков для проекта go-locate.

## Установка pre-commit

### macOS (с Homebrew)
```bash
brew install pre-commit
```

### Python pip
```bash
pip install pre-commit
```

### Проверка установки
```bash
pre-commit --version
```

## Настройка в проекте

### 1. Установка хуков
```bash
# Установить pre-commit хуки в git
pre-commit install

# Установить хук для проверки commit messages
pre-commit install --hook-type commit-msg
```

### 2. Первый запуск
```bash
# Запустить хуки на всех файлах (рекомендуется при первой настройке)
pre-commit run --all-files
```

## Настроенные хуки

### Базовые проверки файлов
- **trailing-whitespace**: Удаляет пробелы в конце строк
- **end-of-file-fixer**: Обеспечивает наличие новой строки в конце файла
- **check-yaml**: Проверяет синтаксис YAML файлов
- **check-toml**: Проверяет синтаксис TOML файлов
- **check-json**: Проверяет синтаксис JSON файлов
- **check-added-large-files**: Предотвращает коммит больших файлов (>1MB)
- **check-merge-conflict**: Проверяет наличие маркеров конфликта слияния
- **mixed-line-ending**: Обеспечивает единообразные окончания строк (LF)

### Go-специфичные хуки
- **go-fmt-repo**: Форматирование кода с помощью gofmt
- **go-imports-repo**: Управление импортами с помощью goimports
- **go-mod-tidy-repo**: Очистка go.mod файлов
- **go-vet-repo-mod**: Статический анализ с помощью go vet
- **go-build-repo-mod**: Проверка компиляции
- **go-test-repo-mod**: Запуск тестов с покрытием кода
- **golangci-lint-repo-mod**: Запуск golangci-lint с существующей конфигурацией

### Безопасность
- **detect-secrets**: Обнаружение секретов в коде

### Commit Messages
- **conventional-pre-commit**: Проверка формата commit messages (conventional commits)

## Использование

### Автоматический запуск
Pre-commit хуки запускаются автоматически при каждом `git commit`.

### Ручной запуск
```bash
# Запустить все хуки
pre-commit run --all-files

# Запустить конкретный хук
pre-commit run go-fmt-repo

# Запустить хуки только на измененных файлах
pre-commit run
```

### Пропуск хуков
```bash
# Пропустить все хуки для одного коммита
git commit --no-verify -m "commit message"

# Пропустить конкретные хуки
SKIP=go-test-repo-mod git commit -m "commit message"

# Пропустить несколько хуков
SKIP=go-test-repo-mod,golangci-lint-repo-mod git commit -m "commit message"
```

## Формат Commit Messages

Проект использует [Conventional Commits](https://www.conventionalcommits.org/) формат:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Примеры типов:
- `feat`: новая функциональность
- `fix`: исправление бага
- `docs`: изменения в документации
- `style`: форматирование, отсутствующие точки с запятой и т.д.
- `refactor`: рефакторинг кода
- `test`: добавление тестов
- `chore`: обновление задач сборки, конфигурации и т.д.

### Примеры commit messages:
```bash
git commit -m "feat: add search functionality"
git commit -m "fix(config): resolve TOML parsing error"
git commit -m "docs: update installation instructions"
git commit -m "test: add unit tests for search module"
```

## Обновление хуков

```bash
# Обновить хуки до последних версий
pre-commit autoupdate

# Очистить кеш pre-commit (если возникают проблемы)
pre-commit clean
```

## Настройка IDE

### VS Code
Добавьте в `.vscode/settings.json`:
```json
{
    "go.formatTool": "goimports",
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
        "source.organizeImports": true
    }
}
```

### GoLand/IntelliJ IDEA
1. Настройте goimports как форматтер по умолчанию
2. Включите автоформатирование при сохранении
3. Настройте golangci-lint как внешний инструмент

## Устранение неполадок

### Проблема с установкой хуков
```bash
# Переустановить хуки
pre-commit uninstall
pre-commit install --install-hooks
```

### Проблемы с golangci-lint
```bash
# Убедитесь, что golangci-lint установлен
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Проверьте конфигурацию
golangci-lint config -h
```

### Проблемы с Go модулями
```bash
# Очистите модули и переустановите зависимости
go clean -modcache
go mod download
go mod tidy
```

## Дополнительные ресурсы

- [Pre-commit документация](https://pre-commit.com/)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [GolangCI-Lint документация](https://golangci-lint.run/)
- [TekWizely/pre-commit-golang](https://github.com/TekWizely/pre-commit-golang)
