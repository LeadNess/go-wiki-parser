package parser

import (
	"fmt"
	"testing"

	"github.com/LeadNess/go-wiki-parser/pkg/parser"
)

func TestRemovingStrong(t *testing.T) {
	strongString := "'''Литва''' , официальное название — '''Литовская Республика'''."
	processedStrongText := "Литва , официальное название — Литовская Республика."
	wikiParser := parser.NewWikiTextProcessor()
	if text, _ := wikiParser.ProcessText(strongString); text != processedStrongText {
		t.Error("Strong text does not processed")
	}
}

func TestRemovingCursive(t *testing.T) {
	cursiveString := "Пример текста про Литву - «''Lytva''» и «''Litua''»."
	processedCursiveText := "Пример текста про Литву - «Lytva» и «Litua»."
	wikiParser := parser.NewWikiTextProcessor()
	if text, _ := wikiParser.ProcessText(cursiveString); text != processedCursiveText {
		t.Error("Cursive text does not processed")
	}
}

func TestRemovingHtml(t *testing.T) {
	HTMLString := "А тут будет ссылка - <!-- КОММЕНТАРИЙ <ref name=\"ВКЛЭ2\">текст ссылки</ref>. -->."
	processedHTMLText := "А тут будет ссылка - ."
	wikiParser := parser.NewWikiTextProcessor()
	if text, _ := wikiParser.ProcessText(HTMLString); text != processedHTMLText {
		t.Error("HTML text does not processed")
		fmt.Printf("%#v\n%#v\n", text, processedHTMLText)
	}
}

func TestRemovingLists(t *testing.T) {
	textWithList := `Территория Литвы разделена на 10 уездов.
{| border=0 cellpadding=5
|
| [[Файл:Литва АТД.png|thumb|250px|Уезды Литвы]]
|
* [[Алитусский уезд]]
* [[Вильнюсский уезд]]
* [[Каунасский уезд]]
* [[Клайпедский уезд]]
* [[Мариямпольский уезд]]
* [[Паневежский уезд]]
* [[Таурагский уезд]]
* [[Тельшяйский уезд]]
* [[Утенский уезд]]
* [[Шяуляйский уезд]]
|}`
	textWithoutList := "Территория Литвы разделена на 10 уездов."
	wikiParser := parser.NewWikiTextProcessor()
	if text, _ := wikiParser.ProcessText(textWithList); text != textWithoutList {
		t.Error("Lists in text are not removed")
	}
}

func TestProcessingFigureBrackets(t *testing.T) {
	textWithFigureBrackets := `Площадь — {{число|65300}} км²{{cite web|url=https://osp.stat.gov.lt/|title=Pradžia}}.{{Государство
| Русское название = Литовская Республика
| Оригинальное название = {{lang-lt|Lietuvos Respublika}}
| Родительный падеж = Литвы
}}`
	processedText := "Площадь — 65300 км²."
	wikiParser := parser.NewWikiTextProcessor()
	if text, _ := wikiParser.ProcessText(textWithFigureBrackets); text != processedText {
		t.Error("Figure brackets in text are not processed")
	}
}

func TestRemovingInternetRefs(t *testing.T) {
	textWithInternetRefs := "Ссылка на сайт[https://osp.stat.gov.lt/]."
	processedText := "Ссылка на сайт."
	wikiParser := parser.NewWikiTextProcessor()
	if text, _ := wikiParser.ProcessText(textWithInternetRefs); text != processedText {
		t.Error("Internet refs in text are not removed")
	}
}

func TestRemovingStresses(t *testing.T) {
	textWithStresses := "Литва́ , официальное название — Лито́вская Респу́блика."
	processedText := "Литва , официальное название — Литовская Республика."
	wikiParser := parser.NewWikiTextProcessor()
	if text, _ := wikiParser.ProcessText(textWithStresses); text != processedText {
		t.Error("Stresses refs in text are not removed")
	}
}

func TestProcessingRefs(t *testing.T) {
	textWithRefs := "Член [[ООН]] с 1991 года. Входит в [[Шенгенская зона|Шенгенскую зону]] и [[Еврозона|Еврозону]]."
	processedText := "Член ООН с 1991 года. Входит в Шенгенскую зону и Еврозону."
	textRefs := []string{"ООН", "Шенгенская зона", "Еврозона"}
	wikiParser := parser.NewWikiTextProcessor()
	if text, refs := wikiParser.ProcessText(textWithRefs); text != processedText {
		t.Error("Refs in text are not processed")
	} else {
		for i := range textRefs {
			if refs[i] != textRefs[i] {
				t.Error("Refs in text are not processed")
			}
		}
	}
}
