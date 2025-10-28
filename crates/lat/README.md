# `lat` - A Latin-esque Spelling's Language Parser

`lat` is a Rust crate that provides a parser for a custom, Latin-inspired spell language. It is designed for game development, particularly for high-paced games that require both fast, simple commands ("Quick Casts") and more complex, powerful ones ("Full Casts").

This parser takes a spell string (`&str`) and transforms it into a structured `LatValue` data type. Your game engine can then interpret this structure to evaluate the spell's final effect, cost, and casting time.

______________________________________________________________________

## 🚀 Usage

The primary API is the `lat::parse` function, which takes a string slice and returns a `Result<LatValue, LatError>`.

```rust
use lat::parser::parse;
use lat::types::{
    LatValue, Action, Element, Modifier, Target, Origin, Emphasis, CoreComponent,
    errors::LatError,
};

fn main() {
    // 1. Quick Cast (Core only)
    // "Burn with great fire"
    let input_quick = "Ure igne magno";
    
    match parse(input_quick) {
        Ok(spell) => {
            assert_eq!(spell.action, Action::Ure);
            assert_eq!(spell.element, Element::Ignis);
            assert_eq!(spell.modifier, Modifier::Magnus);
            assert_eq!(spell.target, None);
            assert!(spell.emphasis_phrases.is_empty());
        },
        Err(e) => panic!("Parse failed: {}", e),
    }

    // 2. Full Cast (Core + Extensions)
    // "Now, from the sky, burn with great fire at the enemy"
    let input_full = "Nunc ex caelo ure igne magno ad hostem";

    match parse(input_full) {
        Ok(spell) => {
            assert_eq!(spell.action, Action::Ure);
            assert_eq!(spell.element, Element::Ignis);
            assert_eq!(spell.modifier, Modifier::Magnus);
            assert_eq!(spell.target, Some(Target::Hostis));
            assert_eq!(spell.origin, Some(Origin::Caelum));
            assert_eq!(spell.emphasis_phrases, vec![Emphasis::Nunc]);
            assert_eq!(spell.source_text, input_full);
        },
        Err(e) => panic!("Parse failed: {}", e),
    }

    // 3. Grammar Error
    // This example is missing the Modifier component
    let input_error = "Ure igne ad hostem";

    match parse(input_error) {
        Ok(_) => panic!("Should have failed"),
        Err(e) => {
            println!("Caught expected error: {}", e);
            // Example error based on a full parser implementation:
            // "Parse Error: Unexpected component. Expected Modifier but found Ad."
        }
    }
}
```

______________________________________________________________________

## 📜 Grammar

The `lat` grammar is built on two concepts: the **Spell Core** (for Quick Casts) and optional **Extension Phrases** (for Full Casts). All words are space-delimited.

### 1. Spell Core (Quick Cast)

This is the minimum set of words required to form a valid spell. It consists of three components in a **fixed order**:

**Syntax:** `[Action] + [Element] + [Modifier]`

- **`[Action]`**: The verb, or what the spell does (e.g., `Ure`).
- **`[Element]`**: The noun, or what element/medium is used (e.g., `igne`).
- **`[Modifier]`**: The adjective, or the scale/form of the spell (e.g., `magno`).

**Example:**

- `Sana aqua celeri` -> (Heal with fast water)
- `Feri umbra forte` -> (Strike with strong shadow)

### 2. Full Cast (Extended Spell)

A Full Cast adds optional "Extension Phrases" to the Spell Core. These phrases add power, specify targets, or add other effects. They can appear **before or after** the Spell Core block.

**Syntax:** `[Extensions...] + [Spell Core] + [Extensions...]`

**Extension Phrase Types:**

- **Target:** A preposition (like `ad`) followed by a target noun.
  - Example: `ad hostem` (at the enemy)
- **Origin:** A preposition (like `ex` or `de`) followed by an origin noun.
  - Example: `ex caelo` (from the sky)
- **Emphasis:** A standalone adverb that modifies the entire spell.
  - Example: `Nunc` (Now), `CumVi` (With force)

**Full Cast Example:**

- `Nunc ex caelo ure igne magno ad hostem`
- **Extensions (Prefix):** `Nunc`, `ex caelo`
- **Spell Core:** `ure igne magno`
- **Extension (Postfix):** `ad hostem`

______________________________________________________________________

## 📦 Data Structures

The parser's output is one of the following two types.

### `LatValue` (Success)

