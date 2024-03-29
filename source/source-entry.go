package source

import (
    "regexp"
    "strings"
)

var entryRgx = regexp.MustCompile(
    "^(deb(?:-src)?)(\\s+\\[[^]]*])?\\s+(\\S+)\\s+(\\S+)(\\s+.*)?$")

type repoType string

const (
    BINARY repoType = "deb"
    SOURCE repoType = "deb-src"
)

type Component string
type ComponentList []Component

type Entry string

func (e Entry) lower() string {
    return strings.ToLower(string(e))
}

func (e Entry) groups() []string {
    return entryRgx.FindStringSubmatch(e.lower())
}

//noinspection GoExportedFuncWithUnexportedType
func (e Entry) RepoType() repoType {
    if e.groups()[1] == string(SOURCE) {
        return SOURCE
    }
    return BINARY
}

func (e Entry) Params() string {
    return strings.TrimSpace(e.groups()[2])
}

func (e Entry) Location() string {
    return e.groups()[3]
}

func (e Entry) DistName() string {
    return e.groups()[4]
}

func (e Entry) Components() ComponentList {
    var components []Component
    for _, comp := range regexp.MustCompile("\\s+").Split(e.groups()[5], -1) {
        components = append(components, Component(comp))
    }
    return components
}

func (e Entry) MergeComponents(components ...Component) Entry {
    s := string(e)
    for _, comp := range components {
        if !e.Components().Contains(comp) {
            s += " " + string(comp)
        }
    }
    return Entry(s)
}

func (e Entry) CanMergeWith(entry Entry) bool {
    g1, g2 := e.groups()[1:], entry.groups()[1:]
    s1, s2 := g1[0] + " ", g2[0] + " "
    s1 += strings.Join(g1[2:4], " ")
    s2 += strings.Join(g2[2:4], " ")
    // TODO let user decide which one to keep if params are different
    return s1 == s2
}

func (cl ComponentList) Contains(comp Component) bool {
    contains := false
    for _, c := range cl {
        contains = contains || comp == c
    }
    return contains
}