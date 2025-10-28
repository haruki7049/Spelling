use crate::types::{LatValue, errors::LatError};

pub fn parse(_input: &str) -> Result<LatValue, LatError> {
    Err(LatError::UnexpectedEof)
}
