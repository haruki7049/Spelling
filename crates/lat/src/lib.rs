wit_bindgen::generate!({
    world: "lat",
    additional_derives: [PartialEq, Eq, Hash, Clone],
});

pub struct Lat;

impl Guest for Lat {
    fn parse(inputs: String) -> Result<LatValue, ParseError> {
        Ok(LatValue::new(inputs))
    }
}

impl LatValue {
    pub fn new(value: String) -> Self {
        Self { value }
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
