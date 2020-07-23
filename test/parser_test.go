package test

import (
	"fmt"
	"testing"

	"../pkg/parser"
)

/*cursiveString := `'''Литва́''' ({{lang-lt|Lietuva}}), официальное название —
		'''Лито́вская Респу́блика''' ({{lang-lt|Lietuvos Respublika}}) — [[государство]],
		расположенное в [[Северная Европа|северной части]] [[Европа|Европы]]. Столица страны — [[Вильнюс]].`
*/
func TestRemovingStrong(t *testing.T) {
	strongString := "'''Литва́''' , официальное название — '''Лито́вская Респу́блика'''."
	processedStrongText := "Литва́ , официальное название — Лито́вская Респу́блика."
	wikiParser := parser.NewWikiTextProcessor()
	if text, _ := wikiParser.ProcessText(strongString); text != processedStrongText {
		t.Errorf("Strong text does not processed")
	}
}

func TestRemovingCursive(t *testing.T) {
	cursiveString := "Пример текста про Литву - «''Lytva''» и «''Litua''»."
	processedCursiveText := "Пример текста про Литву - «Lytva» и «Litua»."
	wikiParser := parser.NewWikiTextProcessor()
	if text, _ := wikiParser.ProcessText(cursiveString); text != processedCursiveText {
		t.Errorf("Cursive text does not processed")
	}
}

func TestRemovingHtml(t *testing.T) {
	HTMLString := "А тут будет ссылка - <ref name=\"ВКЛЭ2\">текст ссылки</ref>. И коммент <!-- КОММЕНТАРИЙ -->."
	processedHTMLText := "А тут будет ссылка - текст ссылки. И коммент ."
	wikiParser := parser.NewWikiTextProcessor()
	if text, _ := wikiParser.ProcessText(HTMLString); text != processedHTMLText {
		t.Errorf("HTML text does not processed")
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
		|}
	`
	textWithoutList := "Территория Литвы разделена на 10 уездов."
	wikiParser := parser.NewWikiTextProcessor()
	if text, _ := wikiParser.ProcessText(textWithList); text != textWithoutList {
		t.Errorf("Lists in text are not removed")
		fmt.Printf("%s\n\n%s\n", text, textWithoutList)
	}
}
