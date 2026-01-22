wit_bindgen::generate!({
    world: "lat",
    additional_derives: [PartialEq, Eq, Hash, Clone],
});

pub struct Lat;

impl Guest for Lat {
    fn parse(inputs: String) -> Result<LatValue, ParseError> {
        Ok(LatValue::new(inputs == "reset"))
    }
}

impl LatValue {
    pub fn new(resetting: bool) -> Self {
        Self { resetting }
    }
}

export!(Lat);

#[cfg(test)]
mod tests {
    use super::{Guest, Lat, LatValue};

    #[test]
    fn lat_parse() {
        let inputs: String = "HELLO".to_string();

        assert_eq!(Lat::parse(inputs.clone()), Ok(LatValue::new(inputs)));
    }
}
