---
title: CSS Test Post
date: 2026-03-24
description: A test of every Pandoc markdown element.
---

# Heading 1

## Heading 2

### Heading 3

#### Heading 4

##### Heading 5

###### Heading 6

---

## Paragraphs & Inline

Plain paragraph text. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus lacinia odio vitae vestibulum habitasse platea dictumst.

**Bold text**, *italic text*, ***bold and italic***, ~~strikethrough~~, and `inline code`.

Superscript: 2^10^ = 1024. Subscript: H~2~O.

A line with a [link to the archive](/archive.html) and a [visited link](/archive.html).

---

## Blockquote

> Poor restless dove, I pity thee;
> And when I hear thy plaintive moan,
> I mourn for thy captivity,
> And in thy woes forget mine own.
>
> — Anne Brontë

---

## Unordered List

- first item
- second item
- third item
    - nested item
    - another nested item
- fourth item

---

## Ordered List

1. first item
2. second item
3. third item
    1. nested item
    2. another nested item
4. fourth item

---

## Code Block

```haskell
postCtx :: Context String
postCtx =
    dateField "date" "%d %b %Y" `mappend`
    defaultContext
```

```rust
fn main() {
    println!("Hello, world!");
}
```

---

## Table

| Header One   | Header Two   | Header Three |
|:-------------|:------------:|-------------:|
| left         | center       | right        |
| aligned      | aligned      | aligned      |
| row three    | row three    | row three    |

---

## Image

![A large test image](https://cdn.britannica.com/13/77413-050-95217C0B/Golden-Gate-Bridge-San-Francisco.jpg)

---

## Horizontal Rule

Above

---

Below

---

## Footnotes

This sentence has a footnote.[^1] This one has another.[^2] And a third.[^3]

[^1]: This is the first footnote.
[^2]: This is the second footnote, with a [link](https://jackiedor.land) inside it.
[^3]: This is the third footnote.

---

## Definition List

Term one
:   Definition of term one.

Term two
:   Definition of term two, which is a bit longer and might wrap depending on viewport width.

---

## Task List

- [x] completed task
- [ ] incomplete task
- [x] another completed task

---

## Nested Blockquote

> This is an outer blockquote.
>
> > This is a nested blockquote inside it.
>
> Back to the outer blockquote.

---

## Mixed List

1. ordered item one
2. ordered item two
    - unordered nested
    - unordered nested
3. ordered item three

---

## Long Paragraph (wrapping test)

Sublime Plaza on 35th & Park was a curious office building. Everything about it was just a little bit *off*. Between the overhang that was a little too long on one side and the somewhat tacky mural of a bewildered Jean Piaget, Sublime Plaza didn't look like a classic medical practice in any way. Despite being filled with doctors and psychiatrists, Sublime Plaza was *almost* dingy. The exterior was clean, but not perfectly so. Some weeds, here and there, would always be growing between the cracks of the sidewalk, waiting to be cut yet again. The most accurate definition of Sublime Plaza, according to the surrounding residents, would be **"a curious building where curious professionals do their curious work."**
