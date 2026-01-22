wasmtime::component::bindgen!({
    world: "lat",
    additional_derives: [Clone],
});

unsafe impl std::marker::Sync for Lat {}
unsafe impl std::marker::Send for Lat {}
impl bevy::ecs::resource::Resource for Lat {}
impl bevy::ecs::component::Component for Lat {
    type Mutability = bevy::ecs::component::Mutable;
    const STORAGE_TYPE: bevy::ecs::component::StorageType = bevy::ecs::component::StorageType::Table;
}

pub use haruki7049::lat::types::{ParseError, LatValue};
