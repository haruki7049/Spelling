use thiserror::Error;

/// Details for the `MissingCoreComponent` error.
#[derive(Debug, Clone, PartialEq, Eq)]
pub enum CoreComponent {
    /// Action (Verb)
    Action,
    /// Element (Noun)
    Element,
    /// Modifier (Adjective)
    Modifier,
}

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
    UnexpectedComponent { expected: String, found: String },

    /// An extension phrase (e.g., "Ad") was started but not followed by its corresponding part.
    #[error("Parse Error: Incomplete extension phrase starting with \"{0}\".")]
    IncompleteExtensionPhrase(String),

    /// The input ended prematurely when more tokens were expected.
    #[error("Parse Error: Unexpected end of input.")]
    UnexpectedEof,
}
