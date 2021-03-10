// Grammar for the vivian row formatting language
package vivian

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

var g = &grammar{
	rules: []*rule{
		{
			name: "Root",
			pos:  position{line: 7, col: 1, offset: 71},
			expr: &choiceExpr{
				pos: position{line: 7, col: 9, offset: 79},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 7, col: 9, offset: 79},
						run: (*parser).callonRoot2,
						expr: &seqExpr{
							pos: position{line: 7, col: 9, offset: 79},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 7, col: 9, offset: 79},
									label: "exprList",
									expr: &oneOrMoreExpr{
										pos: position{line: 7, col: 18, offset: 88},
										expr: &ruleRefExpr{
											pos:  position{line: 7, col: 18, offset: 88},
											name: "Expr",
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 7, col: 24, offset: 94},
									name: "EOF",
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 9, col: 5, offset: 142},
						run: (*parser).callonRoot8,
						expr: &seqExpr{
							pos: position{line: 9, col: 5, offset: 142},
							exprs: []interface{}{
								&zeroOrMoreExpr{
									pos: position{line: 9, col: 5, offset: 142},
									expr: &ruleRefExpr{
										pos:  position{line: 9, col: 5, offset: 142},
										name: "Expr",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 9, col: 11, offset: 148},
									name: "CloseBrace",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Expr",
			pos:  position{line: 13, col: 1, offset: 215},
			expr: &choiceExpr{
				pos: position{line: 13, col: 9, offset: 223},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 13, col: 9, offset: 223},
						name: "Content",
					},
					&ruleRefExpr{
						pos:  position{line: 13, col: 19, offset: 233},
						name: "Input",
					},
					&ruleRefExpr{
						pos:  position{line: 13, col: 27, offset: 241},
						name: "Text",
					},
				},
			},
		},
		{
			name: "Content",
			pos:  position{line: 15, col: 1, offset: 247},
			expr: &choiceExpr{
				pos: position{line: 15, col: 12, offset: 258},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 15, col: 12, offset: 258},
						run: (*parser).callonContent2,
						expr: &seqExpr{
							pos: position{line: 15, col: 12, offset: 258},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 15, col: 12, offset: 258},
									label: "tags",
									expr: &oneOrMoreExpr{
										pos: position{line: 15, col: 17, offset: 263},
										expr: &ruleRefExpr{
											pos:  position{line: 15, col: 17, offset: 263},
											name: "Tag",
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 15, col: 22, offset: 268},
									name: "OpenBrace",
								},
								&labeledExpr{
									pos:   position{line: 15, col: 32, offset: 278},
									label: "children",
									expr: &zeroOrMoreExpr{
										pos: position{line: 15, col: 41, offset: 287},
										expr: &ruleRefExpr{
											pos:  position{line: 15, col: 41, offset: 287},
											name: "Expr",
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 15, col: 47, offset: 293},
									name: "CloseBrace",
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 27, col: 5, offset: 601},
						run: (*parser).callonContent12,
						expr: &seqExpr{
							pos: position{line: 27, col: 5, offset: 601},
							exprs: []interface{}{
								&oneOrMoreExpr{
									pos: position{line: 27, col: 5, offset: 601},
									expr: &ruleRefExpr{
										pos:  position{line: 27, col: 5, offset: 601},
										name: "Tag",
									},
								},
								&notExpr{
									pos: position{line: 27, col: 10, offset: 606},
									expr: &ruleRefExpr{
										pos:  position{line: 27, col: 11, offset: 607},
										name: "OpenBrace",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 29, col: 5, offset: 727},
						run: (*parser).callonContent18,
						expr: &seqExpr{
							pos: position{line: 29, col: 5, offset: 727},
							exprs: []interface{}{
								&oneOrMoreExpr{
									pos: position{line: 29, col: 5, offset: 727},
									expr: &ruleRefExpr{
										pos:  position{line: 29, col: 5, offset: 727},
										name: "Tag",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 29, col: 10, offset: 732},
									name: "OpenBrace",
								},
								&zeroOrMoreExpr{
									pos: position{line: 29, col: 20, offset: 742},
									expr: &anyMatcher{
										line: 29, col: 20, offset: 742,
									},
								},
								&notExpr{
									pos: position{line: 29, col: 23, offset: 745},
									expr: &ruleRefExpr{
										pos:  position{line: 29, col: 24, offset: 746},
										name: "CloseBrace",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 31, col: 5, offset: 821},
						run: (*parser).callonContent27,
						expr: &seqExpr{
							pos: position{line: 31, col: 5, offset: 821},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 31, col: 5, offset: 821},
									name: "TagMarker",
								},
								&ruleRefExpr{
									pos:  position{line: 31, col: 15, offset: 831},
									name: "OpenBrace",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Text",
			pos:  position{line: 36, col: 1, offset: 920},
			expr: &actionExpr{
				pos: position{line: 36, col: 9, offset: 928},
				run: (*parser).callonText1,
				expr: &labeledExpr{
					pos:   position{line: 36, col: 9, offset: 928},
					label: "chunks",
					expr: &oneOrMoreExpr{
						pos: position{line: 36, col: 16, offset: 935},
						expr: &choiceExpr{
							pos: position{line: 36, col: 18, offset: 937},
							alternatives: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 36, col: 18, offset: 937},
									name: "UnescapedChars",
								},
								&ruleRefExpr{
									pos:  position{line: 36, col: 35, offset: 954},
									name: "EscapeSequence",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Input",
			pos:  position{line: 43, col: 1, offset: 1086},
			expr: &choiceExpr{
				pos: position{line: 43, col: 10, offset: 1095},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 43, col: 10, offset: 1095},
						run: (*parser).callonInput2,
						expr: &seqExpr{
							pos: position{line: 43, col: 10, offset: 1095},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 43, col: 10, offset: 1095},
									name: "TagMarker",
								},
								&litMatcher{
									pos:        position{line: 43, col: 20, offset: 1105},
									val:        "-",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 43, col: 24, offset: 1109},
									label: "nonFinalComponents",
									expr: &zeroOrMoreExpr{
										pos: position{line: 43, col: 43, offset: 1128},
										expr: &ruleRefExpr{
											pos:  position{line: 43, col: 43, offset: 1128},
											name: "NonFinalInputComponent",
										},
									},
								},
								&labeledExpr{
									pos:   position{line: 43, col: 67, offset: 1152},
									label: "finalComponent",
									expr: &ruleRefExpr{
										pos:  position{line: 43, col: 82, offset: 1167},
										name: "FinalInputComponent",
									},
								},
								&zeroOrOneExpr{
									pos: position{line: 43, col: 102, offset: 1187},
									expr: &ruleRefExpr{
										pos:  position{line: 43, col: 102, offset: 1187},
										name: "SpaceChomper",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 47, col: 5, offset: 1318},
						run: (*parser).callonInput13,
						expr: &seqExpr{
							pos: position{line: 47, col: 5, offset: 1318},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 47, col: 5, offset: 1318},
									name: "TagMarker",
								},
								&litMatcher{
									pos:        position{line: 47, col: 15, offset: 1328},
									val:        "-",
									ignoreCase: false,
								},
								&notExpr{
									pos: position{line: 47, col: 19, offset: 1332},
									expr: &ruleRefExpr{
										pos:  position{line: 47, col: 20, offset: 1333},
										name: "AnyInputComponent",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 49, col: 5, offset: 1408},
						run: (*parser).callonInput19,
						expr: &seqExpr{
							pos: position{line: 49, col: 5, offset: 1408},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 49, col: 5, offset: 1408},
									name: "TagMarker",
								},
								&litMatcher{
									pos:        position{line: 49, col: 15, offset: 1418},
									val:        "-",
									ignoreCase: false,
								},
								&oneOrMoreExpr{
									pos: position{line: 49, col: 19, offset: 1422},
									expr: &ruleRefExpr{
										pos:  position{line: 49, col: 19, offset: 1422},
										name: "NonFinalInputComponent",
									},
								},
								&litMatcher{
									pos:        position{line: 49, col: 43, offset: 1446},
									val:        ".",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 51, col: 5, offset: 1516},
						run: (*parser).callonInput26,
						expr: &seqExpr{
							pos: position{line: 51, col: 5, offset: 1516},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 51, col: 5, offset: 1516},
									name: "TagMarker",
								},
								&litMatcher{
									pos:        position{line: 51, col: 15, offset: 1526},
									val:        "-",
									ignoreCase: false,
								},
								&oneOrMoreExpr{
									pos: position{line: 51, col: 19, offset: 1530},
									expr: &ruleRefExpr{
										pos:  position{line: 51, col: 19, offset: 1530},
										name: "NonFinalInputComponent",
									},
								},
								&notExpr{
									pos: position{line: 51, col: 43, offset: 1554},
									expr: &ruleRefExpr{
										pos:  position{line: 51, col: 44, offset: 1555},
										name: "FinalInputComponent",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "SpaceChomper",
			pos:  position{line: 55, col: 1, offset: 1716},
			expr: &seqExpr{
				pos: position{line: 55, col: 17, offset: 1732},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 55, col: 17, offset: 1732},
						val:        "~",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 55, col: 21, offset: 1736},
						name: "_",
					},
				},
			},
		},
		{
			name:        "AnyInputComponent",
			displayName: "\"path component\"",
			pos:         position{line: 57, col: 1, offset: 1739},
			expr: &choiceExpr{
				pos: position{line: 57, col: 39, offset: 1777},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 57, col: 39, offset: 1777},
						name: "NonFinalInputComponent",
					},
					&ruleRefExpr{
						pos:  position{line: 57, col: 64, offset: 1802},
						name: "FinalInputComponent",
					},
				},
			},
		},
		{
			name:        "NonFinalInputComponent",
			displayName: "\"path component\"",
			pos:         position{line: 59, col: 1, offset: 1823},
			expr: &actionExpr{
				pos: position{line: 59, col: 44, offset: 1866},
				run: (*parser).callonNonFinalInputComponent1,
				expr: &seqExpr{
					pos: position{line: 59, col: 44, offset: 1866},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 59, col: 44, offset: 1866},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 59, col: 49, offset: 1871},
								name: "Identifier",
							},
						},
						&litMatcher{
							pos:        position{line: 59, col: 60, offset: 1882},
							val:        ".",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name:        "FinalInputComponent",
			displayName: "\"last path component\"",
			pos:         position{line: 63, col: 1, offset: 1919},
			expr: &actionExpr{
				pos: position{line: 63, col: 46, offset: 1964},
				run: (*parser).callonFinalInputComponent1,
				expr: &seqExpr{
					pos: position{line: 63, col: 46, offset: 1964},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 63, col: 46, offset: 1964},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 63, col: 51, offset: 1969},
								name: "Identifier",
							},
						},
						&notExpr{
							pos: position{line: 63, col: 62, offset: 1980},
							expr: &litMatcher{
								pos:        position{line: 63, col: 63, offset: 1981},
								val:        ".",
								ignoreCase: false,
							},
						},
					},
				},
			},
		},
		{
			name: "EscapeSequence",
			pos:  position{line: 67, col: 1, offset: 2018},
			expr: &actionExpr{
				pos: position{line: 67, col: 19, offset: 2036},
				run: (*parser).callonEscapeSequence1,
				expr: &seqExpr{
					pos: position{line: 67, col: 19, offset: 2036},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 67, col: 19, offset: 2036},
							name: "TagMarker",
						},
						&labeledExpr{
							pos:   position{line: 67, col: 29, offset: 2046},
							label: "escapedChar",
							expr: &ruleRefExpr{
								pos:  position{line: 67, col: 42, offset: 2059},
								name: "EscapedChar",
							},
						},
					},
				},
			},
		},
		{
			name: "UnescapedChars",
			pos:  position{line: 71, col: 1, offset: 2119},
			expr: &actionExpr{
				pos: position{line: 71, col: 19, offset: 2137},
				run: (*parser).callonUnescapedChars1,
				expr: &oneOrMoreExpr{
					pos: position{line: 71, col: 19, offset: 2137},
					expr: &seqExpr{
						pos: position{line: 71, col: 20, offset: 2138},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 71, col: 20, offset: 2138},
								expr: &ruleRefExpr{
									pos:  position{line: 71, col: 21, offset: 2139},
									name: "EscapedChar",
								},
							},
							&anyMatcher{
								line: 71, col: 33, offset: 2151,
							},
						},
					},
				},
			},
		},
		{
			name: "EscapedChar",
			pos:  position{line: 75, col: 1, offset: 2189},
			expr: &choiceExpr{
				pos: position{line: 75, col: 16, offset: 2204},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 75, col: 16, offset: 2204},
						name: "TagMarker",
					},
					&ruleRefExpr{
						pos:  position{line: 75, col: 28, offset: 2216},
						name: "CloseBrace",
					},
				},
			},
		},
		{
			name: "Tag",
			pos:  position{line: 77, col: 1, offset: 2228},
			expr: &actionExpr{
				pos: position{line: 77, col: 8, offset: 2235},
				run: (*parser).callonTag1,
				expr: &seqExpr{
					pos: position{line: 77, col: 8, offset: 2235},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 77, col: 8, offset: 2235},
							name: "TagMarker",
						},
						&labeledExpr{
							pos:   position{line: 77, col: 18, offset: 2245},
							label: "tagName",
							expr: &ruleRefExpr{
								pos:  position{line: 77, col: 26, offset: 2253},
								name: "Identifier",
							},
						},
					},
				},
			},
		},
		{
			name: "Identifier",
			pos:  position{line: 81, col: 1, offset: 2300},
			expr: &actionExpr{
				pos: position{line: 81, col: 15, offset: 2314},
				run: (*parser).callonIdentifier1,
				expr: &oneOrMoreExpr{
					pos: position{line: 81, col: 15, offset: 2314},
					expr: &charClassMatcher{
						pos:        position{line: 81, col: 15, offset: 2314},
						val:        "[A-Z0-9]i",
						ranges:     []rune{'a', 'z', '0', '9'},
						ignoreCase: true,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "TagMarker",
			pos:  position{line: 88, col: 1, offset: 2551},
			expr: &litMatcher{
				pos:        position{line: 88, col: 14, offset: 2564},
				val:        "\x01",
				ignoreCase: false,
			},
		},
		{
			name: "OpenBrace",
			pos:  position{line: 89, col: 1, offset: 2571},
			expr: &litMatcher{
				pos:        position{line: 89, col: 14, offset: 2584},
				val:        "\x02",
				ignoreCase: false,
			},
		},
		{
			name: "CloseBrace",
			pos:  position{line: 90, col: 1, offset: 2591},
			expr: &litMatcher{
				pos:        position{line: 90, col: 15, offset: 2605},
				val:        "\x03",
				ignoreCase: false,
			},
		},
		{
			name:        "_",
			displayName: "\"whitespace\"",
			pos:         position{line: 92, col: 1, offset: 2613},
			expr: &zeroOrMoreExpr{
				pos: position{line: 92, col: 18, offset: 2632},
				expr: &charClassMatcher{
					pos:        position{line: 92, col: 18, offset: 2632},
					val:        "[ \\t\\r\\n]",
					chars:      []rune{' ', '\t', '\r', '\n'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 93, col: 1, offset: 2643},
			expr: &notExpr{
				pos: position{line: 93, col: 8, offset: 2650},
				expr: &anyMatcher{
					line: 93, col: 9, offset: 2651,
				},
			},
		},
	},
}

func (c *current) onRoot2(exprList interface{}) (interface{}, error) {
	return makeNodeSlice(exprList), nil
}

func (p *parser) callonRoot2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRoot2(stack["exprList"])
}

func (c *current) onRoot8() (interface{}, error) {
	return nil, fmt.Errorf("Unexpected close brace")
}

func (p *parser) callonRoot8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRoot8()
}

func (c *current) onContent2(tags, children interface{}) (interface{}, error) {
	tagNames := makeStringSlice(tags)
	childNodes := makeNodeSlice(children)
	var topNode *ContentNode
	for i := len(tagNames) - 1; i >= 0; i-- {
		topNode = &ContentNode{
			Tag:      tagNames[i],
			Children: childNodes,
		}
		childNodes = []Node{topNode}
	}
	return topNode, nil
}

func (p *parser) callonContent2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onContent2(stack["tags"], stack["children"])
}

