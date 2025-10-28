pub mod errors;

/// Represents the core action (verb) of the spell.
/// This is a mandatory part of the "Spell Core".
#[derive(Debug, PartialEq, Eq, Clone, Copy)]
pub enum Action {
    /// Burn
    Ure,
    /// Heal
    Sana,
    /// Defend
    Defende,
    /// Strike
    Feri,
}

/// Represents the core element or medium (noun) of the spell.
/// This is a mandatory part of the "Spell Core".
#[derive(Debug, PartialEq, Eq, Clone, Copy)]
pub enum Element {
    /// Fire
    Ignis,
    /// Water
    Aqua,
    /// Wind
    Ventus,
    /// Earth
    Terra,
    /// Light
    Lumen,
    /// Shadow
    Umbra,
}

/// Represents the core modifier (adjective) describing the scale or form of the spell.
/// This is a mandatory part of the "Spell Core".
#[derive(Debug, PartialEq, Eq, Clone, Copy)]
pub enum Modifier {
    /// Large / Great
    Magnus,
    /// Small
    Parvus,
    /// Fast / Quick
    Celer,
    /// Strong
    Fortis,
}

/// Represents an optional target (extension phrase) for the spell.
/// Used in "Full Casts" to specify where the spell goes.
#[derive(Debug, PartialEq, Eq, Clone, Copy)]
pub enum Target {
    /// Enemy
    Hostis,
    /// Ally / Friend
    Amicus,
    /// Self
    Me,
    /// Area
    Area,
}

/// Represents an optional origin point (extension phrase) for the spell.
/// Used in "Full Casts" to specify where the spell comes from.
#[derive(Debug, PartialEq, Eq, Clone, Copy)]
pub enum Origin {
    /// From the sky
    Caelum,
    /// From the earth
    Terra,
}

/// Represents an optional emphasis phrase (extension phrase) for the spell.
/// The parser identifies these, but the game engine evaluates their effect (e.g., increasing power).
#[derive(Debug, PartialEq, Eq, Clone, Copy)]
pub enum Emphasis {
    /// Now
    Nunc,
    /// With force
    CumVi,
    /// At last / Finally
    Tandem,
}

/// Represents the final parsed result of an entire spell string (the Abstract Syntax Tree).
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
