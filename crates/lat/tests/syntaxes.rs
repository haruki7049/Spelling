mod common_syntaxes {
    use lat::{
        parser::parse,
        types::{Action, Element, LatValue, Modifier},
    };

    #[test]
    fn test() -> Result<(), Box<dyn std::error::Error>> {
        const INPUT: &str = "Ure ignis magno";
        let parsed_result: LatValue = parse(INPUT)?;

        assert_eq!(
            parsed_result,
            LatValue {
                action: Action::Ure,
                element: Element::Ignis,
                modifier: Modifier::Magnus,
                target: None,
                origin: None,
                emphasis_phrases: vec![],
                source_text: String::from("Ure ignis magno"),
            }
        );

        Ok(())
    }
}