func (c *current) onContent12() (interface{}, error) {
	return nil, fmt.Errorf("Expected open brace for tag. Escape the tag marker by doubling it, like @@.")
}

func (p *parser) callonContent12() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onContent12()
}

func (c *current) onContent18() (interface{}, error) {
	return nil, fmt.Errorf("Expected close brace for tag.")
}

func (p *parser) callonContent18() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onContent18()
}

func (c *current) onContent27() (interface{}, error) {
	return nil, fmt.Errorf("Expected a tag name")
}

func (p *parser) callonContent27() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onContent27()
}

func (c *current) onText1(chunks interface{}) (interface{}, error) {
	chunkSlice := makeStringSlice(chunks)
	return &TextNode{
		Text: strings.Join(chunkSlice, ""),
	}, nil
}

func (p *parser) callonText1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onText1(stack["chunks"])
}

func (c *current) onInput2(nonFinalComponents, finalComponent interface{}) (interface{}, error) {
	return &InputNode{
		Path: append(makeStringSlice(nonFinalComponents), finalComponent.(string)),
	}, nil
}

func (p *parser) callonInput2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInput2(stack["nonFinalComponents"], stack["finalComponent"])
}

func (c *current) onInput13() (interface{}, error) {
	return nil, fmt.Errorf("Expected an input path")
}

