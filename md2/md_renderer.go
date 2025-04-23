package md2

import (
	"bytes"
	"fmt"
	"github.com/gomarkdown/markdown/parser/latex"
	"io"
	"strings"

	"github.com/gomarkdown/markdown/ast"
)

// Renderer renders to markdown. Allows to convert to a canonnical
// form
type Renderer struct {
	orderedListCounter map[int]int
	// used to keep track of whether a given list item uses a paragraph
	// for large spacing.
	paragraph map[int]bool

	lastOutputLen  int
	listDepth      int
	indentSize     int
	lastNormalText string

	latex *latex.LaTeXToMarkdownV2
}

// NewRenderer returns a Markdown renderer.
func NewRenderer() *Renderer {
	return &Renderer{
		orderedListCounter: map[int]int{},
		paragraph:          map[int]bool{},
		indentSize:         4,
		latex:              latex.NewLaTeXToMarkdownV2(),
	}
}

func (r *Renderer) out(w io.Writer, d []byte) {
	r.lastOutputLen = len(d)
	w.Write(d)
}

func (r *Renderer) outs(w io.Writer, s string) {
	r.lastOutputLen = len(s)
	io.WriteString(w, s)
}

func (r *Renderer) doubleSpace(w io.Writer) {
	// TODO: need to remember number of written bytes
	//if out.Len() > 0 {
	r.outs(w, "\n")
	//}
}

func (r *Renderer) list(w io.Writer, node *ast.List, entering bool) {
	if entering {
		r.listDepth++
		flags := node.ListFlags
		if flags&ast.ListTypeOrdered != 0 {
			r.orderedListCounter[r.listDepth] = 1
		}
	} else {
		r.listDepth--
		fmt.Fprintf(w, "\n")
	}
}

func (r *Renderer) listItem(w io.Writer, node *ast.ListItem, entering bool) {
	flags := node.ListFlags
	bullet := string(node.BulletChar)

	if entering {
		//for i := 1; i < r.listDepth; i++ {
		//	for i := 0; i < r.indentSize; i++ {
		//		fmt.Fprintf(w, " ")
		//	}
		//}

		if r.listDepth >= 1 {
			fmt.Fprintf(w, (strings.Repeat(" ", (r.listDepth-1)*r.indentSize)))
		}

		if flags&ast.ListTypeOrdered != 0 {
			fmt.Fprintf(w, "%d\\. ", r.orderedListCounter[r.listDepth])
			r.orderedListCounter[r.listDepth]++
		} else {
			fmt.Fprintf(w, "\\%s ", bullet)
		}
	}
}

func (r *Renderer) para(w io.Writer, node *ast.Paragraph, entering bool) {
	parent := node.GetParent()
	secondListItem := false

	switch parent.(type) {
	case *ast.ListItem:
		parentList, ok := parent.(*ast.ListItem)
		if !ok {
			break
		}

		if len(parent.GetChildren()) > 1 {
			if node.AsContainer() == parentList.GetChildren()[1].AsContainer() {
				secondListItem = true
			}
		}
	}

	if entering && secondListItem {
		if secondListItem {
			r.outs(w, "\t")
		}
	}

	if !entering && r.lastOutputLen > 0 {
		var br = "\n\n"

		// List items don't need the extra line-break.
		if _, ok := node.Parent.(*ast.ListItem); ok {
			br = "\n"
		}

		r.outs(w, br)
	}
}