This `struct` is returned on a successful parse. It contains the complete AST (Abstract Syntax Tree) for the spell, which your game engine can then evaluate.

```rust
/// Represents the final parsed result of an entire spell string (the AST).
/// This structure is passed from the `lat` parser to the game engine for evaluation.
#[derive(Debug, PartialEq, Eq, Clone)]
pub struct LatValue {
    // --- Spell Core (Mandatory) ---
    
    /// The core action (verb) of the spell.
    pub action: Action,
    /// The core element or medium (noun) of the spell.
    pub element: Element,
    /// The core modifier (adjective) of the spell.
    pub modifier: Modifier,
    
    // --- Spell Extensions (Optional) ---
    
    /// The optional target of the spell (e.g., "ad hostem").
    pub target: Option<Target>,
    /// The optional origin of the spell (e.g., "ex caelo").
    pub origin: Option<Origin>,
    
    /// A list of detected emphasis phrases (e.g., "Nunc", "CumVi").
    /// The game engine interprets this list to modify spell power, cost, etc.
    pub emphasis_phrases: Vec<Emphasis>,
    
    /// The original, unmodified source text of the parsed spell.
    /// Primarily used for debugging.
    pub source_text: String,
}
```

### `LatError` (Failure)

If the input string violates the grammar, a `LatError` is returned. It uses `thiserror` for convenient error messages.

```rust
/// Represents all possible errors that can occur during parsing by the `lat` crate.
#[derive(Error, Debug, Clone, PartialEq, Eq)]
pub enum LatError {
    /// An unknown word (not in the dictionary) was encountered.
    #[error("Parse Error: Unknown word \"{0}\".")]
    UnknownWord(String),
    
    /// A required component for the spell core (Quick Cast) is missing.
    #[error("Parse Error: Missing core spell component {0:?}.")]
    MissingCoreComponent(CoreComponent),
    
    /// A component was detected in an unexpected order.
    #[error("Parse Error: Unexpected component. Expected {expected} but found {found}.")]
    UnexpectedComponent {
        expected: String,
        found: String,
    },
    
    /// An extension phrase (e.g., "Ad") was started but not followed by its corresponding part.
    #[error("Parse Error: Incomplete extension phrase starting with \"{0}\".")]
    IncompleteExtensionPhrase(String),

    /// The input ended prematurely when more tokens were expected.
    #[error("Parse Error: Unexpected end of input.")]
    UnexpectedEof,
}

/// Details for the `MissingCoreComponent` error.
#[derive(Debug, Clone, PartialEq, Eq)]
pub enum CoreComponent {
    Action,
    Element,
    Modifier,
}
```

______________________________________________________________________

## 📖 Dictionary

This is the dictionary of words the parser currently recognizes.

### 1. Action (Verbs)

- `Ure`: Burn
- `Sana`: Heal
- `Defende`: Defend
- `Feri`: Strike

### 2. Element (Nouns)

- `igne`: (with) Fire
- `aqua`: (with) Water
- `vento`: (with) Wind
- `terra`: (with) Earth
- `luce`: (with) Light
- `umbra`: (with) Shadow

*(Note: These words correspond to `Element::Ignis`, `Element::Aqua`, etc.)*

### 3. Modifier (Adjectives)

- `magno` / `magna`: Great, Large (maps to `Modifier::Magnus`)
- `parvo` / `parva`: Small (maps to `Modifier::Parvus`)
- `celeri`: Fast (maps to `Modifier::Celer`)
- `forte`: Strong (maps to `Modifier::Fortis`)

*(Note: The parser can be configured to treat `magno` and `magna` as aliases, or require them based on the Element's grammatical gender for stricter parsing.)*

### 4. Target (Preposition: `ad`)

- `ad hostem`: (at the) Enemy (`Target::Hostis`)
- `ad amicum`: (at the) Ally (`Target::Amicus`)
- `ad me`: (at) Self (`Target::Me`)
- `ad aream`: (at the) Area (`Target::Area`)

### 5. Origin (Prepositions: `ex`, `de`)

- `ex caelo`: (from the) Sky (`Origin::Caelum`)
- `de terra`: (from the) Earth (`Origin::Terra`)

### 6. Emphasis (Adverbs)

- `Nunc`: Now (`Emphasis::Nunc`)
- `CumVi`: With force (`Emphasis::CumVi`)
- `Tandem`: At last (`Emphasis::Tandem`)