func (p *parser) callonInput13() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInput13()
}

func (c *current) onInput19() (interface{}, error) {
	return nil, fmt.Errorf("Unexpected period in input path")
}

func (p *parser) callonInput19() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInput19()
}

func (c *current) onInput26() (interface{}, error) {
	return nil, fmt.Errorf("Unexpected trailing period on input path.\nTip: End a sentence with an input using a chomp, like `@-foo~ .`")
}

func (p *parser) callonInput26() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInput26()
}

func (c *current) onNonFinalInputComponent1(name interface{}) (interface{}, error) {
	return name.(string), nil
}

func (p *parser) callonNonFinalInputComponent1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNonFinalInputComponent1(stack["name"])
}

func (c *current) onFinalInputComponent1(name interface{}) (interface{}, error) {
	return name.(string), nil
}

func (p *parser) callonFinalInputComponent1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFinalInputComponent1(stack["name"])
}

func (c *current) onEscapeSequence1(escapedChar interface{}) (interface{}, error) {
	return string(escapedChar.([]byte)), nil
}

func (p *parser) callonEscapeSequence1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEscapeSequence1(stack["escapedChar"])
}

func (c *current) onUnescapedChars1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonUnescapedChars1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnescapedChars1()
}

func (c *current) onTag1(tagName interface{}) (interface{}, error) {
	return tagName.(string), nil
}

