package latex

import (
	"fmt"
	"regexp"
	"strings"
)

// LaTeXToMarkdownV2 конвертирует LaTeX формулы в формат Markdown V2 для Telegram
type LaTeXToMarkdownV2 struct {
	latexSymbols map[string]string
}

// NewLaTeXToMarkdownV2 создает новый экземпляр конвертера
func NewLaTeXToMarkdownV2() *LaTeXToMarkdownV2 {
	// Словарь для конвертации LaTeX символов в Unicode
	latexSymbols := map[string]string{
		"\\alpha":      "α",
		"\\beta":       "β",
		"\\gamma":      "γ",
		"\\delta":      "δ",
		"\\epsilon":    "ε",
		"\\zeta":       "ζ",
		"\\eta":        "η",
		"\\theta":      "θ",
		"\\lambda":     "λ",
		"\\mu":         "μ",
		"\\pi":         "π",
		"\\rho":        "ρ",
		"\\sigma":      "σ",
		"\\tau":        "τ",
		"\\phi":        "φ",
		"\\omega":      "ω",
		"\\sum":        "∑",
		"\\int":        "∫",
		"\\infty":      "∞",
		"\\partial":    "∂",
		"\\nabla":      "∇",
		"\\pm":         "±",
		"\\times":      "×",
		"\\div":        "÷",
		"\\leq":        "≤",
		"\\geq":        "≥",
		"\\neq":        "≠",
		"\\approx":     "≈",
		"\\cdot":       "·",
		"\\in":         "∈",
		"\\subset":     "⊂",
		"\\cup":        "∪",
		"\\cap":        "∩",
		"\\to":         "→",
		"\\rightarrow": "→",
		"\\leftarrow":  "←",
		"\\Rightarrow": "⇒",
		"\\Leftarrow":  "⇐",
		// Добавьте больше символов по необходимости
	}

	return &LaTeXToMarkdownV2{
		latexSymbols: latexSymbols,
	}
}

// containsLaTeXSymbols проверяет, содержит ли текст LaTeX символы
func (l *LaTeXToMarkdownV2) containsLaTeXSymbols(content string) bool {
	if len(content) < 5 {
		return false
	}

	// Проверяем наличие общих LaTeX команд
	latexCommands := []string{
		"\\frac",
		"\\sqrt",
		"\\begin",
		"\\end",
		"\\left",
		"\\right",
		"\\limits",
		"\\displaystyle",
	}

	// Проверяем команды и символы из словаря
	for _, cmd := range latexCommands {
		if strings.Contains(content, cmd) {
			return true
		}
	}

	for symbol := range l.latexSymbols {
		if strings.Contains(content, symbol) {
			return true
		}
	}

	return false
}

// convertLaTeXToUnicode преобразует LaTeX команды в Unicode символы
func (l *LaTeXToMarkdownV2) convertLaTeXToUnicode(content string) string {
	result := content

	// Заменяем LaTeX символы на Unicode эквиваленты
	for latex, unicode := range l.latexSymbols {
		result = strings.ReplaceAll(result, latex, unicode)
	}

	// Обработка специальных конструкций
	result = l.processSpecialConstructs(result)

	return result
}

// processSpecialConstructs обрабатывает специальные LaTeX конструкции
func (l *LaTeXToMarkdownV2) processSpecialConstructs(content string) string {
	result := content

	// Примеры обработки специальных конструкций
	// \frac{a}{b} -> a/b
	fracRegex := regexp.MustCompile(`\\frac\{([^}]+)\}\{([^}]+)\}`)
	result = fracRegex.ReplaceAllString(result, "$1/$2")

	// \sqrt{x} -> √x
	sqrtRegex := regexp.MustCompile(`\\sqrt\{([^}]+)\}`)
	result = sqrtRegex.ReplaceAllString(result, "√$1")

	// Верхний индекс: x^{2} -> x²
	supRegex := regexp.MustCompile(`\^{([0-9])}`)
	superscripts := map[string]string{
		"0": "⁰", "1": "¹", "2": "²", "3": "³", "4": "⁴",
		"5": "⁵", "6": "⁶", "7": "⁷", "8": "⁸", "9": "⁹",
	}
	result = supRegex.ReplaceAllStringFunc(result, func(match string) string {
		digit := match[2:3]
		if sup, ok := superscripts[digit]; ok {
			return sup
		}
		return match
	})

	// Нижний индекс: x_{i} -> xᵢ
	subRegex := regexp.MustCompile(`_\{([^}]+)\}`)
	subscripts := map[string]string{
		"0": "₀", "1": "₁", "2": "₂", "3": "₃", "4": "₄",
		"5": "₅", "6": "₆", "7": "₇", "8": "₈", "9": "₉",
		"i": "ᵢ", "j": "ⱼ", "k": "ₖ", "n": "ₙ", "m": "ₘ",
	}
	result = subRegex.ReplaceAllStringFunc(result, func(match string) string {
		char := match[2 : len(match)-1]
		if sub, ok := subscripts[char]; ok {
			return sub
		}
		return "_" + char
	})

	return result
}

// escapeLaTeX обрабатывает LaTeX формулы в тексте для Markdown V2
func (l *LaTeXToMarkdownV2) EscapeLaTeX(text []byte) []byte {
	// Регулярные выражения для поиска блочных и инлайн формул
	blockMathRegex := regexp.MustCompile(`(?s)\\\[(.*?)\\\]`)
	inlineMathRegex := regexp.MustCompile(`(?s)\\\((.*?)\\\)`)

	// Функция для обработки match'ей
	processMatch := func(match []string, isBlock bool) string {
		if len(match) < 2 {
			return match[0]
		}

		content := match[1]
		if !l.containsLaTeXSymbols(content) {
			return match[0] // Возвращаем без изменений, если нет LaTeX символов
		}

		// Преобразуем LaTeX в Unicode
		converted := l.convertLaTeXToUnicode(content)

		// Экранируем для Markdown V2
		escaped := escapeMarkdownV2(converted)

		if isBlock {
			return fmt.Sprintf("```\n%s\n```", strings.TrimSpace(escaped))
		}
		return fmt.Sprintf("`%s`", strings.TrimSpace(escaped))
	}

	// Разбиваем текст на параграфы
	lines := strings.Split(string(text), "\n\n")
	processedLines := make([]string, len(lines))

	for i, line := range lines {
		// Обрабатываем блочные формулы
		processedLine := blockMathRegex.ReplaceAllStringFunc(line, func(match string) string {
			matches := blockMathRegex.FindStringSubmatch(match)
			return processMatch(matches, true)
		})

		// Обрабатываем инлайн формулы
		processedLine = inlineMathRegex.ReplaceAllStringFunc(processedLine, func(match string) string {
			matches := inlineMathRegex.FindStringSubmatch(match)
			return processMatch(matches, false)
		})

		processedLines[i] = processedLine
	}

	return []byte(strings.Join(processedLines, "\n\n"))
}

// escapeMarkdownV2 экранирует специальные символы для Markdown V2
func escapeMarkdownV2(text string) string {
	specialChars := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}

	result := text
	for _, char := range specialChars {
		result = strings.ReplaceAll(result, char, "\\"+char)
	}

	return result
}
