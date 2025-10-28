mod common_syntaxes {
    use lat::{LatValue, parse};

    #[test]
    fn test() {
        const INPUT: &str = "";
        let parsed_result: LatValue = parse(INPUT);

        assert_eq!(parsed_result, LatValue {});
    }
}