func (p *parser) callonTag1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTag1(stack["tagName"])
}

func (c *current) onIdentifier1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonIdentifier1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifier1()
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errNoMatch is returned if no match could be found.
	errNoMatch = errors.New("no match found")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match
}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos        position
	val        string
	chars      []rune
	ranges     []rune
	classes    []*unicode.RangeTable
	ignoreCase bool
	inverted   bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner  error
	pos    position
	prefix string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
	}
	p.setOptions(opts)
	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	recover bool
	debug   bool
	depth   int

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// stats
	exprCnt int
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position)
}

func (p *parser) addErrAt(err error, pos position) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, pos: pos, prefix: buf.String()}
	p.errs.add(pe)
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n == 1 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// make sure this doesn't go out silently
			p.addErr(errNoMatch)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint
	var ok bool

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.exprCnt++
	var val interface{}
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position)
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	// can't match EOF
	if cur == utf8.RuneError {
		return nil, false
	}
	start := p.pt
	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(not.expr)
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	var vals []interface{}

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}

func rangeTable(class string) *unicode.RangeTable {
	if rt, ok := unicode.Categories[class]; ok {
		return rt
	}
	if rt, ok := unicode.Properties[class]; ok {
		return rt
	}
	if rt, ok := unicode.Scripts[class]; ok {
		return rt
	}

	// cannot happen
	panic(fmt.Sprintf("invalid Unicode class: %s", class))
}
