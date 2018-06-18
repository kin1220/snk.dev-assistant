package robots_test

import (
	"bytes"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"testing/quick"
	"time"

	"xojoc.pw/crawl/robots"
	"xojoc.pw/must"
)

func TestParse(t *testing.T) {
	rf := `
# robots.txt for http://www.gnu.org/

User-agent: *
Crawl-delay: 2
Disallow: /private/
Disallow: /savannah-checkouts/

Sitemap: http://www.gnu.org/sitemap.xml
`
	got, err := robots.Parse(strings.NewReader(rf))
	must.OK(err)

	want := &robots.Txt{
		CrawlDelay: map[string]int{
			"*": 2,
		},
		Disallow: map[string][]string{
			"*": []string{"/private/", "/savannah-checkouts/"},
		},
		Allow: map[string][]string{},

		Sitemaps: []string{"http://www.gnu.org/sitemap.xml"},
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("# want:\n%#v\n\n# got:\n%#v\n", want, got)
	}
}

func TestAllowed(t *testing.T) {
	rf := `
# robots.txt for http://www.wikipedia.org/ and friends
#
# Please note: There are a lot of pages on this site, and there are
# some misbehaved spiders out there that go _way_ too fast. If you're
# irresponsible, your access to the site may be blocked.
#

# advertising-related bots:
User-agent: Mediapartners-Google*
Disallow: /

# Wikipedia work bots:
User-agent: IsraBot
Disallow:

User-agent: Orthogaffe
Disallow:

# Crawlers that are kind enough to obey, but which we'd rather not have
# unless they're feeding search engines.
User-agent: UbiCrawler
Disallow: /

User-agent: DOC
Disallow: /

User-agent: Zao
Disallow: /

# Some bots are known to be trouble, particularly those designed to copy
# entire sites. Please obey robots.txt.
User-agent: sitecheck.internetseer.com
Disallow: /

User-agent: Zealbot
Disallow: /

User-agent: MSIECrawler
Disallow: /

User-agent: SiteSnagger
Disallow: /

User-agent: WebStripper
Disallow: /

User-agent: WebCopier
Disallow: /

User-agent: Fetch
Disallow: /

User-agent: Offline Explorer
Disallow: /

User-agent: Teleport
Disallow: /

User-agent: TeleportPro
Disallow: /

User-agent: WebZIP
Disallow: /

User-agent: linko
Disallow: /

User-agent: HTTrack
Disallow: /

User-agent: Microsoft.URL.Control
Disallow: /

User-agent: Xenu
Disallow: /

User-agent: larbin
Disallow: /

User-agent: libwww
Disallow: /

User-agent: ZyBORG
Disallow: /

User-agent: Download Ninja
Disallow: /

# Misbehaving: requests much too fast:
User-agent: fast
Disallow: /

#
# Sorry, wget in its recursive mode is a frequent problem.
# Please read the man page and use it properly; there is a
# --wait option you can use to set the delay between hits,
# for instance.
#
User-agent: wget
Disallow: /

#
# The 'grub' distributed client has been *very* poorly behaved.
#
User-agent: grub-client
Disallow: /

#
# Doesn't follow robots.txt anyway, but...
#
User-agent: k2spider
Disallow: /

#
# Hits many times per second, not acceptable
# http://www.nameprotect.com/botinfo.html
User-agent: NPBot
Disallow: /

# A capture bot, downloads gazillions of pages with no public benefit
# http://www.webreaper.net/
User-agent: WebReaper
Disallow: /

# Wayback Machine: defaults and whether to index user-pages
# FIXME: Complete the removal of this block, per T7582.
# User-agent: archive.org_bot
# Allow: /


#
# Friendly, low-speed bots are welcome viewing article pages, but not
# dynamically-generated pages please.
#
# Inktomi's "Slurp" can read a minimum delay between hits; if your
# bot supports such a thing using the 'Crawl-delay' or another
# instruction, please let us know.
#
# There is a special exception for API mobileview to allow dynamic
# mobile web & app views to load section content.
# These views aren't HTTP-cached but use parser cache aggressively
# and don't expose special: pages etc.
#
# Another exception is for REST API documentation, located at
# /api/rest_v1/?doc.
#
User-agent: *
Allow: /w/api.php?action=mobileview&
Allow: /w/load.php?
Allow: /api/rest_v1/?doc
Disallow: /w/
Disallow: /api/
Disallow: /trap/
#
# ar:
Disallow: /wiki/%D8%AE%D8%A7%D8%B5:Search
Disallow: /wiki/%D8%AE%D8%A7%D8%B5%3ASearch
#
# dewiki:
# T6937
# sensible deletion and meta user discussion pages:
Disallow: /wiki/Wikipedia:L%C3%B6schkandidaten/
Disallow: /wiki/Wikipedia:Löschkandidaten/
Disallow: /wiki/Wikipedia:Vandalensperrung/
Disallow: /wiki/Wikipedia:Benutzersperrung/
Disallow: /wiki/Wikipedia:Vermittlungsausschuss/
Disallow: /wiki/Wikipedia:Administratoren/Probleme/
Disallow: /wiki/Wikipedia:Adminkandidaturen/
Disallow: /wiki/Wikipedia:Qualitätssicherung/
Disallow: /wiki/Wikipedia:Qualit%C3%A4tssicherung/
# Search- and random-page
Disallow: /wiki/Spezial:Suche
Disallow: /wiki/Special:Suche
Disallow: /wiki/Spezial:Zufällige_Seite
Disallow: /wiki/Spezial:Zuf%C3%A4llige_Seite
Disallow: /wiki/Special:Zufällige_Seite
Disallow: /wiki/Special:Zuf%C3%A4llige_Seite
# 4937#5
Disallow: /wiki/Wikipedia:Vandalismusmeldung/
Disallow: /wiki/Wikipedia:Gesperrte_Lemmata/
Disallow: /wiki/Wikipedia:Löschprüfung/
Disallow: /wiki/Wikipedia:L%C3%B6schprüfung/
Disallow: /wiki/Wikipedia:Administratoren/Notizen/
Disallow: /wiki/Wikipedia:Schiedsgericht/Anfragen/
Disallow: /wiki/Wikipedia:L%C3%B6schpr%C3%BCfung/
# T14111
Disallow: /wiki/Wikipedia:Checkuser/
Disallow: /wiki/Wikipedia_Diskussion:Checkuser/
Disallow: /wiki/Wikipedia_Diskussion:Adminkandidaturen/
# T15961
Disallow: /wiki/Wikipedia:Spam-Blacklist-Log
Disallow: /wiki/Wikipedia%3ASpam-Blacklist-Log
Disallow: /wiki/Wikipedia_Diskussion:Spam-Blacklist-Log
Disallow: /wiki/Wikipedia_Diskussion%3ASpam-Blacklist-Log
#
# enwiki:
# Folks get annoyed when VfD discussions end up the number 1 google hit for
# their name. See T6776
Disallow: /wiki/Wikipedia:Articles_for_deletion/
Disallow: /wiki/Wikipedia%3AArticles_for_deletion/
Disallow: /wiki/Wikipedia:Votes_for_deletion/
Disallow: /wiki/Wikipedia%3AVotes_for_deletion/
Disallow: /wiki/Wikipedia:Pages_for_deletion/
Disallow: /wiki/Wikipedia%3APages_for_deletion/
Disallow: /wiki/Wikipedia:Miscellany_for_deletion/
Disallow: /wiki/Wikipedia%3AMiscellany_for_deletion/
Disallow: /wiki/Wikipedia:Miscellaneous_deletion/
Disallow: /wiki/Wikipedia%3AMiscellaneous_deletion/
Disallow: /wiki/Wikipedia:Copyright_problems
Disallow: /wiki/Wikipedia%3ACopyright_problems
Disallow: /wiki/Wikipedia:Protected_titles/
Disallow: /wiki/Wikipedia%3AProtected_titles/
# T15398
Disallow: /wiki/Wikipedia:WikiProject_Spam/
Disallow: /wiki/Wikipedia%3AWikiProject_Spam/
# T16075
Disallow: /wiki/MediaWiki:Spam-blacklist
Disallow: /wiki/MediaWiki%3ASpam-blacklist
Disallow: /wiki/MediaWiki_talk:Spam-blacklist
Disallow: /wiki/MediaWiki_talk%3ASpam-blacklist
# T13261
Disallow: /wiki/Wikipedia:Requests_for_arbitration/
Disallow: /wiki/Wikipedia%3ARequests_for_arbitration/
Disallow: /wiki/Wikipedia:Requests_for_comment/
Disallow: /wiki/Wikipedia%3ARequests_for_comment/
Disallow: /wiki/Wikipedia:Requests_for_adminship/
Disallow: /wiki/Wikipedia%3ARequests_for_adminship/
# T12288
Disallow: /wiki/Wikipedia_talk:Articles_for_deletion/
Disallow: /wiki/Wikipedia_talk%3AArticles_for_deletion/
Disallow: /wiki/Wikipedia_talk:Votes_for_deletion/
Disallow: /wiki/Wikipedia_talk%3AVotes_for_deletion/
Disallow: /wiki/Wikipedia_talk:Pages_for_deletion/
Disallow: /wiki/Wikipedia_talk%3APages_for_deletion/
Disallow: /wiki/Wikipedia_talk:Miscellany_for_deletion/
Disallow: /wiki/Wikipedia_talk%3AMiscellany_for_deletion/
Disallow: /wiki/Wikipedia_talk:Miscellaneous_deletion/
Disallow: /wiki/Wikipedia_talk%3AMiscellaneous_deletion/
# T16793
Disallow: /wiki/Wikipedia:Changing_username
Disallow: /wiki/Wikipedia%3AChanging_username
Disallow: /wiki/Wikipedia:Changing_username/
Disallow: /wiki/Wikipedia%3AChanging_username/
Disallow: /wiki/Wikipedia_talk:Changing_username
Disallow: /wiki/Wikipedia_talk%3AChanging_username
Disallow: /wiki/Wikipedia_talk:Changing_username/
Disallow: /wiki/Wikipedia_talk%3AChanging_username/
#
# eswiki:
# T8746
Disallow: /wiki/Wikipedia:Consultas_de_borrado/
Disallow: /wiki/Wikipedia%3AConsultas_de_borrado/
#
# fiwiki:
# T10695
Disallow: /wiki/Wikipedia:Poistettavat_sivut
Disallow: /wiki/K%C3%A4ytt%C3%A4j%C3%A4:
Disallow: /wiki/Käyttäjä:
Disallow: /wiki/Keskustelu_k%C3%A4ytt%C3%A4j%C3%A4st%C3%A4:
Disallow: /wiki/Keskustelu_käyttäjästä:
Disallow: /wiki/Wikipedia:Yll%C3%A4pit%C3%A4j%C3%A4t/
Disallow: /wiki/Wikipedia:Ylläpitäjät/
#
# frwiki:
Disallow: /wiki/Wikip%C3%A9dia:Pages_%C3%A0_supprimer/
Disallow: /wiki/Wikip%C3%A9dia:Pages_soup%C3%A7onn%C3%A9es_de_violation_de_copyright/
#
# hewiki:
Disallow: /wiki/%D7%9E%D7%99%D7%95%D7%97%D7%93:Search
Disallow: /wiki/%D7%9E%D7%99%D7%95%D7%97%D7%93%3ASearch
#T11517
Disallow: /wiki/ויקיפדיה:רשימת_מועמדים_למחיקה/
Disallow: /wiki/ויקיפדיה%3Aרשימת_מועמדים_למחיקה/
Disallow: /wiki/%D7%95%D7%99%D7%A7%D7%99%D7%A4%D7%93%D7%99%D7%94:%D7%A8%D7%A9%D7%99%D7%9E%D7%AA_%D7%9E%D7%95%D7%A2%D7%9E%D7%93%D7%99%D7%9D_%D7%9C%D7%9E%D7%97%D7%99%D7%A7%D7%94/
Disallow: /wiki/%D7%95%D7%99%D7%A7%D7%99%D7%A4%D7%93%D7%99%D7%94%3A%D7%A8%D7%A9%D7%99%D7%9E%D7%AA_%D7%9E%D7%95%D7%A2%D7%9E%D7%93%D7%99%D7%9D_%D7%9C%D7%9E%D7%97%D7%99%D7%A7%D7%94/
Disallow: /wiki/ויקיפדיה:ערכים_לא_קיימים_ומוגנים
Disallow: /wiki/ויקיפדיה%3Aערכים_לא_קיימים_ומוגנים
Disallow: /wiki/%D7%95%D7%99%D7%A7%D7%99%D7%A4%D7%93%D7%99%D7%94:%D7%A2%D7%A8%D7%9B%D7%99%D7%9D_%D7%9C%D7%90_%D7%A7%D7%99%D7%99%D7%9E%D7%99%D7%9D_%D7%95%D7%9E%D7%95%D7%92%D7%A0%D7%99%D7%9D
Disallow: /wiki/%D7%95%D7%99%D7%A7%D7%99%D7%A4%D7%93%D7%99%D7%94%3A%D7%A2%D7%A8%D7%9B%D7%99%D7%9D_%D7%9C%D7%90_%D7%A7%D7%99%D7%99%D7%9E%D7%99%D7%9D_%D7%95%D7%9E%D7%95%D7%92%D7%A0%D7%99%D7%9D
Disallow: /wiki/ויקיפדיה:דפים_לא_קיימים_ומוגנים
Disallow: /wiki/ויקיפדיה%3Aדפים_לא_קיימים_ומוגנים
Disallow: /wiki/%D7%95%D7%99%D7%A7%D7%99%D7%A4%D7%93%D7%99%D7%94:%D7%93%D7%A4%D7%99%D7%9D_%D7%9C%D7%90_%D7%A7%D7%99%D7%99%D7%9E%D7%99%D7%9D_%D7%95%D7%9E%D7%95%D7%92%D7%A0%D7%99%D7%9D
Disallow: /wiki/%D7%95%D7%99%D7%A7%D7%99%D7%A4%D7%93%D7%99%D7%94%3A%D7%93%D7%A4%D7%99%D7%9D_%D7%9C%D7%90_%D7%A7%D7%99%D7%99%D7%9E%D7%99%D7%9D_%D7%95%D7%9E%D7%95%D7%92%D7%A0%D7%99%D7%9D
#
# huwiki:
Disallow: /wiki/Speci%C3%A1lis:Search
Disallow: /wiki/Speci%C3%A1lis%3ASearch
#
# itwiki:
# T7545
Disallow: /wiki/Wikipedia:Pagine_da_cancellare
Disallow: /wiki/Wikipedia%3APagine_da_cancellare
Disallow: /wiki/Wikipedia:Utenti_problematici
Disallow: /wiki/Wikipedia%3AUtenti_problematici
Disallow: /wiki/Wikipedia:Vandalismi_in_corso
Disallow: /wiki/Wikipedia%3AVandalismi_in_corso
Disallow: /wiki/Wikipedia:Amministratori
Disallow: /wiki/Wikipedia%3AAmministratori
Disallow: /wiki/Wikipedia:Proposte_di_cancellazione_semplificata
Disallow: /wiki/Wikipedia%3AProposte_di_cancellazione_semplificata
Disallow: /wiki/Categoria:Da_cancellare_subito
Disallow: /wiki/Categoria%3ADa_cancellare_subito
Disallow: /wiki/Wikipedia:Sospette_violazioni_di_copyright
Disallow: /wiki/Wikipedia%3ASospette_violazioni_di_copyright
Disallow: /wiki/Categoria:Da_controllare_per_copyright
Disallow: /wiki/Categoria%3ADa_controllare_per_copyright
Disallow: /wiki/Progetto:Rimozione_contributi_sospetti
Disallow: /wiki/Progetto%3ARimozione_contributi_sospetti
Disallow: /wiki/Categoria:Da_cancellare_subito_per_violazione_integrale_copyright
Disallow: /wiki/Categoria%3ADa_cancellare_subito_per_violazione_integrale_copyright
Disallow: /wiki/Progetto:Cococo
Disallow: /wiki/Progetto%3ACococo
Disallow: /wiki/Discussioni_progetto:Cococo
Disallow: /wiki/Discussioni_progetto%3ACococo
#
# jawiki
Disallow: /wiki/%E7%89%B9%E5%88%A5:Search
Disallow: /wiki/%E7%89%B9%E5%88%A5%3ASearch
# T7239
Disallow: /wiki/Wikipedia:%E5%89%8A%E9%99%A4%E4%BE%9D%E9%A0%BC/
Disallow: /wiki/Wikipedia%3A%E5%89%8A%E9%99%A4%E4%BE%9D%E9%A0%BC/
Disallow: /wiki/Wikipedia:%E5%88%A9%E7%94%A8%E8%80%85%E3%83%9A%E3%83%BC%E3%82%B8%E3%81%AE%E5%89%8A%E9%99%A4%E4%BE%9D%E9%A0%BC
Disallow: /wiki/Wikipedia%3A%E5%88%A9%E7%94%A8%E8%80%85%E3%83%9A%E3%83%BC%E3%82%B8%E3%81%AE%E5%89%8A%E9%99%A4%E4%BE%9D%E9%A0%BC
# nowiki
# T13432
Disallow: /wiki/Bruker:
Disallow: /wiki/Bruker%3A
Disallow: /wiki/Brukerdiskusjon
Disallow: /wiki/Wikipedia:Administratorer
Disallow: /wiki/Wikipedia%3AAdministratorer
Disallow: /wiki/Wikipedia-diskusjon:Administratorer
Disallow: /wiki/Wikipedia-diskusjon%3AAdministratorer
Disallow: /wiki/Wikipedia:Sletting
Disallow: /wiki/Wikipedia%3ASletting
Disallow: /wiki/Wikipedia-diskusjon:Sletting
Disallow: /wiki/Wikipedia-diskusjon%3ASletting
Disallow: /wiki/Spesial:
Disallow: /wiki/Spesial%3A
#
# plwiki
# T10067
Disallow: /wiki/Wikipedia:Strony_do_usuni%C4%99cia
Disallow: /wiki/Wikipedia%3AStrony_do_usuni%C4%99cia
Disallow: /wiki/Wikipedia:Do_usuni%C4%99cia
Disallow: /wiki/Wikipedia%3ADo_usuni%C4%99cia
Disallow: /wiki/Wikipedia:SDU/
Disallow: /wiki/Wikipedia%3ASDU/
Disallow: /wiki/Wikipedia:Strony_podejrzane_o_naruszenie_praw_autorskich
Disallow: /wiki/Wikipedia%3AStrony_podejrzane_o_naruszenie_praw_autorskich
#
# ptwiki:
# T7394
Disallow: /wiki/Wikipedia:Páginas_para_eliminar/
Disallow: /wiki/Wikipedia:P%C3%A1ginas_para_eliminar/
Disallow: /wiki/Wikipedia%3AP%C3%A1ginas_para_eliminar/
Disallow: /wiki/Wikipedia_Discussão:Páginas_para_eliminar/
Disallow: /wiki/Wikipedia_Discuss%C3%A3o:P%C3%A1ginas_para_eliminar/
Disallow: /wiki/Wikipedia_Discuss%C3%A3o%3AP%C3%A1ginas_para_eliminar/
#
# rowiki:
# T14546
Disallow: /wiki/Wikipedia:Pagini_de_%C5%9Fters
Disallow: /wiki/Wikipedia%3APagini_de_%C5%9Fters
Disallow: /wiki/Discu%C5%A3ie_Wikipedia:Pagini_de_%C5%9Fters
Disallow: /wiki/Discu%C5%A3ie_Wikipedia%3APagini_de_%C5%9Fters
#
# ruwiki:
Disallow: /wiki/%D0%A1%D0%BF%D0%B5%D1%86%D0%B8%D0%B0%D0%BB%D1%8C%D0%BD%D1%8B%D0%B5:Search
Disallow: /wiki/%D0%A1%D0%BF%D0%B5%D1%86%D0%B8%D0%B0%D0%BB%D1%8C%D0%BD%D1%8B%D0%B5%3ASearch
#
# svwiki:
# T12229
Disallow: /wiki/Wikipedia%3ASidor_f%C3%B6reslagna_f%C3%B6r_radering
Disallow: /wiki/Wikipedia:Sidor_f%C3%B6reslagna_f%C3%B6r_radering
Disallow: /wiki/Wikipedia:Sidor_föreslagna_för_radering
Disallow: /wiki/Användare
Disallow: /wiki/Anv%C3%A4ndare
Disallow: /wiki/Användardiskussion
Disallow: /wiki/Anv%C3%A4ndardiskussion
Disallow: /wiki/Wikipedia:Skyddade_sidnamn
Disallow: /wiki/Wikipedia%3ASkyddade_sidnamn
# T13291
Disallow: /wiki/Wikipedia:Sidor_som_bör_raderas
Disallow: /wiki/Wikipedia:Sidor_som_b%C3%B6r_raderas
Disallow: /wiki/Wikipedia%3ASidor_som_b%C3%B6r_raderas
#
# zhwiki:
# T7104
Disallow: /wiki/Wikipedia:删除投票/侵权
Disallow: /wiki/Wikipedia:%E5%88%A0%E9%99%A4%E6%8A%95%E7%A5%A8/%E4%BE%B5%E6%9D%83
Disallow: /wiki/Wikipedia:删除投票和请求
Disallow: /wiki/Wikipedia:%E5%88%A0%E9%99%A4%E6%8A%95%E7%A5%A8%E5%92%8C%E8%AF%B7%E6%B1%82
Disallow: /wiki/Category:快速删除候选
Disallow: /wiki/Category:%E5%BF%AB%E9%80%9F%E5%88%A0%E9%99%A4%E5%80%99%E9%80%89
Disallow: /wiki/Category:维基百科需要翻译的文章
Disallow: /wiki/Category:%E7%BB%B4%E5%9F%BA%E7%99%BE%E7%A7%91%E9%9C%80%E8%A6%81%E7%BF%BB%E8%AF%91%E7%9A%84%E6%96%87%E7%AB%A0
#
# sister projects
#
# enwikinews:
# T7340
Disallow: /wiki/Portal:Prepared_stories/
Disallow: /wiki/Portal%3APrepared_stories/
#
# itwikinews
# T11138
Disallow: /wiki/Wikinotizie:Richieste_di_cancellazione
Disallow: /wiki/Wikinotizie:Sospette_violazioni_di_copyright
Disallow: /wiki/Categoria:Da_cancellare_subito
Disallow: /wiki/Categoria:Da_cancellare_subito_per_violazione_integrale_copyright
Disallow: /wiki/Wikinotizie:Storie_in_preparazione
#
# enwikiquote:
# T17095
Disallow: /wiki/Wikiquote:Votes_for_deletion/
Disallow: /wiki/Wikiquote%3AVotes_for_deletion/
Disallow: /wiki/Wikiquote_talk:Votes_for_deletion/
Disallow: /wiki/Wikiquote_talk%3AVotes_for_deletion/
Disallow: /wiki/Wikiquote:Votes_for_deletion_archive/
Disallow: /wiki/Wikiquote%3AVotes_for_deletion_archive/
Disallow: /wiki/Wikiquote_talk:Votes_for_deletion_archive/
Disallow: /wiki/Wikiquote_talk%3AVotes_for_deletion_archive/
#
# enwikibooks
Disallow: /wiki/Wikibooks:Votes_for_deletion
#
# working...
Disallow: /wiki/Fundraising_2007/comments
#
Disallow: /wiki/Special:Maintenance
# Do not show banner content or record hits
Disallow: /wiki/Special:BannerLoader
Disallow: /wiki/Special:RecordImpression
#
#
#----------------------------------------------------------#
#
#
#
 # <!-- Please do not remove the space at the start of this line, it breaks the rendering.  http://www.robotstxt.org/orig.html says spaces before comments are OK. --><syntaxhighlight lang="robots">
#
# Localisable part of robots.txt for en.wikipedia.org
#
# Edit at https://en.wikipedia.org/w/index.php?title=MediaWiki:Robots.txt&action=edit
# Don't add newlines here. All rules set here are active for every user-agent.
#
# Please check any changes using a syntax validator such as http://tool.motoricerca.info/robots-checker.phtml
# Enter https://en.wikipedia.org/robots.txt as the URL to check.
#
# https://bugzilla.wikimedia.org/show_bug.cgi?id=14075
Disallow: /wiki/MediaWiki:Spam-blacklist
Disallow: /wiki/MediaWiki%3ASpam-blacklist
Disallow: /wiki/MediaWiki_talk:Spam-blacklist
Disallow: /wiki/MediaWiki_talk%3ASpam-blacklist
Disallow: /wiki/Wikipedia:WikiProject_Spam
Disallow: /wiki/Wikipedia_talk:WikiProject_Spam
#
# Folks get annoyed when XfD discussions end up the number 1 google hit for
# their name. 
# https://phabricator.wikimedia.org/T16075
Disallow: /wiki/Wikipedia:Articles_for_deletion
Disallow: /wiki/Wikipedia%3AArticles_for_deletion
Disallow: /wiki/Wikipedia:Votes_for_deletion
Disallow: /wiki/Wikipedia%3AVotes_for_deletion
Disallow: /wiki/Wikipedia:Pages_for_deletion
Disallow: /wiki/Wikipedia%3APages_for_deletion
Disallow: /wiki/Wikipedia:Miscellany_for_deletion
Disallow: /wiki/Wikipedia%3AMiscellany_for_deletion
Disallow: /wiki/Wikipedia:Miscellaneous_deletion
Disallow: /wiki/Wikipedia%3AMiscellaneous_deletion
Disallow: /wiki/Wikipedia:Categories_for_discussion
Disallow: /wiki/Wikipedia%3ACategories_for_discussion
Disallow: /wiki/Wikipedia:Templates_for_deletion
Disallow: /wiki/Wikipedia%3ATemplates_for_deletion
Disallow: /wiki/Wikipedia:Redirects_for_discussion
Disallow: /wiki/Wikipedia%3ARedirects_for_discussion
Disallow: /wiki/Wikipedia:Deletion_review
Disallow: /wiki/Wikipedia%3ADeletion_review
Disallow: /wiki/Wikipedia:WikiProject_Deletion_sorting
Disallow: /wiki/Wikipedia%3AWikiProject_Deletion_sorting
Disallow: /wiki/Wikipedia:Files_for_deletion
Disallow: /wiki/Wikipedia%3AFiles_for_deletion
Disallow: /wiki/Wikipedia:Files_for_discussion
Disallow: /wiki/Wikipedia%3AFiles_for_discussion
Disallow: /wiki/Wikipedia:Possibly_unfree_files
Disallow: /wiki/Wikipedia%3APossibly_unfree_files
#
# https://phabricator.wikimedia.org/T12288
Disallow: /wiki/Wikipedia_talk:Articles_for_deletion
Disallow: /wiki/Wikipedia_talk%3AArticles_for_deletion
Disallow: /wiki/Wikipedia_talk:Votes_for_deletion
Disallow: /wiki/Wikipedia_talk%3AVotes_for_deletion
Disallow: /wiki/Wikipedia_talk:Pages_for_deletion
Disallow: /wiki/Wikipedia_talk%3APages_for_deletion
Disallow: /wiki/Wikipedia_talk:Miscellany_for_deletion
Disallow: /wiki/Wikipedia_talk%3AMiscellany_for_deletion
Disallow: /wiki/Wikipedia_talk:Miscellaneous_deletion
Disallow: /wiki/Wikipedia_talk%3AMiscellaneous_deletion
Disallow: /wiki/Wikipedia_talk:Templates_for_deletion
Disallow: /wiki/Wikipedia_talk%3ATemplates_for_deletion
Disallow: /wiki/Wikipedia_talk:Categories_for_discussion
Disallow: /wiki/Wikipedia_talk%3ACategories_for_discussion
Disallow: /wiki/Wikipedia_talk:Deletion_review
Disallow: /wiki/Wikipedia_talk%3ADeletion_review
Disallow: /wiki/Wikipedia_talk:WikiProject_Deletion_sorting
Disallow: /wiki/Wikipedia_talk%3AWikiProject_Deletion_sorting
Disallow: /wiki/Wikipedia_talk:Files_for_deletion
Disallow: /wiki/Wikipedia_talk%3AFiles_for_deletion
Disallow: /wiki/Wikipedia_talk:Files_for_discussion
Disallow: /wiki/Wikipedia_talk%3AFiles_for_discussion
Disallow: /wiki/Wikipedia_talk:Possibly_unfree_files
Disallow: /wiki/Wikipedia_talk%3APossibly_unfree_files
#
Disallow: /wiki/Wikipedia:Copyright_problems
Disallow: /wiki/Wikipedia%3ACopyright_problems
Disallow: /wiki/Wikipedia_talk:Copyright_problems
Disallow: /wiki/Wikipedia_talk%3ACopyright_problems
Disallow: /wiki/Wikipedia:Suspected_copyright_violations
Disallow: /wiki/Wikipedia%3ASuspected_copyright_violations
Disallow: /wiki/Wikipedia_talk:Suspected_copyright_violations
Disallow: /wiki/Wikipedia_talk%3ASuspected_copyright_violations
Disallow: /wiki/Wikipedia:Contributor_copyright_investigations
Disallow: /wiki/Wikipedia%3AContributor_copyright_investigations
Disallow: /wiki/Wikipedia:Contributor_copyright_investigations
Disallow: /wiki/Wikipedia%3AContributor_copyright_investigations
Disallow: /wiki/Wikipedia_talk:Contributor_copyright_investigations
Disallow: /wiki/Wikipedia_talk%3AContributor_copyright_investigations
Disallow: /wiki/Wikipedia_talk:Contributor_copyright_investigations
Disallow: /wiki/Wikipedia_talk%3AContributor_copyright_investigations
Disallow: /wiki/Wikipedia:Protected_titles
Disallow: /wiki/Wikipedia%3AProtected_titles
Disallow: /wiki/Wikipedia_talk:Protected_titles
Disallow: /wiki/Wikipedia_talk%3AProtected_titles
Disallow: /wiki/Wikipedia:Articles_for_creation
Disallow: /wiki/Wikipedia%3AArticles_for_creation
Disallow: /wiki/Wikipedia_talk:Articles_for_creation
Disallow: /wiki/Wikipedia_talk%3AArticles_for_creation
Disallow: /wiki/Wikipedia_talk:Article_wizard
Disallow: /wiki/Wikipedia_talk%3AArticle_wizard
#
# https://phabricator.wikimedia.org/T13261
Disallow: /wiki/Wikipedia:Requests_for_arbitration
Disallow: /wiki/Wikipedia%3ARequests_for_arbitration
Disallow: /wiki/Wikipedia_talk:Requests_for_arbitration
Disallow: /wiki/Wikipedia_talk%3ARequests_for_arbitration
Disallow: /wiki/Wikipedia:Requests_for_comment
Disallow: /wiki/Wikipedia%3ARequests_for_comment
Disallow: /wiki/Wikipedia_talk:Requests_for_comment
Disallow: /wiki/Wikipedia_talk%3ARequests_for_comment
Disallow: /wiki/Wikipedia:Requests_for_adminship
Disallow: /wiki/Wikipedia%3ARequests_for_adminship
Disallow: /wiki/Wikipedia_talk:Requests_for_adminship
Disallow: /wiki/Wikipedia_talk%3ARequests_for_adminship
#
# https://phabricator.wikimedia.org/T14111
Disallow: /wiki/Wikipedia:Requests_for_checkuser
Disallow: /wiki/Wikipedia%3ARequests_for_checkuser
Disallow: /wiki/Wikipedia_talk:Requests_for_checkuser
Disallow: /wiki/Wikipedia_talk%3ARequests_for_checkuser
#
# https://phabricator.wikimedia.org/T15398
Disallow: /wiki/Wikipedia:WikiProject_Spam
Disallow: /wiki/Wikipedia%3AWikiProject_Spam
#
# https://phabricator.wikimedia.org/T16793
Disallow: /wiki/Wikipedia:Changing_username
Disallow: /wiki/Wikipedia%3AChanging_username
Disallow: /wiki/Wikipedia:Changing_username
Disallow: /wiki/Wikipedia%3AChanging_username
Disallow: /wiki/Wikipedia_talk:Changing_username
Disallow: /wiki/Wikipedia_talk%3AChanging_username
Disallow: /wiki/Wikipedia_talk:Changing_username
Disallow: /wiki/Wikipedia_talk%3AChanging_username
#
Disallow: /wiki/Wikipedia:Administrators%27_noticeboard
Disallow: /wiki/Wikipedia%3AAdministrators%27_noticeboard
Disallow: /wiki/Wikipedia_talk:Administrators%27_noticeboard
Disallow: /wiki/Wikipedia_talk%3AAdministrators%27_noticeboard
Disallow: /wiki/Wikipedia:Community_sanction_noticeboard
Disallow: /wiki/Wikipedia%3ACommunity_sanction_noticeboard
Disallow: /wiki/Wikipedia_talk:Community_sanction_noticeboard
Disallow: /wiki/Wikipedia_talk%3ACommunity_sanction_noticeboard
Disallow: /wiki/Wikipedia:Bureaucrats%27_noticeboard
Disallow: /wiki/Wikipedia%3ABureaucrats%27_noticeboard
Disallow: /wiki/Wikipedia_talk:Bureaucrats%27_noticeboard
Disallow: /wiki/Wikipedia_talk%3ABureaucrats%27_noticeboard
#
Disallow: /wiki/Wikipedia:Sockpuppet_investigations
Disallow: /wiki/Wikipedia%3ASockpuppet_investigations
Disallow: /wiki/Wikipedia_talk:Sockpuppet_investigations
Disallow: /wiki/Wikipedia_talk%3ASockpuppet_investigations
#
Disallow: /wiki/Wikipedia:Neutral_point_of_view/Noticeboard
Disallow: /wiki/Wikipedia%3ANeutral_point_of_view/Noticeboard
Disallow: /wiki/Wikipedia_talk:Neutral_point_of_view/Noticeboard
Disallow: /wiki/Wikipedia_talk%3ANeutral_point_of_view/Noticeboard
#
Disallow: /wiki/Wikipedia:No_original_research/noticeboard
Disallow: /wiki/Wikipedia%3ANo_original_research/noticeboard
Disallow: /wiki/Wikipedia_talk:No_original_research/noticeboard
Disallow: /wiki/Wikipedia_talk%3ANo_original_research/noticeboard
#
Disallow: /wiki/Wikipedia:Fringe_theories/Noticeboard
Disallow: /wiki/Wikipedia%3AFringe_theories/Noticeboard
Disallow: /wiki/Wikipedia_talk:Fringe_theories/Noticeboard
Disallow: /wiki/Wikipedia_talk%3AFringe_theories/Noticeboard
#
Disallow: /wiki/Wikipedia:Conflict_of_interest/Noticeboard
Disallow: /wiki/Wikipedia%3AConflict_of_interest/Noticeboard
Disallow: /wiki/Wikipedia_talk:Conflict_of_interest/Noticeboard
Disallow: /wiki/Wikipedia_talk%3AConflict_of_interest/Noticeboard
#
Disallow: /wiki/Wikipedia:Long-term_abuse
Disallow: /wiki/Wikipedia%3ALong-term_abuse
Disallow: /wiki/Wikipedia_talk:Long-term_abuse
Disallow: /wiki/Wikipedia_talk%3ALong-term_abuse
Disallow: /wiki/Wikipedia:Long_term_abuse
Disallow: /wiki/Wikipedia%3ALong_term_abuse
Disallow: /wiki/Wikipedia_talk:Long_term_abuse
Disallow: /wiki/Wikipedia_talk%3ALong_term_abuse
#
Disallow: /wiki/Wikipedia:Wikiquette_assistance
Disallow: /wiki/Wikipedia%3AWikiquette_assistance
#
Disallow: /wiki/Wikipedia:Abuse_reports
Disallow: /wiki/Wikipedia%3AAbuse_reports
Disallow: /wiki/Wikipedia_talk:Abuse_reports
Disallow: /wiki/Wikipedia_talk%3AAbuse_reports
Disallow: /wiki/Wikipedia:Abuse_response
Disallow: /wiki/Wikipedia%3AAbuse_response
Disallow: /wiki/Wikipedia_talk:Abuse_response
Disallow: /wiki/Wikipedia_talk%3AAbuse_response
#
Disallow: /wiki/Wikipedia:Reliable_sources/Noticeboard
Disallow: /wiki/Wikipedia%3AReliable_sources/Noticeboard
Disallow: /wiki/Wikipedia_talk:Reliable_sources/Noticeboard
Disallow: /wiki/Wikipedia_talk%3AReliable_sources/Noticeboard
#
Disallow: /wiki/Wikipedia:Suspected_sock_puppets
Disallow: /wiki/Wikipedia%3ASuspected_sock_puppets
Disallow: /wiki/Wikipedia_talk:Suspected_sock_puppets
Disallow: /wiki/Wikipedia_talk%3ASuspected_sock_puppets
#
Disallow: /wiki/Wikipedia:Biographies_of_living_persons/Noticeboard
Disallow: /wiki/Wikipedia%3ABiographies_of_living_persons/Noticeboard
Disallow: /wiki/Wikipedia_talk:Biographies_of_living_persons/Noticeboard
Disallow: /wiki/Wikipedia_talk%3ABiographies_of_living_persons/Noticeboard
#
Disallow: /wiki/Wikipedia:Content_noticeboard
Disallow: /wiki/Wikipedia%3AContent_noticeboard
Disallow: /wiki/Wikipedia_talk:Content_noticeboard
Disallow: /wiki/Wikipedia_talk%3AContent_noticeboard
#
Disallow: /wiki/Template:Editnotices
Disallow: /wiki/Template%3AEditnotices
#
Disallow: /wiki/Wikipedia:Arbitration
Disallow: /wiki/Wikipedia%3AArbitration
Disallow: /wiki/Wikipedia_talk:Arbitration
Disallow: /wiki/Wikipedia_talk%3AArbitration
#
Disallow: /wiki/Wikipedia:Arbitration_Committee
Disallow: /wiki/Wikipedia%3AArbitration_Committee
Disallow: /wiki/Wikipedia_talk:Arbitration_Committee
Disallow: /wiki/Wikipedia_talk%3AArbitration_Committee
#
Disallow: /wiki/Wikipedia:Arbitration_Committee_Elections
Disallow: /wiki/Wikipedia%3AArbitration_Committee_Elections
Disallow: /wiki/Wikipedia_talk:Arbitration_Committee_Elections
Disallow: /wiki/Wikipedia_talk%3AArbitration_Committee_Elections
#
Disallow: /wiki/Wikipedia:Mediation_Committee
Disallow: /wiki/Wikipedia%3AMediation_Committee
Disallow: /wiki/Wikipedia_talk:Mediation_Committee
Disallow: /wiki/Wikipedia_talk%3AMediation_Committee
#
Disallow: /wiki/Wikipedia:Mediation_Cabal/Cases
Disallow: /wiki/Wikipedia%3AMediation_Cabal/Cases
#
Disallow: /wiki/Wikipedia:Requests_for_bureaucratship
Disallow: /wiki/Wikipedia%3ARequests_for_bureaucratship
Disallow: /wiki/Wikipedia_talk:Requests_for_bureaucratship
Disallow: /wiki/Wikipedia_talk%3ARequests_for_bureaucratship
#
Disallow: /wiki/Wikipedia:Administrator_review
Disallow: /wiki/Wikipedia%3AAdministrator_review
Disallow: /wiki/Wikipedia_talk:Administrator_review
Disallow: /wiki/Wikipedia_talk%3AAdministrator_review
#
Disallow: /wiki/Wikipedia:Editor_review
Disallow: /wiki/Wikipedia%3AEditor_review
Disallow: /wiki/Wikipedia_talk:Editor_review
Disallow: /wiki/Wikipedia_talk%3AEditor_review
#
Disallow: /wiki/Wikipedia:Article_Incubator
Disallow: /wiki/Wikipedia%3AArticle_Incubator
Disallow: /wiki/Wikipedia_talk:Article_Incubator
Disallow: /wiki/Wikipedia_talk%3AArticle_Incubator
#
Disallow: /wiki/Category:Noindexed_pages
Disallow: /wiki/Category%3ANoindexed_pages
#
# </syntaxhighlight>
`
	txt, err := robots.Parse(strings.NewReader(rf))
	must.OK(err)

	if !txt.Allowed("IsraBot", "/") {
		t.Fatal("expected true")
	}
	if txt.Allowed("DOC", "/") {
		t.Fatal("expected false")
	}
	if !txt.Allowed("anything", "/w/load.php?") {
		t.Fatal("expected true")
	}
	if txt.Allowed("anything", "/w/load.p") {
		t.Fatal("expected false")
	}
}

var cfg = &quick.Config{MaxCount: 100000, Rand: rand.New(rand.NewSource(time.Now().UTC().UnixNano()))}

func TestQuick(t *testing.T) {
	f := func(b []byte) bool {
		_, err := robots.Parse(bytes.NewReader(b))
		if err != nil {
			return false
		}
		return true
	}

	if err := quick.Check(f, cfg); err != nil {
		t.Error(err)
	}
}