// escape replaces instances of backslash with escaped backslash in text.
func escape(text []byte) []byte {
	return bytes.Replace(text, []byte(`\`), []byte(`\\`), -1)
}

func isNumber(data []byte) bool {
	for _, b := range data {
		if b < '0' || b > '9' {
			return false
		}
	}
	return true
}

// "_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!",
func needsEscaping(text []byte, lastNormalText string) bool {
	switch string(text) {
	case `\`,
		"`",
		"*",
		"_",
		"{", "}",
		"[", "]",
		"(", ")",
		"#",
		"+",
		"~",
		"=",
		"|",
		"-":
		return true
	case "!":
		return false
	case ".":
		// Return true if number, because a period after a number must be escaped to not get parsed as an ordered list.
		return isNumber([]byte(lastNormalText))
	case "<", ">":
		return true
	default:
		return false
	}
}

// cleanWithoutTrim is like clean, but doesn't trim blanks.
func cleanWithoutTrim(s string) string {
	var b []byte
	var p byte
	for i := 0; i < len(s); i++ {
		q := s[i]
		if q == '\n' || q == '\r' || q == '\t' {
			q = ' '
		}
		if q != ' ' || p != ' ' {
			b = append(b, q)
			p = q
		}
	}
	return string(b)
}

func (r *Renderer) skipSpaceIfNeededNormalText(w io.Writer, cleanString string) bool {
	if cleanString[0] != ' ' {
		return false
	}

	return false
	//  TODO: what did it mean to do?
	// we no longer use *bytes.Buffer for out, so whatever this tracked,
	// it has to be done in a different wy
	/*
		if _, ok := r.normalTextMarker[out]; !ok {
			r.normalTextMarker[out] = -1
		}
		return r.normalTextMarker[out] == out.Len()
	*/
}

func (r *Renderer) text(w io.Writer, text *ast.Text) {
	lit := []byte(escapeMarkdownV2(string(text.Literal)))
	normalText := string(text.Literal)
	if needsEscaping(lit, r.lastNormalText) {
		lit = append([]byte("\\"), lit...)
	}
	r.lastNormalText = normalText
	if r.listDepth > 0 && string(lit) == "\n" {
		// TODO: See if this can be cleaned up... It's needed for lists.
		return
	}
	cleanString := cleanWithoutTrim(string(lit))
	if cleanString == "" {
		return
	}
	// Skip first space if last character is already a space (i.e., no need for a 2nd space in a row).
	if r.skipSpaceIfNeededNormalText(w, cleanString) {
		cleanString = cleanString[1:]
	}
	r.outs(w, cleanString)
	// If it ends with a space, make note of that.
	//if len(cleanString) >= 1 && cleanString[len(cleanString)-1] == ' ' {
	// TODO: write equivalent of this
	// r.normalTextMarker[out] = out.Len()
	//}
}

func (r *Renderer) surround(w io.Writer, symbol string) {
	r.outs(w, symbol)
}

func (r *Renderer) htmlSpan(w io.Writer, node *ast.HTMLSpan) {
	r.out(w, node.Literal)
}

func (r *Renderer) htmlBlock(w io.Writer, node *ast.HTMLBlock) {
	r.doubleSpace(w)
	r.out(w, node.Literal)
	r.outs(w, "\n\n")
}

func (r *Renderer) codeBlock(w io.Writer, node *ast.CodeBlock) {
	r.doubleSpace(w)
	text := node.Literal
	lang := string(node.Info)
	// Parse out the language name.
	count := 0
	for _, elt := range strings.Fields(lang) {
		if elt[0] == '.' {
			elt = elt[1:]
		}
		if len(elt) == 0 {
			continue
		}
		r.outs(w, "```")
		r.outs(w, elt)
		count++
		break
	}

	if count == 0 {
		r.outs(w, "```")
	}
	r.outs(w, "\n")
	r.out(w, text)
	r.outs(w, "```\n\n")
}

func (r *Renderer) code(w io.Writer, node *ast.Code) {
	r.outs(w, "`")
	r.out(w, node.Literal)
	r.outs(w, "`")
}

//// Рендерер для заголовка таблицы
//func (r *Renderer) tableHeader(w io.Writer, node *ast.TableHeader, entering bool) {
//	if !entering {
//		// Определяем количество столбцов по первой строке заголовка
//		if len(node.Children) > 0 {
//			if headerRow, ok := node.Children[0].(*ast.TableRow); ok {
//				numColumns := len(headerRow.Children)
//
//				// Для каждого столбца добавляем разделитель нужной длины
//				for i := 0; i < numColumns; i++ {
//					// Получаем ячейку заголовка
//					if i < len(headerRow.Children) {
//						if cell, ok := headerRow.Children[i].(*ast.TableCell); ok {
//							// Измеряем длину содержимого ячейки
//							var cellBuf bytes.Buffer
//							for _, child := range cell.Children {
//								ast.WalkFunc(child, func(n ast.Node, entering bool) ast.WalkStatus {
//									return r.RenderNode(&cellBuf, n, entering)
//								})
//							}
//
//							content := cellBuf.String()
//							// Минимальная ширина разделителя
//							width := len(content)
//							if width < 3 {
//								width = 3
//							}
//
//							// Добавляем "-" соответствующей длины
//							w.Write([]byte(" " + strings.Repeat("\\-", width*3)))
//						}
//					} else {
//						// Если ячейки нет, используем дефолтную ширину
//						w.Write([]byte(" \\-\\-\\-\\-\\-\\-"))
//					}
//				}
//			}
//		}
//	}
//}
//
//// Рендерер для тела таблицы
//func (r *Renderer) tableBody(w io.Writer, node *ast.TableBody) {
//}
//
//// Рендерер для строки таблицы
//func (r *Renderer) tableRow(w io.Writer, node *ast.TableRow) {
//	r.out(w, []byte("\n"))
//}
//
//// Рендерер для подвала таблицы
//func (r *Renderer) tableFooter(w io.Writer, node *ast.TableFooter) {
//
//}
//
//// Рендерер для ячейки таблицы
//func (r *Renderer) tableCell(w io.Writer, node *ast.TableCell, entering bool) {
//	r.outs(w, "\\| ")
//
//	r.outs(w, escapeMarkdownV2(string(node.Literal)))
//}

// Функция для извлечения текстового содержимого из узла
func extractTextContent(w io.Writer, node ast.Node) {
	ast.WalkFunc(node, func(n ast.Node, entering bool) ast.WalkStatus {
		if entering && n.GetChildren() == nil {
			if text, ok := n.(*ast.Text); ok {
				// Заменяем переводы строк на пробелы
				content := strings.ReplaceAll(string(text.Literal), "\n", " ")
				content = strings.ReplaceAll(content, "\r", "")
				w.Write([]byte(content))
			}
		}
		return ast.GoToNext
	})
}

// Рендерер для таблицы с исправлением текста после таблицы
func (r *Renderer) table(w io.Writer, node *ast.Table) {
	// Собираем информацию о структуре таблицы
	var rows []*ast.TableRow
	headerRowIdx := -1

	// Собираем все строки и определяем, какие из них - заголовки
	for i, child := range node.Children {
		switch c := child.(type) {
		case *ast.TableHeader:
			for _, headerRow := range c.Children {
				if tr, ok := headerRow.(*ast.TableRow); ok {
					rows = append(rows, tr)
					if headerRowIdx == -1 {
						headerRowIdx = len(rows) - 1
					}
				}
			}
		case *ast.TableBody:
			for _, bodyRow := range c.Children {
				if tr, ok := bodyRow.(*ast.TableRow); ok {
					rows = append(rows, tr)
				}
			}
		case *ast.TableFooter:
			for _, footerRow := range c.Children {
				if tr, ok := footerRow.(*ast.TableRow); ok {
					rows = append(rows, tr)
				}
			}
		case *ast.TableRow:
			if i == 0 { // Если первый элемент TableRow, считаем его заголовком
				headerRowIdx = 0
			}
			rows = append(rows, c)
		}
	}

	// Если строк нет, выходим
	if len(rows) == 0 {
		return
	}

	// Определяем количество столбцов
	numColumns := 0
	for _, row := range rows {
		if len(row.Children) > numColumns {
			numColumns = len(row.Children)
		}
	}

	if numColumns == 0 {
		return
	}

	// Собираем содержимое ячеек
	cellContents := make([][]string, len(rows))

	for i, row := range rows {
		cellContents[i] = make([]string, numColumns)

		for j := 0; j < numColumns; j++ {
			// Если ячейка существует, рендерим её содержимое
			if j < len(row.Children) {
				if cell, ok := row.Children[j].(*ast.TableCell); ok {
					var buf bytes.Buffer

					for _, child := range cell.Children {
						renderNodeText(&buf, child)
					}

					cellContents[i][j] = strings.TrimSpace(buf.String())
				}
			}
		}
	}

	// Определяем максимальную ширину для каждого столбца
	columnWidths := make([]int, numColumns)

	for _, row := range cellContents {
		for j, content := range row {
			if len(content) > columnWidths[j] {
				columnWidths[j] = len(content)
			}
		}
	}

	// Обеспечиваем минимальную ширину столбца
	for j := range columnWidths {
		if columnWidths[j] < 3 {
			columnWidths[j] = 3
		}
	}

	// Начинаем рендеринг таблицы с блока кода для Telegram
	fmt.Fprint(w, "```\n")

	// Рендерим строки таблицы
	for i, row := range cellContents {
		// Формируем строку таблицы
		fmt.Fprint(w, "| ")

		for j, content := range row {
			// Выравнивание по левому краю с добавлением пробелов
			paddedContent := content
			if len(paddedContent) < columnWidths[j] {
				paddedContent += strings.Repeat(" ", columnWidths[j]-len(paddedContent))
			} else if len(paddedContent) > columnWidths[j] {
				// Обрезаем длинное содержимое
				paddedContent = paddedContent[:columnWidths[j]-3] + "..."
			}

			fmt.Fprint(w, paddedContent)

			if j < numColumns-1 {
				fmt.Fprint(w, " | ")
			} else {
				fmt.Fprint(w, " |")
			}
		}

		fmt.Fprint(w, "\n")

		// Добавляем разделитель после заголовка
		if i == headerRowIdx {
			fmt.Fprint(w, "| ")

			for j, width := range columnWidths {
				fmt.Fprint(w, strings.Repeat("-", width))

				if j < numColumns-1 {
					fmt.Fprint(w, " | ")
				} else {
					fmt.Fprint(w, " |")
				}
			}

			fmt.Fprint(w, "\n")
		}
	}

	// Завершаем блок кода
	fmt.Fprint(w, "```\n")
}

// Вспомогательная функция для извлечения текста из узла
func renderNodeText(w io.Writer, node ast.Node) {
	if node == nil {
		return
	}

	switch c := node.(type) {
	case *ast.Text:
		// Заменяем переводы строк и табуляции на пробелы
		text := string(c.Literal)
		text = strings.ReplaceAll(text, "\n", " ")
		text = strings.ReplaceAll(text, "\r", "")
		text = strings.ReplaceAll(text, "\t", " ")

		// Сжимаем множественные пробелы
		for strings.Contains(text, "  ") {
			text = strings.ReplaceAll(text, "  ", " ")
		}

		w.Write([]byte(text))

	case *ast.Link:
		// Для ссылок выводим текст ссылки
		for _, child := range c.Children {
			renderNodeText(w, child)
		}

	case *ast.Strong, *ast.Emph:
		// Для форматированного текста рекурсивно обрабатываем содержимое
		for _, child := range c.GetChildren() {
			renderNodeText(w, child)
		}

	case *ast.Code:
		// Для кода выводим его содержимое
		text := string(c.Literal)
		text = strings.ReplaceAll(text, "\n", " ")
		w.Write([]byte(text))

	case *ast.Image:
		// Для изображений выводим альтернативный текст
		fmt.Fprint(w, "[Image]")

	case *ast.Math:
		// Для математических выражений
		w.Write(c.Literal)

	default:
		// Для других типов узлов рекурсивно обрабатываем дочерние элементы
		if children := node.GetChildren(); children != nil {
			for _, child := range children {
				renderNodeText(w, child)
			}
		}
	}
}

// Чтобы избежать дублирования таблиц, нужно отслеживать обработанные таблицы
// Для этого можно использовать глобальную переменную или добавить метод в рендерер

// Карта для отслеживания обработанных таблиц
var processedTables = make(map[*ast.Table]bool)

// Метод для проверки, была ли таблица уже обработана
func (r *Renderer) wasTableProcessed(table *ast.Table) bool {
	if processedTables[table] {
		return true
	}
	processedTables[table] = true
	return false
}

// Проверяет, находится ли узел внутри таблицы
func isInsideTable(node ast.Node) bool {
	for node != nil {
		if _, ok := node.(*ast.Table); ok {
			return true
		}
		node = node.GetParent()
	}
	return false
}

// Получает предыдущий узел того же уровня
func getPreviousSibling(parent ast.Node, current ast.Node) ast.Node {
	children := parent.GetChildren()
	if len(children) <= 1 {
		return nil
	}

	for i, child := range children {
		if child == current && i > 0 {
			return children[i-1]
		}
	}

	return nil
}

// Проверяет, является ли узел таблицей
func isTable(node ast.Node) bool {
	_, ok := node.(*ast.Table)
	return ok
}

// Эти функции могут остаться пустыми, так как обработка происходит в table
func (r *Renderer) tableHeader(w io.Writer, node *ast.TableHeader) {
	// Обработка происходит в методе table
}

func (r *Renderer) tableBody(w io.Writer, node *ast.TableBody) {
	// Обработка происходит в методе table
}

func (r *Renderer) tableRow(w io.Writer, node *ast.TableRow) {
	// Обработка происходит в методе table
}

func (r *Renderer) tableFooter(w io.Writer, node *ast.TableFooter) {
	// Обработка происходит в методе table
}

func (r *Renderer) tableCell(w io.Writer, node *ast.TableCell) {
	// Обработка происходит в методе table
}

var headings = map[int]string{
	0: "📌",
	1: "✏️",
	2: "📚",
	3: "🔖",
}

func escapeMarkdownV2(text string) string {
	specialChars := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	result := text
	for _, char := range specialChars {
		result = strings.ReplaceAll(result, char, "\\"+char)
	}
	return result
}

func (r *Renderer) heading(w io.Writer, node *ast.Heading, entering bool) {
	if entering {
		r.out(w, []byte("*"))
		r.outs(w, headings[node.Level])
		r.outs(w, " ")
		r.out(w, []byte(escapeMarkdownV2(string(node.Literal))))
	} else {
		r.out(w, []byte("*"))
		r.outs(w, "\n\n")
	}
}

func (r *Renderer) image(w io.Writer, node *ast.Image, entering bool) {
	if entering {
		// alt := node. ??
		var alt []byte
		r.outs(w, "![")
		r.out(w, alt)
	} else {
		link := node.Destination
		title := node.Title
		r.outs(w, "](")
		r.out(w, escape(link))
		if len(title) != 0 {
			r.outs(w, ` "`)
			r.out(w, title)
			r.outs(w, `"`)
		}
		r.outs(w, ")")
	}
}

func (r *Renderer) link(w io.Writer, node *ast.Link, entering bool) {
	if entering {
		r.outs(w, "[")
	} else {
		link := string(escape(node.Destination))
		title := string(node.Title)
		r.outs(w, "](")
		r.outs(w, link)
		if len(title) != 0 {
			r.outs(w, ` "`)
			r.outs(w, title)
			r.outs(w, `"`)
		}
		r.outs(w, ")")
	}
}

// RenderNode renders markdown node
func (r *Renderer) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	// Проверка дублирования для таблиц
	if table, ok := node.(*ast.Table); ok {
		if r.wasTableProcessed(table) {
			return ast.GoToNext
		}
	}

	switch node := node.(type) {
	case *ast.Text:
		//// Проверяем, не находится ли текст сразу после таблицы
		//// Если находится, добавляем пустую строку перед ним
		//parent := node.GetParent()
		//if parent != nil && !isInsideTable(parent) {
		//	prevSibling := getPreviousSibling(parent, node)
		//	if prevSibling != nil && isTable(prevSibling) {
		//		fmt.Fprint(w, "\n")
		//	}
		//}
		//
		//// Рендерим текст с правильным экранированием
		//if !isInsideTable(node) { // Не экранируем текст внутри таблицы (уже в блоке кода)
		//	fmt.Fprint(w, escapeMarkdownV2(string(n.Literal)))
		//} else {
		//	fmt.Fprint(w, string(node.Literal))
		//}

		r.text(w, node)
	case *ast.Softbreak:
		panic(fmt.Sprintf("node %T NYI", node))
	case *ast.Hardbreak:
		panic(fmt.Sprintf("node %T NYI", node))
	case *ast.Emph:
		r.surround(w, "_")
	case *ast.Strong:
		r.surround(w, "*")
	case *ast.Del:
		r.surround(w, "~")
	case *ast.BlockQuote:
		r.blockQuote(w, node)
	case *ast.Aside:
		panic(fmt.Sprintf("node %T NYI", node))
	case *ast.Link:
		r.link(w, node, entering)
	case *ast.CrossReference:
		panic(fmt.Sprintf("node %T NYI", node))
	case *ast.Citation:
		panic(fmt.Sprintf("node %T NYI", node))
	case *ast.Image:
		r.image(w, node, entering)
	case *ast.Code:
		r.code(w, node)
	case *ast.CodeBlock:
		r.codeBlock(w, node)
	case *ast.Caption:
		panic(fmt.Sprintf("node %T NYI", node))
	case *ast.CaptionFigure:
		panic(fmt.Sprintf("node %T NYI", node))
	case *ast.Document:
		// do nothing
	case *ast.Paragraph:
		r.para(w, node, entering)
	case *ast.HTMLSpan:
		r.htmlSpan(w, node)
	case *ast.HTMLBlock:
		r.htmlBlock(w, node)
	case *ast.Heading:
		r.heading(w, node, entering)
	case *ast.HorizontalRule:
	case *ast.List:
		r.list(w, node, entering)
	case *ast.ListItem:
		r.listItem(w, node, entering)
	case *ast.Table:
		r.table(w, node)
		return ast.SkipChildren // Пропускаем внутренний рендеринг таблицы

	case *ast.TableCell:
		//r.tableCell(w, node, entering)
	case *ast.TableHeader:
		//r.tableHeader(w, node, entering)
	case *ast.TableBody:
		r.tableBody(w, node)
	case *ast.TableRow:
		r.tableRow(w, node)
	case *ast.TableFooter:
		r.tableFooter(w, node)
	case *ast.Math:
		r.math(w, node)
	case *ast.MathBlock:
		r.mathBlock(w, node, entering)
	case *ast.DocumentMatter:
		panic(fmt.Sprintf("node %T NYI", node))
	case *ast.Callout:
		panic(fmt.Sprintf("node %T NYI", node))
	case *ast.Index:
		panic(fmt.Sprintf("node %T NYI", node))
	case *ast.Subscript:
		panic(fmt.Sprintf("node %T NYI", node))
	case *ast.Superscript:
		panic(fmt.Sprintf("node %T NYI", node))
	case *ast.Footnotes:
		// nothing by default; just output the list.
	default:
		panic(fmt.Sprintf("Unknown node %T", node))
	}
	return ast.GoToNext
}

