# Vivian Tech Details

This document describes the internal details of how the Vivian template language is parsed. It is not intended to be documentation for Vivian's usage.

# Stages

There are several stages in the Vivian parsing pipeline, which go beyond the simple grammar described in grammar.peg.

## Header

Because Vivian is meant to be flexible enough to use with many types of arbitrary static copy, some escaping of delimiters is inevitable. However, any choice of delimiters will tend to be particularly annoying with text related to some domain. The most trivial example of this would be laying out text _about_ Vivian itself in a Vivian template.

To circumvent this, the first 3-4 characters of a template can specify a "header". This is denoted by a leading `@`, followed by an optional tag marker (if omitted, `@` remains the tag marker), and then an open and close brace. The characters chosen will then be used throughout the template, including for escapes.

The header is handled as the very first step. Once the delimiters are chosen (or left alone), the entire string is processed to replace those characters with the following ascii characters:

1. Tag marker: 0x01
2. Open brace: 0x02
3. Close brace: 0x03

The delimiters are then recorded in the AST, since they will be needed after parsing to restore any escaped instnaces of those characters.

## PEG Parser

The grammar in grammar.peg specifies what might be called "loose Vivian". In fact, you might consider this to be a sort of "heavyweight tokenization" rather than a full parse. It produces an AST, but additional processing is needed before this AST can be used to render anything.

In particular, this parser is completely tag-agnostic. As far as it's concerned, any alphanumeric string is a valid tag name. As such, it also has no notion of the rules governing tag nesting. It will happily allow you to subdivide a cell-creating tag, for example, even though this is not actually valid.

## Hoisting

Vivid uses a non-nested cell-based layout, but there is no one-to-one correspondence between cells and tags. For example `@wrap[foo]` describes a single cell, as does `@fixed50[bar]`. But `@wrap[@fixed50[baz]]` (and its shorthand form `@wrap@fixed50[baz]`) also describes a single cell.

In other words, a cell is described by an _uninterrupted lineage_ of "cell creator" tags, all of which are only-children except for the root. A Vivian tree in which all cell creating nodes obey this property can be considered "normal Vivian", and this is the format we want to work with for rendering. Once a template is in this form, we can split it into distinct cells, each of which only have to worry about local styling rather than layout.

However, this form is not always very convenient to write. For example, if you want to give an entire row a white background, the natural way to write that would be, for example:

    @bgWhite[@fixed[Name:] @max30@green[@-name] @fixed[Age:] @fixed3[@-age]]

However, the normal form of this template is a mess:

    @bgWhite[@fixed[Name:]] @bgWhite[@max30@green[@-name]] @bgWhite[@fixed[Age:]] @bgWhite[@fixed3[@-age]]]

There's clear value in allowing unrestricted styling at the top level since, by definition, styling tags don't affect layout.

The only property that the template _must_ obey is that cell creators never "fork". That is, cell creators cannot be siblings within the lineage, nor can they have style siblings aside from at the root level. For example, `@wrap[@fixed[foo]]` is fine, but `@wrap[@fixed[foo] @min50[bar] @green[baz]]` is not, since it implies nesting of cells.

To enable non-normal Vivian templates, we normalize them through the process of "hoisting". Style nodes are split up, and their pieces are injected into their child cell creators. This process recurses from the leaves of the AST. We also validate the "no forking lineages" rule, rejecting the AST if any exceptions are found. At this time, we also discard any top-level runs of style nodes whose text consists of only whitespace.

After hoisting, the AST has a normal _structure_, but it may still be invalid, because tag names are still unconstrained.

## Cell production

The next step in the pipeline is to take all cell creator lineages (and bare style nodes at the root level) and turn them into cells. This is the first stage at which we actually parse tag names. The parse rule is extremely basic: a cell creator tag always consists of an alphabetic name followed by zero or more digits.

Cell creator lineages each become one cell, with all of the tag attributes being applied to that cell. Consecutive runs of style nodes at the top level are coalesced into their own cells with default layout properties.

After this stage, the cell creator tags are no longer relevant; their information is encoded in the cells, and from this point on, we are only concerned with styling the contents of those cells, and then rendering them.

## Styling

Finally, we create the styled text which lives inside the cells. At this point, we finally validate the style tags. The following tags are valid:

-   bold/b
-   italic/i
-   underline/u
-   Any html color name (e.g. green/blue/salmon/chartreuse). See [here](golang.org/x/image/colornames) for a full listing (case insensitive)
-   Any three or six digit hexadecimal RGB color (case insensitive)
-   "bg" followed by an html color name or hex color as above

# Future work

Validating the style tags only at render time is not ideal. We should be able to validate a template without having content to put inside it. This could be done rather inefficiently by simply passing a nil context to the template and asking it to render, but ideally, we'd do this as its own step.
