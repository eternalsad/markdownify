package latex

import (
	"fmt"
	"regexp"
	"sort"
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
		// Греческие буквы строчные
		"\\alpha":      "α",
		"\\beta":       "β",
		"\\gamma":      "γ",
		"\\delta":      "δ",
		"\\epsilon":    "ε",
		"\\varepsilon": "ε",
		"\\zeta":       "ζ",
		"\\eta":        "η",
		"\\theta":      "θ",
		"\\vartheta":   "ϑ",
		"\\iota":       "ι",
		"\\kappa":      "κ",
		"\\lambda":     "λ",
		"\\mu":         "μ",
		"\\nu":         "ν",
		"\\xi":         "ξ",
		"\\pi":         "π",
		"\\varpi":      "ϖ",
		"\\rho":        "ρ",
		"\\varrho":     "ϱ",
		"\\sigma":      "σ",
		"\\varsigma":   "ς",
		"\\tau":        "τ",
		"\\upsilon":    "υ",
		"\\phi":        "φ",
		"\\varphi":     "φ",
		"\\chi":        "χ",
		"\\psi":        "ψ",
		"\\omega":      "ω",

		// Греческие буквы заглавные
		"\\Gamma":   "Γ",
		"\\Delta":   "Δ",
		"\\Theta":   "Θ",
		"\\Lambda":  "Λ",
		"\\Xi":      "Ξ",
		"\\Pi":      "Π",
		"\\Sigma":   "Σ",
		"\\Upsilon": "Υ",
		"\\Phi":     "Φ",
		"\\Psi":     "Ψ",
		"\\Omega":   "Ω",

		// Математические операторы
		"\\sum":      "∑",
		"\\prod":     "∏",
		"\\coprod":   "∐",
		"\\int":      "∫",
		"\\oint":     "∮",
		"\\iint":     "∬",
		"\\iiint":    "∭",
		"\\partial":  "∂",
		"\\nabla":    "∇",
		"\\pm":       "±",
		"\\mp":       "∓",
		"\\times":    "×",
		"\\div":      "÷",
		"\\setminus": "\\",
		"\\cdot":     "·",
		"\\ast":      "∗",
		"\\star":     "★",
		"\\circ":     "∘",
		"\\bullet":   "•",

		// Отношения
		"\\leq":    "≤",
		"\\geq":    "≥",
		"\\neq":    "≠",
		"\\approx": "≈",
		"\\equiv":  "≡",
		"\\cong":   "≅",
		"\\sim":    "∼",
		"\\propto": "∝",
		"\\prec":   "≺",
		"\\succ":   "≻",
		"\\preceq": "⪯",
		"\\succeq": "⪰",
		"\\ll":     "≪",
		"\\gg":     "≫",

		// Теория множеств
		"\\in":         "∈",
		"\\notin":      "∉",
		"\\ni":         "∋",
		"\\subset":     "⊂",
		"\\supset":     "⊃",
		"\\subseteq":   "⊆",
		"\\supseteq":   "⊇",
		"\\cup":        "∪",
		"\\cap":        "∩",
		"\\emptyset":   "∅",
		"\\varnothing": "∅",

		// Логические символы
		"\\land":      "∧",
		"\\lor":       "∨",
		"\\lnot":      "¬",
		"\\forall":    "∀",
		"\\exists":    "∃",
		"\\nexists":   "∄",
		"\\therefore": "∴",
		"\\because":   "∵",

		// Стрелки
		"\\to":             "→",
		"\\rightarrow":     "→",
		"\\leftarrow":      "←",
		"\\Rightarrow":     "⇒",
		"\\Leftarrow":      "⇐",
		"\\mapsto":         "↦",
		"\\uparrow":        "↑",
		"\\downarrow":      "↓",
		"\\updownarrow":    "↕",
		"\\leftrightarrow": "↔",
		"\\Leftrightarrow": "⇔",

		// Разное
		"\\infty":         "∞",
		"\\aleph":         "ℵ",
		"\\hbar":          "ℏ",
		"\\ell":           "ℓ",
		"\\wp":            "℘",
		"\\Re":            "ℜ",
		"\\Im":            "ℑ",
		"\\angle":         "∠",
		"\\measuredangle": "∡",
		"\\triangle":      "△",
		"\\square":        "□",
		"\\overline":      "‾", // Приближение, может не работать правильно
		"\\dots":          "…",
		"\\cdots":         "⋯",
		"\\vdots":         "⋮",
		"\\ddots":         "⋱",

		// Специальные комбинации
		"\\_0": "₀",
		"\\_1": "₁",
		"\\_2": "₂",
		"\\_3": "₃",
		"\\_4": "₄",
		"\\_5": "₅",
		"\\_6": "₆",
		"\\_7": "₇",
		"\\_8": "₈",
		"\\_9": "₉",
		"\\^0": "⁰",
		"\\^1": "¹",
		"\\^2": "²",
		"\\^3": "³",
		"\\^4": "⁴",
		"\\^5": "⁵",
		"\\^6": "⁶",
		"\\^7": "⁷",
		"\\^8": "⁸",
		"\\^9": "⁹",

		// Векторное исчисление
		"\\vec":  "→", // Обычно отображается над символом, здесь упрощенно
		"\\dot":  "·", // Точка над символом (упрощенно)
		"\\ddot": "¨", // Две точки над символом (упрощенно)

		// Дроби и корни
		"\\frac": "/", // Упрощенное представление дроби
		"\\sqrt": "√",

		// Скобки
		"\\lbrace": "{",
		"\\rbrace": "}",
		"\\langle": "⟨",
		"\\rangle": "⟩",
		"\\lceil":  "⌈",
		"\\rceil":  "⌉",
		"\\lfloor": "⌊",
		"\\rfloor": "⌋",

		// Физические константы
		"\\varepsilon_0": "ε₀",
		"\\mu_0":         "μ₀",

		// Дополнительные операторы
		"\\oplus":    "⊕",
		"\\otimes":   "⊗",
		"\\perp":     "⊥",
		"\\parallel": "∥",
		"\\surd":     "√",
	}

	return &LaTeXToMarkdownV2{
		latexSymbols: latexSymbols,
	}
}

// containsLaTeXSymbols проверяет, содержит ли текст LaTeX символы
func (l *LaTeXToMarkdownV2) ContainsLaTeXSymbols(content string) bool {
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

func (l *LaTeXToMarkdownV2) convertLaTeXToUnicode(content string) string {
	// 1. Сортируем символы по длине (от длинных к коротким)
	type symbolPair struct {
		latex   string
		unicode string
	}

	pairs := make([]symbolPair, 0, len(l.latexSymbols))
	for latex, unicode := range l.latexSymbols {
		pairs = append(pairs, symbolPair{latex, unicode})
	}

	// Сортировка от длинных к коротким
	sort.Slice(pairs, func(i, j int) bool {
		return len(pairs[i].latex) > len(pairs[j].latex)
	})

	// 2. Используем builder для создания результата
	var builder strings.Builder
	i := 0

	for i < len(content) {
		matched := false

		// Проверяем на совпадение с символами (начиная с самых длинных)
		for _, pair := range pairs {
			if i+len(pair.latex) <= len(content) && content[i:i+len(pair.latex)] == pair.latex {
				builder.WriteString(pair.unicode)
				i += len(pair.latex)
				matched = true
				break
			}
		}

		if !matched {
			// Если символ не найден, просто добавляем текущий символ
			builder.WriteByte(content[i])
			i++
		}
	}

	result := builder.String()

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
		if !l.ContainsLaTeXSymbols(content) {
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
