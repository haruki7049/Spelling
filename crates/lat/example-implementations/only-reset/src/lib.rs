// Generate bindings from the WIT file
wit_bindgen::generate!({
    world: "lat",
    generate_all,
});

struct MyLat;

impl exports::haruki7049::lat::parser::Guest for MyLat {
    fn parse(inputs: String) -> Result<haruki7049::lat::types::LatValue, haruki7049::lat::types::ParseError> {
        if inputs == "reset" {
            Ok(haruki7049::lat::types::LatValue { resetting: true })
        } else if inputs.is_empty() {
            Err(haruki7049::lat::types::ParseError {
                value: inputs,
                pos: haruki7049::lat::types::Position { line: 1, column: 0 },
                message: "Input cannot be empty".to_string(),
            })
        } else {
            Ok(haruki7049::lat::types::LatValue { resetting: false })
        }
    }
}

export!(MyLat);