// RenderHeader renders header
func (r *Renderer) RenderHeader(w io.Writer, ast ast.Node) {
	// do nothing
}

// RenderFooter renders footer
func (r *Renderer) RenderFooter(w io.Writer, ast ast.Node) {
	// do nothing
}

func (r *Renderer) math(w io.Writer, node *ast.Math) {
	r.outs(w, "` ")
	r.out(w, []byte(escapeMarkdownV2(string(node.Literal))))
	r.outs(w, " `")
}

func (r *Renderer) mathBlock(w io.Writer, node *ast.MathBlock, entering bool) {
	if r.latex.ContainsLaTeXSymbols(string(node.Literal)) {
		return
	}

	if entering {
		r.outs(w, "```\n")
		r.out(w, []byte(escapeMarkdownV2(string(node.Literal))))
	} else {
		r.outs(w, "\n```")
	}
}

func (r *Renderer) blockQuote(w io.Writer, node *ast.BlockQuote) {
	// blockquote могут быть только детьми шлюх и документа
	if _, ok := node.GetParent().(*ast.Document); !ok {
		r.outs(w, strings.TrimLeft(escapeMarkdownV2(string((node.Literal))), " "))
		return
	}

	r.outs(w, ">")
	r.outs(w, strings.TrimLeft(escapeMarkdownV2(string((node.Literal))), " "))
}
