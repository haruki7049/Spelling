wit_bindgen::generate!({
    world: "lat",
});

pub struct Lat;

impl Guest for Lat {
    fn parse(inputs: String) -> Result<LatValue, ParseError> {
        dbg!(inputs);
        todo!();
    }
}

export!(Lat);
