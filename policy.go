package halalrandomstrings

import (
	"errors"
	"fmt"
	"strings"
)

type blockedNode struct {
	children map[string]*blockedNode
	terminal bool
}

func (e *Engine) buildBlocked(expressions []string) error {
	e.blockedSet = make(map[string]struct{})
	e.blockedTrie = &blockedNode{children: make(map[string]*blockedNode)}
	for _, expression := range expressions {
		normalized := normalizeTokenString(expression)
		if normalized == "" {
			return errors.New("empty blocked expression")
		}
		e.blockedSet[normalized] = struct{}{}
		node := e.blockedTrie
		for _, token := range strings.Split(normalized, "-") {
			if node.children[token] == nil {
				node.children[token] = &blockedNode{children: make(map[string]*blockedNode)}
			}
			node = node.children[token]
		}
		node.terminal = true
	}
	return nil
}

func (e *Engine) buildIndexes() error {
	p := e.words.Policy
	if p.Version != 1 || p.MaxWords != 64 || p.MaxRepeat != 1000000 || len(p.CategoryGroups) != len(e.words.Categories) || len(e.words.Blocked) == 0 {
		return errors.New("missing or unsupported mixing policy")
	}
	if strings.Join(p.AllowedSeparators, "") != "-_./" {
		return errors.New("unexpected separator policy")
	}
	if err := e.buildBlocked(e.words.Blocked); err != nil {
		return err
	}
	e.normWords = make(map[string][]wordToken)
	e.catBuckets = make(map[string]wordBuckets)
	e.groupBuckets = make(map[string]wordBuckets)
	e.groupPhrases = make(map[string]map[string]struct{})
	e.maxPhraseTokens = make(map[string]int)
	owners := make(map[string]string)
	for _, cat := range e.GetCategories() {
		group := p.CategoryGroups[cat]
		switch group {
		case "divine", "prophets", "names", "islamic", "general":
		default:
			return fmt.Errorf("unclassified category %q", cat)
		}
		if group == "divine" && cat != "asma_allah" || group == "prophets" && cat != "prophets" || group == "names" && cat != "muslim_names_male" && cat != "muslim_names_female" {
			return fmt.Errorf("invalid isolated group for %q", cat)
		}
		if _, disabled := p.DisabledCategories[cat]; disabled {
			return fmt.Errorf("disabled category %q present", cat)
		}
		if len(e.words.Categories[cat]) == 0 {
			return fmt.Errorf("empty category %q", cat)
		}
		if e.groupPhrases[group] == nil {
			e.groupPhrases[group] = make(map[string]struct{})
			e.groupBuckets[group] = make(wordBuckets, p.MaxWords+1)
		}
		e.catBuckets[cat] = make(wordBuckets, p.MaxWords+1)
		seen := make(map[string]bool)
		for _, s := range e.words.Categories[cat] {
			for _, r := range s {
				if r > 127 {
					return fmt.Errorf("non-ASCII dictionary entry %q", s)
				}
			}
			n := normalizeTokenString(s)
			count := strings.Count(n, "-") + 1
			if n == "" || count > p.MaxWords || seen[n] || !e.passesBlocklist(n) {
				return fmt.Errorf("invalid, duplicate or blocked entry %q in %s", s, cat)
			}
			for _, r := range n {
				if !(r >= 'a' && r <= 'z' || r == '-') {
					return fmt.Errorf("nonalphabetic entry %q", s)
				}
			}
			if prev, ok := owners[n]; ok && prev != group {
				return fmt.Errorf("ambiguous expression %q in %s and %s", n, prev, group)
			}
			owners[n] = group
			seen[n] = true
			wt := wordToken{n, count}
			e.normWords[cat] = append(e.normWords[cat], wt)
			e.catBuckets[cat][count] = append(e.catBuckets[cat][count], wt)
			if _, exists := e.groupPhrases[group][n]; !exists {
				e.groupBuckets[group][count] = append(e.groupBuckets[group][count], wt)
			}
			e.groupPhrases[group][n] = struct{}{}
			if count > e.maxPhraseTokens[group] {
				e.maxPhraseTokens[group] = count
			}
		}
	}
	for _, group := range []string{"divine", "prophets", "names", "islamic", "general"} {
		if len(e.groupPhrases[group]) == 0 {
			return fmt.Errorf("missing group %s", group)
		}
	}
	if p.CategoryGroups["asma_allah"] != "divine" || p.CategoryGroups["prophets"] != "prophets" || p.CategoryGroups["muslim_names_male"] != "names" || p.CategoryGroups["muslim_names_female"] != "names" {
		return errors.New("protected category assigned to wrong group")
	}
	for _, rule := range e.words.Rules {
		if len(rule.Pattern) == 0 {
			return errors.New("empty rule")
		}
		group := ""
		placeholders := make([]string, len(rule.Pattern))
		for i, cat := range rule.Pattern {
			g, ok := p.CategoryGroups[cat]
			if !ok || group != "" && g != group {
				return fmt.Errorf("incompatible rule %v", rule.Pattern)
			}
			group = g
			placeholders[i] = "{" + cat + "}"
		}
		if rule.Template != strings.Join(placeholders, "-") {
			return fmt.Errorf("unsupported template %q", rule.Template)
		}
	}
	return nil
}
